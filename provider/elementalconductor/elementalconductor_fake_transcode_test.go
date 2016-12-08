package elementalconductor

import (
	"strings"

	"github.com/NYTimes/encoding-wrapper/elementalconductor"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/provider"
)

type fakeElementalConductorClient struct {
	*elementalconductor.Client
	jobs         map[string]elementalconductor.Job
	canceledJobs []string
}

func newFakeElementalConductorClient(cfg *config.ElementalConductor) *fakeElementalConductorClient {
	return &fakeElementalConductorClient{
		jobs: make(map[string]elementalconductor.Job),
		Client: &elementalconductor.Client{
			Host:            cfg.Host,
			UserLogin:       cfg.UserLogin,
			APIKey:          cfg.APIKey,
			AuthExpires:     cfg.AuthExpires,
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
			Destination:     cfg.Destination,
		},
	}
}

func (c *fakeElementalConductorClient) GetPreset(presetID string) (*elementalconductor.Preset, error) {
	container := elementalconductor.MPEG4
	if strings.Contains(presetID, "hls") {
		container = elementalconductor.AppleHTTPLiveStreaming
	}
	return &elementalconductor.Preset{
		Name:      presetID,
		Container: string(container),
	}, nil
}

func (c *fakeElementalConductorClient) CreatePreset(preset *elementalconductor.Preset) (*elementalconductor.Preset, error) {
	return &elementalconductor.Preset{
		Name: preset.Name,
	}, nil
}

func (c *fakeElementalConductorClient) GetJob(jobID string) (*elementalconductor.Job, error) {
	job := c.jobs[jobID]
	return &job, nil
}

func (c *fakeElementalConductorClient) CancelJob(jobID string) (*elementalconductor.Job, error) {
	c.canceledJobs = append(c.canceledJobs, jobID)
	return &elementalconductor.Job{}, nil
}

func fakeElementalConductorFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.ElementalConductor.Host == "" || cfg.ElementalConductor.UserLogin == "" ||
		cfg.ElementalConductor.APIKey == "" || cfg.ElementalConductor.AuthExpires == 0 {
		return nil, errElementalConductorInvalidConfig
	}
	client := newFakeElementalConductorClient(cfg.ElementalConductor)
	return &elementalConductorProvider{client: client, config: cfg.ElementalConductor}, nil
}
