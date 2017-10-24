package zencoder

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateJob(t *testing.T) {
	expectedStatus := http.StatusCreated
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatus)
		if expectedStatus != http.StatusCreated {
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{"id": 1234,"outputs": [{"id": 4321}]}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	var settings EncodingSettings
	resp, err := zc.CreateJob(&settings)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if resp == nil {
		t.Fatal("Expected a response")
	}

	if resp.Id != 1234 {
		t.Fatal("Expected Id=1234", resp.Id)
	}

	if len(resp.Outputs) != 1 {
		t.Fatal("Expected one output", len(resp.Outputs))
	}

	if resp.Outputs[0].Id != 4321 {
		t.Fatal("Expected Id=4321", resp.Outputs[0].Id)
	}

	expectedStatus = http.StatusInternalServerError
	resp, err = zc.CreateJob(&settings)
	if err == nil {
		t.Fatal("Expected error")
	}

	if resp != nil {
		t.Fatal("Expected no response")
	}

	returnBody = false
	expectedStatus = http.StatusOK

	resp, err = zc.CreateJob(&settings)
	if err == nil {
		t.Fatal("Expected error")
	}

	if resp != nil {
		t.Fatal("Expected no response")
	}

	returnBody = true
	srv.Close()

	resp, err = zc.CreateJob(&settings)
	if err == nil {
		t.Fatal("Expected error")
	}

	if resp != nil {
		t.Fatal("Expected no response")
	}
}

func TestListJobs(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs.json", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `[{
  "job": {
    "created_at": "2010-01-01T00:00:00Z",
    "finished_at": "2010-01-01T00:00:00Z",
    "updated_at": "2010-01-01T00:00:00Z",
    "submitted_at": "2010-01-01T00:00:00Z",
    "pass_through": null,
    "id": 1,
    "input_media_file": {
      "format": "mpeg4",
      "created_at": "2010-01-01T00:00:00Z",
      "frame_rate": 29,
      "finished_at": "2010-01-01T00:00:00Z",
      "updated_at": "2010-01-01T00:00:00Z",
      "duration_in_ms": 24883,
      "audio_sample_rate": 48000,
      "url": "s3://bucket/test.mp4",
      "id": 1,
      "error_message": null,
      "error_class": null,
      "audio_bitrate_in_kbps": 95,
      "audio_codec": "aac",
      "height": 352,
      "file_size_bytes": 1862748,
      "video_codec": "h264",
      "test": false,
      "total_bitrate_in_kbps": 593,
      "channels": "2",
      "width": 624,
      "video_bitrate_in_kbps": 498,
      "state": "finished"
    },
    "test": false,
    "output_media_files": [{
      "format": "mpeg4",
      "created_at": "2010-01-01T00:00:00Z",
      "frame_rate": 29,
      "finished_at": "2010-01-01T00:00:00Z",
      "updated_at": "2010-01-01T00:00:00Z",
      "duration_in_ms": 24883,
      "audio_sample_rate": 44100,
      "url": "http://s3.amazonaws.com/bucket/video.mp4",
      "id": 1,
      "error_message": null,
      "error_class": null,
      "audio_bitrate_in_kbps": 92,
      "audio_codec": "aac",
      "height": 352,
      "file_size_bytes": 1386663,
      "video_codec": "h264",
      "test": false,
      "total_bitrate_in_kbps": 443,
      "channels": "2",
      "width": 624,
      "video_bitrate_in_kbps": 351,
      "state": "finished",
      "label": "Web"
    }],
    "thumbnails": [{
      "created_at": "2010-01-01T00:00:00Z",
      "updated_at": "2010-01-01T00:00:00Z",
      "url": "http://s3.amazonaws.com/bucket/video/frame_0000.png",
      "id": 1
    }],
    "state": "finished"
  }
}]`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	jobs, err := zc.ListJobs()
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if jobs == nil {
		t.Fatal("Expected jobs")
	}

	if len(jobs) != 1 {
		t.Fatal("Expected 1 job", len(jobs))
	}

	if jobs[0].Job.CreatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.CreatedAt)
	}
	if jobs[0].Job.FinishedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.FinishedAt)
	}
	if jobs[0].Job.UpdatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.UpdatedAt)
	}
	if jobs[0].Job.SubmittedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.SubmittedAt)
	}
	if jobs[0].Job.PassThrough != nil {
		t.Fatal("Expected nil, got", jobs[0].Job.PassThrough)
	}
	if jobs[0].Job.Id != 1 {
		t.Fatal("Expected 1, got", jobs[0].Job.Id)
	}
	if jobs[0].Job.InputMediaFile.Format != "mpeg4" {
		t.Fatal("Expected mpeg4, got", jobs[0].Job.InputMediaFile.Format)
	}
	if jobs[0].Job.InputMediaFile.CreatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.InputMediaFile.CreatedAt)
	}
	if jobs[0].Job.InputMediaFile.FrameRate != 29 {
		t.Fatal("Expected 29, got", jobs[0].Job.InputMediaFile.FrameRate)
	}
	if jobs[0].Job.InputMediaFile.FinishedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.InputMediaFile.FinishedAt)
	}
	if jobs[0].Job.InputMediaFile.UpdatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.InputMediaFile.UpdatedAt)
	}
	if jobs[0].Job.InputMediaFile.DurationInMs != 24883 {
		t.Fatal("Expected 24883, got", jobs[0].Job.InputMediaFile.DurationInMs)
	}
	if jobs[0].Job.InputMediaFile.AudioSampleRate != 48000 {
		t.Fatal("Expected 48000, got", jobs[0].Job.InputMediaFile.AudioSampleRate)
	}
	if jobs[0].Job.InputMediaFile.Url != "s3://bucket/test.mp4" {
		t.Fatal("Expected s3://bucket/test.mp4, got", jobs[0].Job.InputMediaFile.Url)
	}
	if jobs[0].Job.InputMediaFile.Id != 1 {
		t.Fatal("Expected 1, got", jobs[0].Job.InputMediaFile.Id)
	}
	if jobs[0].Job.InputMediaFile.ErrorMessage != nil {
		t.Fatal("Expected nil, got", jobs[0].Job.InputMediaFile.ErrorMessage)
	}
	if jobs[0].Job.InputMediaFile.ErrorClass != nil {
		t.Fatal("Expected nil, got", jobs[0].Job.InputMediaFile.ErrorClass)
	}
	if jobs[0].Job.InputMediaFile.AudioBitrateInKbps != 95 {
		t.Fatal("Expected 95, got", jobs[0].Job.InputMediaFile.AudioBitrateInKbps)
	}
	if jobs[0].Job.InputMediaFile.AudioCodec != "aac" {
		t.Fatal("Expected aac, got", jobs[0].Job.InputMediaFile.AudioCodec)
	}
	if jobs[0].Job.InputMediaFile.Height != 352 {
		t.Fatal("Expected 352, got", jobs[0].Job.InputMediaFile.Height)
	}
	if jobs[0].Job.InputMediaFile.FileSizeInBytes != 1862748 {
		t.Fatal("Expected 1862748, got", jobs[0].Job.InputMediaFile.FileSizeInBytes)
	}
	if jobs[0].Job.InputMediaFile.VideoCodec != "h264" {
		t.Fatal("Expected h264, got", jobs[0].Job.InputMediaFile.VideoCodec)
	}
	if jobs[0].Job.InputMediaFile.Test != false {
		t.Fatal("Expected false, got", jobs[0].Job.InputMediaFile.Test)
	}
	if jobs[0].Job.InputMediaFile.TotalBitrateInKbps != 593 {
		t.Fatal("Expected 593, got", jobs[0].Job.InputMediaFile.TotalBitrateInKbps)
	}
	if jobs[0].Job.InputMediaFile.Channels != "2" {
		t.Fatal("Expected 2, got", jobs[0].Job.InputMediaFile.Channels)
	}
	if jobs[0].Job.InputMediaFile.Width != 624 {
		t.Fatal("Expected 624, got", jobs[0].Job.InputMediaFile.Width)
	}
	if jobs[0].Job.InputMediaFile.VideoBitrateInKbps != 498 {
		t.Fatal("Expected 498, got", jobs[0].Job.InputMediaFile.VideoBitrateInKbps)
	}
	if jobs[0].Job.InputMediaFile.State != "finished" {
		t.Fatal("Expected finished, got", jobs[0].Job.InputMediaFile.State)
	}
	if jobs[0].Job.Test != false {
		t.Fatal("Expected false, got", jobs[0].Job.Test)
	}
	if len(jobs[0].Job.OutputMediaFiles) != 1 {
		t.Fatal("Expected 1 outputs, got", len(jobs[0].Job.OutputMediaFiles))
	}
	if jobs[0].Job.OutputMediaFiles[0].Format != "mpeg4" {
		t.Fatal("Expected mpeg4, got", jobs[0].Job.OutputMediaFiles[0].Format)
	}
	if jobs[0].Job.OutputMediaFiles[0].CreatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.OutputMediaFiles[0].CreatedAt)
	}
	if jobs[0].Job.OutputMediaFiles[0].FrameRate != 29 {
		t.Fatal("Expected 29, got", jobs[0].Job.OutputMediaFiles[0].FrameRate)
	}
	if jobs[0].Job.OutputMediaFiles[0].FinishedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.OutputMediaFiles[0].FinishedAt)
	}
	if jobs[0].Job.OutputMediaFiles[0].UpdatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.OutputMediaFiles[0].UpdatedAt)
	}
	if jobs[0].Job.OutputMediaFiles[0].DurationInMs != 24883 {
		t.Fatal("Expected 24883, got", jobs[0].Job.OutputMediaFiles[0].DurationInMs)
	}
	if jobs[0].Job.OutputMediaFiles[0].AudioSampleRate != 44100 {
		t.Fatal("Expected 44100, got", jobs[0].Job.OutputMediaFiles[0].AudioSampleRate)
	}
	if jobs[0].Job.OutputMediaFiles[0].Url != "http://s3.amazonaws.com/bucket/video.mp4" {
		t.Fatal("Expected http://s3.amazonaws.com/bucket/video.mp4, got", jobs[0].Job.OutputMediaFiles[0].Url)
	}
	if jobs[0].Job.OutputMediaFiles[0].Id != 1 {
		t.Fatal("Expected 1, got", jobs[0].Job.OutputMediaFiles[0].Id)
	}
	if jobs[0].Job.OutputMediaFiles[0].ErrorMessage != nil {
		t.Fatal("Expected nil, got", jobs[0].Job.OutputMediaFiles[0].ErrorMessage)
	}
	if jobs[0].Job.OutputMediaFiles[0].ErrorClass != nil {
		t.Fatal("Expected nil, got", jobs[0].Job.OutputMediaFiles[0].ErrorClass)
	}
	if jobs[0].Job.OutputMediaFiles[0].AudioBitrateInKbps != 92 {
		t.Fatal("Expected 92, got", jobs[0].Job.OutputMediaFiles[0].AudioBitrateInKbps)
	}
	if jobs[0].Job.OutputMediaFiles[0].AudioCodec != "aac" {
		t.Fatal("Expected aac, got", jobs[0].Job.OutputMediaFiles[0].AudioCodec)
	}
	if jobs[0].Job.OutputMediaFiles[0].Height != 352 {
		t.Fatal("Expected 352, got", jobs[0].Job.OutputMediaFiles[0].Height)
	}
	if jobs[0].Job.OutputMediaFiles[0].FileSizeInBytes != 1386663 {
		t.Fatal("Expected 1386663, got", jobs[0].Job.OutputMediaFiles[0].FileSizeInBytes)
	}
	if jobs[0].Job.OutputMediaFiles[0].VideoCodec != "h264" {
		t.Fatal("Expected h264, got", jobs[0].Job.OutputMediaFiles[0].VideoCodec)
	}
	if jobs[0].Job.OutputMediaFiles[0].Test != false {
		t.Fatal("Expected false, got", jobs[0].Job.OutputMediaFiles[0].Test)
	}
	if jobs[0].Job.OutputMediaFiles[0].TotalBitrateInKbps != 443 {
		t.Fatal("Expected 443, got", jobs[0].Job.OutputMediaFiles[0].TotalBitrateInKbps)
	}
	if jobs[0].Job.OutputMediaFiles[0].Channels != "2" {
		t.Fatal("Expected 2, got", jobs[0].Job.OutputMediaFiles[0].Channels)
	}
	if jobs[0].Job.OutputMediaFiles[0].Width != 624 {
		t.Fatal("Expected 624, got", jobs[0].Job.OutputMediaFiles[0].Width)
	}
	if jobs[0].Job.OutputMediaFiles[0].VideoBitrateInKbps != 351 {
		t.Fatal("Expected 351, got", jobs[0].Job.OutputMediaFiles[0].VideoBitrateInKbps)
	}
	if jobs[0].Job.OutputMediaFiles[0].State != "finished" {
		t.Fatal("Expected finished, got", jobs[0].Job.OutputMediaFiles[0].State)
	}
	if jobs[0].Job.OutputMediaFiles[0].Label == nil || *jobs[0].Job.OutputMediaFiles[0].Label != "Web" {
		t.Fatal("Expected Web, got", jobs[0].Job.OutputMediaFiles[0].Label)
	}
	if len(jobs[0].Job.Thumbnails) != 1 {
		t.Fatal("Expected 1 thumbnail, got", len(jobs[0].Job.Thumbnails))
	}
	if jobs[0].Job.Thumbnails[0].CreatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.Thumbnails[0].CreatedAt)
	}
	if jobs[0].Job.Thumbnails[0].UpdatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", jobs[0].Job.Thumbnails[0].UpdatedAt)
	}
	if jobs[0].Job.Thumbnails[0].Url != "http://s3.amazonaws.com/bucket/video/frame_0000.png" {
		t.Fatal("Expected http://s3.amazonaws.com/bucket/video/frame_0000.png, got", jobs[0].Job.Thumbnails[0].Url)
	}
	if jobs[0].Job.Thumbnails[0].Id != 1 {
		t.Fatal("Expected 1, got", jobs[0].Job.Thumbnails[0].Id)
	}
	if jobs[0].Job.State != "finished" {
		t.Fatal("Expected finished, got", jobs[0].Job.State)
	}

	expectedStatus = http.StatusInternalServerError
	jobs, err = zc.ListJobs()
	if err == nil {
		t.Fatal("Expected error")
	}

	if jobs != nil {
		t.Fatal("Expected no response")
	}

	expectedStatus = http.StatusOK
	returnBody = false
	jobs, err = zc.ListJobs()
	if err == nil {
		t.Fatal("Expected error")
	}

	if jobs != nil {
		t.Fatal("Expected no response")
	}

	srv.Close()
	returnBody = true
	jobs, err = zc.ListJobs()
	if err == nil {
		t.Fatal("Expected error")
	}

	if jobs != nil {
		t.Fatal("Expected no response")
	}
}

func TestGetJobDetails(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs/123.json", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "job": {
    "created_at": "2010-01-01T00:00:00Z",
    "finished_at": "2010-01-01T00:00:00Z",
    "updated_at": "2010-01-01T00:00:00Z",
    "submitted_at": "2010-01-01T00:00:00Z",
    "pass_through": null,
    "id": 1,
    "input_media_file": {
      "format": "mpeg4",
      "created_at": "2010-01-01T00:00:00Z",
      "frame_rate": 29,
      "finished_at": "2010-01-01T00:00:00Z",
      "updated_at": "2010-01-01T00:00:00Z",
      "duration_in_ms": 24883,
      "audio_sample_rate": 48000,
      "url": "s3://bucket/test.mp4",
      "id": 1,
      "error_message": null,
      "error_class": null,
      "audio_bitrate_in_kbps": 95,
      "audio_codec": "aac",
      "height": 352,
      "file_size_bytes": 1862748,
      "video_codec": "h264",
      "test": false,
      "total_bitrate_in_kbps": 593,
      "channels": "2",
      "width": 624,
      "video_bitrate_in_kbps": 498,
      "state": "finished",
      "md5_checksum":"7f106918e02a69466afa0ee014174143"
    },
    "test": false,
    "output_media_files": [{
      "format": "mpeg4",
      "created_at": "2010-01-01T00:00:00Z",
      "frame_rate": 29,
      "finished_at": "2010-01-01T00:00:00Z",
      "updated_at": "2010-01-01T00:00:00Z",
      "duration_in_ms": 24883,
      "audio_sample_rate": 44100,
      "url": "http://s3.amazonaws.com/bucket/video.mp4",
      "id": 1,
      "error_message": null,
      "error_class": null,
      "audio_bitrate_in_kbps": 92,
      "audio_codec": "aac",
      "height": 352,
      "file_size_bytes": 1386663,
      "video_codec": "h264",
      "test": false,
      "total_bitrate_in_kbps": 443,
      "channels": "2",
      "width": 624,
      "video_bitrate_in_kbps": 351,
      "state": "finished",
      "label": "Web",
      "md5_checksum":"7f106918e02a69466afa0ee014172496"
    }],
    "thumbnails": [{
      "created_at": "2010-01-01T00:00:00Z",
      "updated_at": "2010-01-01T00:00:00Z",
      "url": "http://s3.amazonaws.com/bucket/video/frame_0000.png",
      "id": 1
    }],
    "state": "finished"
  }
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	details, err := zc.GetJobDetails(123)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if details == nil {
		t.Fatal("Expected details")
	}

	if details.Job.CreatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.CreatedAt)
	}
	if details.Job.FinishedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.FinishedAt)
	}
	if details.Job.UpdatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.UpdatedAt)
	}
	if details.Job.SubmittedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.SubmittedAt)
	}
	if details.Job.PassThrough != nil {
		t.Fatal("Expected nil, got", details.Job.PassThrough)
	}
	if details.Job.Id != 1 {
		t.Fatal("Expected 1, got", details.Job.Id)
	}
	if details.Job.Test != false {
		t.Fatal("Expected false, got", details.Job.Test)
	}
	if details.Job.State != "finished" {
		t.Fatal("Expected finished, got", details.Job.State)
	}

	if details.Job.InputMediaFile.Format != "mpeg4" {
		t.Fatal("Expected mpeg4, got", details.Job.InputMediaFile.Format)
	}
	if details.Job.InputMediaFile.CreatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.InputMediaFile.CreatedAt)
	}
	if details.Job.InputMediaFile.FrameRate != 29 {
		t.Fatal("Expected 29, got", details.Job.InputMediaFile.FrameRate)
	}
	if details.Job.InputMediaFile.FinishedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.InputMediaFile.FinishedAt)
	}
	if details.Job.InputMediaFile.UpdatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.InputMediaFile.UpdatedAt)
	}
	if details.Job.InputMediaFile.DurationInMs != 24883 {
		t.Fatal("Expected 24883, got", details.Job.InputMediaFile.DurationInMs)
	}
	if details.Job.InputMediaFile.AudioSampleRate != 48000 {
		t.Fatal("Expected 48000, got", details.Job.InputMediaFile.AudioSampleRate)
	}
	if details.Job.InputMediaFile.Url != "s3://bucket/test.mp4" {
		t.Fatal("Expected s3://bucket/test.mp4, got", details.Job.InputMediaFile.Url)
	}
	if details.Job.InputMediaFile.Id != 1 {
		t.Fatal("Expected 1, got", details.Job.InputMediaFile.Id)
	}
	if details.Job.InputMediaFile.ErrorMessage != nil {
		t.Fatal("Expected nil, got", details.Job.InputMediaFile.ErrorMessage)
	}
	if details.Job.InputMediaFile.ErrorClass != nil {
		t.Fatal("Expected nil, got", details.Job.InputMediaFile.ErrorClass)
	}
	if details.Job.InputMediaFile.AudioBitrateInKbps != 95 {
		t.Fatal("Expected 95, got", details.Job.InputMediaFile.AudioBitrateInKbps)
	}
	if details.Job.InputMediaFile.AudioCodec != "aac" {
		t.Fatal("Expected aac, got", details.Job.InputMediaFile.AudioCodec)
	}
	if details.Job.InputMediaFile.Height != 352 {
		t.Fatal("Expected 352, got", details.Job.InputMediaFile.Height)
	}
	if details.Job.InputMediaFile.FileSizeInBytes != 1862748 {
		t.Fatal("Expected 1862748, got", details.Job.InputMediaFile.FileSizeInBytes)
	}
	if details.Job.InputMediaFile.VideoCodec != "h264" {
		t.Fatal("Expected h264, got", details.Job.InputMediaFile.VideoCodec)
	}
	if details.Job.InputMediaFile.Test != false {
		t.Fatal("Expected false, got", details.Job.InputMediaFile.Test)
	}
	if details.Job.InputMediaFile.TotalBitrateInKbps != 593 {
		t.Fatal("Expected 593, got", details.Job.InputMediaFile.TotalBitrateInKbps)
	}
	if details.Job.InputMediaFile.Channels != "2" {
		t.Fatal("Expected 2, got", details.Job.InputMediaFile.Channels)
	}
	if details.Job.InputMediaFile.Width != 624 {
		t.Fatal("Expected 624, got", details.Job.InputMediaFile.Width)
	}
	if details.Job.InputMediaFile.VideoBitrateInKbps != 498 {
		t.Fatal("Expected 498, got", details.Job.InputMediaFile.VideoBitrateInKbps)
	}
	if details.Job.InputMediaFile.State != "finished" {
		t.Fatal("Expected finished, got", details.Job.InputMediaFile.State)
	}
	if details.Job.InputMediaFile.MD5Checksum != "7f106918e02a69466afa0ee014174143" {
		t.Fatal("Expected 7f106918e02a69466afa0ee014174143, got", details.Job.InputMediaFile.MD5Checksum)
	}

	if details.Job.OutputMediaFiles[0].Format != "mpeg4" {
		t.Fatal("Expected mpeg4, got", details.Job.OutputMediaFiles[0].Format)
	}
	if details.Job.OutputMediaFiles[0].CreatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.OutputMediaFiles[0].CreatedAt)
	}
	if details.Job.OutputMediaFiles[0].FrameRate != 29 {
		t.Fatal("Expected 29, got", details.Job.OutputMediaFiles[0].FrameRate)
	}
	if details.Job.OutputMediaFiles[0].FinishedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.OutputMediaFiles[0].FinishedAt)
	}
	if details.Job.OutputMediaFiles[0].UpdatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.OutputMediaFiles[0].UpdatedAt)
	}
	if details.Job.OutputMediaFiles[0].DurationInMs != 24883 {
		t.Fatal("Expected 24883, got", details.Job.OutputMediaFiles[0].DurationInMs)
	}
	if details.Job.OutputMediaFiles[0].AudioSampleRate != 44100 {
		t.Fatal("Expected 44100, got", details.Job.OutputMediaFiles[0].AudioSampleRate)
	}
	if details.Job.OutputMediaFiles[0].Url != "http://s3.amazonaws.com/bucket/video.mp4" {
		t.Fatal("Expected http://s3.amazonaws.com/bucket/video.mp4, got", details.Job.OutputMediaFiles[0].Url)
	}
	if details.Job.OutputMediaFiles[0].Id != 1 {
		t.Fatal("Expected 1, got", details.Job.OutputMediaFiles[0].Id)
	}
	if details.Job.OutputMediaFiles[0].ErrorMessage != nil {
		t.Fatal("Expected nil, got", details.Job.OutputMediaFiles[0].ErrorMessage)
	}
	if details.Job.OutputMediaFiles[0].ErrorClass != nil {
		t.Fatal("Expected nil, got", details.Job.OutputMediaFiles[0].ErrorClass)
	}
	if details.Job.OutputMediaFiles[0].AudioBitrateInKbps != 92 {
		t.Fatal("Expected 92, got", details.Job.OutputMediaFiles[0].AudioBitrateInKbps)
	}
	if details.Job.OutputMediaFiles[0].AudioCodec != "aac" {
		t.Fatal("Expected aac, got", details.Job.OutputMediaFiles[0].AudioCodec)
	}
	if details.Job.OutputMediaFiles[0].Height != 352 {
		t.Fatal("Expected 352, got", details.Job.OutputMediaFiles[0].Height)
	}
	if details.Job.OutputMediaFiles[0].FileSizeInBytes != 1386663 {
		t.Fatal("Expected 1386663, got", details.Job.OutputMediaFiles[0].FileSizeInBytes)
	}
	if details.Job.OutputMediaFiles[0].VideoCodec != "h264" {
		t.Fatal("Expected h264, got", details.Job.OutputMediaFiles[0].VideoCodec)
	}
	if details.Job.OutputMediaFiles[0].Test != false {
		t.Fatal("Expected false, got", details.Job.OutputMediaFiles[0].Test)
	}
	if details.Job.OutputMediaFiles[0].TotalBitrateInKbps != 443 {
		t.Fatal("Expected 443, got", details.Job.OutputMediaFiles[0].TotalBitrateInKbps)
	}
	if details.Job.OutputMediaFiles[0].Channels != "2" {
		t.Fatal("Expected 2, got", details.Job.OutputMediaFiles[0].Channels)
	}
	if details.Job.OutputMediaFiles[0].Width != 624 {
		t.Fatal("Expected 624, got", details.Job.OutputMediaFiles[0].Width)
	}
	if details.Job.OutputMediaFiles[0].VideoBitrateInKbps != 351 {
		t.Fatal("Expected 351, got", details.Job.OutputMediaFiles[0].VideoBitrateInKbps)
	}
	if details.Job.OutputMediaFiles[0].State != "finished" {
		t.Fatal("Expected finished, got", details.Job.OutputMediaFiles[0].State)
	}
	if details.Job.OutputMediaFiles[0].Label == nil || *details.Job.OutputMediaFiles[0].Label != "Web" {
		t.Fatal("Expected Web, got", details.Job.OutputMediaFiles[0].Label)
	}
	if details.Job.OutputMediaFiles[0].MD5Checksum != "7f106918e02a69466afa0ee014172496" {
		t.Fatal("Expected 7f106918e02a69466afa0ee014172496, got", details.Job.OutputMediaFiles[0].MD5Checksum)
	}
	if details.Job.Thumbnails[0].CreatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.Thumbnails[0].CreatedAt)
	}
	if details.Job.Thumbnails[0].UpdatedAt != "2010-01-01T00:00:00Z" {
		t.Fatal("Expected 2010-01-01T00:00:00Z, got", details.Job.Thumbnails[0].UpdatedAt)
	}
	if details.Job.Thumbnails[0].Url != "http://s3.amazonaws.com/bucket/video/frame_0000.png" {
		t.Fatal("Expected http://s3.amazonaws.com/bucket/video/frame_0000.png, got", details.Job.Thumbnails[0].Url)
	}
	if details.Job.Thumbnails[0].Id != 1 {
		t.Fatal("Expected 1, got", details.Job.Thumbnails[0].Id)
	}

	expectedStatus = http.StatusNotFound
	details, err = zc.GetJobDetails(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no response")
	}

	expectedStatus = http.StatusOK
	returnBody = false
	details, err = zc.GetJobDetails(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no response")
	}

	returnBody = false
	srv.Close()
	details, err = zc.GetJobDetails(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no response")
	}
}

func TestGetJobProgress(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs/123/progress.json", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "state": "processing",
  "progress": 32.34567345,
  "input": {
    "id": 1234,
    "state": "finished"
  },
  "outputs": [
    {
      "id": 4567,
      "state": "processing",
      "current_event": "Transcoding",
      "current_event_progress": 25.0323,
      "progress": 35.23532
    },
    {
      "id": 4568,
      "state": "processing",
      "current_event": "Uploading",
      "current_event_progress": 82.32,
      "progress": 95.3223
    }
  ]
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	progress, err := zc.GetJobProgress(123)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if progress == nil {
		t.Fatal("Expected response")
	}

	if progress.State != "processing" {
		t.Fatal("Expected processing, got", progress.State)
	}
	if progress.JobProgress != 32.34567345 {
		t.Fatal("Expected 32.34567345, got", progress.JobProgress)
	}
	if progress.InputProgress.Id != 1234 {
		t.Fatal("Expected 1234, got", progress.InputProgress.Id)
	}
	if progress.InputProgress.State != "finished" {
		t.Fatal("Expected finished, got", progress.InputProgress.State)
	}

	if len(progress.OutputProgress) != 2 {
		t.Fatal("Expected 2 outputs, got", len(progress.OutputProgress))
	}

	if progress.OutputProgress[0].Id != 4567 {
		t.Fatal("Expected 4567, got", progress.OutputProgress[0].Id)
	}
	if progress.OutputProgress[0].State != "processing" {
		t.Fatal("Expected processing, got", progress.OutputProgress[0].State)
	}
	if progress.OutputProgress[0].CurrentEvent != "Transcoding" {
		t.Fatal("Expected Transcoding, got", progress.OutputProgress[0].CurrentEvent)
	}
	if progress.OutputProgress[0].CurrentEventProgress != 25.0323 {
		t.Fatal("Expected 25.0323, got", progress.OutputProgress[0].CurrentEventProgress)
	}
	if progress.OutputProgress[0].OverallProgress != 35.23532 {
		t.Fatal("Expected 35.23532, got", progress.OutputProgress[0].OverallProgress)
	}

	expectedStatus = http.StatusNotFound
	progress, err = zc.GetJobProgress(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if progress != nil {
		t.Fatal("Expected no response", progress)
	}

	expectedStatus = http.StatusOK
	returnBody = false
	progress, err = zc.GetJobProgress(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if progress != nil {
		t.Fatal("Expected no response", progress)
	}

	returnBody = true
	srv.Close()
	progress, err = zc.GetJobProgress(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if progress != nil {
		t.Fatal("Expected no response", progress)
	}
}

func TestResubmitJob(t *testing.T) {
	expectedStatus := http.StatusNoContent

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs/123/resubmit.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatus)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	err := zc.ResubmitJob(123)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	expectedStatus = http.StatusConflict

	err = zc.ResubmitJob(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	expectedStatus = http.StatusOK
	srv.Close()

	err = zc.ResubmitJob(123)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCancelJob(t *testing.T) {
	expectedStatus := http.StatusNoContent

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs/123/cancel.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatus)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	err := zc.CancelJob(123)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	expectedStatus = http.StatusConflict

	err = zc.CancelJob(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	expectedStatus = http.StatusOK
	srv.Close()

	err = zc.CancelJob(123)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestFinishLiveJob(t *testing.T) {
	expectedStatus := http.StatusNoContent

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs/123/finish", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatus)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	err := zc.FinishLiveJob(123)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	expectedStatus = http.StatusConflict

	err = zc.FinishLiveJob(123)
	if err == nil {
		t.Fatal("Expected error")
	}
}
