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
	"path"
	"strconv"
	"strings"
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
	outputs, err := z.buildOutputs(job, transcodeProfile)
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

func (z *zencoderProvider) buildOutputs(job *db.Job, transcodeProfile provider.TranscodeProfile) ([]*zencoder.OutputSettings, error) {
	zencoderOutputs := make([]*zencoder.OutputSettings, 0, len(transcodeProfile.Outputs))
	for _, output := range transcodeProfile.Outputs {
		localPresetOutput, err := z.GetPreset(output.Preset.Name)
		if err != nil {
			return nil, fmt.Errorf("Error getting localpreset: %s", err.Error())
		}
		localPresetStruct := localPresetOutput.(*db.LocalPreset)
		zencoderOutput, err := z.buildOutput(job, localPresetStruct.Preset, output.FileName)
		if err != nil {
			return nil, fmt.Errorf("Error building output: %s", err.Error())
		}
		zencoderOutputs = append(zencoderOutputs, &zencoderOutput)
	}
	return zencoderOutputs, nil
}

func (z *zencoderProvider) getResolution(preset db.Preset) (int32, int32) {
	var width, height int64
	width, err := strconv.ParseInt(preset.Video.Width, 10, 32)
	if err != nil || preset.Video.Width == "0" || preset.Video.Width == "" {
		width = 0
	}
	height, err = strconv.ParseInt(preset.Video.Height, 10, 32)
	if err != nil || preset.Video.Height == "0" || preset.Video.Height == "" {
		height = 0
	}
	return int32(width), int32(height)
}

func (z *zencoderProvider) buildOutput(job *db.Job, preset db.Preset, outputFileName string) (zencoder.OutputSettings, error) {
	zencoderOutput := zencoder.OutputSettings{
		Label:      preset.Name + ":" + preset.Description,
		VideoCodec: preset.Video.Codec,
		AudioCodec: preset.Audio.Codec,
		Filename:   outputFileName,
	}
	destinationURL, err := url.Parse(z.config.Zencoder.Destination)
	if err != nil {
		return zencoder.OutputSettings{}, fmt.Errorf("error parsing destination (%q)", z.config.Zencoder.Destination)
	}
	destinationURL.Path = path.Join(destinationURL.Path, job.ID) + "/"
	zencoderOutput.BaseUrl = destinationURL.String()
	zencoderOutput.Width, zencoderOutput.Height = z.getResolution(preset)
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
		zencoderOutput.H264Profile = strings.ToLower(preset.Profile)
		zencoderOutput.H264Level = strings.ToLower(preset.ProfileLevel)
	}
	if preset.RateControl == "CBR" {
		zencoderOutput.ConstantBitrate = true
	}
	if preset.Container == "m3u8" {
		zencoderOutput.Type = "segmented"
		zencoderOutput.Format = "ts"
		zencoderOutput.SegmentSeconds = int32(job.StreamingParams.SegmentDuration)
		zencoderOutput.PrepareForSegmenting = job.StreamingParams.Protocol
		zencoderOutput.HLSOptimizedTS = true
	} else {
		zencoderOutput.Format = preset.Container
	}

	zencoderOutput.Deinterlace = "on"
	return zencoderOutput, nil
}

func (z *zencoderProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	jobID, err := strconv.ParseInt(job.ProviderJobID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting job ID (%q): %s", job.ID, err)
	}
	jobDetails, err := z.client.GetJobDetails(jobID)
	if err != nil {
		return nil, fmt.Errorf("error getting job details: %s", err)
	}
	jobOutputs, err := z.getJobOutputs(job, jobDetails.Job.OutputMediaFiles)
	if err != nil {
		return nil, fmt.Errorf("error getting job outputs: %s", err)
	}
	progress, err := z.client.GetJobProgress(jobID)
	if err != nil {
		return nil, fmt.Errorf("error getting job progress: %s", err)
	}
	inputMediaFile := jobDetails.Job.InputMediaFile
	return &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: job.ProviderJobID,
		Status:        z.statusMap(progress.State),
		Progress:      progress.JobProgress,
		Output:        jobOutputs,
		SourceInfo: provider.SourceInfo{
			Duration:   time.Duration(inputMediaFile.DurationInMs * 1000),
			Height:     int64(inputMediaFile.Height),
			Width:      int64(inputMediaFile.Width),
			VideoCodec: inputMediaFile.VideoCodec,
		},
		ProviderStatus: map[string]interface{}{
			"sourcefile": jobDetails.Job.InputMediaFile.Url,
			"started":    jobDetails.Job.CreatedAt,
			"finished":   jobDetails.Job.FinishedAt,
			"updated":    jobDetails.Job.UpdatedAt,
			"created":    jobDetails.Job.SubmittedAt,
		},
	}, nil
}

func (z *zencoderProvider) statusMap(zencoderStatus string) provider.Status {
	switch zencoderStatus {
	case "waiting":
		return provider.StatusQueued
	case "pending":
		return provider.StatusQueued
	case "assigning":
		return provider.StatusQueued
	case "processing":
		return provider.StatusStarted
	case "finished":
		return provider.StatusFinished
	case "cancelled":
		return provider.StatusCanceled
	default:
		return provider.StatusFailed
	}
}

func (z *zencoderProvider) getJobOutputs(job *db.Job, outputMediaFiles []*zencoder.MediaFile) (provider.JobOutput, error) {
	files := make([]provider.OutputFile, 0, len(outputMediaFiles))
	for _, mediaFile := range outputMediaFiles {
		file := provider.OutputFile{
			Path:       mediaFile.Url,
			Container:  mediaFile.Format,
			VideoCodec: mediaFile.VideoCodec,
			Width:      int64(mediaFile.Width),
			Height:     int64(mediaFile.Height),
		}
		files = append(files, file)
	}
	destinationURL, err := url.Parse(z.config.Zencoder.Destination)
	if err != nil {
		return provider.JobOutput{}, fmt.Errorf("error parsing destination (%q)", z.config.Zencoder.Destination)
	}

	destinationURL.Path = path.Join(destinationURL.Path, job.ID) + "/"
	return provider.JobOutput{
		Files:       files,
		Destination: destinationURL.String(),
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
