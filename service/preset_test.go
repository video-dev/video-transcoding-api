package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
)

func TestNewPreset(t *testing.T) {
	tests := []struct {
		givenTestCase       string
		givenRequestData    map[string]interface{}
		givenTriggerDBError bool

		wantCode int
		wantBody map[string]interface{}
	}{
		{
			"New preset",
			map[string]interface{}{
				"name": "abc-123",
				"providerMapping": map[string]string{
					"elementalconductor": "18",
					"elastictranscoder":  "18384284-0002",
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
			},
		},
		{
			"New preset duplicate name",
			map[string]interface{}{
				"name": "abc-321",
				"providerMapping": map[string]string{
					"elementalconductor": "18",
					"elastictranscoder":  "18384284-0002",
				},
			},
			false,

			http.StatusConflict,
			map[string]interface{}{
				"error": db.ErrPresetAlreadyExists.Error(),
			},
		},
		{
			"New preset missing name",
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
		srvr := server.NewSimpleServer(nil)
		fakeDB := newFakeDB(test.givenTriggerDBError)
		fakeDB.CreatePreset(&db.Preset{Name: "abc-321"})
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
		if test.wantCode == http.StatusOK {
			preset, err := fakeDB.GetPreset(got["name"].(string))
			if err != nil {
				t.Error(err)
			} else if !reflect.DeepEqual(preset.ProviderMapping, test.givenRequestData["providerMapping"]) {
				t.Errorf("%s: didn't save the preset in the database. Want %#v. Got %#v", test.givenTestCase, test.givenRequestData, preset.ProviderMapping)
			}
		}
	}
}

func TestGetPreset(t *testing.T) {
	tests := []struct {
		givenTestCase   string
		givenPresetName string

		wantBody *db.Preset
		wantCode int
	}{
		{
			"Get preset",
			"preset-1",
			&db.Preset{Name: "preset-1"},
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
		srvr := server.NewSimpleServer(nil)
		fakeDB := newFakeDB(false)
		fakeDB.CreatePreset(&db.Preset{Name: "preset-1"})
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDB,
		})
		r, _ := http.NewRequest("GET", "/presets/"+test.givenPresetName, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		if test.wantBody != nil {
			var gotPreset db.Preset
			err := json.NewDecoder(w.Body).Decode(&gotPreset)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotPreset, *test.wantBody) {
				t.Errorf("%s: wrong body. Want %#v. Got %#v", test.givenTestCase, *test.wantBody, gotPreset)
			}
		}
	}
}

func TestUpdatePreset(t *testing.T) {
	tests := []struct {
		givenTestCase    string
		givenPresetName  string
		givenRequestData map[string]interface{}

		wantBody *db.Preset
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
			&db.Preset{
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
		srvr := server.NewSimpleServer(nil)
		fakeDB := newFakeDB(false)
		fakeDB.CreatePreset(&db.Preset{
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
		r, _ := http.NewRequest("PUT", "/presets/"+test.givenPresetName, bytes.NewReader(data))
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		if test.wantBody != nil {
			var gotPreset db.Preset
			err := json.NewDecoder(w.Body).Decode(&gotPreset)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotPreset, *test.wantBody) {
				t.Errorf("%s: wrong body. Want %#v. Got %#v", test.givenTestCase, *test.wantBody, gotPreset)
			}
			preset, err := fakeDB.GetPreset(gotPreset.Name)
			if err != nil {
				t.Error(err)
			} else if !reflect.DeepEqual(*preset, gotPreset) {
				t.Errorf("%s: didn't update the preset in the database. Want %#v. Got %#v", test.givenTestCase, gotPreset, *preset)
			}
		}
	}
}

func TestDeletePreset(t *testing.T) {
	tests := []struct {
		givenTestCase   string
		givenPresetName string
		wantCode        int
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
		srvr := server.NewSimpleServer(nil)
		fakeDB := newFakeDB(false)
		fakeDB.CreatePreset(&db.Preset{Name: "preset-1"})
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDB,
		})
		r, _ := http.NewRequest("DELETE", "/presets/"+test.givenPresetName, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		if test.wantCode == http.StatusOK {
			_, err := fakeDB.GetPreset(test.givenPresetName)
			if err != db.ErrPresetNotFound {
				t.Errorf("%s: didn't delete the job in the database", test.givenTestCase)
			}
		}
	}
}

func TestListPresets(t *testing.T) {
	tests := []struct {
		givenTestCase string
		givenPresets  []db.Preset

		wantCode int
		wantBody map[string]db.Preset
	}{
		{
			"List presets",
			[]db.Preset{
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
			map[string]db.Preset{
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
			map[string]db.Preset{},
		},
	}
	for _, test := range tests {
		srvr := server.NewSimpleServer(nil)
		fakeDB := newFakeDB(false)
		for i := range test.givenPresets {
			fakeDB.CreatePreset(&test.givenPresets[i])
		}
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDB,
		})
		r, _ := http.NewRequest("GET", "/presets", nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		var got map[string]db.Preset
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Errorf("%s: unable to JSON decode response body: %s", test.givenTestCase, err)
		}
		if !reflect.DeepEqual(got, test.wantBody) {
			t.Errorf("%s: expected response body of\n%#v;\ngot\n%#v", test.givenTestCase, test.wantBody, got)
		}
	}
}
