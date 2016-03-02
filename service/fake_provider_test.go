package service

import (
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

func init() {
	provider.Register("fake", fakeProviderFactory)
}

type fakeProvider struct{}

func (e *fakeProvider) JobStatus(id string) (*provider.JobStatus, error) {
	if id == "provider-job-123" {
		return &provider.JobStatus{
			ProviderJobID: "provider-job-123",
			Status:        provider.StatusFinished,
			StatusMessage: "The job is finished",
			ProviderStatus: map[string]interface{}{
				"progress":   100.0,
				"sourcefile": "http://some.source.file",
			},
		}, nil
	}
	return nil, provider.JobNotFoundError{ID: id}
}

func (e *fakeProvider) Healthcheck() error {
	return nil
}

func (e *fakeProvider) Transcode(transcodeProfile provider.TranscodeProfile) (*provider.JobStatus, error) {
	for _, preset := range transcodeProfile.Presets {
		if _, ok := preset.ProviderMapping["fake"]; !ok {
			return nil, provider.ErrPresetNotFound
		}
	}
	return &provider.JobStatus{
		ProviderJobID: "provider-preset-job-123",
		Status:        provider.StatusFinished,
		StatusMessage: "The job is finished",
		ProviderStatus: map[string]interface{}{
			"progress":   100.0,
			"sourcefile": "http://some.source.file",
		},
	}, nil
}

func fakeProviderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	return &fakeProvider{}, nil
}
