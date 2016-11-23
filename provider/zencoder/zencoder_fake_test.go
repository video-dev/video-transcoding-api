package zencoder

import (
	zencoderClient "github.com/flavioribeiro/zencoder"
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
		State:       "processing",
		JobProgress: 10,
	}, nil
}

func (z *FakeZencoder) GetJobDetails(id int64) (*zencoderClient.JobDetails, error) {
	return &zencoderClient.JobDetails{
		Job: &zencoderClient.Job{
			InputMediaFile: &zencoderClient.MediaFile{
				Url:          "http://nyt.net/input.mov",
				Format:       "mov",
				VideoCodec:   "ProRes422",
				Width:        1920,
				Height:       1080,
				DurationInMs: 10000,
			},
			CreatedAt:   "2016-11-05T05:02:57Z",
			FinishedAt:  "2016-11-05T05:02:57Z",
			UpdatedAt:   "2016-11-05T05:02:57Z",
			SubmittedAt: "2016-11-05T05:02:57Z",
			OutputMediaFiles: []*zencoderClient.MediaFile{
				{
					Url:          "https://mybucket.s3.amazonaws.com/destination-dir/output1.mp4",
					Format:       "mp4",
					VideoCodec:   "h264",
					Width:        1920,
					Height:       1080,
					DurationInMs: 10000,
				},
				{
					Url:          "https://mybucket.s3.amazonaws.com/destination-dir/output2.webm",
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
