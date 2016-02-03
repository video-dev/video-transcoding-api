package provider

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/aws/aws-sdk-go/service/elastictranscoder/elastictranscoderiface"
	"github.com/nytm/video-transcoding-api/config"
)

const defaultAWSRegion = "us-east-1"

var errAWSInvalidConfig = errors.New("invalid Elastic Transcoder config. Please define the configuration entries in the config file or environment variables")

type awsProvider struct {
	c      elastictranscoderiface.ElasticTranscoderAPI
	config *config.ElasticTranscoder
}

func (p *awsProvider) TranscodeWithPresets(source string, presets []string) (*JobStatus, error) {
	input := elastictranscoder.CreateJobInput{
		PipelineId: aws.String(p.config.PipelineID),
		Input:      &elastictranscoder.JobInput{Key: aws.String(source)},
	}
	input.Outputs = make([]*elastictranscoder.CreateJobOutput, len(presets))
	for i, preset := range presets {
		input.Outputs[i] = &elastictranscoder.CreateJobOutput{
			PresetId: aws.String(preset),
			Key:      p.outputKey(source, preset),
		}
	}
	resp, err := p.c.CreateJob(&input)
	if err != nil {
		return nil, err
	}
	return &JobStatus{
		ProviderJobID: *resp.Job.Id,
		Status:        StatusQueued,
	}, nil
}

func (p *awsProvider) outputKey(source, preset string) *string {
	parts := strings.Split(source, "/")
	lastIndex := len(parts) - 1
	parts = append(parts[0:lastIndex], preset, parts[lastIndex])
	return aws.String(strings.Join(parts, "/"))
}

func (p *awsProvider) JobStatus(id string) (*JobStatus, error) {
	resp, err := p.c.ReadJob(&elastictranscoder.ReadJobInput{Id: aws.String(id)})
	if err != nil {
		return nil, err
	}
	return &JobStatus{
		ProviderJobID: *resp.Job.Id,
		Status:        p.statusMap(*resp.Job.Status),
	}, nil
}

func (p *awsProvider) statusMap(awsStatus string) status {
	switch awsStatus {
	case "Submitted":
		return StatusQueued
	case "Progressing":
		return StatusStarted
	case "Complete":
		return StatusFinished
	case "Canceled":
		return StatusCanceled
	default:
		return StatusFailed
	}
}

// ElasticTranscoderProvider is the factory function for the AWS Elastic
// Transcoder provider.
func ElasticTranscoderProvider(cfg *config.Config) (TranscodingProvider, error) {
	if cfg.ElasticTranscoder.AccessKeyID == "" || cfg.ElasticTranscoder.SecretAccessKey == "" || cfg.ElasticTranscoder.PipelineID == "" {
		return nil, errAWSInvalidConfig
	}
	creds := credentials.NewStaticCredentials(cfg.ElasticTranscoder.AccessKeyID, cfg.ElasticTranscoder.SecretAccessKey, "")
	region := cfg.ElasticTranscoder.Region
	if region == "" {
		region = defaultAWSRegion
	}
	awsSession := session.New(aws.NewConfig().WithCredentials(creds).WithRegion(region))
	return &awsProvider{
		c:      elastictranscoder.New(awsSession),
		config: cfg.ElasticTranscoder,
	}, nil
}
