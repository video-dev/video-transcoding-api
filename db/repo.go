package db

import "errors"

var (
	// ErrJobNotFound is the error returned when the job is not found on GetJob or
	// DeleteJob.
	ErrJobNotFound = errors.New("job not found")

	// ErrPresetNotFound is the error returned when the preset is not found
	// on GetPreset, UpdatePreset or DeletePreset.
	ErrPresetNotFound = errors.New("preset not found")

	// ErrPresetAlreadyExists is the error returned when the preset already
	// exists.
	ErrPresetAlreadyExists = errors.New("preset already exists")
)

// Repository represents the repository for persisting types of the API.
type Repository interface {
	JobRepository
	PresetRepository
}

// JobRepository is the interface that defines the set of methods for managing Job
// persistence.
type JobRepository interface {
	SaveJob(*Job) error
	DeleteJob(*Job) error
	GetJob(id string) (*Job, error)
}

// PresetRepository is the interface that defines the set of methods for
// managing Preset persistence.
type PresetRepository interface {
	SavePreset(*Preset) error
	UpdatePreset(*Preset) error
	DeletePreset(*Preset) error
	GetPreset(id string) (*Preset, error)
	ListPresets() ([]Preset, error)
}
