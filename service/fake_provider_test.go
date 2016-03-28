package service

import (
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

func init() {
	provider.Register("fake", fakeProviderFactory)
}

type fakeProvider struct{}

func (p *fakeProvider) Transcode(transcodeProfile provider.TranscodeProfile) (*provider.JobStatus, error) {
	for _, preset := range transcodeProfile.Presets {
		if _, ok := preset.ProviderMapping["fake"]; !ok {
			return nil, provider.ErrPresetMapNotFound
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

func (*fakeProvider) CreatePreset(preset provider.Preset) (string, error) {
	return "presetID_here", nil
}

func (p *fakeProvider) JobStatus(id string) (*provider.JobStatus, error) {
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

func (p *fakeProvider) Healthcheck() error {
	return nil
}

func (p *fakeProvider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "webm", "hls"},
		Destinations:  []string{"akamai", "s3"},
	}
}

func fakeProviderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	return &fakeProvider{}, nil
}
