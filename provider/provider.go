package provider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
)

const defaultAWSRegion = "us-east-1"

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
	Transcode(transcodeProfile TranscodeProfile) (*JobStatus, error)
	JobStatus(id string) (*JobStatus, error)
	CreatePreset(preset Preset) (string, error)

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
	Width         string `json:"width,omitempty"`
	Height        string `json:"height,omitempty"`
	VideoCodec    string `json:"videoCodec,omitempty"`
	VideoBitrate  string `json:"videoBitrate,omitempty"`
	GopSize       string `json:"gopSize,omitempty"`
	GopMode       string `json:"gopMode,omitempty"`
	Profile       string `json:"profile,omitempty"`
	ProfileLevel  string `json:"profileLevel,omitempty"`
	RateControl   string `json:"rateControl,omitempty"`
	InterlaceMode string `json:"interlaceMode,omitempty"`
	AudioCodec    string `json:"audioCodec,omitempty"`
	AudioBitrate  string `json:"audioBitrate,omitempty"`
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
func ListProviders() []string {
	providerNames := make([]string, 0, len(providers))
	for name := range providers {
		providerNames = append(providerNames, name)
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

// ExtractSegmentsFromPresetPlaylist downloads an HLSv3 (.m3u8) playlist from
// the S3 location specified by playlistS3URL, parses its contents to create
// a list of related segment (.ts) filenames and returns that list.
func ExtractSegmentsFromPresetPlaylist(accessKeyID string, secretAccessKey string, playlistS3URL string) []string {
	var segments []string
	playlistContents, err := loadFileFromS3(
		accessKeyID,
		secretAccessKey,
		playlistS3URL,
	)
	if err != nil {
		fmt.Printf("Error loading Playlist to get segments for job status: %s", err.Error())
		return segments
	}
	// Validate that this is a HLSv3 playlist
	if !strings.HasPrefix(playlistContents, "#EXTM3U") {
		fmt.Printf("Invalid Playlist found to get segments for job status; It does not start with #EXTM3U: %s", playlistS3URL)
		return segments
	}
	// Remove # comments in playlist contents
	re := regexp.MustCompile("(?s)#.*?\n")
	playlistContents = re.ReplaceAllString(playlistContents, "")
	// Parse out .ts segments
	segmentsFromFile := strings.Split(playlistContents, "\n")
	segmentsFromFile = segmentsFromFile[:len(segmentsFromFile)-1]
	// Prepend segments with full S3 path to them
	playlistS3URLParts := strings.Split(playlistS3URL, "/")
	playlistS3URLParts = playlistS3URLParts[:len(playlistS3URLParts)-1]
	playListS3URLPath := strings.Join(playlistS3URLParts[:], "/")
	for _, segment := range segmentsFromFile {
		segments = append(segments, fmt.Sprintf("%s/%s", playListS3URLPath, segment))
	}
	return segments
}

func loadFileFromS3(accessKeyID string, secretAccessKey string, pathToFile string) (string, error) {
	parsedURL, err := url.Parse(pathToFile)
	if err != nil {
		return "", fmt.Errorf("Could not parse S3 URL: %s", err.Error())
	}
	scheme := parsedURL.Scheme
	bucket := parsedURL.Host
	key := strings.TrimPrefix(parsedURL.Path, "/")
	if bucket == "" || key == "" {
		return "", fmt.Errorf("Invalid S3 URL: %s", pathToFile)
	}
	if scheme == "http" || scheme == "https" {
		bucket = strings.TrimSuffix(bucket, ".s3.amazonaws.com")
	}
	creds := credentials.NewStaticCredentials(
		accessKeyID,
		secretAccessKey,
		"",
	)
	region := defaultAWSRegion
	awsSession := session.New(aws.NewConfig().WithCredentials(creds).WithRegion(region))
	service := s3.New(awsSession)
	getObjectOutput, err := service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("Could not open S3 File: %s : %s", pathToFile, err.Error())
	}
	var fileContents string
	if b, err := ioutil.ReadAll(getObjectOutput.Body); err == nil {
		fileContents = string(b)
	}
	if err != nil {
		return "", fmt.Errorf("Could not load S3 File: %s", err.Error())
	}
	return fileContents, nil
}
