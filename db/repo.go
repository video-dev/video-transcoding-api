package db

// ErrJobNotFound is the error returned when the job is not found on GetJob or
// DeleteJob.
type ErrJobNotFound string

func (err ErrJobNotFound) Error() string {
	return string(err)
}

// JobRepository is the interface that defines the method for managing Job
// persistence.
type JobRepository interface {
	SaveJob(*Job) error
	DeleteJob(*Job) error
	GetJob(id string) (*Job, error)
}
