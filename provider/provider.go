package provider

import (
	"errors"
	"fmt"

	"github.com/nytm/video-transcoding-api/config"
)

var (
	// ErrProviderAlreadyRegistered is the error returned when trying to register a
	// provider twice.
	ErrProviderAlreadyRegistered = errors.New("provider is already registered")

	// ErrProviderNotFound is the error returned when asking for a provider
	// that is not registered.
	ErrProviderNotFound = errors.New("provider not found")
)

// TranscodingProvider represents a provider of transcoding.
//
// It defines a basic API for transcoding a media and query the status of a
// Job. The underlying provider should handle the profileSpec as deisired (it
// might be a JSON, or an XML, or anything else.
type TranscodingProvider interface {
	JobStatus(id string) (*JobStatus, error)
}

// PresetTranscodingProvider is a transcoding provider that supports
// transcoding media using preset names.
type PresetTranscodingProvider interface {
	TranscodeWithPresets(sourceMedia string, presets []string) (*JobStatus, error)
}

// ProfileTranscodingProvider is a transcsoding provider that suppports
// transcoding media using provided profiles.
type ProfileTranscodingProvider interface {
	TranscodeWithProfiles(sourceMedia string, profiles []Profile) (*JobStatus, error)
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
	Status         Status                 `json:"status,omitempty"`
	ProviderName   string                 `json:"providerName,omitempty"`
	StatusMessage  string                 `json:"statusMessage,omitempty"`
	ProviderStatus map[string]interface{} `json:"providerStatus,omitempty"`
}

// Status is the status of a transcoding job.
type Status string

const (
	// StatusQueued is the status for a job that is in the queue for
	// execution.
	StatusQueued = Status("queued")

	// StatusStarted is the status for a job that is being executed.
	StatusStarted = Status("started")

	// StatusFinished is the status for a job that finished successfully.
	StatusFinished = Status("finished")

	// StatusFailed is the status for a job that has failed.
	StatusFailed = Status("failed")

	// StatusCanceled is the status for a job that has been canceled.
	StatusCanceled = Status("canceled")

	// StatusArchived is the status for a job that has been archived.
	StatusArchived = Status("archived")
)

var providers map[string]Factory

// Register register a new provider in the internal list of providers.
func Register(name string, provider Factory) error {
	if providers == nil {
		providers = make(map[string]Factory)
	}
	if _, ok := providers[name]; ok {
		return ErrProviderAlreadyRegistered
	}
	providers[name] = provider
	return nil
}

// GetProviderFactory looks up the list of registered providers and returns the
// factory function for the given provider name, if it's available.
func GetProviderFactory(name string) (Factory, error) {
	factory, ok := providers[name]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return factory, nil
}
