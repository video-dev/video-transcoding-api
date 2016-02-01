package db

// Job represents the job that is persisted in the database of the Transcoding
// API.
type Job struct {
	ID            string
	ProviderName  string
	ProviderJobID string
}
