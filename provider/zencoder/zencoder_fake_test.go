package zencoder

import (
	zencoderClient "github.com/brandscreen/zencoder"
)

type FakeZencoder struct {
}

func (z *FakeZencoder) CreateJob(settings *zencoderClient.EncodingSettings) (*zencoderClient.CreateJobResponse, error) {
	return &zencoderClient.CreateJobResponse{
		Id: 123,
	}, nil
}

func (z *FakeZencoder) ListJobs() ([]*zencoderClient.JobDetails, error) {
	return []*zencoderClient.JobDetails{}, nil
}
