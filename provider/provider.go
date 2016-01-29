package provider

import (
	"fmt"

	"github.com/nytm/video-transcoding-api/config"
)

// TranscodingProvider represents a provider of transcoding.
//
// It defines a basic API for transcoding a media and query the status of a
// Job. The underlying provider should handle the profileSpec as deisired (it
// might be a JSON, or an XML, or anything else.
type TranscodingProvider interface {
	Transcode(sourceMedia string, profiles []Profile) (*JobStatus, error)
	JobStatus(id string) (*JobStatus, error)
}

// Factory is the function responsible for creating the instance of a
// provider.
type Factory func(cfg *config.Config) (TranscodingProvider, error)

// InvalidConfigError is returned if a provider could not be configured properly
type InvalidConfigError string

// JobNotFoundError is returned if a job with a given id could not be found by the provider
type JobNotFoundError struct {
	ID string
}

func (err InvalidConfigError) Error() string {
	return string(err)
}

func (err JobNotFoundError) Error() string {
	return fmt.Sprintf("could not found job with id: %s", err.ID)
}

// JobStatus is the representation of the status as the provide sees it. The
// provider is able to add customized information in the ProviderStatus field.
type JobStatus struct {
	ProviderJobID  string                 `json:"providerJobId,omitempty"`
	Status         status                 `json:"status,omitempty"`
	ProviderName   string                 `json:"providerName,omitempty"`
	StatusMessage  string                 `json:"statusMessage,omitempty"`
	ProviderStatus map[string]interface{} `json:"providerStatus,omitempty"`
}

type status string

const (
	// StatusQueued is the status for a job that is in the queue for
	// execution.
	StatusQueued = status("queued")

	// StatusStarted is the status for a job that is being executed.
	StatusStarted = status("started")

	// StatusFinished is the status for a job that finished successfully.
	StatusFinished = status("finished")

	// StatusFailed is the status for a job that has failed.
	StatusFailed = status("failed")
)
