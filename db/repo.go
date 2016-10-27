package db

import (
	"errors"
	"time"
)

var (
	// ErrJobNotFound is the error returned when the job is not found on GetJob or
	// DeleteJob.
	ErrJobNotFound = errors.New("job not found")

	// ErrPresetMapNotFound is the error returned when the presetmap is not found
	// on GetPresetMap, UpdatePresetMap or DeletePresetMap.
	ErrPresetMapNotFound = errors.New("presetmap not found")

	// ErrPresetMapAlreadyExists is the error returned when the presetmap already
	// exists.
	ErrPresetMapAlreadyExists = errors.New("presetmap already exists")

	// ErrLocalPresetNotFound is the error returned when the local preset is not found
	// on GetPresetMap, UpdatePresetMap or DeletePresetMap.
	ErrLocalPresetNotFound = errors.New("local preset not found")

	// ErrLocalPresetAlreadyExists is the error returned when the local preset already
	// exists.
	ErrLocalPresetAlreadyExists = errors.New("local preset already exists")
)

// Repository represents the repository for persisting types of the API.
type Repository interface {
	JobRepository
	PresetMapRepository
	LocalPresetRepository
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

// LocalPresetRepository provides an interface that defines the set of methods for
// managing presets when the provider don't have the ability to store/manage it.
type LocalPresetRepository interface {
	CreateLocalPreset(*LocalPreset) error
	UpdateLocalPreset(*LocalPreset) error
	DeleteLocalPreset(*LocalPreset) error
	GetLocalPreset(name string) (*LocalPreset, error)
}
