// Package encodingcom provides a implementation of the provider that uses the
// Encoding.com API for transcoding media files.
//
// It doesn't expose any public type. In order to use the provider, one must
// import this package and then grab the factory from the provider package:
//
//     import (
//         "github.com/nytm/video-transcoding-api/provider"
//         "github.com/nytm/video-transcoding-api/provider/encodingcom"
//     )
//
//     func UseProvider() {
//         factory, err := provider.GetProviderFactory(encodingcom.Name)
//         // handle err and use factory to get an instance of the provider.
//     }
package encodingcom

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// Name is the name used for registering the Encoding.com provider in the
// registry of providers.
const Name = "encodingcom"

var (
	kregexp      = regexp.MustCompile(`000$`)
	s3regexp     = regexp.MustCompile(`^s3://([^/_.]+)/(.+)$`)
	httpS3Regexp = regexp.MustCompile(`https?://([^/_.]+)\.s3\.amazonaws\.com/(.+)$`)
)

var errEncodingComInvalidConfig = provider.InvalidConfigError("missing Encoding.com user id or key. Please define the environment variables ENCODINGCOM_USER_ID and ENCODINGCOM_USER_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, encodingComFactory)
}

type encodingComClient interface {
	AddMedia(source []string, format []encodingcom.Format, Region string) (*encodingcom.AddMediaResponse, error)
	CancelMedia(string) (*encodingcom.Response, error)
	GetStatus(mediaIDs []string) ([]encodingcom.StatusResponse, error)
	SavePreset(name string, format encodingcom.Format) (*encodingcom.SavePresetResponse, error)
	GetPreset(name string) (*encodingcom.Preset, error)
	DeletePreset(name string) (*encodingcom.Response, error)
}

type encodingComProvider struct {
	config *config.Config
	client encodingComClient
}

func (e *encodingComProvider) Transcode(job *db.Job, transcodeProfile provider.TranscodeProfile) (*provider.JobStatus, error) {
	formats, err := e.presetsToFormats(job, transcodeProfile)
	if err != nil {
		return nil, fmt.Errorf("Error converting presets to formats on Transcode operation: %s", err.Error())
	}
	resp, err := e.client.AddMedia([]string{e.sourceMedia(transcodeProfile.SourceMedia)}, formats, e.config.EncodingCom.Region)
	if err != nil {
		return nil, fmt.Errorf("Error making AddMedia request for Transcode operation: %s", err.Error())
	}
	return &provider.JobStatus{
		ProviderJobID: resp.MediaID,
		StatusMessage: resp.Message,
		ProviderName:  Name,
	}, nil
}

func (e *encodingComProvider) CreatePreset(preset provider.Preset) (string, error) {
	resp, err := e.client.SavePreset(preset.Name, e.presetToFormat(preset))
	if err != nil {
		return "", err
	}
	return resp.SavedPreset, nil
}

func (e *encodingComProvider) sourceMedia(original string) string {
	parts := s3regexp.FindStringSubmatch(original)
	if len(parts) > 0 {
		return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", parts[1], parts[2])
	}
	return original
}

func (e *encodingComProvider) presetToFormat(preset provider.Preset) encodingcom.Format {
	falseYesNoBoolean := encodingcom.YesNoBoolean(false)
	format := encodingcom.Format{
		Output:      []string{preset.Container},
		Destination: []string{"ftp://username:password@yourftphost.com/video/encoded/test.flv"},
	}
	if preset.Container == "m3u8" {
		format.Output = []string{"advanced_hls"}
		format.PackFiles = &falseYesNoBoolean
		stream := encodingcom.Stream{
			Profile:      preset.Profile,
			Keyframe:     preset.Video.GopSize,
			Bitrate:      kregexp.ReplaceAllString(preset.Video.Bitrate, "k"),
			VideoCodec:   preset.Video.Codec,
			AudioBitrate: kregexp.ReplaceAllString(preset.Audio.Bitrate, "k"),
			AudioCodec:   preset.Audio.Codec,
			AudioVolume:  100,
		}
		if stream.AudioCodec == "aac" {
			stream.AudioCodec = "dolby_aac"
		}
		if stream.VideoCodec == "h264" {
			stream.VideoCodec = "libx264"
		}
		if preset.RateControl == "VBR" {
			stream.TwoPass = true
		}
		width := preset.Video.Width
		height := preset.Video.Height
		if width == "" {
			width = "0"
		}
		if height == "" {
			height = "0"
		}
		stream.Size = width + "x" + height
		format.Stream = []encodingcom.Stream{stream}
	} else {
		format.Bitrate = kregexp.ReplaceAllString(preset.Video.Bitrate, "k")
		format.AudioBitrate = kregexp.ReplaceAllString(preset.Audio.Bitrate, "k")
		format.AudioCodec = preset.Audio.Codec
		format.VideoCodec = preset.Video.Codec
		format.Profile = preset.Profile
		format.Gop = "cgop"
		format.Keyframe = []string{preset.Video.GopSize}
		format.AudioVolume = 100

		if format.AudioCodec == "aac" {
			format.AudioCodec = "dolby_aac"
		}
		if format.VideoCodec == "h264" {
			format.VideoCodec = "libx264"
		}
		if preset.RateControl == "VBR" {
			format.TwoPass = true
		}
		width := preset.Video.Width
		height := preset.Video.Height
		if width == "" {
			width = "0"
		}
		if height == "" {
			height = "0"
		}
		format.Size = width + "x" + height
	}
	return format
}

func (e *encodingComProvider) GetPreset(presetID string) (interface{}, error) {
	preset, err := e.client.GetPreset(presetID)
	if err != nil {
		return nil, err
	}
	return preset, nil
}

func (e *encodingComProvider) DeletePreset(presetID string) error {
	_, err := e.client.DeletePreset(presetID)
	return err
}

func (e *encodingComProvider) getDestinations(jobID string, transcodeProfile provider.TranscodeProfile, preset db.PresetMap) []string {
	destination := e.buildDestination(
		e.config.EncodingCom.Destination,
		jobID,
		transcodeProfile.OutputPath,
		transcodeProfile.OutputFilePrefix,
		preset.OutputOpts.Label,
		preset.OutputOpts.Extension,
	)
	return []string{destination}
}

func (e *encodingComProvider) buildDestination(outputDestination string, jobID string, outputDestinationPath string, outputFilePrefix string, presetLabel string, extension string) string {
	outputPath := strings.TrimRight(outputDestinationPath, "/")
	destinationPathWithPrefix := path.Join(jobID, outputPath, outputFilePrefix)
	outputFile := destinationPathWithPrefix + "_" + presetLabel + "." + extension
	if extension == "m3u8" {
		outputFile = destinationPathWithPrefix + "_" + presetLabel + "/video.m3u8"
	}
	return strings.TrimRight(outputDestination, "/") + "/" + outputFile
}

func (e *encodingComProvider) presetsToFormats(job *db.Job, transcodeProfile provider.TranscodeProfile) ([]encodingcom.Format, error) {
	streams := []encodingcom.Stream{}
	streamingPresetDestinations := []string{}
	formats := make([]encodingcom.Format, 0, len(transcodeProfile.Presets))
	for _, preset := range transcodeProfile.Presets {
		presetName := preset.Name
		presetID, ok := preset.ProviderMapping[Name]
		if !ok {
			return nil, provider.ErrPresetMapNotFound
		}
		presetOutput, err := e.GetPreset(presetID)
		if err != nil {
			return nil, fmt.Errorf("Error getting preset info: %s", err.Error())
		}
		presetStruct := presetOutput.(*encodingcom.Preset)
		if presetStruct.Output == "advanced_hls" {
			for _, stream := range presetStruct.Format.Stream() {
				stream.SubPath = presetName
				streams = append(streams, stream)
			}
			destination := e.getDestinations(job.ID, transcodeProfile, preset)
			streamingPresetDestinations = append(streamingPresetDestinations, destination[0])
		} else {
			extension := preset.OutputOpts.Extension
			if extension == "" {
				extension = "." + filepath.Ext(transcodeProfile.SourceMedia)
			}
			format := encodingcom.Format{
				Output:      []string{presetID},
				Destination: e.getDestinations(job.ID, transcodeProfile, preset),
			}
			formats = append(formats, format)
		}
	}
	if len(streams) > 0 {
		falseValue := encodingcom.YesNoBoolean(false)
		format := encodingcom.Format{
			Output:          []string{"advanced_hls"},
			Destination:     streamingPresetDestinations,
			SegmentDuration: transcodeProfile.StreamingParams.SegmentDuration,
			Stream:          streams,
			PackFiles:       &falseValue,
		}
		formats = append(formats, format)
	}
	return formats, nil
}

func (e *encodingComProvider) JobStatus(id string) (*provider.JobStatus, error) {
	resp, err := e.client.GetStatus([]string{id})
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, errors.New("invalid value returned by the Encoding.com API: []")
	}
	return &provider.JobStatus{
		ProviderJobID: id,
		ProviderName:  "encoding.com",
		Status:        e.statusMap(resp[0].MediaStatus),
		ProviderStatus: map[string]interface{}{
			"progress":          resp[0].Progress,
			"sourcefile":        resp[0].SourceFile,
			"timeleft":          resp[0].TimeLeft,
			"created":           resp[0].CreateDate,
			"started":           resp[0].StartDate,
			"finished":          resp[0].FinishDate,
			"formatStatus":      e.getFormatStatus(resp),
			"destinationStatus": e.getOutputDestinationStatus(resp),
		},
		OutputDestination: e.getOutputDestination(resp),
	}, nil
}

func (e *encodingComProvider) getFormatStatus(status []encodingcom.StatusResponse) []string {
	formatStatusList := []string{}
	formats := status[0].Formats
	for _, formatStatus := range formats {
		formatStatusList = append(formatStatusList, formatStatus.Status)
	}
	return formatStatusList
}

func (e *encodingComProvider) getOutputDestinationStatus(status []encodingcom.StatusResponse) []encodingcom.DestinationStatus {
	destinationStatusList := []encodingcom.DestinationStatus{}
	formats := status[0].Formats
	for _, formatStatus := range formats {
		for _, destinationStatus := range formatStatus.Destinations {
			destinationStatus.Name = e.destinationMedia(destinationStatus.Name)
			destinationStatusList = append(destinationStatusList, destinationStatus)
		}
	}
	return destinationStatusList
}

func (e *encodingComProvider) getOutputDestination(status []encodingcom.StatusResponse) string {
	formats := status[0].Formats
	for _, formatStatus := range formats {
		for _, destinationStatus := range formatStatus.Destinations {
			if destinationStatus.Status == "Saved" {
				destination := strings.Split(destinationStatus.Name, "/")
				return e.destinationMedia(strings.Join(destination[:len(destination)-1], "/"))
			}
		}
	}
	return ""
}

func (e *encodingComProvider) destinationMedia(input string) string {
	parts := httpS3Regexp.FindStringSubmatch(input)
	if len(parts) > 0 {
		return fmt.Sprintf("s3://%s/%s", parts[1], parts[2])
	}
	return input
}

func (e *encodingComProvider) statusMap(encodingComStatus string) provider.Status {
	switch strings.ToLower(encodingComStatus) {
	case "new":
		return provider.StatusQueued
	case "downloading", "ready to process", "waiting for encoder", "processing", "saving":
		return provider.StatusStarted
	case "finished":
		return provider.StatusFinished
	case "error":
		return provider.StatusFailed
	default:
		return provider.StatusUnknown
	}
}

func (e *encodingComProvider) CancelJob(id string) error {
	_, err := e.client.CancelMedia(id)
	return err
}

func (e *encodingComProvider) Healthcheck() error {
	status, err := encodingcom.APIStatus(e.config.EncodingCom.StatusEndpoint)
	if err != nil {
		return err
	}
	if !status.OK() {
		return fmt.Errorf("Status code: %s.\nIncident: %s\nStatus: %s", status.StatusCode, status.Incident, status.Status)
	}
	return nil
}

func (e *encodingComProvider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls", "webm"},
		Destinations:  []string{"akamai", "s3"},
	}
}

func encodingComFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.EncodingCom.UserID == "" || cfg.EncodingCom.UserKey == "" {
		return nil, errEncodingComInvalidConfig
	}
	client, err := encodingcom.NewClient("https://manage.encoding.com", cfg.EncodingCom.UserID, cfg.EncodingCom.UserKey)
	if err != nil {
		return nil, err
	}
	return &encodingComProvider{client: client, config: cfg}, nil
}
