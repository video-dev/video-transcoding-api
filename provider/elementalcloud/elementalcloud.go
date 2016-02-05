// Package elementalcloud provides a implementation of the provider that uses the
// Elemental Cloud API for transcoding media files.
//
// It doesn't expose any public type. In order to use the provider, one must
// import this package and then grab the factory from the provider package:
//
//     import (
//         "github.com/nytm/video-transcoding-api/provider"
//         "github.com/nytm/video-transcoding-api/provider/elementalcloud"
//     )
//
//     func UseProvider() {
//         factory, err := provider.GetProviderFactory(elementalcloud.Name)
//         // handle err and use factory to get an instance of the provider.
//     }
package elementalcloud

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/NYTimes/encoding-wrapper/elementalcloud"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

// Name is the name used for registering the Elemental Cloud provider in the
// registry of providers.
const Name = "elementalcloud"

const defaultJobPriority = 50
const defaultOutputGroupOrder = 1
const defaultExtension = ".mp4"

var errElementalCloudInvalidConfig = provider.InvalidConfigError("missing Elemental user login or api key. Please define the environment variables ELEMENTALCLOUD_USER_LOGIN and ELEMENTALCLOUD_API_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, elementalCloudFactory)
}

type elementalCloudProvider struct {
	config *config.Config
	client *elementalcloud.Client
}

func (p *elementalCloudProvider) TranscodeWithPresets(source string, presets []string) (*provider.JobStatus, error) {
	newJob := p.NewJob(source, presets)
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

func (p *elementalCloudProvider) JobStatus(id string) (*provider.JobStatus, error) {
	_, err := p.client.GetJob(id)
	if err != nil {
		return nil, err
	}
	return &provider.JobStatus{ProviderName: Name}, nil
}

func (p *elementalCloudProvider) buildFullDestination(source string) string {
	sourceParts := strings.Split(source, "/")
	sourceFilenamePart := sourceParts[len(sourceParts)-1]
	sourceFileName := strings.TrimSuffix(sourceFilenamePart, filepath.Ext(sourceFilenamePart))
	destination := strings.TrimRight(p.client.Destination, "/")
	return destination + "/" + sourceFileName
}

func buildOutputsAndStreamAssemblies(presets []string) ([]elementalcloud.Output, []elementalcloud.StreamAssembly) {
	var outputList []elementalcloud.Output
	var streamAssemblyList []elementalcloud.StreamAssembly
	for index, preset := range presets {
		indexString := strconv.Itoa(index)
		streamAssemblyName := "stream_" + indexString
		output := elementalcloud.Output{
			StreamAssemblyName: streamAssemblyName,
			Order:              index,
			Extension:          defaultExtension,
		}
		streamAssembly := elementalcloud.StreamAssembly{
			Name:   streamAssemblyName,
			Preset: preset,
		}
		outputList = append(outputList, output)
		streamAssemblyList = append(streamAssemblyList, streamAssembly)
	}
	return outputList, streamAssemblyList
}

// NewJob constructs a job spec from the given source and presets
func (p *elementalCloudProvider) NewJob(source string, presets []string) *elementalcloud.Job {
	inputLocation := elementalcloud.Location{
		URI:      source,
		Username: p.client.AccessKeyID,
		Password: p.client.SecretAccessKey,
	}
	outputLocation := elementalcloud.Location{
		URI:      p.buildFullDestination(source),
		Username: p.client.AccessKeyID,
		Password: p.client.SecretAccessKey,
	}
	outputList, streamAssemblyList :=
		buildOutputsAndStreamAssemblies(presets)
	newJob := elementalcloud.Job{
		Input: elementalcloud.Input{
			FileInput: inputLocation,
		},
		Priority: defaultJobPriority,
		OutputGroup: elementalcloud.OutputGroup{
			Order: defaultOutputGroupOrder,
			FileGroupSettings: elementalcloud.FileGroupSettings{
				Destination: outputLocation,
			},
			Type:   "file_group_settings",
			Output: outputList,
		},
		StreamAssembly: streamAssemblyList,
	}
	return &newJob
}

func elementalCloudFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.ElementalCloud.Host == "" || cfg.ElementalCloud.UserLogin == "" ||
		cfg.ElementalCloud.APIKey == "" || cfg.ElementalCloud.AuthExpires == 0 {
		return nil, errElementalCloudInvalidConfig
	}
	client := elementalcloud.NewClient(
		cfg.ElementalCloud.Host,
		cfg.ElementalCloud.UserLogin,
		cfg.ElementalCloud.APIKey,
		cfg.ElementalCloud.AuthExpires,
		cfg.ElementalCloud.AccessKeyID,
		cfg.ElementalCloud.SecretAccessKey,
		cfg.ElementalCloud.Destination,
	)
	return &elementalCloudProvider{client: client, config: cfg}, nil
}
