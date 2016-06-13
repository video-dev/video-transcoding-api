package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	gizmoConfig "github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/db/dbtest"
)

func TestNewPreset(t *testing.T) {
	tests := []struct {
		givenTestCase    string
		givenRequestData map[string]interface{}
		wantBody         map[string]interface{}
		wantCode         int
	}{
		{
			"Create new preset",
			map[string]interface{}{
				"providers": []string{"fake", "encodingcom"},
				"preset": map[string]interface{}{
					"name":         "nyt_test_here_2wq",
					"description":  "testing creation from api",
					"container":    "mp4",
					"profile":      "Main",
					"profileLevel": "3.1",
					"rateControl":  "VBR",
					"video": map[string]string{
						"height":        "720",
						"codec":         "h264",
						"bitrate":       "1000",
						"gopSize":       "90",
						"gopMode":       "fixed",
						"interlaceMode": "progressive",
					},
					"audio": map[string]string{
						"codec":   "aac",
						"bitrate": "64000",
					},
				},
			},
			map[string]interface{}{
				"Results": map[string]interface{}{
					"fake": map[string]interface{}{
						"PresetID": "presetID_here",
						"Error":    "",
					},
					"encodingcom": map[string]interface{}{
						"PresetID": "",
						"Error":    "getting factory: provider not found",
					},
				},
				"PresetMap": "nyt_test_here_2wq",
			},
			http.StatusOK,
		},
		{
			"Error creating preset in all providers",
			map[string]interface{}{
				"providers": []string{"elastictranscoder", "encodingcom"},
				"preset": map[string]interface{}{
					"name":         "nyt_test_here_3wq",
					"description":  "testing creation from api",
					"container":    "mp4",
					"profile":      "Main",
					"profileLevel": "3.1",
					"rateControl":  "VBR",
					"video": map[string]string{
						"height":        "720",
						"codec":         "h264",
						"bitrate":       "1000",
						"gopSize":       "90",
						"gopMode":       "fixed",
						"interlaceMode": "progressive",
					},
					"audio": map[string]string{
						"codec":   "aac",
						"bitrate": "64000",
					},
				},
			},
			map[string]interface{}{
				"Results": map[string]interface{}{
					"elastictranscoder": map[string]interface{}{
						"PresetID": "",
						"Error":    "getting factory: provider not found",
					},
					"encodingcom": map[string]interface{}{
						"PresetID": "",
						"Error":    "getting factory: provider not found",
					},
				},
				"PresetMap": "",
			},
			http.StatusOK,
		},
	}

	for _, test := range tests {
		srvr := server.NewSimpleServer(&gizmoConfig.Server{RouterType: "fast"})
		fakeDB := dbtest.NewFakeRepository(false)

		srvr.Register(&TranscodingService{config: &config.Config{}, db: fakeDB})
		body, _ := json.Marshal(test.givenRequestData)
		r, _ := http.NewRequest("POST", "/presets", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		var got map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Errorf("%s: unable to JSON decode response body: %s", test.givenTestCase, err)
		}
		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("%s: expected response body of\n%#v;\ngot\n%#v", test.givenTestCase, test.wantBody, got)
		}
	}
}

func TestDeletePreset(t *testing.T) {
	tests := []struct {
		givenTestCase string
		wantBody      map[string]interface{}
		wantCode      int
	}{
		{
			"Delete a preset",
			map[string]interface{}{
				"results": map[string]interface{}{
					"fake": map[string]interface{}{
						"presetId": "presetID_here",
					},
				},
				"presetMap": "removed successfully",
			},
			http.StatusOK,
		},
	}

	for _, test := range tests {
		srvr := server.NewSimpleServer(&gizmoConfig.Server{RouterType: "fast"})
		fakeDB := dbtest.NewFakeRepository(false)
		fakeProviderMapping := make(map[string]string)
		fakeProviderMapping["fake"] = "presetID_here"
		fakeDB.CreatePresetMap(&db.PresetMap{Name: "abc-321", ProviderMapping: fakeProviderMapping})
		srvr.Register(&TranscodingService{config: &config.Config{}, db: fakeDB})
		r, _ := http.NewRequest("DELETE", "/presets/abc-321", nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		var got map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Errorf("%s: unable to JSON decode response body: %s", test.givenTestCase, err)
		}
		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("%s: expected response body of\n%#v;\ngot\n%#v", test.givenTestCase, test.wantBody, got)
		}
	}
}
