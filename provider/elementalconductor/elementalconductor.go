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
	"fmt"
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

var errElementalConductorInvalidConfig = provider.InvalidConfigError("missing Elemental user login or api key. Please define the environment variables ELEMENTALCONDUCTOR_USER_LOGIN and ELEMENTALCONDUCTOR_API_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, elementalConductorFactory)
}

type elementalConductorProvider struct {
	config *config.Config
	client *elementalconductor.Client
}

func (p *elementalConductorProvider) CreatePreset(preset provider.Preset) (string, error) {
	elementalConductorPreset := elementalconductor.Preset{
		XMLName: xml.Name{Local: "preset"},
	}
	elementalConductorPreset.Name = preset.Name
	elementalConductorPreset.Description = preset.Description
	elementalConductorPreset.Container = preset.Container
	elementalConductorPreset.Width = preset.Width
	elementalConductorPreset.Height = preset.Height
	elementalConductorPreset.VideoCodec = preset.VideoCodec
	elementalConductorPreset.VideoBitrate = preset.VideoBitrate
	elementalConductorPreset.GopSize = preset.GopSize
	elementalConductorPreset.GopMode = preset.GopMode
	elementalConductorPreset.Profile = preset.Profile
	elementalConductorPreset.ProfileLevel = preset.ProfileLevel
	elementalConductorPreset.RateControl = preset.RateControl
	elementalConductorPreset.InterlaceMode = preset.InterlaceMode
	elementalConductorPreset.AudioCodec = preset.AudioCodec
	elementalConductorPreset.AudioBitrate = preset.AudioBitrate

	result, err := p.client.CreatePreset(&elementalConductorPreset)
	if err != nil {
		return "", err
	}

	return result.Name, nil
}

func (p *elementalConductorProvider) Transcode(transcodeProfile provider.TranscodeProfile) (*provider.JobStatus, error) {
	newJob, err := p.newJob(transcodeProfile)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.CreateJob(newJob)
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
	case "preprocessing", "running", "postprocessing":
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

func buildOutputGroupAndStreamAssemblies(outputLocation elementalconductor.Location, transcodeProfile provider.TranscodeProfile) (elementalconductor.OutputGroup, []elementalconductor.StreamAssembly, error) {
	var outputList []elementalconductor.Output
	var streamAssemblyList []elementalconductor.StreamAssembly
	var outputGroup elementalconductor.OutputGroup
	for index, preset := range transcodeProfile.Presets {
		indexString := strconv.Itoa(index)
		streamAssemblyName := "stream_" + indexString
		output := elementalconductor.Output{
			StreamAssemblyName: streamAssemblyName,
			NameModifier:       "_" + preset.Name,
			Order:              index,
		}

		if transcodeProfile.StreamingParams.Protocol == "hls" {
			output.Container = elementalconductor.AppleHTTPLiveStreaming
		} else {
			ext := strings.TrimLeft(preset.OutputOpts.Extension, ".")
			output.Container = elementalconductor.Container(ext)
		}

		presetID, ok := preset.ProviderMapping[Name]
		if !ok {
			return outputGroup, nil, provider.ErrPresetMapNotFound
		}
		streamAssembly := elementalconductor.StreamAssembly{
			Name:   streamAssemblyName,
			Preset: presetID,
		}
		outputList = append(outputList, output)
		streamAssemblyList = append(streamAssemblyList, streamAssembly)
	}
	if transcodeProfile.StreamingParams.Protocol == "hls" {
		outputGroup = elementalconductor.OutputGroup{
			Order: defaultOutputGroupOrder,
			AppleLiveGroupSettings: &elementalconductor.AppleLiveGroupSettings{
				Destination:     &outputLocation,
				SegmentDuration: transcodeProfile.StreamingParams.SegmentDuration,
			},
			Type:   elementalconductor.AppleLiveOutputGroupType,
			Output: outputList,
		}
	} else {
		outputGroup = elementalconductor.OutputGroup{
			Order: defaultOutputGroupOrder,
			FileGroupSettings: &elementalconductor.FileGroupSettings{
				Destination: &outputLocation,
			},
			Type:   elementalconductor.FileOutputGroupType,
			Output: outputList,
		}
	}
	return outputGroup, streamAssemblyList, nil
}

// newJob constructs a job spec from the given source and presets
func (p *elementalConductorProvider) newJob(transcodeProfile provider.TranscodeProfile) (*elementalconductor.Job, error) {
	inputLocation := elementalconductor.Location{
		URI:      transcodeProfile.SourceMedia,
		Username: p.client.AccessKeyID,
		Password: p.client.SecretAccessKey,
	}
	outputLocation := elementalconductor.Location{
		URI:      p.buildFullDestination(transcodeProfile.SourceMedia),
		Username: p.client.AccessKeyID,
		Password: p.client.SecretAccessKey,
	}
	outputGroup, streamAssemblyList, err := buildOutputGroupAndStreamAssemblies(outputLocation, transcodeProfile)
	if err != nil {
		return nil, err
	}
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
	return &newJob, nil
}

func (p *elementalConductorProvider) Healthcheck() error {
	nodes, err := p.client.GetNodes()
	if err != nil {
		return err
	}
	cloudConfig, err := p.client.GetCloudConfig()
	if err != nil {
		return err
	}
	var serverCount int
	for _, node := range nodes {
		if node.Product == elementalconductor.ProductServer && node.Status == "active" {
			serverCount++
		}
	}
	if serverCount < cloudConfig.MinNodes {
		return fmt.Errorf("there are not enough active nodes. %d nodes required to be active, but found only %d", cloudConfig.MinNodes, serverCount)
	}
	return nil
}

func (p *elementalConductorProvider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls"},
		Destinations:  []string{"akamai", "s3"},
	}
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
