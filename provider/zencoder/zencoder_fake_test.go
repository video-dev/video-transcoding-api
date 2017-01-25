package zencoder

import "github.com/flavioribeiro/zencoder"

type FakeZencoder struct {
}

func (z *FakeZencoder) CreateJob(settings *zencoder.EncodingSettings) (*zencoder.CreateJobResponse, error) {
	return &zencoder.CreateJobResponse{
		Id: 123,
	}, nil
}

func (z *FakeZencoder) CancelJob(id int64) error {
	return nil
}

func (z *FakeZencoder) GetJobProgress(id int64) (*zencoder.JobProgress, error) {
	if id == 1234567890 {
		return &zencoder.JobProgress{State: "processing", JobProgress: 10}, nil
	}
	return &zencoder.JobProgress{State: "finished", JobProgress: 0}, nil
}

func (z *FakeZencoder) GetJobDetails(id int64) (*zencoder.JobDetails, error) {
	state := "finished"
	if id == 1234567890 || id == 837958345 {
		state = "processing"
	}
	return &zencoder.JobDetails{
		Job: &zencoder.Job{
			State: state,
			InputMediaFile: &zencoder.MediaFile{
				Url:          "http://nyt.net/input.mov",
				Format:       "mov",
				VideoCodec:   "ProRes422",
				Width:        1920,
				Height:       1080,
				DurationInMs: 50000,
			},
			CreatedAt:   "2016-11-05T05:02:57Z",
			FinishedAt:  "2016-11-05T05:02:57Z",
			UpdatedAt:   "2016-11-05T05:02:57Z",
			SubmittedAt: "2016-11-05T05:02:57Z",
			OutputMediaFiles: []*zencoder.MediaFile{
				{
					Url:             "https://mybucket.s3.amazonaws.com/destination-dir/output1.mp4",
					Format:          "mp4",
					VideoCodec:      "h264",
					Width:           1920,
					Height:          1080,
					DurationInMs:    10000,
					FileSizeInBytes: 66885256,
				},
				{
					Url:             "https://mybucket.s3.amazonaws.com/destination-dir/output2.webm",
					Format:          "webm",
					VideoCodec:      "vp8",
					Width:           1080,
					Height:          720,
					DurationInMs:    10000,
					FileSizeInBytes: 92140022,
				},
			},
		},
	}, nil
}

func (z *FakeZencoder) GetVodUsage(settings *zencoder.ReportSettings) (*zencoder.VodUsage, error) {
	return &zencoder.VodUsage{}, nil
}
