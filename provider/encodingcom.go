package provider

import (
	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
)

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
	// add validation
	// create client
	var client encodingcom.Client
	return &encodingComProvider{
		client: &client,
	}, nil
}
