package mediaconvert

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
	"github.com/pkg/errors"
	"github.com/video-dev/video-transcoding-api/v2/config"
	"github.com/video-dev/video-transcoding-api/v2/db"
	"github.com/video-dev/video-transcoding-api/v2/internal/provider"
)

const (
	// Name identifies the MediaConvert provider by name
	Name = "mediaconvert"

	defaultAudioSampleRate = 48000
)

func init() {
	provider.Register(Name, mediaconvertFactory)
}

type mediaconvertClient interface {
	CreateJob(context.Context, *mediaconvert.CreateJobInput, ...func(*mediaconvert.Options)) (*mediaconvert.CreateJobOutput, error)
	GetJob(context.Context, *mediaconvert.GetJobInput, ...func(*mediaconvert.Options)) (*mediaconvert.GetJobOutput, error)
	ListJobs(context.Context, *mediaconvert.ListJobsInput, ...func(*mediaconvert.Options)) (*mediaconvert.ListJobsOutput, error)
	CancelJob(context.Context, *mediaconvert.CancelJobInput, ...func(*mediaconvert.Options)) (*mediaconvert.CancelJobOutput, error)
	CreatePreset(context.Context, *mediaconvert.CreatePresetInput, ...func(*mediaconvert.Options)) (*mediaconvert.CreatePresetOutput, error)
	GetPreset(context.Context, *mediaconvert.GetPresetInput, ...func(*mediaconvert.Options)) (*mediaconvert.GetPresetOutput, error)
	DeletePreset(context.Context, *mediaconvert.DeletePresetInput, ...func(*mediaconvert.Options)) (*mediaconvert.DeletePresetOutput, error)
}

type mcProvider struct {
	client mediaconvertClient
	cfg    *config.MediaConvert
}

func (p *mcProvider) Transcode(job *db.Job) (*provider.JobStatus, error) {
	presets, err := p.outputPresetsFrom(job.Outputs)
	if err != nil {
		return nil, errors.Wrap(err, "building map of output presetID to MediaConvert preset")
	}

	outputGroups, err := p.outputGroupsFrom(job, presets)
	if err != nil {
		return nil, errors.Wrap(err, "generating Mediaconvert output groups")
	}
	createJobInput := mediaconvert.CreateJobInput{
		Queue: aws.String(p.cfg.Queue),
		Role:  aws.String(p.cfg.Role),
		Priority: int32(job.QueuePriority),
		Settings: &types.JobSettings{
			Inputs: []types.Input{
				{
					FileInput: aws.String(job.SourceMedia),
					AudioSelectors: map[string]types.AudioSelector{
						"Audio Selector 1": {DefaultSelection: types.AudioDefaultSelectionDefault},
					},
					VideoSelector: &types.VideoSelector{
						ColorSpace: types.ColorSpaceFollow,
					},
				},
			},
			OutputGroups: outputGroups,
		},
	}

	resp, err := p.client.CreateJob(context.Background(), &createJobInput)
	if err != nil {
		return nil, err
	}
	return &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: *resp.Job.Id,
		Status:        provider.StatusQueued,
	}, nil
}

func (p *mcProvider) outputPresetsFrom(outputs []db.TranscodeOutput) (map[string]types.Preset, error) {
	presetCh := make(chan *presetResult)
	presets := map[string]types.Preset{}

	var wg sync.WaitGroup
	for _, output := range outputs {
		presetID, ok := output.Preset.ProviderMapping[Name]
		if !ok {
			return nil, provider.ErrPresetMapNotFound
		}

		wg.Add(1)
		go p.makeGetPresetRequest(presetID, presetCh, &wg)
	}

	go func() {
		wg.Wait()
		close(presetCh)
	}()

	for res := range presetCh {
		if res.err != nil {
			return nil, fmt.Errorf("error getting preset info: %s", res.err)
		}

		presets[res.presetID] = res.preset
	}

	return presets, nil
}

func (p *mcProvider) outputGroupsFrom(job *db.Job, presets map[string]types.Preset) ([]types.OutputGroup, error) {
	outputGroups := map[types.ContainerType][]db.TranscodeOutput{}
	for _, output := range job.Outputs {
		presetID, ok := output.Preset.ProviderMapping[Name]
		if !ok {
			return nil, provider.ErrPresetMapNotFound
		}

		preset, ok := presets[presetID]
		if !ok {
			return nil, errors.New("mediaconvert preset not found in preset results")
		}

		container := preset.Settings.ContainerSettings.Container
		outputGroups[container] = append(outputGroups[container], output)
	}

	mcOutputGroups := []types.OutputGroup{}
	for container, outputs := range outputGroups {
		mcOutputGroup := types.OutputGroup{}

		var mcOutputs []types.Output
		for _, output := range outputs {
			presetID, ok := output.Preset.ProviderMapping[Name]
			if !ok {
				return nil, provider.ErrPresetMapNotFound
			}

			rawExtension := path.Ext(output.FileName)
			filename := strings.Replace(path.Base(output.FileName), rawExtension, "", 1)
			extension := strings.Replace(rawExtension, ".", "", -1)

			mcOutputs = append(mcOutputs, types.Output{
				Preset:       aws.String(presetID),
				NameModifier: aws.String("_" + filename),
				Extension:    aws.String(extension),
			})
		}
		mcOutputGroup.Outputs = mcOutputs

		destination := destinationPathFrom(p.cfg.Destination, job.ID)

		switch container {
		case types.ContainerTypeM3u8:
			mcOutputGroup.OutputGroupSettings = &types.OutputGroupSettings{
				Type: types.OutputGroupTypeHlsGroupSettings,
				HlsGroupSettings: &types.HlsGroupSettings{
					Destination:            aws.String(destination + "/hls/index"),
					SegmentLength:          int32(job.StreamingParams.SegmentDuration),
					MinSegmentLength:       1,
					DirectoryStructure:     types.HlsDirectoryStructureSingleDirectory,
					ManifestDurationFormat: types.HlsManifestDurationFormatFloatingPoint,
					OutputSelection:        types.HlsOutputSelectionManifestsAndSegments,
					SegmentControl:         types.HlsSegmentControlSegmentedFiles,
				},
			}
		case types.ContainerTypeMp4:
			mcOutputGroup.OutputGroupSettings = &types.OutputGroupSettings{
				Type: types.OutputGroupTypeFileGroupSettings,
				FileGroupSettings: &types.FileGroupSettings{
					Destination: aws.String(destination),
				},
			}
		case types.ContainerTypeRaw:
			mcOutputGroup.OutputGroupSettings = &types.OutputGroupSettings{
				Type: types.OutputGroupTypeFileGroupSettings,
				FileGroupSettings: &types.FileGroupSettings{
					Destination: aws.String(destination + "/thumb/thumbnail"),
				},
			}
		default:
			return nil, fmt.Errorf("container %s is not yet supported with mediaconvert", string(container))
		}

		mcOutputGroups = append(mcOutputGroups, mcOutputGroup)
	}

	return mcOutputGroups, nil
}

func destinationPathFrom(destBase string, jobID string) string {
	return fmt.Sprintf("%s/%s/", strings.TrimRight(destBase, "/"), jobID)
}

type presetResult struct {
	presetID string
	preset   types.Preset
	err      error
}

func (p *mcProvider) makeGetPresetRequest(presetID string, ch chan *presetResult, wg *sync.WaitGroup) {
	defer wg.Done()
	result := &presetResult{presetID: presetID}

	presetResp, err := p.fetchPreset(presetID)
	if err != nil {
		result.err = err
		ch <- result
		return
	}

	result.preset = presetResp
	ch <- result
}

func (p *mcProvider) CreatePreset(preset db.Preset) (string, error) {

	if preset.Video != (db.VideoPreset{}) {
		// call video function
		return p.CreateVideoPreset(preset)
	} else if preset.Thumbnail != (db.ThumbnailPreset{}) {
		return p.CreateThumbnailPreset(preset)
	}

	return "", fmt.Errorf("missing video description settings")
}

func (p *mcProvider) CreateVideoPreset(preset db.Preset) (string, error) {
	
	container, err := containerFrom(preset.Container)
	if err != nil {
		return "", errors.Wrap(err, "mapping preset container to MediaConvert container")
	}

	videoPreset, err := videoPresetFrom(preset)
	if err != nil {
		return "", errors.Wrap(err, "generating video preset")
	}
	
	audioPreset, err := audioPresetFrom(preset)
	if err != nil {
		return "", errors.Wrap(err, "generating audio preset")
	}

	presetInput := mediaconvert.CreatePresetInput{
		Name:        &preset.Name,
		Description: &preset.Description,
		Settings: &types.PresetSettings{
			ContainerSettings: &types.ContainerSettings{
				Container: container,
			},
			VideoDescription:  videoPreset,
			AudioDescriptions: []types.AudioDescription{*audioPreset},
		},
	}

	resp, err := p.client.CreatePreset(context.Background(), &presetInput)
	if err != nil {
		return "", err
	}
	if resp == nil || resp.Preset == nil || resp.Preset.Name == nil {
		return "", fmt.Errorf("unexpected response from MediaConvert: %v", resp)
	}

	return *resp.Preset.Name, nil
	

}

func (p *mcProvider) CreateThumbnailPreset(preset db.Preset) (string, error) {

	container, err := containerFrom(preset.Container)
	if err != nil {
		return "", errors.Wrap(err, "mapping preset container to MediaConvert container")
	}

	thumbnailPreset, err := thumbnailPresetFrom(preset)
	if err != nil {
		return "", errors.Wrap(err, "generating thumbnail preset")
	}

	presetInput := mediaconvert.CreatePresetInput{
						Name:        &preset.Name,
						Description: &preset.Description,
						Settings: &types.PresetSettings{
							ContainerSettings: &types.ContainerSettings{
								Container: container,
							},
							VideoDescription:  thumbnailPreset,
						},
					}

	resp, err := p.client.CreatePreset(context.Background(), &presetInput)
	if err != nil {
		return "", err
	}
	if resp == nil || resp.Preset == nil || resp.Preset.Name == nil {
		return "", fmt.Errorf("unexpected response from MediaConvert: %v", resp)
	}

	return *resp.Preset.Name, nil
}


func (p *mcProvider) GetPreset(presetID string) (interface{}, error) {
	preset, err := p.fetchPreset(presetID)
	if err != nil {
		return nil, err
	}

	return preset, err
}

func (p *mcProvider) fetchPreset(presetID string) (types.Preset, error) {
	preset, err := p.client.GetPreset(context.Background(), &mediaconvert.GetPresetInput{
		Name: aws.String(presetID),
	})
	if err != nil {
		return types.Preset{}, err
	}
	if preset == nil || preset.Preset == nil {
		return types.Preset{}, fmt.Errorf("unexpected response from MediaConvert: %v", preset)
	}

	return *preset.Preset, err
}

func (p *mcProvider) DeletePreset(presetID string) error {
	_, err := p.client.DeletePreset(context.Background(), &mediaconvert.DeletePresetInput{
		Name: aws.String(presetID),
	})

	return err
}

func (p *mcProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	jobResp, err := p.client.GetJob(context.Background(), &mediaconvert.GetJobInput{
		Id: aws.String(job.ProviderJobID),
	})

	if err != nil {
		return &provider.JobStatus{}, errors.Wrap(err, "fetching job info with the mediaconvert API")
	}

	return p.jobStatusFrom(job.ProviderJobID, job.ID, jobResp.Job), nil
}

func (p *mcProvider) jobStatusFrom(providerJobID string, jobID string, job *types.Job) *provider.JobStatus {
	status := &provider.JobStatus{
		ProviderJobID: providerJobID,
		ProviderName:  Name,
		Status:        providerStatusFrom(job.Status),
		StatusMessage: statusMsgFrom(job),
		Output: provider.JobOutput{
			Destination: destinationPathFrom(p.cfg.Destination, jobID),
		},
	}

	if status.Status == provider.StatusFinished {
		status.Progress = 100
	} else {
		status.Progress = float64(job.JobPercentComplete)
	}

	var files []provider.OutputFile
	for _, groupDetails := range job.OutputGroupDetails {
		for _, outputDetails := range groupDetails.OutputDetails {
			if outputDetails.VideoDetails == nil {
				continue
			}

			files = append(files, provider.OutputFile{
				Height: 	int64(outputDetails.VideoDetails.HeightInPx),
				Width:  	int64(outputDetails.VideoDetails.WidthInPx),
				Path:   	jobID + "/hls/index.m3u8",
				Thumbnail:  jobID + "/thumb/thumbnail_",
			})
		}
	}
	status.Output.Files = files

	return status
}

func statusMsgFrom(job *types.Job) string {
	if job.ErrorMessage != nil {
		return *job.ErrorMessage
	}

	return string(job.CurrentPhase)
}

func (p *mcProvider) CancelJob(id string) error {
	_, err := p.client.CancelJob(context.Background(), &mediaconvert.CancelJobInput{
		Id: aws.String(id),
	})

	return err
}

func (p *mcProvider) Healthcheck() error {
	_, err := p.client.ListJobs(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "listing jobs")
	}
	return nil
}

func (p *mcProvider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		InputFormats:  []string{"h264"},
		OutputFormats: []string{"mp4", "hls"},
		Destinations:  []string{"s3"},
	}
}

func mediaconvertFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.MediaConvert.Endpoint == "" || cfg.MediaConvert.Queue == "" || cfg.MediaConvert.Role == "" {
		return nil, errors.New("incomplete MediaConvert config")
	}

	mcCfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "loading default aws config")
	}

	if cfg.MediaConvert.AccessKeyID+cfg.MediaConvert.SecretAccessKey != "" {
		mcCfg.Credentials = &credentials.StaticCredentialsProvider{Value: aws.Credentials{
			AccessKeyID:     cfg.MediaConvert.AccessKeyID,
			SecretAccessKey: cfg.MediaConvert.SecretAccessKey,
		}}
	}

	if cfg.MediaConvert.Region != "" {
		mcCfg.Region = cfg.MediaConvert.Region
	}

	return &mcProvider{
		client: mediaconvert.New(mediaconvert.Options{
			EndpointResolver: mediaconvert.EndpointResolverFromURL(cfg.MediaConvert.Endpoint),
			Region:           cfg.MediaConvert.Region,
			Credentials: mcCfg.Credentials,
		}),
		cfg: cfg.MediaConvert,
	}, nil
}
