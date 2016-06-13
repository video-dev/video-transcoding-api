package elastictranscoder

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
)

type failure struct {
	op  string
	err error
}

type fakeElasticTranscoder struct {
	*elastictranscoder.ElasticTranscoder
	jobs     map[string]*elastictranscoder.CreateJobInput
	failures chan failure
}

func newFakeElasticTranscoder() *fakeElasticTranscoder {
	return &fakeElasticTranscoder{
		ElasticTranscoder: &elastictranscoder.ElasticTranscoder{},
		failures:          make(chan failure, 1),
		jobs:              make(map[string]*elastictranscoder.CreateJobInput),
	}
}

func (c *fakeElasticTranscoder) CreateJob(input *elastictranscoder.CreateJobInput) (*elastictranscoder.CreateJobResponse, error) {
	if err := c.getError("CreateJob"); err != nil {
		return nil, err
	}
	id := fmt.Sprintf("job-%x", generateID())
	c.jobs[id] = input
	return &elastictranscoder.CreateJobResponse{
		Job: &elastictranscoder.Job{
			Id:         aws.String(id),
			Input:      input.Input,
			PipelineId: input.PipelineId,
			Status:     aws.String("Submitted"),
		},
	}, nil
}

func (c *fakeElasticTranscoder) CreatePreset(input *elastictranscoder.CreatePresetInput) (*elastictranscoder.CreatePresetOutput, error) {
	var presetID = *input.Name + "-abc123"
	return &elastictranscoder.CreatePresetOutput{
		Preset: &elastictranscoder.Preset{
			Audio:       input.Audio,
			Container:   input.Container,
			Description: input.Description,
			Name:        input.Name,
			Id:          &presetID,
			Thumbnails:  input.Thumbnails,
			Video:       input.Video,
		},
	}, nil
}

func (c *fakeElasticTranscoder) ReadPreset(input *elastictranscoder.ReadPresetInput) (*elastictranscoder.ReadPresetOutput, error) {
	container := "mp4"
	if strings.Contains(*input.Id, "hls") {
		container = "ts"
	}
	return &elastictranscoder.ReadPresetOutput{
		Preset: &elastictranscoder.Preset{
			Id:        input.Id,
			Name:      input.Id,
			Container: aws.String(container),
		},
	}, nil
}

func (c *fakeElasticTranscoder) ReadJob(input *elastictranscoder.ReadJobInput) (*elastictranscoder.ReadJobOutput, error) {
	if err := c.getError("ReadJob"); err != nil {
		return nil, err
	}
	createJobInput, ok := c.jobs[*input.Id]
	if !ok {
		return nil, errors.New("job not found")
	}
	outputs := make([]*elastictranscoder.JobOutput, len(createJobInput.Outputs))
	for i, createJobOutput := range createJobInput.Outputs {
		outputs[i] = &elastictranscoder.JobOutput{
			Key:          createJobOutput.Key,
			StatusDetail: aws.String("it's finished!"),
		}
	}
	return &elastictranscoder.ReadJobOutput{
		Job: &elastictranscoder.Job{
			Id:         input.Id,
			Input:      createJobInput.Input,
			PipelineId: createJobInput.PipelineId,
			Status:     aws.String("Complete"),
			Outputs:    outputs,
		},
	}, nil
}

func (c *fakeElasticTranscoder) ReadPipeline(input *elastictranscoder.ReadPipelineInput) (*elastictranscoder.ReadPipelineOutput, error) {
	if err := c.getError("ReadPipeline"); err != nil {
		return nil, err
	}
	return &elastictranscoder.ReadPipelineOutput{
		Pipeline: &elastictranscoder.Pipeline{
			Id:           input.Id,
			Name:         aws.String("nice pipeline"),
			OutputBucket: aws.String("some bucket"),
		},
	}, nil
}

func (c *fakeElasticTranscoder) prepareFailure(op string, err error) {
	c.failures <- failure{op: op, err: err}
}

func (c *fakeElasticTranscoder) getError(op string) error {
	select {
	case prepFailure := <-c.failures:
		if prepFailure.op == op {
			return prepFailure.err
		}
		c.failures <- prepFailure
	default:
	}
	return nil
}

func generateID() []byte {
	var b [4]byte
	rand.Read(b[:])
	return b[:]
}
