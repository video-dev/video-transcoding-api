package bitmovin

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/bitmovin/bitmovin-go/services"
	"net/url"
)

// Name is the name used for registering the bitmovin provider in the
// registry of providers.
const Name = "bitmovin"

const bitmovinAPIErrorMsg = "ERROR"

// Just to double check the interface is properly implemented
var _ provider.TranscodingProvider = (*bitmovinProvider)(nil)

var cloudRegions = map[string]struct{}{
	"AWS_US_EAST_1":        {},
	"AWS_US_WEST_1":        {},
	"AWS_US_WEST_2":        {},
	"AWS_EU_WEST_1":        {},
	"AWS_EU_CENTRAL_1":     {},
	"AWS_AP_SOUTH_1":       {},
	"AWS_AP_NORTHEAST_1":   {},
	"AWS_AP_NORTHEAST_2":   {},
	"AWS_AP_SOUTHEAST_1":   {},
	"AWS_AP_SOUTHEAST_2":   {},
	"AWS_SA_EAST_1":        {},
	"GOOGLE_EUROPE_WEST_1": {},
	"GOOGLE_US_EAST_1":     {},
	"GOOGLE_US_CENTRAL_1":  {},
	"GOOGLE_ASIA_EAST_1":   {},
}

var awsCloudRegions = map[string]struct{}{
	"US_EAST_1":      {},
	"US_WEST_1":      {},
	"US_WEST_2":      {},
	"EU_WEST_1":      {},
	"EU_CENTRAL_1":   {},
	"AP_SOUTH_1":     {},
	"AP_NORTHEAST_1": {},
	"AP_NORTHEAST_2": {},
	"AP_SOUTHEAST_1": {},
	"AP_SOUTHEAST_2": {},
	"SA_EAST_1":      {},
}

func init() {
	provider.Register(Name, bitmovinFactory)
}

var h264Levels = []bitmovintypes.H264Level{
	bitmovintypes.H264Level1,
	bitmovintypes.H264Level1b,
	bitmovintypes.H264Level1_1,
	bitmovintypes.H264Level1_2,
	bitmovintypes.H264Level1_3,
	bitmovintypes.H264Level2,
	bitmovintypes.H264Level2_1,
	bitmovintypes.H264Level2_2,
	bitmovintypes.H264Level3,
	bitmovintypes.H264Level3_1,
	bitmovintypes.H264Level3_2,
	bitmovintypes.H264Level4,
	bitmovintypes.H264Level4_1,
	bitmovintypes.H264Level4_2,
	bitmovintypes.H264Level5,
	bitmovintypes.H264Level5_1,
	bitmovintypes.H264Level5_2}

var errBitmovinInvalidConfig = provider.InvalidConfigError("Invalid configuration")

var s3Pattern = regexp.MustCompile(`^s3://`)
var httpPattern = regexp.MustCompile(`^http://`)
var httpsPattern = regexp.MustCompile(`^https://`)

type bitmovinProvider struct {
	client *bitmovin.Bitmovin
	config *config.Bitmovin
}

type bitmovinH264Preset struct {
	Video models.H264CodecConfiguration
	Audio models.AACCodecConfiguration
}

type bitmovinVP8Preset struct {
	Video models.VP8CodecConfiguration
	Audio models.VorbisCodecConfiguration
}

func (p *bitmovinProvider) CreatePreset(preset db.Preset) (string, error) {
	fmt.Println("HERE")
	if strings.ToLower(preset.Audio.Codec) == "aac" && strings.ToLower(preset.Video.Codec) == "h264" {
		aac := services.NewAACCodecConfigurationService(p.client)
		var audioConfigID string
		bitrate, err := strconv.Atoi(preset.Audio.Bitrate)
		if err != nil {
			return "", err
		}
		temp := int64(bitrate)
		audioConfig := &models.AACCodecConfiguration{
			Name:         stringToPtr(preset.Name),
			Bitrate:      &temp,
			SamplingRate: floatToPtr(48000.0),
		}
		audioResp, err := aac.Create(audioConfig)
		if err != nil {
			return "", err
		}
		if audioResp.Status == bitmovinAPIErrorMsg {
			return "", errors.New("Error in creating audio portion of Preset")
		}

		audioConfigID = *audioResp.Data.Result.ID

		customData := make(map[string]interface{})
		customData["audio"] = audioConfigID
		customData["container"] = preset.Container
		h264Config, err := p.createH264VideoPreset(preset, customData)
		if err != nil {
			return "", err
		}

		h264 := services.NewH264CodecConfigurationService(p.client)
		videoResp, err := h264.Create(h264Config)
		if err != nil {
			return "", err
		}
		if videoResp.Status == bitmovinAPIErrorMsg {
			return "", errors.New("error in creating video portion of Preset")
		}
		return *videoResp.Data.Result.ID, nil
	}
	if strings.ToLower(preset.Audio.Codec) == "vorbis" && strings.ToLower(preset.Video.Codec) == "vp8" {
		vorbis := services.NewVorbisCodecConfigurationService(p.client)
		var audioConfigID string
		bitrate, err := strconv.Atoi(preset.Audio.Bitrate)
		if err != nil {
			return "", err
		}
		temp := int64(bitrate)
		audioConfig := &models.VorbisCodecConfiguration{
			Name:         stringToPtr(preset.Name),
			Bitrate:      &temp,
			SamplingRate: floatToPtr(48000.0),
		}
		audioResp, err := vorbis.Create(audioConfig)
		if err != nil {
			return "", err
		}
		if audioResp.Status == bitmovinAPIErrorMsg {
			return "", errors.New("Error in creating audio portion of Preset")
		}

		audioConfigID = *audioResp.Data.Result.ID

		customData := make(map[string]interface{})
		customData["audio"] = audioConfigID
		customData["container"] = preset.Container
		vp8Config, err := p.createVP8VideoPreset(preset, customData)
		if err != nil {
			return "", err
		}

		vp8 := services.NewVP8CodecConfigurationService(p.client)
		videoResp, err := vp8.Create(vp8Config)
		if err != nil {
			return "", err
		}
		if videoResp.Status == bitmovinAPIErrorMsg {
			return "", errors.New("error in creating video portion of Preset")
		}
		return *videoResp.Data.Result.ID, nil
	}
	return "", fmt.Errorf("Unsupported Audio codec: %v", preset.Audio.Codec)
}

func (p *bitmovinProvider) createAACAudioPreset(preset db.Preset) (*models.AACCodecConfiguration, error) {
	bitrate, err := strconv.Atoi(preset.Audio.Bitrate)
	if err != nil {
		return nil, err
	}
	temp := int64(bitrate)
	audioConfig := &models.AACCodecConfiguration{
		Name:         stringToPtr(preset.Name),
		Bitrate:      &temp,
		SamplingRate: floatToPtr(48000.0),
	}
	return audioConfig, nil
}

func (p *bitmovinProvider) createVorbisAudioPreset(preset db.Preset) (*models.VorbisCodecConfiguration, error) {
	bitrate, err := strconv.Atoi(preset.Audio.Bitrate)
	if err != nil {
		return nil, err
	}
	temp := int64(bitrate)
	audioConfig := &models.VorbisCodecConfiguration{
		Name:    stringToPtr(preset.Name),
		Bitrate: &temp,
		// FIXME: double check bitrate sampling rate comibnatins for vorbis
		SamplingRate: floatToPtr(48000.0),
	}
	return audioConfig, nil
}

func (p *bitmovinProvider) createH264VideoPreset(preset db.Preset, customData map[string]interface{}) (*models.H264CodecConfiguration, error) {
	h264 := &models.H264CodecConfiguration{
		CustomData: customData,
	}
	profile := strings.ToLower(preset.Video.Profile)
	h264.Name = stringToPtr(preset.Name)
	switch profile {
	case "high":
		h264.Profile = bitmovintypes.H264ProfileHigh
	case "main":
		h264.Profile = bitmovintypes.H264ProfileMain
	case "baseline":
		h264.Profile = bitmovintypes.H264ProfileBaseline
	case "":
		h264.Profile = bitmovintypes.H264ProfileMain
	default:
		return nil, fmt.Errorf("Unrecognized H264 Profile: %v", preset.Video.Profile)
	}
	foundLevel := false
	for _, l := range h264Levels {
		if l == bitmovintypes.H264Level(preset.Video.ProfileLevel) {
			h264.Level = l
			foundLevel = true
			break
		}
	}
	if !foundLevel {
		// Just set it to the highest Level
		h264.Level = bitmovintypes.H264Level5_2
	}
	if preset.Video.Width != "" {
		width, err := strconv.Atoi(preset.Video.Width)
		if err != nil {
			return nil, err
		}
		h264.Width = intToPtr(int64(width))
	}
	if preset.Video.Height != "" {
		height, err := strconv.Atoi(preset.Video.Height)
		if err != nil {
			return nil, err
		}
		h264.Height = intToPtr(int64(height))
	}

	if preset.Video.Bitrate == "" {
		return nil, errors.New("Video Bitrate must be set")
	}
	bitrate, err := strconv.Atoi(preset.Video.Bitrate)
	if err != nil {
		return nil, err
	}
	h264.Bitrate = intToPtr(int64(bitrate))
	if preset.Video.GopSize != "" {
		gopSize, err := strconv.Atoi(preset.Video.GopSize)
		if err != nil {
			return nil, err
		}
		h264.MaxGOP = intToPtr(int64(gopSize))
	}

	return h264, nil
}

func (p *bitmovinProvider) createVP8VideoPreset(preset db.Preset, customData map[string]interface{}) (*models.VP8CodecConfiguration, error) {
	vp8 := &models.VP8CodecConfiguration{
		CustomData: customData,
	}
	vp8.Name = stringToPtr(preset.Name)

	if preset.Video.Width != "" {
		width, err := strconv.Atoi(preset.Video.Width)
		if err != nil {
			return nil, err
		}
		vp8.Width = intToPtr(int64(width))
	}
	if preset.Video.Height != "" {
		height, err := strconv.Atoi(preset.Video.Height)
		if err != nil {
			return nil, err
		}
		vp8.Height = intToPtr(int64(height))
	}

	if preset.Video.Bitrate == "" {
		return nil, errors.New("Video Bitrate must be set")
	}
	bitrate, err := strconv.Atoi(preset.Video.Bitrate)
	if err != nil {
		return nil, err
	}
	vp8.Bitrate = intToPtr(int64(bitrate))

	return vp8, nil
}

func (p *bitmovinProvider) DeletePreset(presetID string) error {
	// Delete both the audio and video preset
	h264 := services.NewH264CodecConfigurationService(p.client)
	_, err := h264.Retrieve(presetID)
	if err == nil {
		cdResp, err := h264.RetrieveCustomData(presetID)
		if err != nil {
			return err
		}
		if cdResp.Status == bitmovinAPIErrorMsg {
			return errors.New("Video Preset must contain custom data to hold audio and container information")
		}
		var audioPresetID string
		if cdResp.Data.Result.CustomData != nil {
			cd := cdResp.Data.Result.CustomData
			i, ok := cd["audio"]
			if !ok {
				return errors.New("No Audio configuration found for Video Preset")
			}
			audioPresetID, ok = i.(string)
			if !ok {
				return errors.New("Audio Configuration somehow not a string")
			}
		} else {
			return errors.New("No Audio configuration found for Video Preset")
		}

		aac := services.NewAACCodecConfigurationService(p.client)
		audioDeleteResp, err := aac.Delete(audioPresetID)
		if err != nil {
			return err
		}
		if audioDeleteResp.Status == bitmovinAPIErrorMsg {
			return errors.New("Error in deleting audio portion of Preset")
		}

		videoDeleteResp, err := h264.Delete(presetID)
		if err != nil {
			return err
		}
		if videoDeleteResp.Status == bitmovinAPIErrorMsg {
			return errors.New("Error in deleting video portion of Preset")
		}
		return nil
	}
	vp8 := services.NewVP8CodecConfigurationService(p.client)
	_, err = vp8.Retrieve(presetID)
	if err == nil {
		cdResp, err := vp8.RetrieveCustomData(presetID)
		if err != nil {
			return err
		}
		if cdResp.Status == bitmovinAPIErrorMsg {
			return errors.New("Video Preset must contain custom data to hold audio and container information")
		}
		var audioPresetID string
		if cdResp.Data.Result.CustomData != nil {
			cd := cdResp.Data.Result.CustomData
			i, ok := cd["audio"]
			if !ok {
				return errors.New("No Audio configuration found for Video Preset")
			}
			audioPresetID, ok = i.(string)
			if !ok {
				return errors.New("Audio Configuration somehow not a string")
			}
		} else {
			return errors.New("No Audio configuration found for Video Preset")
		}

		vorbis := services.NewVorbisCodecConfigurationService(p.client)
		audioDeleteResp, err := vorbis.Delete(audioPresetID)
		if err != nil {
			return err
		}
		if audioDeleteResp.Status == bitmovinAPIErrorMsg {
			return errors.New("Error in deleting audio portion of Preset")
		}

		videoDeleteResp, err := vp8.Delete(presetID)
		if err != nil {
			return err
		}
		if videoDeleteResp.Status == bitmovinAPIErrorMsg {
			return errors.New("Error in deleting video portion of Preset")
		}
		return nil
	}
	return errors.New("Could not find H.264 or VP8 configuration to delete")
}

// func (p *bitmovinProvider) DeletePreset(presetID string) error {
// 	// Delete both the audio and video preset
// 	h264 := services.NewH264CodecConfigurationService(p.client)
// 	cdResp, err := h264.RetrieveCustomData(presetID)
// 	if err != nil {
// 		return err
// 	}
// 	if cdResp.Status == bitmovinAPIErrorMsg {
// 		return errors.New("Video Preset must contain custom data to hold audio and container information")
// 	}
// 	var audioPresetID string
// 	if cdResp.Data.Result.CustomData != nil {
// 		cd := cdResp.Data.Result.CustomData
// 		i, ok := cd["audio"]
// 		if !ok {
// 			return errors.New("No Audio configuration found for Video Preset")
// 		}
// 		audioPresetID, ok = i.(string)
// 		if !ok {
// 			return errors.New("Audio Configuration somehow not a string")
// 		}
// 	} else {
// 		return errors.New("No Audio configuration found for Video Preset")
// 	}

// 	aac := services.NewAACCodecConfigurationService(p.client)
// 	audioDeleteResp, err := aac.Delete(audioPresetID)
// 	if err != nil {
// 		return err
// 	}
// 	if audioDeleteResp.Status == bitmovinAPIErrorMsg {
// 		return errors.New("Error in deleting audio portion of Preset")
// 	}

// 	videoDeleteResp, err := h264.Delete(presetID)
// 	if err != nil {
// 		return err
// 	}
// 	if videoDeleteResp.Status == bitmovinAPIErrorMsg {
// 		return errors.New("Error in deleting video portion of Preset")
// 	}
// 	return nil
// }

func (p *bitmovinProvider) GetPreset(presetID string) (interface{}, error) {
	h264 := services.NewH264CodecConfigurationService(p.client)
	response, err := h264.Retrieve(presetID)
	if err == nil {
		// It is H.264 and AAC
		if response.Status == bitmovinAPIErrorMsg {
			return nil, errors.New("Error in retrieving video portion of Preset")
		}
		h264Config := response.Data.Result
		cd, err := h264.RetrieveCustomData(presetID)
		if err != nil {
			return nil, err
		}
		if cd.Status == bitmovinAPIErrorMsg {
			return nil, errors.New("")
		}
		if cd.Data.Result.CustomData != nil {
			h264Config.CustomData = cd.Data.Result.CustomData
			i, ok := h264Config.CustomData["audio"]
			if !ok {
				return nil, errors.New("No Audio configuration found for Video Preset")
			}
			s, ok := i.(string)
			if !ok {
				return nil, errors.New("Audio Configuration somehow not a string")
			}
			aac := services.NewAACCodecConfigurationService(p.client)
			audioResponse, err := aac.Retrieve(s)
			if err != nil {
				return nil, err
			}
			if audioResponse.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in retrieving audio portion of Preset")
			}
			aacConfig := audioResponse.Data.Result
			preset := bitmovinH264Preset{
				Video: h264Config,
				Audio: aacConfig,
			}
			return preset, nil
		}
		return nil, errors.New("No Audio configuration found for Video Preset")
	}

	vp8 := services.NewVP8CodecConfigurationService(p.client)
	vp8Response, err := vp8.Retrieve(presetID)
	if err == nil {
		// It is VP8 and Vorbis
		if vp8Response.Status == bitmovinAPIErrorMsg {
			return nil, errors.New("Error in retrieving video portion of Preset")
		}
		vp8Config := vp8Response.Data.Result
		cd, err := vp8.RetrieveCustomData(presetID)
		if err != nil {
			return nil, err
		}
		if cd.Status == bitmovinAPIErrorMsg {
			return nil, errors.New("")
		}
		if cd.Data.Result.CustomData != nil {
			vp8Config.CustomData = cd.Data.Result.CustomData
			i, ok := vp8Config.CustomData["audio"]
			if !ok {
				return nil, errors.New("No Audio configuration found for Video Preset")
			}
			s, ok := i.(string)
			if !ok {
				return nil, errors.New("Audio Configuration somehow not a string")
			}
			vorbis := services.NewVorbisCodecConfigurationService(p.client)
			audioResponse, err := vorbis.Retrieve(s)
			if err != nil {
				return nil, err
			}
			if audioResponse.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in retrieving audio portion of Preset")
			}
			vorbisConfig := audioResponse.Data.Result
			preset := bitmovinVP8Preset{
				Video: vp8Config,
				Audio: vorbisConfig,
			}
			return preset, nil
		}
	}

	return nil, errors.New("No Audio configuration found for Video Preset")
}

// func (p *bitmovinProvider) GetPreset(presetID string) (interface{}, error) {
// 	h264 := services.NewH264CodecConfigurationService(p.client)
// 	response, err := h264.Retrieve(presetID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if response.Status == bitmovinAPIErrorMsg {
// 		return nil, errors.New("Error in retrieving video portion of Preset")
// 	}
// 	h264Config := response.Data.Result
// 	cd, err := h264.RetrieveCustomData(presetID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if cd.Status == bitmovinAPIErrorMsg {
// 		return nil, errors.New("")
// 	}
// 	if cd.Data.Result.CustomData != nil {
// 		h264Config.CustomData = cd.Data.Result.CustomData
// 		i, ok := h264Config.CustomData["audio"]
// 		if !ok {
// 			return nil, errors.New("No Audio configuration found for Video Preset")
// 		}
// 		s, ok := i.(string)
// 		if !ok {
// 			return nil, errors.New("Audio Configuration somehow not a string")
// 		}
// 		aac := services.NewAACCodecConfigurationService(p.client)
// 		audioResponse, err := aac.Retrieve(s)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if audioResponse.Status == bitmovinAPIErrorMsg {
// 			return nil, errors.New("Error in retrieving audio portion of Preset")
// 		}
// 		aacConfig := audioResponse.Data.Result
// 		preset := bitmovinPreset{
// 			Video: h264Config,
// 			Audio: aacConfig,
// 		}
// 		return preset, nil
// 	}
// 	return nil, errors.New("No Audio configuration found for Video Preset")
// }

// func (p *bitmovinProvider) Transcode(job *db.Job) (*provider.JobStatus, error) {
// 	aclEntry := models.ACLItem{
// 		Permission: bitmovintypes.ACLPermissionPublicRead,
// 	}
// 	acl := []models.ACLItem{aclEntry}

// 	cloudRegion := bitmovintypes.AWSCloudRegion(p.config.AWSStorageRegion)
// 	outputBucketName, err := grabBucketNameFromS3Destination(p.config.Destination)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s3OS := services.NewS3OutputService(p.client)
// 	s3Output := &models.S3Output{
// 		BucketName:  stringToPtr(outputBucketName),
// 		AccessKey:   stringToPtr(p.config.AccessKeyID),
// 		SecretKey:   stringToPtr(p.config.SecretAccessKey),
// 		CloudRegion: cloudRegion,
// 	}

// 	s3OSResponse, err := s3OS.Create(s3Output)
// 	if err != nil {
// 		return nil, err
// 	} else if s3OSResponse.Status == bitmovinAPIErrorMsg {
// 		return nil, errors.New("Error in setting up S3 input")
// 	}

// 	inputID, inputFullPath, err := createInput(p, job.SourceMedia)
// 	if err != nil {
// 		return nil, err
// 	}

// 	videoInputStream := models.InputStream{
// 		InputID:       stringToPtr(inputID),
// 		InputPath:     stringToPtr(inputFullPath),
// 		SelectionMode: bitmovintypes.SelectionModeAuto,
// 	}

// 	audioInputStream := models.InputStream{
// 		InputID:       stringToPtr(inputID),
// 		InputPath:     stringToPtr(inputFullPath),
// 		SelectionMode: bitmovintypes.SelectionModeAuto,
// 	}

// 	viss := []models.InputStream{videoInputStream}
// 	aiss := []models.InputStream{audioInputStream}

// 	h264S := services.NewH264CodecConfigurationService(p.client)

// 	var masterManifestPath string
// 	var masterManifestFile string
// 	outputtingHLS := false
// 	manifestID := ""

// 	//create the master manifest if needed so we can add it to the customData of the encoding response
// 	for _, output := range job.Outputs {
// 		videoPresetID := output.Preset.ProviderMapping[Name]
// 		customDataResp, cdErr := h264S.RetrieveCustomData(videoPresetID)

// 		if cdErr != nil {
// 			return nil, cdErr
// 		}
// 		if customDataResp.Status == bitmovinAPIErrorMsg {
// 			return nil, errors.New("")
// 		}
// 		containerInterface, ok := customDataResp.Data.Result.CustomData["container"]
// 		if !ok {
// 			return nil, errors.New("")
// 		}
// 		container, ok := containerInterface.(string)
// 		if !ok {
// 			return nil, errors.New("")
// 		}
// 		if container == "m3u8" {
// 			outputtingHLS = true
// 			break
// 		}
// 	}

// 	hlsService := services.NewHLSManifestService(p.client)

// 	if outputtingHLS {
// 		masterManifestPath = filepath.Dir(job.StreamingParams.PlaylistFileName)
// 		masterManifestFile = filepath.Base(job.StreamingParams.PlaylistFileName)
// 		manifestOutput := models.Output{
// 			OutputID:   s3OSResponse.Data.Result.ID,
// 			OutputPath: stringToPtr(masterManifestPath),
// 			ACL:        acl,
// 		}
// 		hlsMasterManifest := &models.HLSManifest{
// 			ManifestName: stringToPtr(masterManifestFile),
// 			Outputs:      []models.Output{manifestOutput},
// 		}
// 		hlsMasterManifestResp, manErr := hlsService.Create(hlsMasterManifest)
// 		if manErr != nil {
// 			return nil, manErr
// 		} else if hlsMasterManifestResp.Status == bitmovinAPIErrorMsg {
// 			return nil, errors.New("Error in HLS Master Manifest creation")
// 		}
// 		manifestID = *hlsMasterManifestResp.Data.Result.ID
// 	}

// 	encodingS := services.NewEncodingService(p.client)
// 	customData := make(map[string]interface{})
// 	if outputtingHLS {
// 		customData["manifest"] = manifestID
// 	}
// 	encodingRegion := bitmovintypes.CloudRegion(p.config.EncodingRegion)
// 	encodingVersion := bitmovintypes.EncoderVersion(p.config.EncodingVersion)
// 	encoding := &models.Encoding{
// 		Name:           stringToPtr("encoding"),
// 		CustomData:     customData,
// 		CloudRegion:    encodingRegion,
// 		EncoderVersion: encodingVersion,
// 	}

// 	encodingResp, err := encodingS.Create(encoding)
// 	if err != nil {
// 		return nil, err
// 	} else if encodingResp.Status == bitmovinAPIErrorMsg {
// 		return nil, errors.New("Error in Encoding Creation")
// 	}

// 	for _, output := range job.Outputs {
// 		videoPresetID := output.Preset.ProviderMapping[Name]
// 		videoResponse, h264Err := h264S.Retrieve(videoPresetID)
// 		if err != nil {
// 			return nil, h264Err
// 		}
// 		if videoResponse.Status == bitmovinAPIErrorMsg {
// 			return nil, errors.New("Error in retrieving video portion of preset")
// 		}
// 		customDataResp, h264CDErr := h264S.RetrieveCustomData(videoPresetID)
// 		if err != nil {
// 			return nil, h264CDErr
// 		}
// 		if customDataResp.Status == bitmovinAPIErrorMsg {
// 			return nil, errors.New("Error in retrieving video custom data where the audio ID and container type is stored")
// 		}
// 		audioPresetIDInterface, ok := customDataResp.Data.Result.CustomData["audio"]
// 		if !ok {
// 			return nil, errors.New("Audio ID not found in video custom data")
// 		}
// 		audioPresetID, ok := audioPresetIDInterface.(string)
// 		if !ok {
// 			return nil, errors.New("Audio ID somehow not a string")
// 		}

// 		var audioStreamID, videoStreamID string
// 		audioStream := &models.Stream{
// 			CodecConfigurationID: &audioPresetID,
// 			InputStreams:         aiss,
// 		}
// 		audioStreamResp, audioErr := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream)
// 		if audioErr != nil {
// 			return nil, audioErr
// 		}
// 		if audioStreamResp.Status == bitmovinAPIErrorMsg {
// 			return nil, errors.New("Error in adding audio stream to Encoding")
// 		}
// 		audioStreamID = *audioStreamResp.Data.Result.ID

// 		videoStream := &models.Stream{
// 			CodecConfigurationID: &videoPresetID,
// 			InputStreams:         viss,
// 		}
// 		videoStreamResp, vsErr := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream)
// 		if vsErr != nil {
// 			return nil, vsErr
// 		}
// 		if videoStreamResp.Status == bitmovinAPIErrorMsg {
// 			return nil, errors.New("Error in adding video stream to Encoding")
// 		}
// 		videoStreamID = *videoStreamResp.Data.Result.ID

// 		audioMuxingStream := models.StreamItem{
// 			StreamID: &audioStreamID,
// 		}
// 		videoMuxingStream := models.StreamItem{
// 			StreamID: &videoStreamID,
// 		}

// 		containerInterface, ok := customDataResp.Data.Result.CustomData["container"]
// 		if !ok {
// 			return nil, errors.New("Container type not found in video custom data")
// 		}
// 		container, ok := containerInterface.(string)
// 		if !ok {
// 			return nil, errors.New("Container type somehow not a string")
// 		}
// 		if container == "m3u8" {
// 			audioMuxingOutput := models.Output{
// 				OutputID:   s3OSResponse.Data.Result.ID,
// 				OutputPath: stringToPtr(filepath.Join(masterManifestPath, audioPresetID)),
// 				ACL:        acl,
// 			}
// 			audioMuxing := &models.TSMuxing{
// 				SegmentLength: floatToPtr(float64(job.StreamingParams.SegmentDuration)),
// 				SegmentNaming: stringToPtr("seg_%number%.ts"),
// 				Streams:       []models.StreamItem{audioMuxingStream},
// 				Outputs:       []models.Output{audioMuxingOutput},
// 			}
// 			audioMuxingResp, muxErr := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, audioMuxing)
// 			if muxErr != nil {
// 				return nil, muxErr
// 			}
// 			if audioMuxingResp.Status == bitmovinAPIErrorMsg {
// 				return nil, errors.New("Error in adding TS Muxing for audio")
// 			}

// 			// create the MediaInfo
// 			audioMediaInfo := &models.MediaInfo{
// 				Type:            bitmovintypes.MediaTypeAudio,
// 				URI:             stringToPtr(audioPresetID + ".m3u8"),
// 				GroupID:         stringToPtr(audioPresetID),
// 				Language:        stringToPtr("en"),
// 				Name:            stringToPtr(audioPresetID),
// 				IsDefault:       boolToPtr(false),
// 				Autoselect:      boolToPtr(false),
// 				Forced:          boolToPtr(false),
// 				SegmentPath:     stringToPtr(audioPresetID),
// 				Characteristics: []string{"public.accessibility.describes-video"},
// 				EncodingID:      encodingResp.Data.Result.ID,
// 				StreamID:        audioStreamResp.Data.Result.ID,
// 				MuxingID:        audioMuxingResp.Data.Result.ID,
// 			}

// 			// Add to Master manifest, we will set the m3u8 and segments relative to the master

// 			audioMediaInfoResp, miErr := hlsService.AddMediaInfo(manifestID, audioMediaInfo)
// 			if miErr != nil {
// 				return nil, miErr
// 			}
// 			if audioMediaInfoResp.Status == bitmovinAPIErrorMsg {
// 				return nil, errors.New("Error in adding EXT-X-MEDIA")
// 			}

// 			videoMuxingOutput := models.Output{
// 				OutputID:   s3OSResponse.Data.Result.ID,
// 				OutputPath: stringToPtr(filepath.Join(masterManifestPath, videoPresetID)),
// 				ACL:        acl,
// 			}
// 			videoMuxing := &models.TSMuxing{
// 				SegmentLength: floatToPtr(float64(job.StreamingParams.SegmentDuration)),
// 				SegmentNaming: stringToPtr("seg_%number%.ts"),
// 				Streams:       []models.StreamItem{videoMuxingStream},
// 				Outputs:       []models.Output{videoMuxingOutput},
// 			}
// 			videoMuxingResp, vmuxErr := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoMuxing)
// 			if err != nil {
// 				return nil, vmuxErr
// 			}
// 			if videoMuxingResp.Status == bitmovinAPIErrorMsg {
// 				return nil, errors.New("Error in adding TS Muxing for video")
// 			}

// 			videoStreamInfo := &models.StreamInfo{
// 				Audio:       stringToPtr(audioPresetID),
// 				SegmentPath: stringToPtr(videoPresetID),
// 				URI:         stringToPtr(filepath.Base(output.FileName)),
// 				EncodingID:  encodingResp.Data.Result.ID,
// 				StreamID:    videoStreamResp.Data.Result.ID,
// 				MuxingID:    videoMuxingResp.Data.Result.ID,
// 			}

// 			videoStreamInfoResp, vsiErr := hlsService.AddStreamInfo(manifestID, videoStreamInfo)
// 			if vsiErr != nil {
// 				return nil, vsiErr
// 			}
// 			if videoStreamInfoResp.Status == bitmovinAPIErrorMsg {
// 				return nil, errors.New("Error in adding EXT-X-STREAM-INF")
// 			}
// 		} else if container == "mp4" {
// 			videoMuxingOutput := models.Output{
// 				OutputID:   s3OSResponse.Data.Result.ID,
// 				ACL:        acl,
// 				OutputPath: stringToPtr(filepath.Dir(output.FileName)),
// 			}
// 			videoMuxing := &models.MP4Muxing{
// 				Filename: stringToPtr(filepath.Base(output.FileName)),
// 				Outputs:  []models.Output{videoMuxingOutput},
// 				Streams:  []models.StreamItem{videoMuxingStream, audioMuxingStream},
// 			}
// 			videoMuxingResp, vmErr := encodingS.AddMP4Muxing(*encodingResp.Data.Result.ID, videoMuxing)
// 			if err != nil {
// 				return nil, vmErr
// 			}
// 			if videoMuxingResp.Status == bitmovinAPIErrorMsg {
// 				return nil, errors.New("Error in adding MP4 Muxing")
// 			}
// 		}
// 	}

// 	startResp, err := encodingS.Start(*encodingResp.Data.Result.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if startResp.Status == bitmovinAPIErrorMsg {
// 		return nil, errors.New("Error in starting encoding")
// 	}

// 	jobStatus := &provider.JobStatus{
// 		ProviderName:  Name,
// 		ProviderJobID: *encodingResp.Data.Result.ID,
// 		Status:        provider.StatusQueued,
// 	}

// 	return jobStatus, nil
// }

func (p *bitmovinProvider) Transcode(job *db.Job) (*provider.JobStatus, error) {
	aclEntry := models.ACLItem{
		Permission: bitmovintypes.ACLPermissionPublicRead,
	}
	acl := []models.ACLItem{aclEntry}

	cloudRegion := bitmovintypes.AWSCloudRegion(p.config.AWSStorageRegion)
	outputBucketName, err := grabBucketNameFromS3Destination(p.config.Destination)
	if err != nil {
		return nil, err
	}

	s3OS := services.NewS3OutputService(p.client)
	s3Output := &models.S3Output{
		BucketName:  stringToPtr(outputBucketName),
		AccessKey:   stringToPtr(p.config.AccessKeyID),
		SecretKey:   stringToPtr(p.config.SecretAccessKey),
		CloudRegion: cloudRegion,
	}

	s3OSResponse, err := s3OS.Create(s3Output)
	if err != nil {
		return nil, err
	} else if s3OSResponse.Status == bitmovinAPIErrorMsg {
		return nil, errors.New("Error in setting up S3 input")
	}

	inputID, inputFullPath, err := createInput(p, job.SourceMedia)
	if err != nil {
		return nil, err
	}

	videoInputStream := models.InputStream{
		InputID:       stringToPtr(inputID),
		InputPath:     stringToPtr(inputFullPath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}

	audioInputStream := models.InputStream{
		InputID:       stringToPtr(inputID),
		InputPath:     stringToPtr(inputFullPath),
		SelectionMode: bitmovintypes.SelectionModeAuto,
	}

	viss := []models.InputStream{videoInputStream}
	aiss := []models.InputStream{audioInputStream}

	h264S := services.NewH264CodecConfigurationService(p.client)
	vp8S := services.NewVP8CodecConfigurationService(p.client)

	var masterManifestPath string
	var masterManifestFile string
	outputtingHLS := false
	manifestID := ""

	//create the master manifest if needed so we can add it to the customData of the encoding response
	for _, output := range job.Outputs {
		videoPresetID := output.Preset.ProviderMapping[Name]
		customDataResp, cdErr := h264S.RetrieveCustomData(videoPresetID)

		if cdErr != nil {
			return nil, cdErr
		}
		if customDataResp.Status == bitmovinAPIErrorMsg {
			return nil, errors.New("")
		}
		containerInterface, ok := customDataResp.Data.Result.CustomData["container"]
		if !ok {
			return nil, errors.New("")
		}
		container, ok := containerInterface.(string)
		if !ok {
			return nil, errors.New("")
		}
		if container == "m3u8" {
			outputtingHLS = true
			break
		}
	}

	hlsService := services.NewHLSManifestService(p.client)

	if outputtingHLS {
		masterManifestPath = filepath.Dir(job.StreamingParams.PlaylistFileName)
		masterManifestFile = filepath.Base(job.StreamingParams.PlaylistFileName)
		manifestOutput := models.Output{
			OutputID:   s3OSResponse.Data.Result.ID,
			OutputPath: stringToPtr(masterManifestPath),
			ACL:        acl,
		}
		hlsMasterManifest := &models.HLSManifest{
			ManifestName: stringToPtr(masterManifestFile),
			Outputs:      []models.Output{manifestOutput},
		}
		hlsMasterManifestResp, manErr := hlsService.Create(hlsMasterManifest)
		if manErr != nil {
			return nil, manErr
		} else if hlsMasterManifestResp.Status == bitmovinAPIErrorMsg {
			return nil, errors.New("Error in HLS Master Manifest creation")
		}
		manifestID = *hlsMasterManifestResp.Data.Result.ID
	}

	encodingS := services.NewEncodingService(p.client)
	customData := make(map[string]interface{})
	if outputtingHLS {
		customData["manifest"] = manifestID
	}
	encodingRegion := bitmovintypes.CloudRegion(p.config.EncodingRegion)
	encodingVersion := bitmovintypes.EncoderVersion(p.config.EncodingVersion)
	encoding := &models.Encoding{
		Name:           stringToPtr("encoding"),
		CustomData:     customData,
		CloudRegion:    encodingRegion,
		EncoderVersion: encodingVersion,
	}

	encodingResp, err := encodingS.Create(encoding)
	if err != nil {
		return nil, err
	} else if encodingResp.Status == bitmovinAPIErrorMsg {
		return nil, errors.New("Error in Encoding Creation")
	}

	for _, output := range job.Outputs {
		videoPresetID := output.Preset.ProviderMapping[Name]
		h264VideoResponse, h264Err := h264S.Retrieve(videoPresetID)
		if h264Err == nil {
			if h264VideoResponse.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in retrieving video portion of preset")
			}
			customDataResp, h264CDErr := h264S.RetrieveCustomData(videoPresetID)
			if h264CDErr != nil {
				return nil, h264CDErr
			}
			if customDataResp.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in retrieving video custom data where the audio ID and container type is stored")
			}
			audioPresetIDInterface, ok := customDataResp.Data.Result.CustomData["audio"]
			if !ok {
				return nil, errors.New("Audio ID not found in video custom data")
			}
			audioPresetID, ok := audioPresetIDInterface.(string)
			if !ok {
				return nil, errors.New("Audio ID somehow not a string")
			}

			var audioStreamID, videoStreamID string
			audioStream := &models.Stream{
				CodecConfigurationID: &audioPresetID,
				InputStreams:         aiss,
			}
			audioStreamResp, audioErr := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream)
			if audioErr != nil {
				return nil, audioErr
			}
			if audioStreamResp.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in adding audio stream to Encoding")
			}
			audioStreamID = *audioStreamResp.Data.Result.ID

			videoStream := &models.Stream{
				CodecConfigurationID: &videoPresetID,
				InputStreams:         viss,
			}
			videoStreamResp, vsErr := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream)
			if vsErr != nil {
				return nil, vsErr
			}
			if videoStreamResp.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in adding video stream to Encoding")
			}
			videoStreamID = *videoStreamResp.Data.Result.ID

			audioMuxingStream := models.StreamItem{
				StreamID: &audioStreamID,
			}
			videoMuxingStream := models.StreamItem{
				StreamID: &videoStreamID,
			}

			containerInterface, ok := customDataResp.Data.Result.CustomData["container"]
			if !ok {
				return nil, errors.New("Container type not found in video custom data")
			}
			container, ok := containerInterface.(string)
			if !ok {
				return nil, errors.New("Container type somehow not a string")
			}
			if container == "m3u8" {
				audioMuxingOutput := models.Output{
					OutputID:   s3OSResponse.Data.Result.ID,
					OutputPath: stringToPtr(filepath.Join(masterManifestPath, audioPresetID)),
					ACL:        acl,
				}
				audioMuxing := &models.TSMuxing{
					SegmentLength: floatToPtr(float64(job.StreamingParams.SegmentDuration)),
					SegmentNaming: stringToPtr("seg_%number%.ts"),
					Streams:       []models.StreamItem{audioMuxingStream},
					Outputs:       []models.Output{audioMuxingOutput},
				}
				audioMuxingResp, muxErr := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, audioMuxing)
				if muxErr != nil {
					return nil, muxErr
				}
				if audioMuxingResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding TS Muxing for audio")
				}

				// create the MediaInfo
				audioMediaInfo := &models.MediaInfo{
					Type:            bitmovintypes.MediaTypeAudio,
					URI:             stringToPtr(audioPresetID + ".m3u8"),
					GroupID:         stringToPtr(audioPresetID),
					Language:        stringToPtr("en"),
					Name:            stringToPtr(audioPresetID),
					IsDefault:       boolToPtr(false),
					Autoselect:      boolToPtr(false),
					Forced:          boolToPtr(false),
					SegmentPath:     stringToPtr(audioPresetID),
					Characteristics: []string{"public.accessibility.describes-video"},
					EncodingID:      encodingResp.Data.Result.ID,
					StreamID:        audioStreamResp.Data.Result.ID,
					MuxingID:        audioMuxingResp.Data.Result.ID,
				}

				// Add to Master manifest, we will set the m3u8 and segments relative to the master

				audioMediaInfoResp, miErr := hlsService.AddMediaInfo(manifestID, audioMediaInfo)
				if miErr != nil {
					return nil, miErr
				}
				if audioMediaInfoResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding EXT-X-MEDIA")
				}

				videoMuxingOutput := models.Output{
					OutputID:   s3OSResponse.Data.Result.ID,
					OutputPath: stringToPtr(filepath.Join(masterManifestPath, videoPresetID)),
					ACL:        acl,
				}
				videoMuxing := &models.TSMuxing{
					SegmentLength: floatToPtr(float64(job.StreamingParams.SegmentDuration)),
					SegmentNaming: stringToPtr("seg_%number%.ts"),
					Streams:       []models.StreamItem{videoMuxingStream},
					Outputs:       []models.Output{videoMuxingOutput},
				}
				videoMuxingResp, vmuxErr := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoMuxing)
				if err != nil {
					return nil, vmuxErr
				}
				if videoMuxingResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding TS Muxing for video")
				}

				videoStreamInfo := &models.StreamInfo{
					Audio:       stringToPtr(audioPresetID),
					SegmentPath: stringToPtr(videoPresetID),
					URI:         stringToPtr(filepath.Base(output.FileName)),
					EncodingID:  encodingResp.Data.Result.ID,
					StreamID:    videoStreamResp.Data.Result.ID,
					MuxingID:    videoMuxingResp.Data.Result.ID,
				}

				videoStreamInfoResp, vsiErr := hlsService.AddStreamInfo(manifestID, videoStreamInfo)
				if vsiErr != nil {
					return nil, vsiErr
				}
				if videoStreamInfoResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding EXT-X-STREAM-INF")
				}
			} else if container == "mp4" {
				videoMuxingOutput := models.Output{
					OutputID:   s3OSResponse.Data.Result.ID,
					ACL:        acl,
					OutputPath: stringToPtr(filepath.Dir(output.FileName)),
				}
				videoMuxing := &models.MP4Muxing{
					Filename: stringToPtr(filepath.Base(output.FileName)),
					Outputs:  []models.Output{videoMuxingOutput},
					Streams:  []models.StreamItem{videoMuxingStream, audioMuxingStream},
				}
				videoMuxingResp, vmErr := encodingS.AddMP4Muxing(*encodingResp.Data.Result.ID, videoMuxing)
				if err != nil {
					return nil, vmErr
				}
				if videoMuxingResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding MP4 Muxing")
				}
			}
			return nil, errors.New("unknown container format")
		}
		vp8VideoResponse, vp8Err := vp8S.Retrieve(videoPresetID)
		if vp8Err == nil {
			if vp8VideoResponse.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in retrieving video portion of preset")
			}
			customDataResp, vp8CDErr := vp8S.RetrieveCustomData(videoPresetID)
			if vp8CDErr != nil {
				return nil, vp8CDErr
			}
			if customDataResp.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in retrieving video custom data where the audio ID and container type is stored")
			}
			audioPresetIDInterface, ok := customDataResp.Data.Result.CustomData["audio"]
			if !ok {
				return nil, errors.New("Audio ID not found in video custom data")
			}
			audioPresetID, ok := audioPresetIDInterface.(string)
			if !ok {
				return nil, errors.New("Audio ID somehow not a string")
			}
			var audioStreamID, videoStreamID string
			audioStream := &models.Stream{
				CodecConfigurationID: &audioPresetID,
				InputStreams:         aiss,
			}
			audioStreamResp, audioErr := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream)
			if audioErr != nil {
				return nil, audioErr
			}
			if audioStreamResp.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in adding audio stream to Encoding")
			}
			audioStreamID = *audioStreamResp.Data.Result.ID

			videoStream := &models.Stream{
				CodecConfigurationID: &videoPresetID,
				InputStreams:         viss,
			}
			videoStreamResp, vsErr := encodingS.AddStream(*encodingResp.Data.Result.ID, videoStream)
			if vsErr != nil {
				return nil, vsErr
			}
			if videoStreamResp.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in adding video stream to Encoding")
			}
			videoStreamID = *videoStreamResp.Data.Result.ID

			audioMuxingStream := models.StreamItem{
				StreamID: &audioStreamID,
			}
			videoMuxingStream := models.StreamItem{
				StreamID: &videoStreamID,
			}
			containerInterface, ok := customDataResp.Data.Result.CustomData["container"]
			if !ok {
				return nil, errors.New("Container type not found in video custom data")
			}
			container, ok := containerInterface.(string)
			if !ok {
				return nil, errors.New("Container type somehow not a string")
			}
			if container != "webm" {
				return nil, errors.New("unknown container for vp8 encoding")
			}
			videoMuxingOutput := models.Output{
				OutputID:   s3OSResponse.Data.Result.ID,
				ACL:        acl,
				OutputPath: stringToPtr(filepath.Dir(output.FileName)),
			}
			videoMuxing := &models.ProgressiveWebMMuxing{
				Filename: stringToPtr(filepath.Base(output.FileName)),
				Outputs:  []models.Output{videoMuxingOutput},
				Streams:  []models.StreamItem{videoMuxingStream, audioMuxingStream},
			}
			videoMuxingResp, vmErr := encodingS.AddProgressiveWebMMuxing(*encodingResp.Data.Result.ID, videoMuxing)
			if err != nil {
				return nil, vmErr
			}
			if videoMuxingResp.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error in adding MP4 Muxing")
			}
		}
		return nil, errors.New("No H264 or VP8 codec configuration found")
	}

	startResp, err := encodingS.Start(*encodingResp.Data.Result.ID)
	if err != nil {
		return nil, err
	}
	if startResp.Status == bitmovinAPIErrorMsg {
		return nil, errors.New("Error in starting encoding")
	}

	jobStatus := &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: *encodingResp.Data.Result.ID,
		Status:        provider.StatusQueued,
	}

	return jobStatus, nil
}

func (p *bitmovinProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	encodingS := services.NewEncodingService(p.client)
	statusResp, err := encodingS.RetrieveStatus(job.ProviderJobID)

	if err != nil {
		return nil, err
	}
	if statusResp.Status == bitmovinAPIErrorMsg {
		return &provider.JobStatus{
			ProviderName:  Name,
			ProviderJobID: job.ProviderJobID,
			Status:        provider.StatusFailed,
		}, nil
	}
	if *statusResp.Data.Result.Status == "FINISHED" {
		// see if manifest generation needs to happen
		cdResp, err := encodingS.RetrieveCustomData(job.ProviderJobID)
		if err != nil {
			return nil, err
		}
		if cdResp.Status == bitmovinAPIErrorMsg {
			return nil, errors.New("No Custom Data on Encoding, there should at least be container information here")
		}
		cd := cdResp.Data.Result.CustomData
		i, ok := cd["manifest"]
		if !ok {
			// No manifest generation needed, we are finished
			return &provider.JobStatus{
				ProviderName:  Name,
				ProviderJobID: job.ProviderJobID,
				Status:        provider.StatusFinished,
			}, nil
		}
		manifestID, ok := i.(string)
		if !ok {
			return nil, errors.New("Audio Configuration somehow not a string")
		}
		manifestS := services.NewHLSManifestService(p.client)
		manifestStatusResp, err := manifestS.RetrieveStatus(manifestID)
		if err != nil {
			return nil, err
		}
		if *manifestStatusResp.Data.Result.Status == bitmovinAPIErrorMsg {
			return &provider.JobStatus{
				ProviderName:  Name,
				ProviderJobID: job.ProviderJobID,
				Status:        provider.StatusFailed,
			}, nil
		}

		if *manifestStatusResp.Data.Result.Status == "CREATED" {
			// start the manifest generation
			startResp, err := manifestS.Start(manifestID)
			if err != nil {
				return nil, err
			} else if startResp.Status == bitmovinAPIErrorMsg {
				return &provider.JobStatus{
					ProviderName:  Name,
					ProviderJobID: job.ProviderJobID,
					Status:        provider.StatusFailed,
				}, nil
			}
			return &provider.JobStatus{
				ProviderName:  Name,
				ProviderJobID: job.ProviderJobID,
				Status:        provider.StatusStarted,
			}, nil
		}

		if *manifestStatusResp.Data.Result.Status == "QUEUED" || *manifestStatusResp.Data.Result.Status == "RUNNING" {
			return &provider.JobStatus{
				ProviderName:  Name,
				ProviderJobID: job.ProviderJobID,
				Status:        provider.StatusStarted,
			}, nil
		}

		if *manifestStatusResp.Data.Result.Status == "FINISHED" {
			return &provider.JobStatus{
				ProviderName:  Name,
				ProviderJobID: job.ProviderJobID,
				Status:        provider.StatusFinished,
			}, nil
		}
	} else if *statusResp.Data.Result.Status == "CREATED" || *statusResp.Data.Result.Status == "QUEUED" {
		return &provider.JobStatus{
			ProviderName:  Name,
			ProviderJobID: job.ProviderJobID,
			Status:        provider.StatusQueued,
		}, nil
	} else if *statusResp.Data.Result.Status == "RUNNING" {
		return &provider.JobStatus{
			ProviderName:  Name,
			ProviderJobID: job.ProviderJobID,
			Status:        provider.StatusStarted,
		}, nil
	}

	return &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: job.ProviderJobID,
		Status:        provider.StatusFailed,
	}, nil
}

func (p *bitmovinProvider) CancelJob(jobID string) error {
	// stop the job
	encodingS := services.NewEncodingService(p.client)
	resp, err := encodingS.Stop(jobID)
	if err != nil {
		return err
	}
	if resp.Status == bitmovinAPIErrorMsg {
		return errors.New("Error in canceling Job")
	}
	return nil
}

func (p *bitmovinProvider) Healthcheck() error {
	// Just going to call list encodings, and if it errors, then clearly it is unhealthy
	encodingS := services.NewEncodingService(p.client)
	resp, err := encodingS.List(int64(0), int64(1))
	if err != nil {
		return err
	}
	if resp.Status == bitmovinAPIErrorMsg {
		return errors.New("Bitmovin service unavailable")
	}
	return nil
}

func (p *bitmovinProvider) Capabilities() provider.Capabilities {
	// FIXME ?
	return provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls"},
		Destinations:  []string{"s3"},
	}
}

func bitmovinFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.Bitmovin.APIKey == "" {
		return nil, errBitmovinInvalidConfig
	}
	if _, ok := cloudRegions[cfg.Bitmovin.EncodingRegion]; !ok {
		return nil, errBitmovinInvalidConfig
	}
	if _, ok := awsCloudRegions[cfg.Bitmovin.AWSStorageRegion]; !ok {
		return nil, errBitmovinInvalidConfig
	}
	client := bitmovin.NewBitmovin(cfg.Bitmovin.APIKey, cfg.Bitmovin.Endpoint, int64(cfg.Bitmovin.Timeout))
	return &bitmovinProvider{client: client, config: cfg.Bitmovin}, nil
}

func parseS3URL(input string) (bucketName string, path string, fileName string, err error) {
	if s3Pattern.MatchString(input) {
		truncatedInput := strings.TrimPrefix(input, "s3://")
		splitTruncatedInput := strings.Split(truncatedInput, "/")
		bucketName = splitTruncatedInput[0]
		fileName = splitTruncatedInput[len(splitTruncatedInput)-1]
		truncatedInput = strings.TrimPrefix(truncatedInput, bucketName+"/")
		path = strings.TrimSuffix(truncatedInput, fileName)
		return
	}
	return "", "", "", errors.New("Could not parse S3 URL")
}

func grabBucketNameFromS3Destination(input string) (bucketName string, err error) {
	if s3Pattern.MatchString(input) {
		name := strings.TrimPrefix(input, "s3://")
		name = strings.TrimSuffix(name, "/")
		return name, err
	}
	return "", errors.New("Could not parse S3 URL")
}

func createInput(provider *bitmovinProvider, input string) (inputID string, path string, err error) {
	if s3Pattern.MatchString(input) {
		inputMediaBucketName, inputMediaPath, inputMediaFileName, err := parseS3URL(input)
		if err != nil {
			return "", "", err
		}
		cloudRegion := bitmovintypes.AWSCloudRegion(provider.config.AWSStorageRegion)

		s3IS := services.NewS3InputService(provider.client)
		s3Input := &models.S3Input{
			BucketName:  stringToPtr(inputMediaBucketName),
			AccessKey:   stringToPtr(provider.config.AccessKeyID),
			SecretKey:   stringToPtr(provider.config.SecretAccessKey),
			CloudRegion: cloudRegion,
		}

		s3ISResponse, err := s3IS.Create(s3Input)
		if err != nil {
			return "", "", err
		} else if s3ISResponse.Status == bitmovinAPIErrorMsg {
			return "", "", errors.New("Error in setting up S3 input")
		}
		var inputFullPath string
		if inputMediaPath == "" {
			inputFullPath = inputMediaFileName
		} else {
			inputFullPath = inputMediaPath + inputMediaFileName
		}
		return *s3ISResponse.Data.Result.ID, inputFullPath, nil
	} else if httpPattern.MatchString(input) {
		u, err := url.Parse(input)
		if err != nil {
			return "", "", errors.New("Error in setting up HTTP input")
		}
		httpIS := services.NewHTTPInputService(provider.client)
		httpInput := &models.HTTPInput{
			Host: stringToPtr(u.Host),
		}
		httpISResponse, err := httpIS.Create(httpInput)
		if err != nil {
			return "", "", errors.New("Error in setting up HTTP input")
		} else if httpISResponse.Status == bitmovinAPIErrorMsg {
			return "", "", errors.New("Error in setting up HTTP input")
		}
		return *httpISResponse.Data.Result.ID, u.Path, nil
	} else if httpsPattern.MatchString(input) {
		u, err := url.Parse(input)
		if err != nil {
			return "", "", errors.New("Error in setting up HTTPS input")
		}
		httpsIS := services.NewHTTPSInputService(provider.client)
		httpsInput := &models.HTTPSInput{
			Host: stringToPtr(u.Host),
		}
		httpsISResponse, err := httpsIS.Create(httpsInput)
		if err != nil {
			return "", "", errors.New("Error in setting up HTTPS input")
		} else if httpsISResponse.Status == bitmovinAPIErrorMsg {
			return "", "", errors.New("Error in setting up HTTPS input")
		}
		return *httpsISResponse.Data.Result.ID, u.Path, nil
	}
	return "", "", errors.New("Only S3, http, and https URLS are supported")
}

func stringToPtr(s string) *string {
	return &s
}

func intToPtr(i int64) *int64 {
	return &i
}

func boolToPtr(b bool) *bool {
	return &b
}

func floatToPtr(f float64) *float64 {
	return &f
}
