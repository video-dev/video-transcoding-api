package zencoder

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInputDetails(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/inputs/123.json", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "audio_bitrate_in_kbps": 96,
  "audio_codec": "aac",
  "audio_sample_rate": 48000,
  "channels": "2",
  "duration_in_ms": 24892,
  "file_size_in_bytes": 1862748,
  "format": "mpeg4",
  "frame_rate": 29.98,
  "height": 352,
  "id": 6816,
  "job_id": 6816,
  "privacy": false,
  "state": "finished",
  "total_bitrate_in_kbps": 594,
  "url": "s3://example/file.mp4",
  "video_bitrate_in_kbps": 498,
  "video_codec": "h264",
  "width": 624,
  "md5_checksum": "7f106918e02a69466afa0ee014174143"
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	details, err := zc.GetInputDetails(123)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if details == nil {
		t.Fatal("Expected details")
	}

	if details.AudioBitrateInKbps != 96 {
		t.Fatal("Expected 96, got", details.AudioBitrateInKbps)
	}
	if details.AudioCodec != "aac" {
		t.Fatal("Expected aac, got", details.AudioCodec)
	}
	if details.AudioSampleRate != 48000 {
		t.Fatal("Expected 48000, got", details.AudioSampleRate)
	}
	if details.Channels != "2" {
		t.Fatal("Expected 2, got", details.Channels)
	}
	if details.DurationInMs != 24892 {
		t.Fatal("Expected 24892, got", details.DurationInMs)
	}
	if details.FileSizeInBytes != 1862748 {
		t.Fatal("Expected 1862748, got", details.FileSizeInBytes)
	}
	if details.Format != "mpeg4" {
		t.Fatal("Expected mpeg4, got", details.Format)
	}
	if details.FrameRate != 29.98 {
		t.Fatal("Expected 29.98, got", details.FrameRate)
	}
	if details.Height != 352 {
		t.Fatal("Expected 352, got", details.Height)
	}
	if details.Id != 6816 {
		t.Fatal("Expected 6816, got", details.Id)
	}
	if details.JobId != 6816 {
		t.Fatal("Expected 6816, got", details.JobId)
	}
	if details.Privacy != false {
		t.Fatal("Expected false, got", details.Privacy)
	}
	if details.State != "finished" {
		t.Fatal("Expected finished, got", details.State)
	}
	if details.TotalBitrateInKbps != 594 {
		t.Fatal("Expected 594, got", details.TotalBitrateInKbps)
	}
	if details.Url != "s3://example/file.mp4" {
		t.Fatal("Expected s3://example/file.mp4, got", details.Url)
	}
	if details.VideoBitrateInKbps != 498 {
		t.Fatal("Expected 498, got", details.VideoBitrateInKbps)
	}
	if details.VideoCodec != "h264" {
		t.Fatal("Expected h264, got", details.VideoCodec)
	}
	if details.Width != 624 {
		t.Fatal("Expected 624, got", details.Width)
	}
	if details.MD5Checksum != "7f106918e02a69466afa0ee014174143" {
		t.Fatal("Expected 7f106918e02a69466afa0ee014174143, got", details.MD5Checksum)
	}

	expectedStatus = http.StatusInternalServerError

	details, err = zc.GetInputDetails(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	expectedStatus = http.StatusOK
	returnBody = false

	details, err = zc.GetInputDetails(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	srv.Close()
	returnBody = true

	details, err = zc.GetInputDetails(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}
}

func TestGetInputProgress(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/inputs/123/progress.json", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "state": "processing",
  "current_event": "Downloading",
  "current_event_progress": 32.34567345,
  "progress": 45.2353255
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	progress, err := zc.GetInputProgress(123)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if progress == nil {
		t.Fatal("Expected progress")
	}

	if progress.State != "processing" {
		t.Fatal("Expected processing, got", progress.State)
	}

	if progress.CurrentEvent != "Downloading" {
		t.Fatal("Expected Downloading, got", progress.CurrentEvent)
	}

	if progress.CurrentEventProgress != 32.34567345 {
		t.Fatal("Expected 32.34567345, got", progress.CurrentEventProgress)
	}

	if progress.OverallProgress != 45.2353255 {
		t.Fatal("Expected 45.2353255, got", progress.OverallProgress)
	}

	expectedStatus = http.StatusInternalServerError

	progress, err = zc.GetInputProgress(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if progress != nil {
		t.Fatal("Expected no progress", progress)
	}

	expectedStatus = http.StatusOK
	returnBody = false

	progress, err = zc.GetInputProgress(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if progress != nil {
		t.Fatal("Expected no progress", progress)
	}

	srv.Close()
	returnBody = true

	progress, err = zc.GetInputProgress(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if progress != nil {
		t.Fatal("Expected no progress", progress)
	}
}
