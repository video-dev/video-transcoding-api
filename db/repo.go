package db

import (
	"errors"
	"time"
)

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
	CreateJob(*Job) error
	DeleteJob(*Job) error
	GetJob(id string) (*Job, error)
	ListJobs(JobFilter) ([]Job, error)
}

// JobFilter contains a set of parameters for filtering the list of jobs in
// JobRepository.
type JobFilter struct {
	// Filter jobs since the given time.
	Since time.Time

	// Limit the number of jobs in the result. 0 means no limit.
	Limit uint
}

// PresetRepository is the interface that defines the set of methods for
// managing Preset persistence.
type PresetRepository interface {
	CreatePreset(*Preset) error
	UpdatePreset(*Preset) error
	DeletePreset(*Preset) error
	GetPreset(name string) (*Preset, error)
	ListPresets() ([]Preset, error)
}
