// Package zencoder provides a implementation of the provider that uses the
// Zencoder API for transcoding media files.
//
// It doesn't expose any public type. In order to use the provider, one must
// import this package and then grab the factory from the provider package:
//
//     import (
//         "github.com/NYTimes/video-transcoding-api/provider"
//         "github.com/NYTimes/video-transcoding-api/provider/zencoder"
//     )
//
//     func UseProvider() {
//         factory, err := provider.GetProviderFactory(Zencoder.Name)
//         // handle err and use factory to get an instance of the provider.
//     }
package zencoder

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/brandscreen/zencoder"
)

// Name is the name used for registering the Encoding.com provider in the
// registry of providers.
const Name = "zencoder"

var errZencoderInvalidConfig = provider.InvalidConfigError("missing Zencoder API key. Please define the environment variables ZENCODER_API_KEY or set these values in the configuration file")

func init() {
	provider.Register(Name, zencoderFactory)
}

// Client is a interface that both
// brandscreen/zencoder and fakeZencoder implements
type Client interface {
	CreateJob(*zencoder.EncodingSettings) (*zencoder.CreateJobResponse, error)
	ListJobs() ([]*zencoder.JobDetails, error)
	CancelJob(id int64) error
	GetJobProgress(id int64) (*zencoder.JobProgress, error)
	GetJobDetails(id int64) (*zencoder.JobDetails, error)
}

type zencoderProvider struct {
	config *config.Config
	client Client
	db     db.Repository
}

func (z *zencoderProvider) Transcode(job *db.Job, transcodeProfile provider.TranscodeProfile) (*provider.JobStatus, error) {
	outputs, err := z.buildOutputs(transcodeProfile)
	if err != nil {
		return nil, err
	}
	encodingSettings := zencoder.EncodingSettings{
		Input:      transcodeProfile.SourceMedia,
		Outputs:    outputs,
		LiveStream: false,
		Region:     "US",
	}
	response, err := z.client.CreateJob(&encodingSettings)
	if err != nil {
		return nil, err
	}
	return &provider.JobStatus{
		ProviderJobID: strconv.FormatInt(response.Id, 10),
		StatusMessage: "created",
		ProviderName:  Name,
	}, nil
}

func (z *zencoderProvider) buildOutputs(transcodeProfile provider.TranscodeProfile) ([]*zencoder.OutputSettings, error) {
	zencoderOutputs := make([]*zencoder.OutputSettings, 0, len(transcodeProfile.Outputs))
	for _, output := range transcodeProfile.Outputs {
		localPresetOutput, err := z.GetPreset(output.Preset.Name)
		if err != nil {
			return nil, fmt.Errorf("Error getting localpreset: %s", err.Error())
		}
		localPresetStruct := localPresetOutput.(*db.LocalPreset)
		zencoderOutput, err := z.buildOutput(localPresetStruct.Preset, output.FileName)
		if err != nil {
			return nil, fmt.Errorf("Error building output: %s", err.Error())
		}
		zencoderOutputs = append(zencoderOutputs, &zencoderOutput)
	}
	return zencoderOutputs, nil
}

func (z *zencoderProvider) buildOutput(preset db.Preset, outputFileName string) (zencoder.OutputSettings, error) {
	zencoderOutput := zencoder.OutputSettings{
		Label:      preset.Name + ":" + preset.Description,
		Format:     preset.Container,
		VideoCodec: preset.Video.Codec,
		AudioCodec: preset.Audio.Codec,
		Filename:   outputFileName,
	}
	destinationURL, err := url.Parse(z.config.Zencoder.Destination)
	if err != nil {
		return zencoder.OutputSettings{}, fmt.Errorf("error parsing destination (%q)", z.config.Zencoder.Destination)
	}
	zencoderOutput.BaseUrl = destinationURL.String()

	width, err := strconv.ParseInt(preset.Video.Width, 10, 32)
	if err != nil {
		return zencoder.OutputSettings{}, fmt.Errorf("error converting preset width (%q): %s", preset.Video.Width, err)
	}
	zencoderOutput.Width = int32(width)

	height, err := strconv.ParseInt(preset.Video.Height, 10, 32)
	if err != nil {
		return zencoder.OutputSettings{}, fmt.Errorf("error converting preset height (%q): %s", preset.Video.Height, err)
	}
	zencoderOutput.Height = int32(height)

	videoBitrate, err := strconv.ParseInt(preset.Video.Bitrate, 10, 32)
	if err != nil {
		return zencoder.OutputSettings{}, fmt.Errorf("error converting preset video bitrate (%q): %s", preset.Video.Bitrate, err)
	}
	zencoderOutput.VideoBitrate = int32(videoBitrate) / 1000

	keyframeInterval, err := strconv.ParseInt(preset.Video.GopSize, 10, 32)
	if err != nil {
		return zencoder.OutputSettings{}, fmt.Errorf("error converting preset keyframe interval (%q): %s", preset.Video.GopSize, err)
	}
	zencoderOutput.KeyframeInterval = int32(keyframeInterval)

	audioBitrate, err := strconv.ParseInt(preset.Audio.Bitrate, 10, 32)
	if err != nil {
		return zencoder.OutputSettings{}, fmt.Errorf("error converting preset audio bitrate (%q): %s", preset.Audio.Bitrate, err)
	}
	zencoderOutput.AudioBitrate = int32(audioBitrate) / 1000

	if preset.Video.GopMode == "fixed" {
		zencoderOutput.FixedKeyframeInterval = true
	}
	if preset.Video.Codec == "h264" {
		zencoderOutput.H264Profile = preset.Profile
		zencoderOutput.H264Level = preset.ProfileLevel
	}
	if preset.RateControl == "CBR" {
		zencoderOutput.ConstantBitrate = true
	}
	zencoderOutput.Deinterlace = "on"
	return zencoderOutput, nil
}

func (z *zencoderProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	jobID, err := strconv.ParseInt(job.ProviderJobID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting job ID (%q): %s", job.ID, err)
	}
	jobOutputs, err := z.getJobOutputs(jobID)
	if err != nil {
		return nil, fmt.Errorf("error getting job outputs: %s", err)
	}
	providerStatus, err := z.getProviderStatus(jobID)
	if err != nil {
		return nil, fmt.Errorf("error getting provider status: %s", err)
	}
	progress, err := z.client.GetJobProgress(jobID)
	if err != nil {
		return nil, fmt.Errorf("error getting job progress: %s", err)
	}
	sourceInfo, err := z.getSourceInfo(jobID)
	if err != nil {
		return nil, fmt.Errorf("error getting media info: %s", err)
	}
	return &provider.JobStatus{
		ProviderName:   Name,
		ProviderJobID:  job.ProviderJobID,
		Status:         provider.Status(progress.State),
		Progress:       progress.JobProgress,
		Output:         jobOutputs,
		SourceInfo:     sourceInfo,
		ProviderStatus: providerStatus,
	}, nil
}

func (z *zencoderProvider) getProviderStatus(jobID int64) (map[string]interface{}, error) {
	jobDetails, err := z.client.GetJobDetails(jobID)
	if err != nil {
		return nil, fmt.Errorf("error getting job details: %s", err)
	}
	return map[string]interface{}{
		"source":    jobDetails.Job.InputMediaFile.Url,
		"created":   jobDetails.Job.CreatedAt,
		"finished":  jobDetails.Job.FinishedAt,
		"updated":   jobDetails.Job.UpdatedAt,
		"submitted": jobDetails.Job.SubmittedAt,
	}, nil
}

func (z *zencoderProvider) getSourceInfo(jobID int64) (provider.SourceInfo, error) {
	jobDetails, err := z.client.GetJobDetails(jobID)
	if err != nil {
		return provider.SourceInfo{}, fmt.Errorf("error getting job details: %s", err)
	}
	inputMediaFile := jobDetails.Job.InputMediaFile
	return provider.SourceInfo{
		Duration:   time.Duration(inputMediaFile.DurationInMs * 1000),
		Height:     int64(inputMediaFile.Height),
		Width:      int64(inputMediaFile.Width),
		VideoCodec: inputMediaFile.VideoCodec,
	}, nil
}

func (z *zencoderProvider) getJobOutputs(jobID int64) (provider.JobOutput, error) {
	jobDetails, err := z.client.GetJobDetails(jobID)
	if err != nil {
		return provider.JobOutput{}, fmt.Errorf("error getting job details: %s", err)
	}
	files := make([]provider.OutputFile, 0, len(jobDetails.Job.OutputMediaFiles))
	for _, mediaFile := range jobDetails.Job.OutputMediaFiles {
		file := provider.OutputFile{
			Path:       mediaFile.Url,
			Container:  mediaFile.Format,
			VideoCodec: mediaFile.VideoCodec,
			Width:      int64(mediaFile.Width),
			Height:     int64(mediaFile.Height),
		}
		files = append(files, file)
	}
	return provider.JobOutput{
		Files: files,
	}, nil
}

func (z *zencoderProvider) CancelJob(id string) error {
	jobID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("error canceling job %s: %s", id, err)
	}
	return z.client.CancelJob(jobID)
}

func (z *zencoderProvider) Healthcheck() error {
	_, err := z.client.ListJobs()
	return err
}

func (z *zencoderProvider) CreatePreset(preset db.Preset) (string, error) {
	err := z.db.CreateLocalPreset(&db.LocalPreset{
		Name:   preset.Name,
		Preset: preset,
	})
	if err != nil {
		return "", err
	}
	return preset.Name, nil
}

func (z *zencoderProvider) GetPreset(presetID string) (interface{}, error) {
	return z.db.GetLocalPreset(presetID)
}

func (z *zencoderProvider) DeletePreset(presetID string) error {
	preset, err := z.GetPreset(presetID)
	if err != nil {
		return err
	}
	return z.db.DeleteLocalPreset(preset.(*db.LocalPreset))
}

func (z *zencoderProvider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls", "webm"},
		Destinations:  []string{"akamai", "s3"},
	}
}

func zencoderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.Zencoder.APIKey == "" {
		return nil, errZencoderInvalidConfig
	}
	client := zencoder.NewZencoder(cfg.Zencoder.APIKey)
	dbRepo, err := redis.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error initializing zencoder wrapper: %s", err)
	}
	return &zencoderProvider{client: client, db: dbRepo, config: cfg}, nil
}
