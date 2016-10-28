package zencoder

import (
	zencoderClient "github.com/brandscreen/zencoder"
)

type FakeZencoder struct {
}

func (z *FakeZencoder) CreateJob(settings *zencoderClient.EncodingSettings) (*zencoderClient.CreateJobResponse, error) {
	return &zencoderClient.CreateJobResponse{}, nil
}
