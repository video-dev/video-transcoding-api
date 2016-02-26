// Package elementalconductor provides a implementation of the provider that uses the
// Elemental Conductor API for transcoding media files.
//
// It doesn't expose any public type. In order to use the provider, one must
// import this package and then grab the factory from the provider package:
//
//     import (
//         "github.com/nytm/video-transcoding-api/provider"
//         "github.com/nytm/video-transcoding-api/provider/elementalconductor"
//     )
//
//     func UseProvider() {
//         factory, err := provider.GetProviderFactory(elementalconductor.Name)
//         // handle err and use factory to get an instance of the provider.
//     }
package elementalconductor

import (
	"encoding/xml"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/NYTimes/encoding-wrapper/elementalconductor"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

// Name is the name used for registering the Elemental Conductor provider in the
// registry of providers.
const Name = "elementalconductor"

const defaultJobPriority = 50
const defaultOutputGroupOrder = 1
const defaultContainer = elementalconductor.MPEG4

var errElementalConductorInvalidConfig = provider.InvalidConfigError("missing Elemental user login or api key. Please define the environment variables ELEMENTALCONDUCTOR_USER_LOGIN and ELEMENTALCONDUCTOR_API_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, elementalConductorFactory)
}

type elementalConductorProvider struct {
	config *config.Config
	client *elementalconductor.Client
}

func (p *elementalConductorProvider) TranscodeWithPresets(source string, presets []string, adaptiveStreaming bool) (*provider.JobStatus, error) {
	newJob := p.newJob(source, presets, adaptiveStreaming)
	resp, err := p.client.PostJob(newJob)
	if err != nil {
		return nil, err
	}
	return &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: resp.GetID(),
		Status:        provider.StatusQueued,
	}, nil
}

func (p *elementalConductorProvider) JobStatus(id string) (*provider.JobStatus, error) {
	resp, err := p.client.GetJob(id)
	if err != nil {
		return nil, err
	}
	providerStatus := map[string]interface{}{
		"status":       resp.Status,
		"pct_complete": strconv.Itoa(resp.PercentComplete),
		"submitted":    resp.Submitted,
	}
	if !resp.StartTime.IsZero() {
		providerStatus["start_time"] = resp.StartTime
	}
	if !resp.CompleteTime.IsZero() {
		providerStatus["complete_time"] = resp.CompleteTime
	}
	if !resp.ErroredTime.IsZero() {
		providerStatus["errored_time"] = resp.ErroredTime
	}
	if len(resp.ErrorMessages) > 0 {
		providerStatus["error_messages"] = resp.ErrorMessages
	}
	return &provider.JobStatus{
		ProviderName:   Name,
		ProviderJobID:  resp.GetID(),
		Status:         p.statusMap(resp.Status),
		ProviderStatus: providerStatus,
	}, nil
}

func (p *elementalConductorProvider) statusMap(elementalConductorStatus string) provider.Status {
	switch strings.ToLower(elementalConductorStatus) {
	case "pending":
		return provider.StatusQueued
	case "preprocessing":
		return provider.StatusStarted
	case "running":
		return provider.StatusStarted
	case "postprocessing":
		return provider.StatusStarted
	case "complete":
		return provider.StatusFinished
	case "cancelled":
		return provider.StatusCanceled
	case "error":
		return provider.StatusFailed
	default:
		return provider.StatusUnknown
	}
}

func (p *elementalConductorProvider) buildFullDestination(source string) string {
	sourceParts := strings.Split(source, "/")
	sourceFilenamePart := sourceParts[len(sourceParts)-1]
	sourceFileName := strings.TrimSuffix(sourceFilenamePart, filepath.Ext(sourceFilenamePart))
	destination := strings.TrimRight(p.client.Destination, "/")
	return destination + "/" + sourceFileName
}

func buildOutputGroupAndStreamAssemblies(outputLocation elementalconductor.Location, adaptiveStreaming bool, presets []string) (elementalconductor.OutputGroup, []elementalconductor.StreamAssembly) {
	var outputList []elementalconductor.Output
	var streamAssemblyList []elementalconductor.StreamAssembly
	for index, preset := range presets {
		indexString := strconv.Itoa(index)
		streamAssemblyName := "stream_" + indexString
		output := elementalconductor.Output{
			StreamAssemblyName: streamAssemblyName,
			NameModifier:       "_" + preset,
			Order:              index,
		}
		if adaptiveStreaming {
			output.Container = elementalconductor.AppleHTTPLiveStreaming
		} else {
			output.Container = defaultContainer
		}
		streamAssembly := elementalconductor.StreamAssembly{
			Name:   streamAssemblyName,
			Preset: preset,
		}
		outputList = append(outputList, output)
		streamAssemblyList = append(streamAssemblyList, streamAssembly)
	}
	var outputGroup elementalconductor.OutputGroup
	if adaptiveStreaming {
		outputGroup = elementalconductor.OutputGroup{
			Order: defaultOutputGroupOrder,
			AppleLiveGroupSettings: elementalconductor.AppleLiveGroupSettings{
				Destination: outputLocation,
			},
			Type:   elementalconductor.AppleLiveOutputGroupType,
			Output: outputList,
		}
	} else {
		outputGroup = elementalconductor.OutputGroup{
			Order: defaultOutputGroupOrder,
			FileGroupSettings: elementalconductor.FileGroupSettings{
				Destination: outputLocation,
			},
			Type:   elementalconductor.FileOutputGroupType,
			Output: outputList,
		}
	}
	return outputGroup, streamAssemblyList
}

// newJob constructs a job spec from the given source and presets
func (p *elementalConductorProvider) newJob(source string, presets []string, adaptiveStreaming bool) *elementalconductor.Job {
	inputLocation := elementalconductor.Location{
		URI:      source,
		Username: p.client.AccessKeyID,
		Password: p.client.SecretAccessKey,
	}
	outputLocation := elementalconductor.Location{
		URI:      p.buildFullDestination(source),
		Username: p.client.AccessKeyID,
		Password: p.client.SecretAccessKey,
	}
	outputGroup, streamAssemblyList := buildOutputGroupAndStreamAssemblies(outputLocation, adaptiveStreaming, presets)
	newJob := elementalconductor.Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Input: elementalconductor.Input{
			FileInput: inputLocation,
		},
		Priority:       defaultJobPriority,
		OutputGroup:    outputGroup,
		StreamAssembly: streamAssemblyList,
	}
	return &newJob
}

func (p *elementalConductorProvider) Healthcheck() error {
	return nil
}

func elementalConductorFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.ElementalConductor.Host == "" || cfg.ElementalConductor.UserLogin == "" ||
		cfg.ElementalConductor.APIKey == "" || cfg.ElementalConductor.AuthExpires == 0 {
		return nil, errElementalConductorInvalidConfig
	}
	client := elementalconductor.NewClient(
		cfg.ElementalConductor.Host,
		cfg.ElementalConductor.UserLogin,
		cfg.ElementalConductor.APIKey,
		cfg.ElementalConductor.AuthExpires,
		cfg.ElementalConductor.AccessKeyID,
		cfg.ElementalConductor.SecretAccessKey,
		cfg.ElementalConductor.Destination,
	)
	return &elementalConductorProvider{client: client, config: cfg}, nil
}
