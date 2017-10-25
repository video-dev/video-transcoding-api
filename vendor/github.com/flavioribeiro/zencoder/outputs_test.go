package zencoder

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOutputDetails(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/outputs/123.json", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "audio_bitrate_in_kbps": 74,
  "audio_codec": "aac",
  "audio_sample_rate": 48000,
  "channels": "2",
  "duration_in_ms": 24892,
  "file_size_in_bytes": 1215110,
  "format": "mpeg4",
  "frame_rate": 29.97,
  "height": 352,
  "id": 13339,
  "label": null,
  "state": "finished",
  "total_bitrate_in_kbps": 387,
  "url": "https://example.com/file.mp4",
  "video_bitrate_in_kbps": 313,
  "video_codec": "h264",
  "width": 624,
  "md5_checksum": "7f106918e02a69466afa0ee014174143"
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	details, err := zc.GetOutputDetails(123)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if details == nil {
		t.Fatal("Expected details")
	}

	if details.AudioBitrateInKbps != 74 {
		t.Fatal("Expected 74, got", details.AudioBitrateInKbps)
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
	if details.FileSizeInBytes != 1215110 {
		t.Fatal("Expected 1215110, got", details.FileSizeInBytes)
	}
	if details.Format != "mpeg4" {
		t.Fatal("Expected mpeg4, got", details.Format)
	}
	if details.FrameRate != 29.97 {
		t.Fatal("Expected 29.97, got", details.FrameRate)
	}
	if details.Height != 352 {
		t.Fatal("Expected 352, got", details.Height)
	}
	if details.Id != 13339 {
		t.Fatal("Expected 13339, got", details.Id)
	}
	if details.Label != nil {
		t.Fatal("Expected nil, got", details.Label)
	}
	if details.State != "finished" {
		t.Fatal("Expected finished, got", details.State)
	}
	if details.TotalBitrateInKbps != 387 {
		t.Fatal("Expected 387, got", details.TotalBitrateInKbps)
	}
	if details.Url != "https://example.com/file.mp4" {
		t.Fatal("Expected https://example.com/file.mp4, got", details.Url)
	}
	if details.VideoBitrateInKbps != 313 {
		t.Fatal("Expected 313, got", details.VideoBitrateInKbps)
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

	details, err = zc.GetOutputDetails(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	expectedStatus = http.StatusOK
	returnBody = false

	details, err = zc.GetOutputDetails(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}

	srv.Close()
	returnBody = true

	details, err = zc.GetOutputDetails(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if details != nil {
		t.Fatal("Expected no details", details)
	}
}

func TestGetOutputProgress(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/outputs/123/progress.json", func(w http.ResponseWriter, r *http.Request) {
		if expectedStatus != http.StatusOK {
			w.WriteHeader(expectedStatus)
			return
		}

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "state": "processing",
  "current_event": "Transcoding",
  "current_event_progress": 45.32525,
  "progress": 32.34567345
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	progress, err := zc.GetOutputProgress(123)
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if progress == nil {
		t.Fatal("Expected progress")
	}

	if progress.State != "processing" {
		t.Fatal("Expected processing, got", progress.State)
	}

	if progress.CurrentEvent != "Transcoding" {
		t.Fatal("Expected Transcoding, got", progress.CurrentEvent)
	}

	if progress.CurrentEventProgress != 45.32525 {
		t.Fatal("Expected 45.32525, got", progress.CurrentEventProgress)
	}

	if progress.OverallProgress != 32.34567345 {
		t.Fatal("Expected 32.34567345, got", progress.OverallProgress)
	}

	expectedStatus = http.StatusInternalServerError

	progress, err = zc.GetOutputProgress(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if progress != nil {
		t.Fatal("Expected no progress", progress)
	}

	expectedStatus = http.StatusOK
	returnBody = false

	progress, err = zc.GetOutputProgress(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if progress != nil {
		t.Fatal("Expected no progress", progress)
	}

	srv.Close()
	returnBody = true

	progress, err = zc.GetOutputProgress(123)
	if err == nil {
		t.Fatal("Expected error")
	}

	if progress != nil {
		t.Fatal("Expected no progress", progress)
	}
}
