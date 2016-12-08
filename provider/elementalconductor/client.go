package elementalconductor

import "github.com/NYTimes/encoding-wrapper/elementalconductor"

type clientInterface interface {
	GetPreset(presetID string) (*elementalconductor.Preset, error)
	CreatePreset(preset *elementalconductor.Preset) (*elementalconductor.Preset, error)
	DeletePreset(presetID string) error
	CreateJob(job *elementalconductor.Job) (*elementalconductor.Job, error)
	GetJob(jobID string) (*elementalconductor.Job, error)
	CancelJob(jobID string) (*elementalconductor.Job, error)
	GetNodes() ([]elementalconductor.Node, error)
	GetCloudConfig() (*elementalconductor.CloudConfig, error)
}
