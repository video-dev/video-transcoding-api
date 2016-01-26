package db

// JobRepository is the interface that defines the method for managing Job
// persistence.
type JobRepository interface {
	SaveJob(*Job) error
	DeleteJob(*Job) error
	GetJob(id string) (*Job, error)
}
