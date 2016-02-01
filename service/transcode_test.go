package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
	"github.com/rcrowley/go-metrics"
)

const testProfileString = `{
   "output":["webm"],
   "size":{"height":360},
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

type fakeProvider struct{}

func (e *fakeProvider) Transcode(sourceMedia string, profiles []provider.Profile) (*provider.JobStatus, error) {
	return &provider.JobStatus{
		ProviderJobID: "provider-job-123",
		Status:        provider.StatusFinished,
		StatusMessage: "The job is finished",
		ProviderStatus: map[string]interface{}{
			"progress":   100.0,
			"sourcefile": "http://some.source.file",
		},
	}, nil
}

func (e *fakeProvider) JobStatus(id string) (*provider.JobStatus, error) {
	if id == "provider-job-123" {
		return &provider.JobStatus{
			ProviderJobID: "provider-job-123",
			Status:        provider.StatusFinished,
			StatusMessage: "The job is finished",
			ProviderStatus: map[string]interface{}{
				"progress":   100.0,
				"sourcefile": "http://some.source.file",
			},
		}, nil
	}
	return nil, provider.JobNotFoundError{ID: id}
}

func fakeProviderFactory(cfg *config.Config) (provider.TranscodingProvider, error) {
	return &fakeProvider{}, nil
}

type fakeDB struct {
	TriggerDBError bool
}

func (d *fakeDB) SaveJob(job *db.Job) error {
	if d.TriggerDBError {
		return fmt.Errorf("Database error")
	}
	job.ID = "12345"
	return nil
}

func (d *fakeDB) DeleteJob(job *db.Job) error {
	return nil
}

func (d *fakeDB) GetJob(id string) (*db.Job, error) {
	if id == "12345" {
		return &db.Job{
			ID:            "12345",
			ProviderName:  "fake",
			ProviderJobID: "provider-job-123",
		}, nil
	}
	return nil, db.ErrJobNotFound
}

func TestTranscode(t *testing.T) {

	tests := []struct {
		givenTestCase       string
		givenURI            string
		givenHTTPMethod     string
		givenRequestBody    string
		givenTriggerDBError bool

		wantCode int
		wantBody interface{}
	}{
		{
			"New job",
			"/jobs",
			"POST",
			fmt.Sprintf(`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "profiles": [%s],
  "provider": "fake"
}`, testProfileString),
			false,

			http.StatusOK,
			map[string]interface{}{
				"jobId": "12345",
			},
		},
		{
			"New job with database error",
			"/jobs",
			"POST",
			fmt.Sprintf(`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "profiles": [%s],
  "provider": "fake"
}`, testProfileString),
			true,

			http.StatusInternalServerError,
			map[string]interface{}{
				"error": "Database error",
			},
		},
		{
			"New job with invalid provider",
			"/jobs",
			"POST",
			fmt.Sprintf(`{
  "source": "http://another.non.existent/video.mp4",
  "destination": "s3://some.bucket.s3.amazonaws.com/some_path",
  "profiles": [%s],
  "provider": "nonexistent-provider"
}`, testProfileString),
			false,

			http.StatusBadRequest,
			map[string]interface{}{
				"error": "Unknown provider found in request: nonexistent-provider",
			},
		},

		{
			"Get job",
			"/jobs/12345",
			"GET",
			"",
			false,

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
			"GET",
			"",
			false,

			http.StatusNotFound,
			map[string]interface{}{
				"error": "Error retrieving job with id 'non_existent_job': job not found",
			},
		},
	}

	for _, test := range tests {

		srvr := server.NewSimpleServer(nil)
		fakeDBObj := db.JobRepository(&fakeDB{
			TriggerDBError: test.givenTriggerDBError,
		})
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDBObj,
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

		// ** THIS IS REQUIRED in order to run the test multiple times.
		metrics.DefaultRegistry.UnregisterAll()
	}

}
