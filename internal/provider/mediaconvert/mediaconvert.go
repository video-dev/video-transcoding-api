package mediaconvert

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
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
	CreateJobRequest(*mediaconvert.CreateJobInput) mediaconvert.CreateJobRequest
	GetJobRequest(*mediaconvert.GetJobInput) mediaconvert.GetJobRequest
	ListJobsRequest(*mediaconvert.ListJobsInput) mediaconvert.ListJobsRequest
	CancelJobRequest(*mediaconvert.CancelJobInput) mediaconvert.CancelJobRequest
	CreatePresetRequest(*mediaconvert.CreatePresetInput) mediaconvert.CreatePresetRequest
	GetPresetRequest(*mediaconvert.GetPresetInput) mediaconvert.GetPresetRequest
	DeletePresetRequest(*mediaconvert.DeletePresetInput) mediaconvert.DeletePresetRequest
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
		Settings: &mediaconvert.JobSettings{
			Inputs: []mediaconvert.Input{
				{
					FileInput: aws.String(job.SourceMedia),
					AudioSelectors: map[string]mediaconvert.AudioSelector{
						"Audio Selector 1": {DefaultSelection: mediaconvert.AudioDefaultSelectionDefault},
					},
					VideoSelector: &mediaconvert.VideoSelector{
						ColorSpace: mediaconvert.ColorSpaceFollow,
					},
				},
			},
			OutputGroups: outputGroups,
		},
	}

	resp, err := p.client.CreateJobRequest(&createJobInput).Send(context.Background())
	if err != nil {
		return nil, err
	}
	return &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: aws.StringValue(resp.Job.Id),
		Status:        provider.StatusQueued,
	}, nil
}

func (p *mcProvider) outputPresetsFrom(outputs []db.TranscodeOutput) (map[string]mediaconvert.Preset, error) {
	presetCh := make(chan *presetResult)
	presets := map[string]mediaconvert.Preset{}

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

func (p *mcProvider) outputGroupsFrom(job *db.Job, presets map[string]mediaconvert.Preset) ([]mediaconvert.OutputGroup, error) {
	outputGroups := map[mediaconvert.ContainerType][]db.TranscodeOutput{}
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

	mcOutputGroups := []mediaconvert.OutputGroup{}
	for container, outputs := range outputGroups {
		mcOutputGroup := mediaconvert.OutputGroup{}

		var mcOutputs []mediaconvert.Output
		for _, output := range outputs {
			presetID, ok := output.Preset.ProviderMapping[Name]
			if !ok {
				return nil, provider.ErrPresetMapNotFound
			}

			rawExtension := path.Ext(output.FileName)
			filename := strings.Replace(path.Base(output.FileName), rawExtension, "", 1)
			extension := strings.Replace(rawExtension, ".", "", -1)

			mcOutputs = append(mcOutputs, mediaconvert.Output{
				Preset:       aws.String(presetID),
				NameModifier: aws.String(filename),
				Extension:    aws.String(extension),
			})
		}
		mcOutputGroup.Outputs = mcOutputs

		destination := destinationPathFrom(p.cfg.Destination, job.ID)

		switch container {
		case mediaconvert.ContainerTypeM3u8:
			mcOutputGroup.OutputGroupSettings = &mediaconvert.OutputGroupSettings{
				Type: mediaconvert.OutputGroupTypeHlsGroupSettings,
				HlsGroupSettings: &mediaconvert.HlsGroupSettings{
					Destination:            aws.String(destination),
					SegmentLength:          aws.Int64(int64(job.StreamingParams.SegmentDuration)),
					MinSegmentLength:       aws.Int64(0),
					DirectoryStructure:     mediaconvert.HlsDirectoryStructureSingleDirectory,
					ManifestDurationFormat: mediaconvert.HlsManifestDurationFormatFloatingPoint,
					OutputSelection:        mediaconvert.HlsOutputSelectionManifestsAndSegments,
					SegmentControl:         mediaconvert.HlsSegmentControlSegmentedFiles,
				},
			}
		case mediaconvert.ContainerTypeMp4:
			mcOutputGroup.OutputGroupSettings = &mediaconvert.OutputGroupSettings{
				Type: mediaconvert.OutputGroupTypeFileGroupSettings,
				FileGroupSettings: &mediaconvert.FileGroupSettings{
					Destination: aws.String(destination),
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
	preset   mediaconvert.Preset
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
		Settings: &mediaconvert.PresetSettings{
			ContainerSettings: &mediaconvert.ContainerSettings{
				Container: container,
			},
			VideoDescription:  videoPreset,
			AudioDescriptions: []mediaconvert.AudioDescription{*audioPreset},
		},
	}

	resp, err := p.client.CreatePresetRequest(&presetInput).Send(context.Background())
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

func (p *mcProvider) fetchPreset(presetID string) (mediaconvert.Preset, error) {
	preset, err := p.client.GetPresetRequest(&mediaconvert.GetPresetInput{
		Name: aws.String(presetID),
	}).Send(context.Background())
	if err != nil {
		return mediaconvert.Preset{}, err
	}
	if preset == nil || preset.Preset == nil {
		return mediaconvert.Preset{}, fmt.Errorf("unexpected response from MediaConvert: %v", preset)
	}

	return *preset.Preset, err
}

func (p *mcProvider) DeletePreset(presetID string) error {
	_, err := p.client.DeletePresetRequest(&mediaconvert.DeletePresetInput{
		Name: aws.String(presetID),
	}).Send(context.Background())

	return err
}

func (p *mcProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	jobResp, err := p.client.GetJobRequest(&mediaconvert.GetJobInput{
		Id: aws.String(job.ProviderJobID),
	}).Send(context.Background())
	if err != nil {
		return &provider.JobStatus{}, errors.Wrap(err, "fetching job info with the mediaconvert API")
	}

	return p.jobStatusFrom(job.ProviderJobID, job.ID, jobResp.Job), nil
}

func (p *mcProvider) jobStatusFrom(providerJobID string, jobID string, job *mediaconvert.Job) *provider.JobStatus {
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
	} else if p := job.JobPercentComplete; p != nil {
		status.Progress = float64(*p)
	}

	var files []provider.OutputFile
	for _, groupDetails := range job.OutputGroupDetails {
		for _, outputDetails := range groupDetails.OutputDetails {
			if outputDetails.VideoDetails == nil {
				continue
			}

			file := provider.OutputFile{}

			if height := outputDetails.VideoDetails.HeightInPx; height != nil {
				file.Height = *height
			}

			if width := outputDetails.VideoDetails.WidthInPx; width != nil {
				file.Width = *width
			}

			files = append(files, file)
		}
	}
	status.Output.Files = files

	return status
}

func statusMsgFrom(job *mediaconvert.Job) string {
	if job.ErrorMessage != nil {
		return *job.ErrorMessage
	}

	return string(job.CurrentPhase)
}

func (p *mcProvider) CancelJob(id string) error {
	_, err := p.client.CancelJobRequest(&mediaconvert.CancelJobInput{
		Id: aws.String(id),
	}).Send(context.Background())

	return err
}

func (p *mcProvider) Healthcheck() error {
	_, err := p.client.ListJobsRequest(nil).Send(context.Background())
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

	mcCfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "loading default aws config")
	}

	if cfg.MediaConvert.AccessKeyID+cfg.MediaConvert.SecretAccessKey != "" {
		mcCfg.Credentials = &aws.StaticCredentialsProvider{Value: aws.Credentials{
			AccessKeyID:     cfg.MediaConvert.AccessKeyID,
			SecretAccessKey: cfg.MediaConvert.SecretAccessKey,
		}}
	}

	if cfg.MediaConvert.Region != "" {
		mcCfg.Region = cfg.MediaConvert.Region
	}

	mcCfg.EndpointResolver = &aws.ResolveWithEndpoint{
		URL: cfg.MediaConvert.Endpoint,
	}

	return &mcProvider{
		client: mediaconvert.New(mcCfg),
		cfg:    cfg.MediaConvert,
	}, nil
}
