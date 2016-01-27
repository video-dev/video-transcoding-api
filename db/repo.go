package db

import "errors"

// ErrJobNotFound is the error returned when the job is not found on GetJob or
// DeleteJob.
var ErrJobNotFound = errors.New("job not found")

// JobRepository is the interface that defines the method for managing Job
// persistence.
type JobRepository interface {
	SaveJob(*Job) error
	DeleteJob(*Job) error
	GetJob(id string) (*Job, error)
}
