package mediaconvert

import (
	"context"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
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

	jobReturnedByGetJob      types.Job
	jobIDReturnedByCreateJob string
	getPresetContainerType   types.ContainerType
}

func (c *testMediaConvertClient) CreatePreset(_ context.Context, input *mediaconvert.CreatePresetInput, _ ...func(*mediaconvert.Options)) (*mediaconvert.CreatePresetOutput, error) {
	c.createPresetCalledWith = input
	return &mediaconvert.CreatePresetOutput{
		Preset: &types.Preset{
			Name: input.Name,
			Settings: &types.PresetSettings{
				ContainerSettings: &types.ContainerSettings{
					Container: input.Settings.ContainerSettings.Container,
				},
			},
		},
	}, nil
}

func (c *testMediaConvertClient) GetJob(context.Context, *mediaconvert.GetJobInput, ...func(*mediaconvert.Options)) (*mediaconvert.GetJobOutput, error) {
	return &mediaconvert.GetJobOutput{
		Job: &c.jobReturnedByGetJob,
	}, nil
}

func (c *testMediaConvertClient) ListJobs(context.Context, *mediaconvert.ListJobsInput, ...func(*mediaconvert.Options)) (*mediaconvert.ListJobsOutput, error) {
	c.listJobsCalled = true
	return &mediaconvert.ListJobsOutput{}, nil
}

func (c *testMediaConvertClient) CreateJob(_ context.Context, input *mediaconvert.CreateJobInput, _ ...func(*mediaconvert.Options)) (*mediaconvert.CreateJobOutput, error) {
	c.createJobCalledWith = *input
	return &mediaconvert.CreateJobOutput{
		Job: &types.Job{
			Id: aws.String(c.jobIDReturnedByCreateJob),
		},
	}, nil
}

func (c *testMediaConvertClient) CancelJob(_ context.Context, input *mediaconvert.CancelJobInput, _ ...func(*mediaconvert.Options)) (*mediaconvert.CancelJobOutput, error) {
	c.cancelJobCalledWith = *input.Id
	return &mediaconvert.CancelJobOutput{}, nil
}

func (c *testMediaConvertClient) GetPreset(_ context.Context, input *mediaconvert.GetPresetInput, _ ...func(*mediaconvert.Options)) (*mediaconvert.GetPresetOutput, error) {
	// atomically set the value of getPresetCalledWith to avoid data races,
	// should probably take a different approach?
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&c.getPresetCalledWith)), unsafe.Pointer(input.Name))

	return &mediaconvert.GetPresetOutput{
		Preset: &types.Preset{
			Name: input.Name,
			Settings: &types.PresetSettings{
				ContainerSettings: &types.ContainerSettings{
					Container: c.getPresetContainerType,
				},
			},
		},
	}, nil
}

func (c *testMediaConvertClient) DeletePreset(_ context.Context, input *mediaconvert.DeletePresetInput, _ ...func(*mediaconvert.Options)) (*mediaconvert.DeletePresetOutput, error) {
	c.deletePresetCalledWith = *input.Name
	return &mediaconvert.DeletePresetOutput{}, nil
}
