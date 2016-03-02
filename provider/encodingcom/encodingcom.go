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
	"path/filepath"
	"strings"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// Name is the name used for registering the Encoding.com provider in the
// registry of providers.
const Name = "encodingcom"

var errEncodingComInvalidConfig = provider.InvalidConfigError("missing Encoding.com user id or key. Please define the environment variables ENCODINGCOM_USER_ID and ENCODINGCOM_USER_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, encodingComFactory)
}

type encodingComProvider struct {
	config *config.Config
	client *encodingcom.Client
}

func (e *encodingComProvider) Transcode(sourceMedia string, presets []db.Preset) (*provider.JobStatus, error) {
	formats, err := e.presetsToFormats(sourceMedia, presets)
	if err != nil {
		return nil, err
	}
	resp, err := e.client.AddMedia([]string{sourceMedia}, formats)
	if err != nil {
		return nil, err
	}
	return &provider.JobStatus{
		ProviderJobID: resp.MediaID,
		StatusMessage: resp.Message,
		ProviderName:  Name,
	}, nil
}

func (e *encodingComProvider) getDestinations(sourceMedia string, preset db.Preset) []string {
	var extension string
	switch v := preset.OutputOpts.Extension; v {
	case "hls", "ts", "m3u8":
		extension = ".m3u8"
	case "":
		extension = "." + filepath.Ext(sourceMedia)
	default:
		extension = "." + v
	}
	sourceParts := strings.Split(sourceMedia, "/")
	sourceFilenamePart := sourceParts[len(sourceParts)-1]
	sourceFileName := strings.TrimSuffix(sourceFilenamePart, filepath.Ext(sourceFilenamePart))
	outputDestination := strings.TrimRight(e.config.EncodingCom.Destination, "/") + "/" + preset.Name + "/"
	return []string{outputDestination + sourceFileName + extension}
}

func (e *encodingComProvider) presetsToFormats(sourceMedia string, presets []db.Preset) ([]encodingcom.Format, error) {
	formats := make([]encodingcom.Format, 0, len(presets))
	for _, preset := range presets {
		presetName, ok := preset.ProviderMapping[Name]
		if !ok {
			return nil, provider.ErrPresetNotFound
		}
		format := encodingcom.Format{
			Output:      []string{presetName},
			Destination: e.getDestinations(sourceMedia, preset),
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
	}, nil
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
