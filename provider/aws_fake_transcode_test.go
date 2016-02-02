package provider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
)

type failure struct {
	op  string
	err error
}

type fakeElasticTranscoder struct {
	*elastictranscoder.ElasticTranscoder
	jobs     []*elastictranscoder.CreateJobInput
	failures chan failure
}

func newFakeElasticTranscoder() *fakeElasticTranscoder {
	return &fakeElasticTranscoder{
		ElasticTranscoder: &elastictranscoder.ElasticTranscoder{},
		failures:          make(chan failure, 1),
	}
}

func (c *fakeElasticTranscoder) CreateJob(input *elastictranscoder.CreateJobInput) (*elastictranscoder.CreateJobResponse, error) {
	if err := c.getError("CreateJob"); err != nil {
		return nil, err
	}
	c.jobs = append(c.jobs, input)
	return &elastictranscoder.CreateJobResponse{
		Job: &elastictranscoder.Job{
			Id:         aws.String("job-" + generateID(4)),
			Input:      input.Input,
			PipelineId: input.PipelineId,
			Status:     aws.String("Submitted"),
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
