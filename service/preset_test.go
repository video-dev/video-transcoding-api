package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/dbtest"
	"github.com/sirupsen/logrus"
)

func TestNewPreset(t *testing.T) {
	tests := []struct {
		givenTestCase    string
		givenRequestData map[string]interface{}
		wantOutputOpts   db.OutputOptions
		wantBody         map[string]interface{}
		wantCode         int
	}{
		{
			"Create new preset",
			map[string]interface{}{
				"providers": []string{"fake", "encodingcom"},
				"outputOptions": map[string]interface{}{
					"extension": "mp5",
				},
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
			db.OutputOptions{
				Extension: "mp4",
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
				"outputOptions": map[string]interface{}{
					"extension": "mp5",
				},
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
			db.OutputOptions{},
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
			http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		name := test.givenRequestData["preset"].(map[string]interface{})["name"].(string)
		srvr := server.NewSimpleServer(&server.Config{RouterType: "fast"})
		fakeDB := dbtest.NewFakeRepository(false)
		service, err := NewTranscodingService(&config.Config{Server: &server.Config{}}, logrus.New())
		if err != nil {
			t.Fatal(err)
		}
		service.db = fakeDB
		srvr.Register(service)
		body, _ := json.Marshal(test.givenRequestData)
		r, _ := http.NewRequest("POST", "/presets", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		var got map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Errorf("%s: unable to JSON decode response body: %s", test.givenTestCase, err)
		}
		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("%s: expected response body of\n%#v;\ngot\n%#v", test.givenTestCase, test.wantBody, got)
		}
		if test.wantCode == http.StatusOK {
			presetMap, err := fakeDB.GetPresetMap(name)
			if err != nil {
				t.Fatalf("%s: %s", test.givenTestCase, err)
			}
			if !reflect.DeepEqual(presetMap.OutputOpts, test.wantOutputOpts) {
				t.Errorf("%s: wrong output options saved.\nWant %#v\nGot  %#v", test.givenTestCase, test.wantOutputOpts, presetMap.OutputOpts)
			}
		}
	}
}

func TestNewPresetWithExistentPresetMap(t *testing.T) {
	data := map[string]interface{}{
		"providers":     []string{"zencoder"},
		"outputOptions": map[string]interface{}{},
		"preset": map[string]interface{}{
			"name":         "presetID_here",
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
	}

	presetMap := db.PresetMap{
		Name: "presetID_here",
		ProviderMapping: map[string]string{
			"fake": "presetID_here",
		},
		OutputOpts: db.OutputOptions{},
	}

	fakeDB := dbtest.NewFakeRepository(false)
	err := fakeDB.CreatePresetMap(&presetMap)
	if err != nil {
		t.Fatal(err)
	}

	srvr := server.NewSimpleServer(&server.Config{RouterType: "fast"})
	service, err := NewTranscodingService(&config.Config{Server: &server.Config{}}, logrus.New())
	if err != nil {
		t.Fatal(err)
	}
	service.db = fakeDB
	srvr.Register(service)
	body, _ := json.Marshal(data)
	r, _ := http.NewRequest("POST", "/presets", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("wrong response code. Want %d. Got %d", http.StatusOK, w.Code)
	}

	var got map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Errorf("%s: unable to JSON decode response body: %s", w.Body, err)
	}

	expectedBody := map[string]interface{}{
		"Results": map[string]interface{}{
			"fake":     map[string]interface{}{"PresetID": "presetID_here", "Error": ""},
			"zencoder": map[string]interface{}{"PresetID": "presetID_here", "Error": ""},
		},
		"PresetMap": "presetID_here",
	}

	if !reflect.DeepEqual(got, expectedBody) {
		t.Errorf("expected response body of\n%#v;\ngot\n%#v", expectedBody, got)
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
		srvr := server.NewSimpleServer(&server.Config{RouterType: "fast"})
		fakeDB := dbtest.NewFakeRepository(false)
		fakeProviderMapping := make(map[string]string)
		fakeProviderMapping["fake"] = "presetID_here"
		fakeDB.CreatePresetMap(&db.PresetMap{Name: "abc-321", ProviderMapping: fakeProviderMapping})
		service, err := NewTranscodingService(&config.Config{Server: &server.Config{}}, logrus.New())
		if err != nil {
			t.Fatal(err)
		}
		service.db = fakeDB
		srvr.Register(service)
		r, _ := http.NewRequest("DELETE", "/presets/abc-321", nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		var got map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Errorf("%s: unable to JSON decode response body: %s", test.givenTestCase, err)
		}
		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("%s: expected response body of\n%#v;\ngot\n%#v", test.givenTestCase, test.wantBody, got)
		}
	}
}
