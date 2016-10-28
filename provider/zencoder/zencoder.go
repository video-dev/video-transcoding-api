// Package encodingcom provides a implementation of the provider that uses the
// Encoding.com API for transcoding media files.
//
// It doesn't expose any public type. In order to use the provider, one must
// import this package and then grab the factory from the provider package:
//
//     import (
//         "github.com/NYTimes/video-transcoding-api/provider"
//         "github.com/NYTimes/video-transcoding-api/provider/encodingcom"
//     )
//
//     func UseProvider() {
//         factory, err := provider.GetProviderFactory(encodingcom.Name)
//         // handle err and use factory to get an instance of the provider.
//     }
package zencoder

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/brandscreen/zencoder"
)

// Name is the name used for registering the Encoding.com provider in the
// registry of providers.
const Name = "zencoder"

var errZencoderInvalidConfig = provider.InvalidConfigError("missing Zencoder API key. Please define the environment variables ZENCODER_API_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, zencoderFactory)
}

type zencoderProvider struct {
	config *config.Config
	client *zencoder.Zencoder
}

func (z *zencoderProvider) Transcode(job *db.Job, transcodeProfile provider.TranscodeProfile) (*provider.JobStatus, error) {
	return &provider.JobStatus{}, nil
}

func (z *zencoderProvider) CreatePreset(preset db.Preset) (string, error) {
	return "", nil
}

func (z *zencoderProvider) GetPreset(presetID string) (interface{}, error) {
	return "", nil
}

func (z *zencoderProvider) DeletePreset(presetID string) error {
	return nil
}

func (z *zencoderProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
}

func (z *zencoderProvider) CancelJob(id string) error {
}

func (z *zencoderProvider) Healthcheck() error {
}

func (z *zencoderProvider) Capabilities() provider.Capabilities {
}

func zencoderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.Zencoder.ApiKey == "" {
		return nil, errZencoderInvalidConfig
	}
	client, err := zencoder.NewZencoder(cfg.Zencoder.ApiKey)
	if err != nil {
		return nil, err
	}
	return &zencoderProvider{client: client, config: cfg}, nil
}
