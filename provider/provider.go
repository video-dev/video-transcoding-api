package provider

import (
	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
)

type TranscodingProvider interface {
	Transcode(sourceMedia string, profileSpec []byte) (*JobStatus, error)
	JobStatus(id string) (*JobStatus, error)
}

type JobStatus struct {
	ProviderJobID  string
	Status         string
	ProviderName   string
	StatusMessage  string
	ProviderStatus map[string]interface{}
}

type encodingComProvider struct {
	client *encodingcom.Client
}

func (e *encodingComProvider) Transcode(sourceMedia string, profileSpec []byte) (*JobStatus, error) {
	return nil, nil
}

func (e *encodingComProvider) JobStatus(id string) (*JobStatus, error) {
	return nil, nil
}

type ProviderFactory func(cfg *config.Config) (TranscodingProvider, error)

func EncodingComProvider(cfg *config.Config) (TranscodingProvider, error) {
	// add validation
	// create client
	var client encodingcom.Client
	return &encodingComProvider{
		client: &client,
	}, nil
}
