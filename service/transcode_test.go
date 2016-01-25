package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
	"github.com/rcrowley/go-metrics"
)

type fakeProvider struct{}

const testProfileString = `{  
   "output":"webm",
   "size":"0x360",
   "bitrate":"900k",
   "audio_bitrate":"64k",
   "audio_sample_rate":"48000",
   "audio_channels_number":"2",
   "framerate":"30",
   "keep_aspect_ratio":"yes",
   "video_codec":"libvpx",
   "profile":"main",
   "audio_codec":"libvorbis",
   "two_pass":"yes",
   "turbo":"no",
   "twin_turbo":"no",
   "cbr":"no",
   "deinterlacing":"auto",
   "keyframe":"90",
   "audio_volume":"100",
   "rotate":0,
   "strip_chapters":"no",
   "hint":"no"
}`

func (e *fakeProvider) Transcode(sourceMedia string, destination string, profile provider.Profile) (*provider.JobStatus, error) {
	return &provider.JobStatus{
		ProviderJobID: "12345",
		Status:        provider.StatusFinished,
		StatusMessage: "The job is finished",
		ProviderStatus: map[string]interface{}{
			"progress":   100.0,
			"sourcefile": "http://some.source.file",
		},
	}, nil
}

func (e *fakeProvider) JobStatus(id string) (*provider.JobStatus, error) {
	return nil, nil
}

func fakeProviderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	return &fakeProvider{}, nil
}

func TestTranscode(t *testing.T) {

	tests := []struct {
		givenURI         string
		givenHTTPMethod  string
		givenRequestBody string

		wantCode int
		wantBody interface{}
	}{
		{
			"/jobs",
			"POST",
			fmt.Sprintf(`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "profile": %s,
  "provider": "fake"
}`, strconv.Quote(testProfileString)),

			http.StatusOK,
			map[string]interface{}{
				"providerJobId": "12345",
				"status":        "finished",
				"providerName":  "fake",
				"statusMessage": "The job is finished",
				"providerStatus": map[string]interface{}{
					"progress":   100.0,
					"sourcefile": "http://some.source.file",
				},
			},
		},
		{
			"/jobs",
			"POST",
			fmt.Sprintf(`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "profile": %s,
  "provider": "nonexistent-provider"
}`, strconv.Quote(testProfileString)),

			http.StatusBadRequest,
			map[string]interface{}{
				"error": "Unknown provider found in request: nonexistent-provider",
			},
		},
	}

	for _, test := range tests {

		srvr := server.NewSimpleServer(nil)
		srvr.Register(&TranscodingService{
			providers: map[string]provider.Factory{
				"fake": fakeProviderFactory,
			},
		})

		r, _ := http.NewRequest(
			test.givenHTTPMethod,
			test.givenURI,
			strings.NewReader(test.givenRequestBody),
		)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)

		if w.Code != test.wantCode {
			t.Errorf("expected response code of %d; got %d", test.wantCode, w.Code)
		}

		var got interface{}
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Error("unable to JSON decode response body: ", err)
		}

		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("expected response body of\n%#v;\ngot\n%#v", test.wantBody, got)
		}

		// ** THIS IS REQUIRED in order to run the test multiple times.
		metrics.DefaultRegistry.UnregisterAll()
	}

}
