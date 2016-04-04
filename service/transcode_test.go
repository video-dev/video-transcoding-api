package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	gizmoConfig "github.com/NYTimes/gizmo/config"
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
		wantBody interface{}
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
			map[string]interface{}{"jobId": "12345"},
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
		srvr := server.NewSimpleServer(&gizmoConfig.Server{RouterType: "fast"})
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
			"/jobs/12345",
			false,
			"hls",
			5,
			http.StatusOK,
			map[string]interface{}{
				"providerJobId": "provider-job-123",
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
		srvr := server.NewSimpleServer(&gizmoConfig.Server{RouterType: "fast"})
		fakeDBObj := dbtest.NewFakeRepository(test.givenTriggerDBError)
		fakeDBObj.CreateJob(&db.Job{ProviderName: "fake",
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
