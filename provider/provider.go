package provider

import "github.com/nytm/video-transcoding-api/config"

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

}

func (e *encodingComProvider) JobStatus(id string) (*JobStatus, error) {

}

type ProviderFactory func(cfg *config.Config) (Provider, error)

func EncodingComProvider(cfg *config.Config) (Provider, error) {
	// add validation
	// create client
	return &encodingComProvider{
		client: &client,
	}, nil
}
