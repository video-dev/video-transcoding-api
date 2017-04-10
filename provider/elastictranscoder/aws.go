// Package elastictranscoder provides a implementation of the provider that
// uses AWS Elastic Transcoder for transcoding media files.
//
// It doesn't expose any public type. In order to use the provider, one must
// import this package and then grab the factory from the provider package:
//
//     import (
//         "github.com/NYTimes/video-transcoding-api/provider"
//         "github.com/NYTimes/video-transcoding-api/provider/elastictranscoder"
//     )
//
//     func UseProvider() {
//         factory, err := provider.GetProviderFactory(elastictranscoder.Name)
//         // handle err and use factory to get an instance of the provider.
//     }
package elastictranscoder

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/aws/aws-sdk-go/service/elastictranscoder/elastictranscoderiface"
)

const (
	// Name is the name used for registering the Elastic Transcoder
	// provider in the registry of providers.
	Name = "elastictranscoder"

	defaultAWSRegion = "us-east-1"
	hlsPlayList      = "HLSv3"
)

var (
	errAWSInvalidConfig = errors.New("invalid Elastic Transcoder config. Please define the configuration entries in the config file or environment variables")
	s3Pattern           = regexp.MustCompile(`^s3://`)
)

func init() {
	provider.Register(Name, elasticTranscoderFactory)
}

type awsProvider struct {
	c      elastictranscoderiface.ElasticTranscoderAPI
	config *config.ElasticTranscoder
}

func (p *awsProvider) Transcode(job *db.Job) (*provider.JobStatus, error) {
	var adaptiveStreamingOutputs []db.TranscodeOutput
	source := p.normalizeSource(job.SourceMedia)
	params := elastictranscoder.CreateJobInput{
		PipelineId: aws.String(p.config.PipelineID),
		Input:      &elastictranscoder.JobInput{Key: aws.String(source)},
	}
	params.Outputs = make([]*elastictranscoder.CreateJobOutput, len(job.Outputs))
	for i, output := range job.Outputs {
		presetID, ok := output.Preset.ProviderMapping[Name]
		if !ok {
			return nil, provider.ErrPresetMapNotFound
		}
		presetQuery := &elastictranscoder.ReadPresetInput{
			Id: aws.String(presetID),
		}
		presetOutput, err := p.c.ReadPreset(presetQuery)
		if err != nil {
			return nil, err
		}
		if presetOutput.Preset == nil || presetOutput.Preset.Container == nil {
			return nil, fmt.Errorf("misconfigured preset: %s", presetID)
		}
		var isAdaptiveStreamingPreset bool
		if *presetOutput.Preset.Container == "ts" {
			isAdaptiveStreamingPreset = true
			adaptiveStreamingOutputs = append(adaptiveStreamingOutputs, output)
		}
		params.Outputs[i] = &elastictranscoder.CreateJobOutput{
			PresetId: aws.String(presetID),
			Key:      p.outputKey(job, output.FileName, isAdaptiveStreamingPreset),
		}
		if isAdaptiveStreamingPreset {
			params.Outputs[i].SegmentDuration = aws.String(strconv.Itoa(int(job.StreamingParams.SegmentDuration)))
		}
	}

	if len(adaptiveStreamingOutputs) > 0 {
		playlistFileName := job.StreamingParams.PlaylistFileName
		playlistFileName = strings.TrimRight(playlistFileName, filepath.Ext(playlistFileName))
		jobPlaylist := elastictranscoder.CreateJobPlaylist{
			Format: aws.String(hlsPlayList),
			Name:   aws.String(job.ID + "/" + playlistFileName),
		}

		jobPlaylist.OutputKeys = make([]*string, len(adaptiveStreamingOutputs))
		for i, output := range adaptiveStreamingOutputs {
			jobPlaylist.OutputKeys[i] = p.outputKey(job, output.FileName, true)
		}

		params.Playlists = []*elastictranscoder.CreateJobPlaylist{&jobPlaylist}
	}
	resp, err := p.c.CreateJob(&params)
	if err != nil {
		return nil, err
	}
	return &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: aws.StringValue(resp.Job.Id),
		Status:        provider.StatusQueued,
	}, nil
}

func (p *awsProvider) normalizeSource(source string) string {
	if s3Pattern.MatchString(source) {
		source = strings.Replace(source, "s3://", "", 1)
		parts := strings.SplitN(source, "/", 2)
		return parts[len(parts)-1]
	}
	return source
}

func (p *awsProvider) outputKey(job *db.Job, fileName string, adaptive bool) *string {
	if adaptive {
		fileName = strings.TrimRight(fileName, filepath.Ext(fileName))
	}
	return aws.String(job.ID + "/" + fileName)
}

func (p *awsProvider) createVideoPreset(preset db.Preset) *elastictranscoder.VideoParameters {
	videoPreset := elastictranscoder.VideoParameters{
		DisplayAspectRatio: aws.String("auto"),
		FrameRate:          aws.String("auto"),
		SizingPolicy:       aws.String("Fill"),
		PaddingPolicy:      aws.String("Pad"),
		Codec:              &preset.Video.Codec,
		KeyframesMaxDist:   &preset.Video.GopSize,
		CodecOptions: map[string]*string{
			"Profile":            aws.String(strings.ToLower(preset.Video.Profile)),
			"Level":              &preset.Video.ProfileLevel,
			"MaxReferenceFrames": aws.String("2"),
		},
	}
	if preset.Video.Width != "" {
		videoPreset.MaxWidth = &preset.Video.Width
	} else {
		videoPreset.MaxWidth = aws.String("auto")
	}
	if preset.Video.Height != "" {
		videoPreset.MaxHeight = &preset.Video.Height
	} else {
		videoPreset.MaxHeight = aws.String("auto")
	}
	normalizedVideoBitRate, _ := strconv.Atoi(preset.Video.Bitrate)
	videoBitrate := strconv.Itoa(normalizedVideoBitRate / 1000)
	videoPreset.BitRate = &videoBitrate
	switch preset.Video.Codec {
	case "h264":
		videoPreset.Codec = aws.String("H.264")
	case "vp8", "vp9":
		videoPreset.Codec = aws.String(preset.Video.Codec)
		delete(videoPreset.CodecOptions, "MaxReferenceFrames")
		delete(videoPreset.CodecOptions, "Level")
		// Recommended profile value is zero, based on:
		// http://www.webmproject.org/docs/encoder-parameters/
		videoPreset.CodecOptions["Profile"] = aws.String("0")
	}
	if preset.Video.GopMode == "fixed" {
		videoPreset.FixedGOP = aws.String("true")
	}
	return &videoPreset
}

func (p *awsProvider) createThumbsPreset(preset db.Preset) *elastictranscoder.Thumbnails {
	thumbsPreset := &elastictranscoder.Thumbnails{
		PaddingPolicy: aws.String("Pad"),
		Format:        aws.String("png"),
		Interval:      aws.String("1"),
		SizingPolicy:  aws.String("Fill"),
		MaxWidth:      aws.String("auto"),
		MaxHeight:     aws.String("auto"),
	}
	return thumbsPreset
}

func (p *awsProvider) createAudioPreset(preset db.Preset) *elastictranscoder.AudioParameters {
	audioPreset := &elastictranscoder.AudioParameters{
		Codec:      &preset.Audio.Codec,
		Channels:   aws.String("auto"),
		SampleRate: aws.String("auto"),
	}

	normalizedAudioBitRate, _ := strconv.Atoi(preset.Audio.Bitrate)
	audioBitrate := strconv.Itoa(normalizedAudioBitRate / 1000)
	audioPreset.BitRate = &audioBitrate

	switch preset.Audio.Codec {
	case "aac":
		audioPreset.Codec = aws.String("AAC")
	case "libvorbis":
		audioPreset.Codec = aws.String("vorbis")
	}

	return audioPreset
}

func (p *awsProvider) CreatePreset(preset db.Preset) (string, error) {
	presetInput := elastictranscoder.CreatePresetInput{
		Name:        &preset.Name,
		Description: &preset.Description,
	}
	if preset.Container == "m3u8" {
		presetInput.Container = aws.String("ts")
	} else {
		presetInput.Container = &preset.Container
	}
	presetInput.Video = p.createVideoPreset(preset)
	presetInput.Audio = p.createAudioPreset(preset)
	presetInput.Thumbnails = p.createThumbsPreset(preset)
	presetOutput, err := p.c.CreatePreset(&presetInput)
	if err != nil {
		return "", err
	}
	return *presetOutput.Preset.Id, nil
}

func (p *awsProvider) GetPreset(presetID string) (interface{}, error) {
	readPresetInput := &elastictranscoder.ReadPresetInput{
		Id: aws.String(presetID),
	}
	readPresetOutput, err := p.c.ReadPreset(readPresetInput)
	if err != nil {
		return nil, err
	}
	return readPresetOutput, err
}

func (p *awsProvider) DeletePreset(presetID string) error {
	presetInput := elastictranscoder.DeletePresetInput{
		Id: &presetID,
	}
	_, err := p.c.DeletePreset(&presetInput)
	return err
}

func (p *awsProvider) JobStatus(job *db.Job) (*provider.JobStatus, error) {
	id := job.ProviderJobID
	resp, err := p.c.ReadJob(&elastictranscoder.ReadJobInput{Id: aws.String(id)})
	if err != nil {
		return nil, err
	}
	totalJobs := len(resp.Job.Outputs)
	completedJobs := float64(0)
	outputs := make(map[string]interface{}, totalJobs)
	for _, output := range resp.Job.Outputs {
		outputStatus := p.statusMap(aws.StringValue(output.Status))
		switch outputStatus {
		case provider.StatusFinished, provider.StatusCanceled, provider.StatusFailed:
			completedJobs++
		}
		outputs[aws.StringValue(output.Key)] = aws.StringValue(output.StatusDetail)
	}
	outputDestination, err := p.getOutputDestination(job, resp.Job)
	if err != nil {
		outputDestination = err.Error()
	}
	outputFiles, err := p.getOutputFiles(resp.Job)
	if err != nil {
		return nil, err
	}
	var sourceInfo provider.SourceInfo
	if resp.Job.Input.DetectedProperties != nil {
		sourceInfo = provider.SourceInfo{
			Duration: time.Duration(aws.Int64Value(resp.Job.Input.DetectedProperties.DurationMillis)) * time.Millisecond,
			Height:   aws.Int64Value(resp.Job.Input.DetectedProperties.Height),
			Width:    aws.Int64Value(resp.Job.Input.DetectedProperties.Width),
		}
	}
	statusMessage := ""
	if len(resp.Job.Outputs) > 0 {
		statusMessage = aws.StringValue(resp.Job.Outputs[0].StatusDetail)
		if strings.Contains(statusMessage, ":") {
			errorMessage := strings.SplitN(statusMessage, ":", 2)[1]
			statusMessage = strings.TrimSpace(errorMessage)
		}
	}
	return &provider.JobStatus{
		ProviderJobID:  aws.StringValue(resp.Job.Id),
		Status:         p.statusMap(aws.StringValue(resp.Job.Status)),
		StatusMessage:  statusMessage,
		Progress:       completedJobs / float64(totalJobs) * 100,
		ProviderStatus: map[string]interface{}{"outputs": outputs},
		SourceInfo:     sourceInfo,
		Output: provider.JobOutput{
			Destination: outputDestination,
			Files:       outputFiles,
		},
	}, nil
}

func (p *awsProvider) getOutputDestination(job *db.Job, awsJob *elastictranscoder.Job) (string, error) {
	readPipelineOutput, err := p.c.ReadPipeline(&elastictranscoder.ReadPipelineInput{
		Id: awsJob.PipelineId,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("s3://%s/%s",
		aws.StringValue(readPipelineOutput.Pipeline.OutputBucket),
		job.ID,
	), nil
}

func (p *awsProvider) getOutputFiles(job *elastictranscoder.Job) ([]provider.OutputFile, error) {
	pipeline, err := p.c.ReadPipeline(&elastictranscoder.ReadPipelineInput{
		Id: job.PipelineId,
	})
	if err != nil {
		return nil, err
	}
	files := make([]provider.OutputFile, 0, len(job.Outputs)+len(job.Playlists))
	for _, output := range job.Outputs {
		preset, err := p.c.ReadPreset(&elastictranscoder.ReadPresetInput{
			Id: output.PresetId,
		})
		if err != nil {
			return nil, err
		}
		filePath := fmt.Sprintf("s3://%s/%s%s",
			aws.StringValue(pipeline.Pipeline.OutputBucket),
			aws.StringValue(job.OutputKeyPrefix),
			aws.StringValue(output.Key),
		)
		container := aws.StringValue(preset.Preset.Container)
		if container == "ts" {
			continue
		}
		file := provider.OutputFile{
			Path:       filePath,
			Container:  container,
			VideoCodec: aws.StringValue(preset.Preset.Video.Codec),
			Width:      aws.Int64Value(output.Width),
			Height:     aws.Int64Value(output.Height),
		}
		files = append(files, file)
	}
	for _, playlist := range job.Playlists {
		filePath := fmt.Sprintf("s3://%s/%s%s",
			aws.StringValue(pipeline.Pipeline.OutputBucket),
			aws.StringValue(job.OutputKeyPrefix),
			aws.StringValue(playlist.Name)+".m3u8",
		)
		files = append(files, provider.OutputFile{Path: filePath, Container: "m3u8"})
	}
	return files, nil
}

func (p *awsProvider) statusMap(awsStatus string) provider.Status {
	switch awsStatus {
	case "Submitted":
		return provider.StatusQueued
	case "Progressing":
		return provider.StatusStarted
	case "Complete":
		return provider.StatusFinished
	case "Canceled":
		return provider.StatusCanceled
	default:
		return provider.StatusFailed
	}
}

func (p *awsProvider) CancelJob(id string) error {
	_, err := p.c.CancelJob(&elastictranscoder.CancelJobInput{Id: aws.String(id)})
	return err
}

func (p *awsProvider) Healthcheck() error {
	_, err := p.c.ReadPipeline(&elastictranscoder.ReadPipelineInput{
		Id: aws.String(p.config.PipelineID),
	})
	return err
}

func (p *awsProvider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		InputFormats:  []string{"h264"},
		OutputFormats: []string{"mp4", "hls", "webm"},
		Destinations:  []string{"s3"},
	}
}

func elasticTranscoderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	if cfg.ElasticTranscoder.AccessKeyID == "" || cfg.ElasticTranscoder.SecretAccessKey == "" || cfg.ElasticTranscoder.PipelineID == "" {
		return nil, errAWSInvalidConfig
	}
	creds := credentials.NewStaticCredentials(cfg.ElasticTranscoder.AccessKeyID, cfg.ElasticTranscoder.SecretAccessKey, "")
	region := cfg.ElasticTranscoder.Region
	if region == "" {
		region = defaultAWSRegion
	}
	awsSession, err := session.NewSession(aws.NewConfig().WithCredentials(creds).WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &awsProvider{
		c:      elastictranscoder.New(awsSession),
		config: cfg.ElasticTranscoder,
	}, nil
}
