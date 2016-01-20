package provider

import (
	"errors"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
)

// ErrMissingData is the error returned by the factory when required data is
// missing.
var ErrMissingData = errors.New("missing Encoding.com user id or key. Please define the environment variables ENCODINGCOM_USER_ID and ENCODINGCOM_USER_KEY")

type encodingComProvider struct {
	client *encodingcom.Client
}

func (e *encodingComProvider) Transcode(sourceMedia, destination string, profile Profile) (*JobStatus, error) {
	return nil, nil
}

func (e *encodingComProvider) JobStatus(id string) (*JobStatus, error) {
	return nil, nil
}

// EncodingComProvider is the factory function for the Encoding.com provider.
func EncodingComProvider(cfg *config.Config) (TranscodingProvider, error) {
	if cfg.EncodingComUserID == "" || cfg.EncodingComUserKey == "" {
		return nil, ErrMissingData
	}
	client, err := encodingcom.NewClient("https://manage.encoding.com", cfg.EncodingComUserID, cfg.EncodingComUserKey)
	if err != nil {
		return nil, err
	}
	return &encodingComProvider{client: client}, nil
}
