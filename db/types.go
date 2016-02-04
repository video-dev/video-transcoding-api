package db

// Job represents the job that is persisted in the repository of the Transcoding
// API.
type Job struct {
	ID            string `redis-hash:"-"`
	ProviderName  string `redis-hash:"providerName"`
	ProviderJobID string `redis-hash:"providerJobID"`
}

// Preset represents the preset that is persisted in the repository of the
// Transcoding API.
type Preset struct {
	ID              string            `redis-hash:"-"`
	ProviderMapping map[string]string `redis-hash:",expand"`
}
