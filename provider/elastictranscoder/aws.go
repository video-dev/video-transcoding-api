// Package elastictranscoder provides a implementation of the provider that
// uses AWS Elastic Transcoder for transcoding media files.
//
// It doesn't expose any public type. In order to use the provider, one must
// import this package and then grab the factory from the provider package:
//
//     import (
//         "github.com/nytm/video-transcoding-api/provider"
//         "github.com/nytm/video-transcoding-api/provider/elastictranscoder"
//     )
//
//     func UseProvider() {
//         factory, err := provider.GetProviderFactory(elastictranscoder.Name)
//         // handle err and use factory to get an instance of the provider.
//     }
package elastictranscoder

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/aws/aws-sdk-go/service/elastictranscoder/elastictranscoderiface"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

const (
	// Name is the name used for registering the Elastic Transcoder
	// provider in the registry of providers.
	Name = "elastictranscoder"

	defaultAWSRegion = "us-east-1"
)

var errAWSInvalidConfig = errors.New("invalid Elastic Transcoder config. Please define the configuration entries in the config file or environment variables")

func init() {
	provider.Register(Name, elasticTranscoderProvider)
}

type awsProvider struct {
	c      elastictranscoderiface.ElasticTranscoderAPI
	config *config.ElasticTranscoder
}

func (p *awsProvider) TranscodeWithPresets(source string, presets []string) (*provider.JobStatus, error) {
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
	return &provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: *resp.Job.Id,
		Status:        provider.StatusQueued,
	}, nil
}

func (p *awsProvider) outputKey(source, preset string) *string {
	parts := strings.Split(source, "/")
	lastIndex := len(parts) - 1
	parts = append(parts[0:lastIndex], preset, parts[lastIndex])
	return aws.String(strings.Join(parts, "/"))
}

func (p *awsProvider) JobStatus(id string) (*provider.JobStatus, error) {
	resp, err := p.c.ReadJob(&elastictranscoder.ReadJobInput{Id: aws.String(id)})
	if err != nil {
		return nil, err
	}
	outputs := make(map[string]interface{})
	for _, output := range resp.Job.Outputs {
		outputs[*output.Key] = *output.StatusDetail
	}
	return &provider.JobStatus{
		ProviderJobID:  *resp.Job.Id,
		Status:         p.statusMap(*resp.Job.Status),
		ProviderStatus: map[string]interface{}{"outputs": outputs},
	}, nil
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

func elasticTranscoderProvider(cfg *config.Config) (provider.TranscodingProvider, error) {
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
