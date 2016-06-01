package provider

import (
	"errors"
	"fmt"
	"sort"

	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
)

var (
	// ErrProviderAlreadyRegistered is the error returned when trying to register a
	// provider twice.
	ErrProviderAlreadyRegistered = errors.New("provider is already registered")

	// ErrProviderNotFound is the error returned when asking for a provider
	// that is not registered.
	ErrProviderNotFound = errors.New("provider not found")

	// ErrPresetMapNotFound is the error returned when the given preset is not
	// found in the provider.
	ErrPresetMapNotFound = errors.New("preset not found in provider")
)

// TranscodingProvider represents a provider of transcoding.
//
// It defines a basic API for transcoding a media and query the status of a
// Job. The underlying provider should handle the profileSpec as desired (it
// might be a JSON, or an XML, or anything else.
type TranscodingProvider interface {
	Transcode(*db.Job, TranscodeProfile) (*JobStatus, error)
	JobStatus(id string) (*JobStatus, error)
	CreatePreset(Preset) (string, error)
	DeletePreset(presetID string) error

	// Healthcheck should return nil if the provider is currently available
	// for transcoding videos, otherwise it should return an error
	// explaining what's going on.
	Healthcheck() error

	// Capabilities describes the capabilities of the provider.
	Capabilities() Capabilities
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
//
// swagger:model
type JobStatus struct {
	ProviderJobID     string                 `json:"providerJobId,omitempty"`
	Status            Status                 `json:"status,omitempty"`
	ProviderName      string                 `json:"providerName,omitempty"`
	StatusMessage     string                 `json:"statusMessage,omitempty"`
	ProviderStatus    map[string]interface{} `json:"providerStatus,omitempty"`
	OutputDestination string                 `json:"outputDestination,omitempty"`
}

// StreamingParams contains all parameters related to the streaming protocol used.
type StreamingParams struct {
	SegmentDuration uint   `json:"segmentDuration,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
}

// Preset define the set of parameters of a given preset
type Preset struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Container     string `json:"container,omitempty"`
	Width         string `json:"video>width,omitempty"`
	Height        string `json:"video>height,omitempty"`
	VideoCodec    string `json:"video>codec,omitempty"`
	VideoBitrate  string `json:"video>bitrate,omitempty"`
	GopSize       string `json:"video>gopSize,omitempty"`
	GopMode       string `json:"video>gopMode,omitempty"`
	InterlaceMode string `json:"video>interlaceMode,omitempty"`
	Profile       string `json:"profile,omitempty"`
	ProfileLevel  string `json:"profileLevel,omitempty"`
	RateControl   string `json:"rateControl,omitempty"`
	AudioCodec    string `json:"audio>codec,omitempty"`
	AudioBitrate  string `json:"audio>bitrate,omitempty"`
}

// TranscodeProfile defines the set of inputs necessary for running a transcoding job.
type TranscodeProfile struct {
	SourceMedia     string
	Presets         []db.PresetMap
	StreamingParams StreamingParams
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

	// StatusUnknown is an unexpected status for a job.
	StatusUnknown = Status("unknown")
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

// ListProviders returns the list of currently registered providers,
// alphabetically ordered.
func ListProviders(c *config.Config) []string {
	providerNames := make([]string, 0, len(providers))
	for name, factory := range providers {
		if _, err := factory(c); err == nil {
			providerNames = append(providerNames, name)
		}
	}
	sort.Strings(providerNames)
	return providerNames
}

// DescribeProvider describes the given provider. It includes information about
// the provider's capabilities and its current health state.
func DescribeProvider(name string, c *config.Config) (*Description, error) {
	factory, err := GetProviderFactory(name)
	if err != nil {
		return nil, err
	}
	description := Description{Name: name}
	provider, err := factory(c)
	if err != nil {
		return &description, nil
	}
	description.Enabled = true
	description.Capabilities = provider.Capabilities()
	description.Health = Health{OK: true}
	if err = provider.Healthcheck(); err != nil {
		description.Health = Health{OK: false, Message: err.Error()}
	}
	return &description, nil
}
