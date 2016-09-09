package elementalconductor

import (
	"strings"

	"github.com/NYTimes/encoding-wrapper/elementalconductor"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

type fakeElementalConductorClient struct {
	*elementalconductor.Client
	jobs         map[string]elementalconductor.Job
	canceledJobs []string
}

func newFakeElementalConductorClient(cfg *config.Config) *fakeElementalConductorClient {
	return &fakeElementalConductorClient{
		jobs: make(map[string]elementalconductor.Job),
		Client: &elementalconductor.Client{
			Host:            cfg.ElementalConductor.Host,
			UserLogin:       cfg.ElementalConductor.UserLogin,
			APIKey:          cfg.ElementalConductor.APIKey,
			AuthExpires:     cfg.ElementalConductor.AuthExpires,
			AccessKeyID:     cfg.ElementalConductor.AccessKeyID,
			SecretAccessKey: cfg.ElementalConductor.SecretAccessKey,
			Destination:     cfg.ElementalConductor.Destination,
		},
	}
}

func (c *fakeElementalConductorClient) GetAccessKeyID() string {
	return c.AccessKeyID
}

func (c *fakeElementalConductorClient) GetSecretAccessKey() string {
	return c.SecretAccessKey
}

func (c *fakeElementalConductorClient) GetDestination() string {
	return c.Destination
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
	client := newFakeElementalConductorClient(cfg)
	return &elementalConductorProvider{client: client, config: cfg}, nil
}
