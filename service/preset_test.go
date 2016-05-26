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

func TestNewPresetMap(t *testing.T) {
	tests := []struct {
		givenTestCase       string
		givenRequestData    map[string]interface{}
		givenTriggerDBError bool

		wantCode int
		wantBody map[string]interface{}
	}{
		{
			"New presetmap",
			map[string]interface{}{
				"name": "abc-123",
				"providerMapping": map[string]string{
					"elementalconductor": "18",
					"elastictranscoder":  "18384284-0002",
				},
				"output": map[string]interface{}{
					"extension": "mp4",
				},
			},
			false,

			http.StatusOK,
			map[string]interface{}{
				"name": "abc-123",
				"providerMapping": map[string]interface{}{
					"elementalconductor": "18",
					"elastictranscoder":  "18384284-0002",
				},
				"output": map[string]interface{}{
					"extension": "mp4",
				},
			},
		},
		{
			"New presetmap duplicate name",
			map[string]interface{}{
				"name": "abc-321",
				"providerMapping": map[string]string{
					"elementalconductor": "18",
					"elastictranscoder":  "18384284-0002",
				},
				"output": map[string]interface{}{
					"extension": "mp4",
				},
			},
			false,

			http.StatusConflict,
			map[string]interface{}{
				"error": db.ErrPresetMapAlreadyExists.Error(),
			},
		},
		{
			"New presetmap missing name",
			map[string]interface{}{
				"providerMapping": map[string]string{
					"elementalconductor": "18",
					"elastictranscoder":  "18384284-0002",
				},
			},
			false,

			http.StatusBadRequest,
			map[string]interface{}{
				"error": "missing field name from the request",
			},
		},
		{
			"New preset missing providers",
			map[string]interface{}{
				"name":            "mypreset",
				"providerMapping": nil,
			},
			false,

			http.StatusBadRequest,
			map[string]interface{}{
				"error": "missing field providerMapping from the request",
			},
		},
		{
			"New preset DB failure",
			map[string]interface{}{
				"name": "super-preset",
				"providerMapping": map[string]string{
					"elementalconductor": "18",
					"elastictranscoder":  "18384284-0002",
				},
			},
			true,

			http.StatusInternalServerError,
			map[string]interface{}{"error": "database error"},
		},
	}
	for _, test := range tests {
		srvr := server.NewSimpleServer(&gizmoConfig.Server{RouterType: "fast"})
		fakeDB := dbtest.NewFakeRepository(test.givenTriggerDBError)
		fakeDB.CreatePresetMap(&db.PresetMap{Name: "abc-321"})
		srvr.Register(&TranscodingService{config: &config.Config{}, db: fakeDB})
		body, _ := json.Marshal(test.givenRequestData)
		r, _ := http.NewRequest("POST", "/presetmaps", bytes.NewReader(body))
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
		if test.wantCode == http.StatusOK {
			presetmap, err := fakeDB.GetPresetMap(got["name"].(string))
			if err != nil {
				t.Error(err)
			} else if !reflect.DeepEqual(presetmap.ProviderMapping, test.givenRequestData["providerMapping"]) {
				t.Errorf("%s: didn't save the preset in the database. Want %#v. Got %#v", test.givenTestCase, test.givenRequestData, presetmap.ProviderMapping)
			}
		}
	}
}

func TestGetPresetMap(t *testing.T) {
	tests := []struct {
		givenTestCase      string
		givenPresetMapName string

		wantBody *db.PresetMap
		wantCode int
	}{
		{
			"Get preset",
			"preset-1",
			&db.PresetMap{Name: "preset-1"},
			http.StatusOK,
		},
		{
			"Get preset not found",
			"preset-unknown",
			nil,
			http.StatusNotFound,
		},
	}
	for _, test := range tests {
		srvr := server.NewSimpleServer(&gizmoConfig.Server{RouterType: "fast"})
		fakeDB := dbtest.NewFakeRepository(false)
		fakeDB.CreatePresetMap(&db.PresetMap{Name: "preset-1"})
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDB,
		})
		r, _ := http.NewRequest("GET", "/presetmaps/"+test.givenPresetMapName, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		if test.wantBody != nil {
			var gotPresetMap db.PresetMap
			err := json.NewDecoder(w.Body).Decode(&gotPresetMap)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotPresetMap, *test.wantBody) {
				t.Errorf("%s: wrong body. Want %#v. Got %#v", test.givenTestCase, *test.wantBody, gotPresetMap)
			}
		}
	}
}

func TestUpdatePresetMap(t *testing.T) {
	tests := []struct {
		givenTestCase      string
		givenPresetMapName string
		givenRequestData   map[string]interface{}

		wantBody *db.PresetMap
		wantCode int
	}{
		{
			"Update preset",
			"preset-1",
			map[string]interface{}{
				"providerMapping": map[string]string{
					"elementalconductor": "abc-123",
					"elastictranscoder":  "def-345",
				},
			},
			&db.PresetMap{
				Name: "preset-1",
				ProviderMapping: map[string]string{
					"elementalconductor": "abc-123",
					"elastictranscoder":  "def-345",
				},
			},
			http.StatusOK,
		},
		{
			"Update preset not found",
			"preset-unknown",
			map[string]interface{}{
				"providerMapping": map[string]string{
					"elementalconductor": "abc-123",
					"elastictranscoder":  "def-345",
				},
			},
			nil,
			http.StatusNotFound,
		},
	}
	for _, test := range tests {
		srvr := server.NewSimpleServer(&gizmoConfig.Server{RouterType: "fast"})
		fakeDB := dbtest.NewFakeRepository(false)
		fakeDB.CreatePresetMap(&db.PresetMap{
			Name: "preset-1",
			ProviderMapping: map[string]string{
				"elementalconductor": "some-id",
			},
		})
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDB,
		})
		data, _ := json.Marshal(test.givenRequestData)
		r, _ := http.NewRequest("PUT", "/presetmaps/"+test.givenPresetMapName, bytes.NewReader(data))
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		if test.wantBody != nil {
			var gotPresetMap db.PresetMap
			err := json.NewDecoder(w.Body).Decode(&gotPresetMap)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotPresetMap, *test.wantBody) {
				t.Errorf("%s: wrong body. Want %#v. Got %#v", test.givenTestCase, *test.wantBody, gotPresetMap)
			}
			preset, err := fakeDB.GetPresetMap(gotPresetMap.Name)
			if err != nil {
				t.Error(err)
			} else if !reflect.DeepEqual(*preset, gotPresetMap) {
				t.Errorf("%s: didn't update the preset in the database. Want %#v. Got %#v", test.givenTestCase, gotPresetMap, *preset)
			}
		}
	}
}

func TestDeletePresetMap(t *testing.T) {
	tests := []struct {
		givenTestCase      string
		givenPresetMapName string
		wantCode           int
	}{
		{
			"Delete preset",
			"preset-1",
			http.StatusOK,
		},
		{
			"Delete preset not found",
			"preset-unknown",
			http.StatusNotFound,
		},
	}
	for _, test := range tests {
		srvr := server.NewSimpleServer(&gizmoConfig.Server{RouterType: "fast"})
		fakeDB := dbtest.NewFakeRepository(false)
		fakeDB.CreatePresetMap(&db.PresetMap{Name: "preset-1"})
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDB,
		})
		r, _ := http.NewRequest("DELETE", "/presetmaps/"+test.givenPresetMapName, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		if test.wantCode == http.StatusOK {
			_, err := fakeDB.GetPresetMap(test.givenPresetMapName)
			if err != db.ErrPresetMapNotFound {
				t.Errorf("%s: didn't delete the job in the database", test.givenTestCase)
			}
		}
	}
}

func TestListPresetMaps(t *testing.T) {
	tests := []struct {
		givenTestCase   string
		givenPresetMaps []db.PresetMap

		wantCode int
		wantBody map[string]db.PresetMap
	}{
		{
			"List presets",
			[]db.PresetMap{
				{
					Name:            "preset-1",
					ProviderMapping: map[string]string{"elementalconductor": "abc123"},
				},
				{
					Name:            "preset-2",
					ProviderMapping: map[string]string{"elementalconductor": "abc124"},
				},
				{
					Name:            "preset-3",
					ProviderMapping: map[string]string{"elementalconductor": "abc125"},
				},
			},
			http.StatusOK,
			map[string]db.PresetMap{
				"preset-1": {
					Name:            "preset-1",
					ProviderMapping: map[string]string{"elementalconductor": "abc123"},
				},
				"preset-2": {
					Name:            "preset-2",
					ProviderMapping: map[string]string{"elementalconductor": "abc124"},
				},
				"preset-3": {
					Name:            "preset-3",
					ProviderMapping: map[string]string{"elementalconductor": "abc125"},
				},
			},
		},
		{
			"Empty list of presets",
			nil,
			http.StatusOK,
			map[string]db.PresetMap{},
		},
	}
	for _, test := range tests {
		srvr := server.NewSimpleServer(&gizmoConfig.Server{RouterType: "fast"})
		fakeDB := dbtest.NewFakeRepository(false)
		for i := range test.givenPresetMaps {
			fakeDB.CreatePresetMap(&test.givenPresetMaps[i])
		}
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDB,
		})
		r, _ := http.NewRequest("GET", "/presetmaps", nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		var got map[string]db.PresetMap
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Errorf("%s: unable to JSON decode response body: %s", test.givenTestCase, err)
		}
		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("%s: expected response body of\n%#v;\ngot\n%#v", test.givenTestCase, test.wantBody, got)
		}
	}
}

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
