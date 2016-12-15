package service

import (
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
)

func init() {
	provider.Register("fake", fakeProviderFactory)
	provider.Register("zencoder", fakeProviderFactory)
}

type fakeProvider struct {
	jobs         []*db.Job
	canceledJobs []string
}

var fprovider fakeProvider

func (p *fakeProvider) Transcode(job *db.Job) (*provider.JobStatus, error) {
	for _, output := range job.Outputs {
		if _, ok := output.Preset.ProviderMapping["fake"]; !ok {
			return nil, provider.ErrPresetMapNotFound
		}
	}
	p.jobs = append(p.jobs, job)
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

func (*fakeProvider) CreatePreset(preset db.Preset) (string, error) {
	return "presetID_here", nil
}

func (*fakeProvider) GetPreset(presetID string) (interface{}, error) {
	return struct{ presetID string }{"presetID_here"}, nil
}

func (*fakeProvider) DeletePreset(presetID string) error {
	return nil
}

func (p *fakeProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	id := job.ProviderJobID
	if id == "provider-job-123" {
		status := provider.StatusFinished
		if len(p.canceledJobs) > 0 {
			status = provider.StatusCanceled
		}
		return &provider.JobStatus{
			ProviderJobID: "provider-job-123",
			Status:        status,
			StatusMessage: "The job is finished",
			Progress:      10.3,
			SourceInfo: provider.SourceInfo{
				Width:      4096,
				Height:     2160,
				Duration:   183e9,
				VideoCodec: "VP9",
			},
			ProviderStatus: map[string]interface{}{
				"progress":   10.3,
				"sourcefile": "http://some.source.file",
			},
			Output: provider.JobOutput{
				Destination: "s3://mybucket/some/dir/job-123",
			},
		}, nil
	}
	return nil, provider.JobNotFoundError{ID: id}
}

func (p *fakeProvider) CancelJob(id string) error {
	if id == "provider-job-123" {
		p.canceledJobs = append(p.canceledJobs, id)
		return nil
	}
	return provider.JobNotFoundError{ID: id}
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
	return &fprovider, nil
}
