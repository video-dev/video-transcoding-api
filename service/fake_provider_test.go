package service

import (
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

func init() {
	provider.Register("fake", fakeProviderFactory)
	provider.Register("profile-fake", profileFakeProviderFactory)
	provider.Register("preset-fake", presetFakeProviderFactory)
}

type baseFakeProvider struct{}

func (e *baseFakeProvider) JobStatus(id string) (*provider.JobStatus, error) {
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

type profileFakeProvider struct {
	*baseFakeProvider
}

func (e *profileFakeProvider) TranscodeWithProfiles(sourceMedia string, profiles []provider.Profile) (*provider.JobStatus, error) {
	return &provider.JobStatus{
		ProviderJobID: "provider-profile-job-123",
		Status:        provider.StatusFinished,
		StatusMessage: "The job is finished",
		ProviderStatus: map[string]interface{}{
			"progress":   100.0,
			"sourcefile": "http://some.source.file",
		},
	}, nil
}

type presetFakeProvider struct {
	*baseFakeProvider
}

func (e *presetFakeProvider) TranscodeWithPresets(sourceMedia string, presets []string, adaptiveStreaming bool) (*provider.JobStatus, error) {
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

type fakeProvider struct {
	*baseFakeProvider
	*profileFakeProvider
	*presetFakeProvider
}

func fakeProviderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	return &fakeProvider{}, nil
}

func profileFakeProviderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	return &profileFakeProvider{}, nil
}

func presetFakeProviderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	return &presetFakeProvider{}, nil
}
