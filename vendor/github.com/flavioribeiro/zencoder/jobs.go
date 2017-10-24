package zencoder

import (
	"fmt"
)

type FileProgress struct {
	Id                   int64   `json:"id,omitempty"`
	State                string  `json:"state,omitempty"`
	CurrentEvent         string  `json:"current_event,omitempty"`
	CurrentEventProgress float64 `json:"current_event_progress,omitempty"`
	OverallProgress      float64 `json:"progress,omitempty"`
}

type JobProgress struct {
	State          JobState        `json:"state,omitempty"`
	JobProgress    float64         `json:"progress,omitempty"`
	InputProgress  *FileProgress   `json:"input,omitempty"`
	OutputProgress []*FileProgress `json:"outputs,omitempty"`
}

type JobState string

const (
	JobStatePending    = JobState("pending")
	JobStateWaiting    = JobState("waiting")
	JobStateProcessing = JobState("processing")
	JobStateAssigning  = JobState("assigning")
	JobStateFinished   = JobState("finished")
	JobStateCancelled  = JobState("cancelled")
	JobStateFailed     = JobState("failed")
)

// Response from CreateJob
type CreateJobResponse struct {
	Id      int64 `json:"id,omitempty"`
	Test    bool  `json:"test,omitempty"`
	Outputs []struct {
		Id    int64   `json:"id,omitempty"`
		Label *string `json:"label,omitempty"`
		Url   string  `json:"url,omitempty"`
	} `json:"outputs,omitempty"`
}

// A MediaFile
type MediaFile struct {
	Id                 int64   `json:"id,omitempty"`
	Url                string  `json:"url,omitempty"`
	Format             string  `json:"format,omitempty"`
	State              string  `json:"state,omitempty"`
	Test               bool    `json:"test,omitempty"`
	Privacy            bool    `json:"privacy"`
	Width              int32   `json:"width,omitempty"`
	Height             int32   `json:"height,omitempty"`
	FrameRate          float64 `json:"frame_rate,omitempty"`
	DurationInMs       int32   `json:"duration_in_ms,omitempty"`
	Channels           string  `json:"channels,omitempty"`
	AudioCodec         string  `json:"audio_codec,omitempty"`
	AudioBitrateInKbps int32   `json:"audio_bitrate_in_kbps,omitempty"`
	AudioSampleRate    int32   `json:"audio_sample_rate,omitempty"`
	VideoCodec         string  `json:"video_codec,omitempty"`
	VideoBitrateInKbps int32   `json:"video_bitrate_in_kbps,omitempty"`
	TotalBitrateInKbps int32   `json:"total_bitrate_in_kbps,omitempty"`
	MD5Checksum        string  `json:"md5_checksum,omitempty"`
	ErrorMessage       *string `json:"error_message,omitempty"`
	ErrorClass         *string `json:"error_class,omitempty"`
	Label              *string `json:"label,omitempty"`
	CreatedAt          string  `json:"created_at,omitempty"`
	FinishedAt         string  `json:"finished_at,omitempty"`
	UpdatedAt          string  `json:"updated_at,omitempty"`
	FileSizeInBytes    int64   `json:"file_size_bytes,omitempty"`
}

type InputMediaFile struct {
	MediaFile
	FileSizeInBytes int64 `json:"file_size_in_bytes,omitempty"`
	JobId           int64 `json:"job_id,omitempty"`
}

type OutputMediaFile struct {
	MediaFile
	FileSizeInBytes int64 `json:"file_size_in_bytes,omitempty"`
	JobId           int64 `json:"job_id,omitempty"`
}

// A Thumbnail
type Thumbnail struct {
	Id        int64  `json:"id,omitempty"`
	Url       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// A Job
type Job struct {
	Id               int64        `json:"id,omitempty"`
	PassThrough      *string      `json:"pass_through,omitempty"`
	State            string       `json:"state,omitempty"`
	InputMediaFile   *MediaFile   `json:"input_media_file,omitempty"`
	Test             bool         `json:"test,omitempty"`
	OutputMediaFiles []*MediaFile `json:"output_media_files,omitempty"`
	Thumbnails       []*Thumbnail `json:"thumbnails,omitempty"`
	CreatedAt        string       `json:"created_at,omitempty"`
	FinishedAt       string       `json:"finished_at,omitempty"`
	UpdatedAt        string       `json:"updated_at,omitempty"`
	SubmittedAt      string       `json:"submitted_at,omitempty"`
}

// Job Details wrapper
type JobDetails struct {
	Job *Job `json:"job,omitempty"`
}

// Create a Job
func (z *Zencoder) CreateJob(settings *EncodingSettings) (*CreateJobResponse, error) {
	var result CreateJobResponse

	if err := z.post("jobs", settings, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// List Jobs
func (z *Zencoder) ListJobs() ([]*JobDetails, error) {
	var result []*JobDetails

	if err := z.getBody("jobs.json", &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Get Job Details
func (z *Zencoder) GetJobDetails(id int64) (*JobDetails, error) {
	var result JobDetails

	if err := z.getBody(fmt.Sprintf("jobs/%d.json", id), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Job Progress
func (z *Zencoder) GetJobProgress(id int64) (*JobProgress, error) {
	var result JobProgress

	if err := z.getBody(fmt.Sprintf("jobs/%d/progress.json", id), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Resubmit a Job
func (z *Zencoder) ResubmitJob(id int64) error {
	return z.putNoContent(fmt.Sprintf("jobs/%d/resubmit.json", id))
}

// Cancel a Job
func (z *Zencoder) CancelJob(id int64) error {
	return z.putNoContent(fmt.Sprintf("jobs/%d/cancel.json", id))
}

// Finish a Live Job
func (z *Zencoder) FinishLiveJob(id int64) error {
	return z.putNoContent(fmt.Sprintf("jobs/%d/finish", id))
}
