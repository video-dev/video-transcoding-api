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
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// Name is the name used for registering the Elemental Conductor provider in the
// registry of providers.
const Name = "elementalconductor"

const defaultJobPriority = 50

var errElementalConductorInvalidConfig = provider.InvalidConfigError("missing Elemental user login or api key. Please define the environment variables ELEMENTALCONDUCTOR_USER_LOGIN and ELEMENTALCONDUCTOR_API_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, elementalConductorFactory)
}

type elementalConductorProvider struct {
	config *config.Config
	client clientInterface
}

func (p *elementalConductorProvider) DeletePreset(presetID string) error {
	return p.client.DeletePreset(presetID)
}

func (p *elementalConductorProvider) CreatePreset(preset provider.Preset) (string, error) {
	elementalConductorPreset := elementalconductor.Preset{
		XMLName: xml.Name{Local: "preset"},
	}
	elementalConductorPreset.Name = preset.Name
	elementalConductorPreset.Description = preset.Description
	elementalConductorPreset.Container = preset.Container
	elementalConductorPreset.Profile = preset.Profile
	elementalConductorPreset.ProfileLevel = preset.ProfileLevel
	elementalConductorPreset.RateControl = preset.RateControl
	elementalConductorPreset.Width = preset.Video.Width
	elementalConductorPreset.Height = preset.Video.Height
	elementalConductorPreset.VideoCodec = preset.Video.Codec
	elementalConductorPreset.VideoBitrate = preset.Video.Bitrate
	elementalConductorPreset.GopSize = preset.Video.GopSize
	elementalConductorPreset.GopMode = preset.Video.GopMode
	elementalConductorPreset.InterlaceMode = preset.Video.InterlaceMode
	elementalConductorPreset.AudioCodec = preset.Audio.Codec
	elementalConductorPreset.AudioBitrate = preset.Audio.Bitrate

	result, err := p.client.CreatePreset(&elementalConductorPreset)
	if err != nil {
		return "", err
	}

	return result.Name, nil
}

func (p *elementalConductorProvider) GetPreset(presetID string) (interface{}, error) {
	preset, err := p.client.GetPreset(presetID)
	if err != nil {
		return nil, err
	}
	return preset, err
}

func (p *elementalConductorProvider) Transcode(job *db.Job, transcodeProfile provider.TranscodeProfile) (*provider.JobStatus, error) {
	newJob, err := p.newJob(job, transcodeProfile)
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

func (p *elementalConductorProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	resp, err := p.client.GetJob(job.ProviderJobID)
	if err != nil {
		return nil, err
	}
	providerStatus := map[string]interface{}{
		"status":    resp.Status,
		"submitted": resp.Submitted,
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
		ProviderJobID:  job.ProviderJobID,
		Progress:       float64(resp.PercentComplete),
		Status:         p.statusMap(resp.Status),
		ProviderStatus: providerStatus,
		Output: provider.JobOutput{
			Destination: p.getOutputDestination(job),
		},
	}, nil
}

func (p *elementalConductorProvider) getOutputDestination(job *db.Job) string {
	return strings.TrimRight(p.config.ElementalConductor.Destination, "/") + "/" + job.ID
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

func (p *elementalConductorProvider) buildOutputGroupAndStreamAssemblies(outputLocation elementalconductor.Location, transcodeProfile provider.TranscodeProfile) ([]elementalconductor.OutputGroup, []elementalconductor.StreamAssembly, error) {
	var streamingOutputList []elementalconductor.Output
	var streamAssemblyList []elementalconductor.StreamAssembly
	var outputGroupList []elementalconductor.OutputGroup
	var outputGroupOrder int
	var streamingGroupOrder int
	for index, output := range transcodeProfile.Outputs {
		indexString := strconv.Itoa(index)
		streamAssemblyName := "stream_" + indexString
		out := elementalconductor.Output{StreamAssemblyName: streamAssemblyName}
		presetID, ok := output.Preset.ProviderMapping[Name]
		if !ok {
			return outputGroupList, nil, provider.ErrPresetMapNotFound
		}
		presetOutput, err := p.GetPreset(presetID)
		if err != nil {
			return outputGroupList, nil, err
		}
		presetStruct := presetOutput.(*elementalconductor.Preset)
		if presetStruct.Container == string(elementalconductor.AppleHTTPLiveStreaming) {
			streamingGroupOrder++
			out.NameModifier = fmt.Sprintf("_%010d", streamingGroupOrder)
			out.Container = elementalconductor.AppleHTTPLiveStreaming
			out.Order = streamingGroupOrder
			streamingOutputList = append(streamingOutputList, out)
		} else {
			outputGroupOrder++
			location := outputLocation
			location.URI += "/" + output.FileName[:len(output.FileName)-len(filepath.Ext(output.FileName))]
			ext := strings.TrimLeft(output.Preset.OutputOpts.Extension, ".")
			out.Container = elementalconductor.Container(ext)
			out.Order = 1
			outputGroupList = append(outputGroupList, elementalconductor.OutputGroup{
				Order:  outputGroupOrder,
				Type:   elementalconductor.FileOutputGroupType,
				Output: []elementalconductor.Output{out},
				FileGroupSettings: &elementalconductor.FileGroupSettings{
					Destination: &location,
				},
			})
		}
		streamAssembly := elementalconductor.StreamAssembly{
			Name:   streamAssemblyName,
			Preset: presetID,
		}
		streamAssemblyList = append(streamAssemblyList, streamAssembly)
	}
	if len(streamingOutputList) > 0 {
		playlistFileName := transcodeProfile.StreamingParams.PlaylistFileName
		location := outputLocation
		location.URI += "/" + strings.TrimRight(playlistFileName, filepath.Ext(playlistFileName))
		outputGroupOrder++
		streamingOutputGroup := elementalconductor.OutputGroup{
			Order: outputGroupOrder,
			AppleLiveGroupSettings: &elementalconductor.AppleLiveGroupSettings{
				Destination:     &location,
				SegmentDuration: transcodeProfile.StreamingParams.SegmentDuration,
				EmitSingleFile:  true,
			},
			Type:   elementalconductor.AppleLiveOutputGroupType,
			Output: streamingOutputList,
		}
		outputGroupList = append(outputGroupList, streamingOutputGroup)
	}
	return outputGroupList, streamAssemblyList, nil
}

// newJob constructs a job spec from the given source and presets
func (p *elementalConductorProvider) newJob(job *db.Job, transcodeProfile provider.TranscodeProfile) (*elementalconductor.Job, error) {
	inputLocation := elementalconductor.Location{
		URI:      transcodeProfile.SourceMedia,
		Username: p.client.GetAccessKeyID(),
		Password: p.client.GetSecretAccessKey(),
	}
	baseLocation := strings.TrimRight(p.config.ElementalConductor.Destination, "/")
	outputLocation := elementalconductor.Location{
		URI:      baseLocation + "/" + job.ID,
		Username: p.client.GetAccessKeyID(),
		Password: p.client.GetSecretAccessKey(),
	}
	outputGroup, streamAssemblyList, err := p.buildOutputGroupAndStreamAssemblies(outputLocation, transcodeProfile)
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

func (p *elementalConductorProvider) CancelJob(id string) error {
	_, err := p.client.CancelJob(id)
	return err
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
