package bitmovin

import (
	"errors"

	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/video-transcoding-api/config"
)

// Name is the name used for registering the bitmovin provider in the
// registry of providers.
const Name = "bitmovin"

var errBitmovinInvalidConfig = provider.InvalidConfigError("missing Bitmovin api key. Please define the environment variable BITMOVIN_API_KEY set this value in the configuration file")

type bitmovinProvider struct {
	client *bitmovin.Bitmovin
	config *config.Bitmovin
}

func (p *bitmovinProvider) CreatePreset(db.Preset) (string, error) {
	return "", errors.New("Not implemented")
}

func (p *bitmovinProvider) DeletePreset(presetID string) error {
	return errors.New("Not implemented")
}

func (p *bitmovinProvider) GetPreset(presetID string) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

func (p *bitmovinProvider) Transcode(*db.Job) (*provider.JobStatus, error) {
	return nil, errors.New("Not implemented")
}

func (p *bitmovinProvider) JobStatus(*db.Job) (*provider.JobStatus, error) {
	return nil, errors.New("Not implemented")
}

func (p *bitmovinProvider) CancelJob(jobID string) error {
	return errors.New("Not implemented")
}

func (p *bitmovinProvider) Healthcheck() error {
	return errors.New("Not implemented")
}

func (p *bitmovinProvider) Capabilities() provider.Capabilities {
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
	client := bitmovin.NewBitmovin(cfg.Bitmovin.APIKey, cfg.Bitmovin.Endpoint, int64(cfg.Bitmovin.Timeout))
	return &bitmovinProvider{client: client, config: cfg.Bitmovin}, nil
}
