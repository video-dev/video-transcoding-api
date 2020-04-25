package mediaconvert

import (
	"net/http"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
)

// testMediaConvertClient is an implementation of the mediaconvertClient interface
// to be used with tests
type testMediaConvertClient struct {
	t *testing.T

	createPresetCalledWith *mediaconvert.CreatePresetInput
	getPresetCalledWith    *string
	deletePresetCalledWith string
	createJobCalledWith    mediaconvert.CreateJobInput
	cancelJobCalledWith    string
	listJobsCalled         bool

	jobReturnedByGetJob      mediaconvert.Job
	jobIDReturnedByCreateJob string
	getPresetContainerType   mediaconvert.ContainerType
}

func (c *testMediaConvertClient) CreatePresetRequest(input *mediaconvert.CreatePresetInput) mediaconvert.CreatePresetRequest {
	c.createPresetCalledWith = input
	return mediaconvert.CreatePresetRequest{
		Request: &aws.Request{
			Retryer:     retry.NewStandard(),
			HTTPRequest: &http.Request{}, Data: &mediaconvert.CreatePresetOutput{
				Preset: &mediaconvert.Preset{
					Name: input.Name,
					Settings: &mediaconvert.PresetSettings{
						ContainerSettings: &mediaconvert.ContainerSettings{
							Container: input.Settings.ContainerSettings.Container,
						},
					},
				},
			},
		},
	}
}

func (c *testMediaConvertClient) GetJobRequest(*mediaconvert.GetJobInput) mediaconvert.GetJobRequest {
	return mediaconvert.GetJobRequest{Request: &aws.Request{
		Retryer:     retry.NewStandard(),
		HTTPRequest: &http.Request{},
		Data: &mediaconvert.GetJobOutput{
			Job: &c.jobReturnedByGetJob,
		},
	}}
}

func (c *testMediaConvertClient) ListJobsRequest(*mediaconvert.ListJobsInput) mediaconvert.ListJobsRequest {
	c.listJobsCalled = true
	return mediaconvert.ListJobsRequest{Request: &aws.Request{
		Retryer:     retry.NewStandard(),
		HTTPRequest: &http.Request{},
		Data:        &mediaconvert.ListJobsOutput{},
	}}
}

func (c *testMediaConvertClient) CreateJobRequest(input *mediaconvert.CreateJobInput) mediaconvert.CreateJobRequest {
	c.createJobCalledWith = *input
	return mediaconvert.CreateJobRequest{
		Request: &aws.Request{
			Retryer:     retry.NewStandard(),
			HTTPRequest: &http.Request{}, Data: &mediaconvert.CreateJobOutput{
				Job: &mediaconvert.Job{
					Id: aws.String(c.jobIDReturnedByCreateJob),
				},
			}},
	}
}

func (c *testMediaConvertClient) CancelJobRequest(input *mediaconvert.CancelJobInput) mediaconvert.CancelJobRequest {
	c.cancelJobCalledWith = *input.Id
	return mediaconvert.CancelJobRequest{Request: &aws.Request{
		Retryer:     retry.NewStandard(),
		HTTPRequest: &http.Request{},
		Data:        &mediaconvert.CancelJobOutput{},
	}}
}

func (c *testMediaConvertClient) GetPresetRequest(input *mediaconvert.GetPresetInput) mediaconvert.GetPresetRequest {
	// atomically set the value of getPresetCalledWith to avoid data races,
	// should probably take a different approach?
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&c.getPresetCalledWith)), unsafe.Pointer(input.Name))

	return mediaconvert.GetPresetRequest{
		Request: &aws.Request{
			Retryer:     retry.NewStandard(),
			HTTPRequest: &http.Request{}, Data: &mediaconvert.GetPresetOutput{
				Preset: &mediaconvert.Preset{
					Name: input.Name,
					Settings: &mediaconvert.PresetSettings{
						ContainerSettings: &mediaconvert.ContainerSettings{
							Container: c.getPresetContainerType,
						},
					},
				},
			},
		},
	}
}

func (c *testMediaConvertClient) DeletePresetRequest(input *mediaconvert.DeletePresetInput) mediaconvert.DeletePresetRequest {
	c.deletePresetCalledWith = *input.Name
	return mediaconvert.DeletePresetRequest{Request: &aws.Request{
		Retryer:     retry.NewStandard(),
		HTTPRequest: &http.Request{},
		Data:        &mediaconvert.DeletePresetOutput{},
	}}
}
