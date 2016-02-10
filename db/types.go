package db

// Job represents the job that is persisted in the repository of the Transcoding
// API.
type Job struct {
	ID            string `redis-hash:"-" json:"jobId"`
	ProviderName  string `redis-hash:"providerName" json:"providerName"`
	ProviderJobID string `redis-hash:"providerJobID" json:"providerJobId"`
}

// Preset represents the preset that is persisted in the repository of the
// Transcoding API.
//
// Each preset is just an aggregator of provider presets, where each preset in
// the API maps to a preset on each provider.
type Preset struct {
	ID              string            `redis-hash:"-" json:"presetId"`
	ProviderMapping map[string]string `redis-hash:",expand" json:"providerMapping"`
}
