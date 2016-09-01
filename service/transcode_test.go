package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/db/dbtest"
	"github.com/nytm/video-transcoding-api/provider"
)

func TestTranscode(t *testing.T) {
	tests := []struct {
		givenTestCase       string
		givenRequestBody    string
		givenTriggerDBError bool

		wantCode int
		wantBody map[string]interface{}
	}{
		{
			"New job",
			`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "presets": ["mp4_1080p"],
  "provider": "fake"
}`,
			false,

			http.StatusOK,
			map[string]interface{}{"jobId": ""},
		},
		{
			"New job with preset not found in provider",
			`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "presets": ["mp4_360p"],
  "provider": "fake"
}`,
			false,

			http.StatusBadRequest,
			map[string]interface{}{"error": provider.ErrPresetMapNotFound.Error()},
		},
		{
			"New job with preset not found in the API",
			`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "presets": ["mp4_720p"],
  "provider": "fake"
}`,
			false,

			http.StatusBadRequest,
			map[string]interface{}{"error": db.ErrPresetMapNotFound.Error()},
		},
		{
			"New job with database error",
			`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "presets": ["mp4_1080p"],
  "provider": "fake"
}`,
			true,

			http.StatusInternalServerError,
			map[string]interface{}{
				"error": "database error",
			},
		},
		{
			"New job with invalid provider",
			`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "presets": ["mp4_1080p"],
  "provider": "nonexistent-provider"
}`,
			false,

			http.StatusBadRequest,
			map[string]interface{}{
				"error": "provider not found",
			},
		},
		{
			"New job missing presets",
			`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "provider": "fake"
}`,
			false,

			http.StatusBadRequest,
			map[string]interface{}{
				"error": "missing preset list from request",
			},
		},
	}

	for _, test := range tests {
		srvr := server.NewSimpleServer(&server.Config{RouterType: "fast"})
		fakeDBObj := dbtest.NewFakeRepository(test.givenTriggerDBError)
		fakeDBObj.CreatePresetMap(&db.PresetMap{Name: "mp4_1080p", ProviderMapping: map[string]string{"fake": "18828"}})
		fakeDBObj.CreatePresetMap(&db.PresetMap{Name: "mp4_360p", ProviderMapping: map[string]string{"elementalconductor": "172712"}})
		srvr.Register(&TranscodingService{config: &config.Config{}, db: fakeDBObj})
		r, _ := http.NewRequest("POST", "/jobs", strings.NewReader(test.givenRequestBody))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: expected response code of %d; got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		var got map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Errorf("%s: unable to JSON decode response body: %s", test.givenTestCase, err)
		}
		if got["jobId"] == "" {
			t.Errorf("%s: missing jobId from the response: %#v", test.givenTestCase, got)
		}
		if _, ok := test.wantBody["jobId"]; ok {
			test.wantBody["jobId"] = got["jobId"]
		}
		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("%s: expected response body of\n%#v;\ngot\n%#v", test.givenTestCase, test.wantBody, got)
		}
		if test.wantCode == http.StatusOK {
			_, err = fakeDBObj.GetJob(got["jobId"].(string))
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestGetTranscodeJob(t *testing.T) {
	tests := []struct {
		givenTestCase        string
		givenURI             string
		givenTriggerDBError  bool
		givenProtocol        string
		givenSegmentDuration uint

		wantCode int
		wantBody interface{}
	}{
		{
			"Get job",
			"/jobs/job-123",
			false,
			"hls",
			5,
			http.StatusOK,
			map[string]interface{}{
				"providerJobId": "provider-job-123",
				"status":        "finished",
				"providerName":  "fake",
				"statusMessage": "The job is finished",
				"progress":      10.3,
				"providerStatus": map[string]interface{}{
					"progress":   10.3,
					"sourcefile": "http://some.source.file",
				},
				"mediaInfo": map[string]interface{}{
					"width":      float64(4096),
					"height":     float64(2160),
					"duration":   183e9,
					"videoCodec": "VP9",
				},
			},
		},
		{
			"Get job with inexistent job id",
			"/jobs/non_existent_job",
			false,
			"",
			0,
			http.StatusNotFound,
			map[string]interface{}{"error": "job not found"},
		},
	}

	for _, test := range tests {
		srvr := server.NewSimpleServer(&server.Config{RouterType: "fast"})
		fakeDBObj := dbtest.NewFakeRepository(test.givenTriggerDBError)
		fakeDBObj.CreateJob(&db.Job{
			ID:            "job-123",
			ProviderName:  "fake",
			ProviderJobID: "provider-job-123",
			StreamingParams: db.StreamingParams{
				SegmentDuration: test.givenSegmentDuration,
				Protocol:        test.givenProtocol,
			}})
		srvr.Register(&TranscodingService{config: &config.Config{}, db: fakeDBObj})
		r, _ := http.NewRequest("GET", test.givenURI, nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: expected response code of %d; got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		var got interface{}
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Errorf("%s: unable to JSON decode response body: %s", test.givenTestCase, err)
		}
		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("%s: expected response body of\n%#v;\ngot\n%#v", test.givenTestCase, test.wantBody, got)
		}
	}
}

func TestCancelTranscodeJob(t *testing.T) {
	var tests = []struct {
		givenTestCase       string
		givenJobID          string
		givenTriggerDBError bool

		wantCode int
		wantBody map[string]interface{}
	}{
		{
			"valid job",
			"job-123",
			false,

			http.StatusOK,
			map[string]interface{}{
				"providerJobId": "provider-job-123",
				"status":        "canceled",
				"providerName":  "fake",
				"statusMessage": "The job is finished",
				"progress":      10.3,
				"providerStatus": map[string]interface{}{
					"progress":   10.3,
					"sourcefile": "http://some.source.file",
				},
				"mediaInfo": map[string]interface{}{
					"width":      float64(4096),
					"height":     float64(2160),
					"duration":   183e9,
					"videoCodec": "VP9",
				},
			},
		},
		{
			"job that doesn't exist in the provider",
			"job-1234",
			false,

			http.StatusGone,
			map[string]interface{}{"error": "could not found job with id: some-job"},
		},
		{
			"non-existing job",
			"some-id",
			false,

			http.StatusNotFound,
			map[string]interface{}{"error": db.ErrJobNotFound.Error()},
		},
		{
			"db error",
			"job-123",
			true,

			http.StatusInternalServerError,
			map[string]interface{}{"error": `error retrieving job with id "job-123": database error`},
		},
	}
	defer func() { fprovider.canceledJobs = nil }()
	for _, test := range tests {
		fprovider.canceledJobs = nil
		srvr := server.NewSimpleServer(&server.Config{RouterType: "fast"})
		fakeDBObj := dbtest.NewFakeRepository(test.givenTriggerDBError)
		fakeDBObj.CreateJob(&db.Job{
			ID:            "job-123",
			ProviderName:  "fake",
			ProviderJobID: "provider-job-123",
		})
		fakeDBObj.CreateJob(&db.Job{
			ID:            "job-1234",
			ProviderName:  "fake",
			ProviderJobID: "some-job",
		})
		srvr.Register(&TranscodingService{config: &config.Config{}, db: fakeDBObj})
		r, _ := http.NewRequest("POST", "/jobs/"+test.givenJobID+"/cancel", bytes.NewReader(nil))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong code returned. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		var body map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &body)
		if err != nil {
			t.Fatalf("%s: %s", test.givenTestCase, err)
		}
		if !reflect.DeepEqual(body, test.wantBody) {
			t.Errorf("%s: wrong body returned.\nWant %#v\nGot  %#v", test.givenTestCase, test.wantBody, body)
		}
		if test.wantCode == http.StatusOK {
			if len(fprovider.canceledJobs) < 1 {
				t.Errorf("%s: did not cancel the job in the provider", test.givenTestCase)
			} else if fprovider.canceledJobs[0] != "provider-job-123" {
				t.Errorf("%s: did not send the correct job id to the provider. Want %q. Got %q", test.givenTestCase, "provider-job-123", fprovider.canceledJobs[0])
			}
		}
	}
}
