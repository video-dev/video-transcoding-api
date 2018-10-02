package bitmovin

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
	"github.com/bitmovin/bitmovin-go/services"
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
	if strings.ToLower(preset.Audio.Codec) == "aac" && strings.ToLower(preset.Video.Codec) == "h264" {
		aac := services.NewAACCodecConfigurationService(p.client)
		bitrate, err := strconv.Atoi(preset.Audio.Bitrate)
		if err != nil {
			return "", err
		}
		audioConfigID, err := getAACConfig(aac, bitrate, 48000)
		if err != nil {
			return "", err
		}

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
		bitrate, err := strconv.Atoi(preset.Audio.Bitrate)
		if err != nil {
			return "", err
		}
		audioConfigID, err := getVorbisConfig(vorbis, bitrate, 48000)
		if err != nil {
			return "", err
		}

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
	h264Response, err := h264.Retrieve(presetID)
	if err == nil {
		if h264Response.Status == bitmovinAPIErrorMsg {
			//if it were merely to not exist, then the err would not be nil
			return errors.New("API Error")
		}
		cdResp, cdErr := h264.RetrieveCustomData(presetID)
		if cdErr != nil {
			return cdErr
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
		audioDeleteResp, adErr := aac.Delete(audioPresetID)
		if adErr != nil {
			return adErr
		}
		if audioDeleteResp.Status == bitmovinAPIErrorMsg {
			return errors.New("Error in deleting audio portion of Preset")
		}

		videoDeleteResp, vdErr := h264.Delete(presetID)
		if vdErr != nil {
			return vdErr
		}
		if videoDeleteResp.Status == bitmovinAPIErrorMsg {
			return errors.New("Error in deleting video portion of Preset")
		}
		return nil
	}
	vp8 := services.NewVP8CodecConfigurationService(p.client)
	vp8Response, err := vp8.Retrieve(presetID)
	if err == nil {
		if vp8Response.Status == bitmovinAPIErrorMsg {
			//if it were merely to not exist, then the err would not be nil
			return errors.New("API Error")
		}
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

func (p *bitmovinProvider) GetPreset(presetID string) (interface{}, error) {
	h264 := services.NewH264CodecConfigurationService(p.client)
	response, err := h264.Retrieve(presetID)
	if err == nil {
		// It is H.264 and AAC
		if response.Status == bitmovinAPIErrorMsg {
			return nil, errors.New("Error in retrieving video portion of Preset")
		}
		h264Config := response.Data.Result
		cd, cdErr := h264.RetrieveCustomData(presetID)
		if cdErr != nil {
			return nil, cdErr
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
			audioResponse, arErr := aac.Retrieve(s)
			if arErr != nil {
				return nil, arErr
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

func (p *bitmovinProvider) Transcode(job *db.Job) (*provider.JobStatus, error) {
	aclEntry := models.ACLItem{
		Permission: bitmovintypes.ACLPermissionPublicRead,
	}
	acl := []models.ACLItem{aclEntry}

	cloudRegion := bitmovintypes.AWSCloudRegion(p.config.AWSStorageRegion)
	outputBucketName, prefix, err := parseS3URL(p.config.Destination)
	if err != nil {
		return nil, err
	}
	prefix = path.Join(prefix, job.ID)

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
		if output.Preset.OutputOpts.Extension != "webm" {
			videoPresetID := output.Preset.ProviderMapping[Name]
			customDataResp, cdErr := h264S.RetrieveCustomData(videoPresetID)

			if cdErr != nil {
				return nil, cdErr
			}
			if customDataResp.Status == bitmovinAPIErrorMsg {
				return nil, errors.New("Error retrieving Custom Data of Video Preset")
			}
			containerInterface, ok := customDataResp.Data.Result.CustomData["container"]
			if !ok {
				return nil, errors.New("Could not find container field in Custom Data")
			}
			container, ok := containerInterface.(string)
			if !ok {
				return nil, errors.New("Container in Custom Data could not be converted into a string")
			}
			if container == "m3u8" {
				outputtingHLS = true
				break
			}
		}
	}

	hlsService := services.NewHLSManifestService(p.client)

	if outputtingHLS {
		masterManifestPath = path.Dir(path.Join(prefix, job.StreamingParams.PlaylistFileName))
		masterManifestFile = path.Base(job.StreamingParams.PlaylistFileName)
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

	// create a map of audioMuxingStreams referenced by AudioPresetID so they are only created once.
	uniqueAudioMuxingStreams := make(map[string]models.StreamItem)
	uniqueAudioStreamResps := make(map[string]*models.StreamResponse)

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

			_, isRepeatedAudio := uniqueAudioMuxingStreams[audioPresetID]

			var audioStreamID, videoStreamID string

			var audioMuxingStream models.StreamItem

			if !isRepeatedAudio {
				audioStream := &models.Stream{
					CodecConfigurationID: &audioPresetID,
					InputStreams:         aiss,
					Conditions:           models.NewAttributeCondition(bitmovintypes.ConditionAttributeInputStream, "==", "true"),
				}
				audioStreamResp, audioErr := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream)
				if audioErr != nil {
					return nil, audioErr
				}
				if audioStreamResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding audio stream to Encoding")
				}
				audioStreamID = *audioStreamResp.Data.Result.ID

				audioMuxingStream = models.StreamItem{
					StreamID: &audioStreamID,
				}

				uniqueAudioMuxingStreams[audioPresetID] = audioMuxingStream
				uniqueAudioStreamResps[audioPresetID] = audioStreamResp
			}

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
				if !isRepeatedAudio {
					audioMuxingOutput := models.Output{
						OutputID:   s3OSResponse.Data.Result.ID,
						OutputPath: stringToPtr(path.Join(masterManifestPath, audioPresetID)),
						ACL:        acl,
					}
					audioMuxing := &models.TSMuxing{
						SegmentLength:        floatToPtr(float64(job.StreamingParams.SegmentDuration)),
						SegmentNaming:        stringToPtr("seg_%number%.ts"),
						Streams:              []models.StreamItem{uniqueAudioMuxingStreams[audioPresetID]},
						Outputs:              []models.Output{audioMuxingOutput},
						StreamConditionsMode: bitmovintypes.ConditionModeDropStream,
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
						StreamID:        uniqueAudioStreamResps[audioPresetID].Data.Result.ID,
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
				}

				videoMuxingOutput := models.Output{
					OutputID:   s3OSResponse.Data.Result.ID,
					OutputPath: stringToPtr(path.Join(masterManifestPath, videoPresetID)),
					ACL:        acl,
				}
				videoMuxing := &models.TSMuxing{
					SegmentLength: floatToPtr(float64(job.StreamingParams.SegmentDuration)),
					SegmentNaming: stringToPtr("seg_%number%.ts"),
					Streams:       []models.StreamItem{videoMuxingStream},
					Outputs:       []models.Output{videoMuxingOutput},
				}
				videoMuxingResp, vmuxErr := encodingS.AddTSMuxing(*encodingResp.Data.Result.ID, videoMuxing)
				if vmuxErr != nil {
					return nil, vmuxErr
				}
				if videoMuxingResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding TS Muxing for video")
				}

				videoManifestURI, err := filepath.Rel(masterManifestPath, path.Join(prefix, output.FileName))
				if err != nil {
					return nil, err
				}
				videoSegmentPath, err := filepath.Rel(path.Dir(path.Join(prefix, output.FileName)), path.Join(masterManifestPath, videoPresetID))
				if err != nil {
					return nil, err
				}

				videoStreamInfo := &models.StreamInfo{
					Audio:       stringToPtr(audioPresetID),
					SegmentPath: stringToPtr(videoSegmentPath),
					URI:         stringToPtr(videoManifestURI),
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
					OutputPath: stringToPtr(path.Dir(path.Join(prefix, output.FileName))),
				}
				videoMuxing := &models.MP4Muxing{
					Filename:             stringToPtr(path.Base(output.FileName)),
					Outputs:              []models.Output{videoMuxingOutput},
					Streams:              []models.StreamItem{videoMuxingStream, uniqueAudioMuxingStreams[audioPresetID]},
					StreamConditionsMode: bitmovintypes.ConditionModeDropStream,
				}
				videoMuxingResp, vmErr := encodingS.AddMP4Muxing(*encodingResp.Data.Result.ID, videoMuxing)
				if err != nil {
					return nil, vmErr
				}
				if videoMuxingResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding MP4 Muxing")
				}
			} else if container == "mov" {
				videoMuxingOutput := models.Output{
					OutputID:   s3OSResponse.Data.Result.ID,
					ACL:        acl,
					OutputPath: stringToPtr(path.Dir(path.Join(prefix, output.FileName))),
				}
				videoMuxing := &models.ProgressiveMOVMuxing{
					Filename:             stringToPtr(path.Base(output.FileName)),
					Outputs:              []models.Output{videoMuxingOutput},
					Streams:              []models.StreamItem{videoMuxingStream, uniqueAudioMuxingStreams[audioPresetID]},
					StreamConditionsMode: bitmovintypes.ConditionModeDropStream,
				}
				videoMuxingResp, vmErr := encodingS.AddProgressiveMOVMuxing(*encodingResp.Data.Result.ID, videoMuxing)
				if err != nil {
					return nil, vmErr
				}
				if videoMuxingResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding MP4 Muxing")
				}
			} else {
				return nil, errors.New("unknown container format")
			}
		} else {
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

				_, isRepeatedAudio := uniqueAudioMuxingStreams[audioPresetID]

				var audioStreamID, videoStreamID string

				if !isRepeatedAudio {

					audioStream := &models.Stream{
						CodecConfigurationID: &audioPresetID,
						InputStreams:         aiss,
						Conditions:           models.NewAttributeCondition(bitmovintypes.ConditionAttributeInputStream, "==", "true"),
					}
					audioStreamResp, audioErr := encodingS.AddStream(*encodingResp.Data.Result.ID, audioStream)
					if audioErr != nil {
						return nil, audioErr
					}
					if audioStreamResp.Status == bitmovinAPIErrorMsg {
						return nil, errors.New("Error in adding audio stream to Encoding")
					}
					audioStreamID = *audioStreamResp.Data.Result.ID

					audioMuxingStream := models.StreamItem{
						StreamID: &audioStreamID,
					}

					uniqueAudioMuxingStreams[audioPresetID] = audioMuxingStream
				}

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
					OutputPath: stringToPtr(path.Dir(path.Join(prefix, output.FileName))),
				}
				videoMuxing := &models.ProgressiveWebMMuxing{
					Filename:             stringToPtr(path.Base(output.FileName)),
					Outputs:              []models.Output{videoMuxingOutput},
					Streams:              []models.StreamItem{videoMuxingStream, uniqueAudioMuxingStreams[audioPresetID]},
					StreamConditionsMode: bitmovintypes.ConditionModeDropStream,
				}
				videoMuxingResp, vmErr := encodingS.AddProgressiveWebMMuxing(*encodingResp.Data.Result.ID, videoMuxing)
				if err != nil {
					return nil, vmErr
				}
				if videoMuxingResp.Status == bitmovinAPIErrorMsg {
					return nil, errors.New("Error in adding MP4 Muxing")
				}
			} else {
				return nil, errors.New("No H264 or VP8 codec configuration found")
			}
		}
	}
	startResp := &models.StartStopResponse{}
	if outputtingHLS {
		vodHLSManifest := models.VodHlsManifest{
			ManifestID: manifestID,
		}
		startOptions := &models.StartOptions{
			VodHlsManifests: []models.VodHlsManifest{vodHLSManifest},
		}
		startResp, err = encodingS.StartWithOptions(*encodingResp.Data.Result.ID, startOptions)
		if err != nil {
			return nil, err
		}
	} else {
		startResp, err = encodingS.Start(*encodingResp.Data.Result.ID)
		if err != nil {
			return nil, err
		}
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
	progress := 0.
	if statusResp.Data.Result.Progress != nil {
		progress = *statusResp.Data.Result.Progress
	}
	jobStatus := provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: job.ProviderJobID,
		Status:        p.mapStatus(stringValue(statusResp.Data.Result.Status)),
		Progress:      progress,
		ProviderStatus: map[string]interface{}{
			"message":        stringValue(statusResp.Data.Message),
			"originalStatus": stringValue(statusResp.Data.Result.Status),
		},
		Output: provider.JobOutput{
			Destination: strings.TrimRight(p.config.Destination, "/") + "/" + job.ID + "/",
		},
	}

	if jobStatus.Status == provider.StatusFinished {
		err = p.addManifestStatusInfo(&jobStatus)
		if err != nil {
			return nil, err
		}
	}

	if jobStatus.Status == provider.StatusFinished {
		err = p.addSourceInfo(job, &jobStatus)
		if err != nil {
			return nil, err
		}

		err = p.addOutputFilesInfo(job, &jobStatus)
		if err != nil {
			return nil, err
		}
	}

	return &jobStatus, nil
}

func (p *bitmovinProvider) addManifestStatusInfo(status *provider.JobStatus) error {
	encodingS := services.NewEncodingService(p.client)
	cdResp, err := encodingS.RetrieveCustomData(status.ProviderJobID)
	if err != nil {
		return err
	}
	if cdResp.Status == bitmovinAPIErrorMsg {
		return errors.New("No Custom Data on Encoding, there should at least be container information here")
	}
	cd := cdResp.Data.Result.CustomData
	_, ok := cd["manifest"]
	if !ok {
		// no manifest requested, no-op
		return nil
	}

	status.Status = p.mapStatus("FINISHED")
	status.ProviderStatus["manifestStatus"] = "FINISHED"

	return nil
}

func (p *bitmovinProvider) addOutputFilesInfo(job *db.Job, status *provider.JobStatus) error {
	err := p.addMP4OutputFilesInfo(job, status)
	if err != nil {
		return err
	}
	err = p.addWebmOutputFilesInfo(job, status)
	if err != nil {
		return err
	}
	return p.addMOVOutputFilesInfo(job, status)
}

func (p *bitmovinProvider) addMP4OutputFilesInfo(job *db.Job, status *provider.JobStatus) error {
	muxings, err := p.listMP4Muxing(job.ProviderJobID)
	if err != nil {
		return err
	}

	encodingS := services.NewEncodingService(p.client)
	var files []provider.OutputFile
	for _, muxing := range muxings {
		resp, err := encodingS.RetrieveMP4MuxingInformation(job.ProviderJobID, stringValue(muxing.ID))
		if err != nil {
			return err
		}

		info := resp.Data.Result
		if len(info.VideoTracks) == 0 {
			return fmt.Errorf("no video track found for encodingID %s muxingID %s", job.ProviderJobID, stringValue(muxing.ID))
		}

		files = append(files, provider.OutputFile{
			Path:       status.Output.Destination + stringValue(muxing.Filename),
			Container:  stringValue(info.ContainerFormat),
			FileSize:   int64Value(info.FileSize),
			VideoCodec: stringValue(info.VideoTracks[0].Codec),
			Width:      int64Value(info.VideoTracks[0].FrameWidth),
			Height:     int64Value(info.VideoTracks[0].FrameHeight),
		})
	}

	status.Output.Files = append(status.Output.Files, files...)
	return nil
}

func (p *bitmovinProvider) listMP4Muxing(jobID string) ([]models.MP4Muxing, error) {
	encodingS := services.NewEncodingService(p.client)

	var totalCount int64 = 1
	var muxings []models.MP4Muxing
	for int64(len(muxings)) < totalCount {
		resp, err := encodingS.ListMP4Muxing(jobID, int64(len(muxings)), 100)
		if err != nil {
			return nil, err
		}
		totalCount = int64Value(resp.Data.Result.TotalCount)
		muxings = append(muxings, resp.Data.Result.Items...)
	}

	return muxings, nil
}

func (p *bitmovinProvider) addWebmOutputFilesInfo(job *db.Job, status *provider.JobStatus) error {
	muxings, err := p.listWebmMuxing(job.ProviderJobID)
	if err != nil {
		return err
	}

	encodingS := services.NewEncodingService(p.client)
	for _, muxing := range muxings {
		resp, err := encodingS.RetrieveProgressiveWebMMuxingInformation(job.ProviderJobID, stringValue(muxing.ID))
		if err != nil {
			return err
		}

		info := resp.Data.Result
		if len(info.VideoTracks) == 0 {
			return fmt.Errorf("no video track found for encodingID %s muxingID %s", job.ProviderJobID, stringValue(muxing.ID))
		}

		status.Output.Files = append(status.Output.Files, provider.OutputFile{
			Path:       status.Output.Destination + stringValue(muxing.Filename),
			Container:  stringValue(info.ContainerFormat),
			FileSize:   int64Value(info.FileSize),
			VideoCodec: stringValue(info.VideoTracks[0].Codec),
			Width:      int64Value(info.VideoTracks[0].FrameWidth),
			Height:     int64Value(info.VideoTracks[0].FrameHeight),
		})
	}
	return nil
}

func (p *bitmovinProvider) listWebmMuxing(jobID string) ([]models.ProgressiveWebMMuxing, error) {
	encodingS := services.NewEncodingService(p.client)

	var totalCount int64 = 1
	var muxings []models.ProgressiveWebMMuxing
	for int64(len(muxings)) < totalCount {
		resp, err := encodingS.ListProgressiveWebMMuxing(jobID, int64(len(muxings)), 100)
		if err != nil {
			return nil, err
		}
		totalCount = int64Value(resp.Data.Result.TotalCount)
		muxings = append(muxings, resp.Data.Result.Items...)
	}

	return muxings, nil
}

func (p *bitmovinProvider) addMOVOutputFilesInfo(job *db.Job, status *provider.JobStatus) error {
	muxings, err := p.listMOVMuxing(job.ProviderJobID)
	if err != nil {
		return err
	}

	encodingS := services.NewEncodingService(p.client)
	for _, muxing := range muxings {
		resp, err := encodingS.RetrieveProgressiveMOVMuxingInformation(job.ProviderJobID, stringValue(muxing.ID))
		if err != nil {
			return err
		}

		info := resp.Data.Result
		if len(info.VideoTracks) == 0 {
			return fmt.Errorf("no video track found for encodingID %s muxingID %s", job.ProviderJobID, stringValue(muxing.ID))
		}

		status.Output.Files = append(status.Output.Files, provider.OutputFile{
			Path:       status.Output.Destination + stringValue(muxing.Filename),
			Container:  stringValue(info.ContainerFormat),
			FileSize:   int64Value(info.FileSize),
			VideoCodec: stringValue(info.VideoTracks[0].Codec),
			Width:      int64Value(info.VideoTracks[0].FrameWidth),
			Height:     int64Value(info.VideoTracks[0].FrameHeight),
		})
	}
	return nil
}

func (p *bitmovinProvider) listMOVMuxing(jobID string) ([]models.ProgressiveMOVMuxing, error) {
	encodingS := services.NewEncodingService(p.client)

	var totalCount int64 = 1
	var muxings []models.ProgressiveMOVMuxing
	for int64(len(muxings)) < totalCount {
		resp, err := encodingS.ListProgressiveMOVMuxing(jobID, int64(len(muxings)), 100)
		if err != nil {
			return nil, err
		}
		totalCount = int64Value(resp.Data.Result.TotalCount)
		muxings = append(muxings, resp.Data.Result.Items...)
	}

	return muxings, nil
}

func (p *bitmovinProvider) addSourceInfo(job *db.Job, status *provider.JobStatus) error {
	encodingS := services.NewEncodingService(p.client)
	resp, err := encodingS.ListStream(job.ProviderJobID, 0, 1)
	if err != nil {
		return err
	}
	if len(resp.Data.Result.Items) == 0 {
		return fmt.Errorf("no stream item found for encodingID %s", job.ProviderJobID)
	}

	streamID := stringValue(resp.Data.Result.Items[0].ID)
	streamInput, err := encodingS.RetrieveStreamInputData(job.ProviderJobID, streamID)
	if err != nil {
		return err
	}
	if len(streamInput.Data.Result.VideoStreams) == 0 {
		return fmt.Errorf("no video stream input found for encodingID %s streamID %s", job.ProviderJobID, streamID)
	}

	videoStream := streamInput.Data.Result.VideoStreams[0]
	status.SourceInfo = provider.SourceInfo{
		Duration:   time.Duration(floatValue(streamInput.Data.Result.Duration) * float64(time.Second)),
		Width:      int64Value(videoStream.Width),
		Height:     int64Value(videoStream.Height),
		VideoCodec: stringValue(videoStream.Codec),
	}
	return nil
}

func (p *bitmovinProvider) mapStatus(status string) provider.Status {
	switch status {
	case "CREATED", "QUEUED":
		return provider.StatusQueued
	case "RUNNING":
		return provider.StatusStarted
	case "FINISHED":
		return provider.StatusFinished
	case "ERROR":
		return provider.StatusFailed
	default:
		return provider.StatusUnknown
	}
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
	return provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "mov", "hls", "webm"},
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

func parseS3URL(input string) (bucketName string, objectKey string, err error) {
	s3URL, err := url.Parse(input)
	if err != nil || s3URL.Scheme != "s3" {
		return "", "", errors.New("Could not parse S3 URL")
	}
	return s3URL.Host, strings.TrimLeft(s3URL.Path, "/"), nil
}

func createInput(provider *bitmovinProvider, input string) (inputID string, path string, err error) {
	if s3Pattern.MatchString(input) {
		inputMediaBucketName, inputMediaPath, err := parseS3URL(input)
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
		return *s3ISResponse.Data.Result.ID, inputMediaPath, nil
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

func getAACConfig(s *services.AACCodecConfigurationService, bitrate int, samplingRate float64) (string, error) {
	l, err := s.List(0, 25)
	if err != nil {
		return "", err
	}

	for _, i := range l.Data.Result.Items {
		if *i.Bitrate == int64(bitrate) && *i.SamplingRate == samplingRate {
			return *i.ID, nil
		}
	}

	audioConfig := &models.AACCodecConfiguration{
		Name:         stringToPtr("unused"),
		Bitrate:      intToPtr(int64(bitrate)),
		SamplingRate: floatToPtr(samplingRate),
	}

	audioResp, err := s.Create(audioConfig)
	if err != nil {
		return "", err
	}
	if audioResp.Status == bitmovinAPIErrorMsg {
		return "", errors.New("Error in creating audio portion of Preset")
	}
	return *audioResp.Data.Result.ID, nil

}

func getVorbisConfig(s *services.VorbisCodecConfigurationService, bitrate int, samplingRate float64) (string, error) {
	l, err := s.List(0, 25)
	if err != nil {
		return "", err
	}

	for _, i := range l.Data.Result.Items {
		if *i.Bitrate == int64(bitrate) && *i.SamplingRate == samplingRate {
			return *i.ID, nil
		}
	}
	audioConfig := &models.VorbisCodecConfiguration{
		Name:         stringToPtr("unused"),
		Bitrate:      intToPtr(int64(bitrate)),
		SamplingRate: floatToPtr(samplingRate),
	}

	audioResp, err := s.Create(audioConfig)
	if err != nil {
		return "", err
	}
	if audioResp.Status == bitmovinAPIErrorMsg {
		return "", errors.New("Error in creating audio portion of Preset")
	}
	return *audioResp.Data.Result.ID, nil
}

func floatValue(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func int64Value(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func stringValue(str *string) string {
	if str == nil {
		return ""
	}
	return *str
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
