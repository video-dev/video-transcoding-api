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

func (z *FakeZencoder) CancelJob(id int64) error {
	return nil
}

func (z *FakeZencoder) GetJobProgress(id int64) (*zencoderClient.JobProgress, error) {
	return &zencoderClient.JobProgress{
		State:       "Transcoding",
		JobProgress: 10,
	}, nil
}

func (z *FakeZencoder) GetJobDetails(id int64) (*zencoderClient.JobDetails, error) {
	return &zencoderClient.JobDetails{
		Job: &zencoderClient.Job{
			OutputMediaFiles: []*zencoderClient.MediaFile{
				{
					Url:          "http://nyt.net/output1.mp4",
					Format:       "mp4",
					VideoCodec:   "h264",
					Width:        1920,
					Height:       1080,
					DurationInMs: 10000,
				},
				{
					Url:          "http://nyt.net/output2.webm",
					Format:       "webm",
					VideoCodec:   "vp8",
					Width:        1080,
					Height:       720,
					DurationInMs: 10000,
				},
			},
		},
	}, nil
}
