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

var kregexp = regexp.MustCompile(`000$`)

var errEncodingComInvalidConfig = provider.InvalidConfigError("missing Encoding.com user id or key. Please define the environment variables ENCODINGCOM_USER_ID and ENCODINGCOM_USER_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, encodingComFactory)
}

type encodingComClient interface {
	AddMedia(source []string, format []encodingcom.Format, Region string) (*encodingcom.AddMediaResponse, error)
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
		return nil, err
	}
	resp, err := e.client.AddMedia([]string{transcodeProfile.SourceMedia}, formats, e.config.EncodingCom.Region)
	if err != nil {
		return nil, fmt.Errorf("Error on AddMedia operation: %s", err.Error())
	}
	return &provider.JobStatus{
		ProviderJobID: resp.MediaID,
		StatusMessage: resp.Message,
		ProviderName:  Name,
	}, nil
}

func (e *encodingComProvider) CreatePreset(preset provider.Preset) (string, error) {
	resp, err := e.client.SavePreset("", e.presetToFormat(preset))
	if err != nil {
		return "", err
	}
	return resp.SavedPreset, nil
}

func (e *encodingComProvider) presetToFormat(preset provider.Preset) encodingcom.Format {
	format := encodingcom.Format{
		Output:       []string{preset.Container},
		Profile:      preset.Profile,
		Bitrate:      kregexp.ReplaceAllString(preset.Video.Bitrate, "k"),
		VideoCodec:   preset.Video.Codec,
		AudioBitrate: kregexp.ReplaceAllString(preset.Audio.Bitrate, "k"),
		AudioCodec:   preset.Audio.Codec,
		AudioVolume:  100,
		Gop:          "cgop",
		Keyframe:     []string{preset.Video.GopSize},
	}
	if preset.Container == "m3u8" {
		format.Output = []string{"advanced_hls"}
	}
	if format.AudioCodec == "aac" {
		format.AudioCodec = "dolby_aac"
	}
	if format.VideoCodec == "h264" {
		format.VideoCodec = "libx264"
	}
	width := preset.Video.Width
	height := preset.Video.Height
	if width == "" {
		width = "0"
	}
	if height == "" {
		height = "0"
	}
	if preset.RateControl == "VBR" {
		format.TwoPass = true
	}
	format.Size = width + "x" + height
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

func (e *encodingComProvider) getDestinations(jobID, sourceMedia string, preset db.PresetMap) []string {
	var extension string

	if preset.OutputOpts.Extension == "" {
		extension = "." + filepath.Ext(sourceMedia)
	} else {
		extension = "." + preset.OutputOpts.Extension
	}

	sourceParts := strings.Split(sourceMedia, "/")
	sourceFilenamePart := sourceParts[len(sourceParts)-1]
	sourceFileName := strings.TrimSuffix(sourceFilenamePart, filepath.Ext(sourceFilenamePart))
	outputDestination := strings.TrimRight(e.config.EncodingCom.Destination, "/") + "/" + path.Join(jobID, preset.Name) + "/"
	if preset.OutputOpts.Extension == "m3u8" {
		return []string{outputDestination + sourceFileName + "/master.m3u8"}
	}
	return []string{outputDestination + sourceFileName + extension}
}

func (e *encodingComProvider) presetsToFormats(job *db.Job, transcodeProfile provider.TranscodeProfile) ([]encodingcom.Format, error) {
	formats := make([]encodingcom.Format, 0, len(transcodeProfile.Presets))
	for _, preset := range transcodeProfile.Presets {
		presetName, ok := preset.ProviderMapping[Name]
		if !ok {
			return nil, provider.ErrPresetMapNotFound
		}
		format := encodingcom.Format{
			Output:          []string{presetName},
			Destination:     e.getDestinations(job.ID, transcodeProfile.SourceMedia, preset),
			SegmentDuration: transcodeProfile.StreamingParams.SegmentDuration,
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
			"destinationStatus": resp[0].Formats[0].Destinations,
		},
		OutputDestination: e.getOutputDestination(resp),
	}, nil
}

func (e *encodingComProvider) getOutputDestination(status []encodingcom.StatusResponse) string {
	formats := status[0].Formats
	for _, formatStatus := range formats {
		for _, destinationStatus := range formatStatus.Destinations {
			if destinationStatus.Status == "Saved" {
				destination := strings.Split(destinationStatus.Name, "/")
				destination = destination[:len(destination)-1]
				return strings.Join(destination, "/")
			}
		}
	}
	return ""
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
