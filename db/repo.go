package db

import (
	"errors"
	"time"
)

var (
	// ErrJobNotFound is the error returned when the job is not found on GetJob or
	// DeleteJob.
	ErrJobNotFound = errors.New("job not found")

	// ErrPresetMapNotFound is the error returned when the preset is not found
	// on GetPresetMap, UpdatePresetMap or DeletePresetMap.
	ErrPresetMapNotFound = errors.New("preset not found")

	// ErrPresetMapAlreadyExists is the error returned when the preset already
	// exists.
	ErrPresetMapAlreadyExists = errors.New("preset already exists")
)

// Repository represents the repository for persisting types of the API.
type Repository interface {
	JobRepository
	PresetMapRepository
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

// PresetMapRepository is the interface that defines the set of methods for
// managing PresetMap persistence.
type PresetMapRepository interface {
	CreatePresetMap(*PresetMap) error
	UpdatePresetMap(*PresetMap) error
	DeletePresetMap(*PresetMap) error
	GetPresetMap(name string) (*PresetMap, error)
	ListPresetMaps() ([]PresetMap, error)
}
