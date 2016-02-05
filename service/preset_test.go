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
		givenRequestData    map[string]string
		givenTriggerDBError bool

		wantCode int
		wantBody map[string]interface{}
	}{
		{
			"New preset",
			map[string]string{"elementalcloud": "18", "elastictranscoder": "18384284-0002"},
			false,

			http.StatusOK,
			map[string]interface{}{"presetId": "12345"},
		},
		{
			"New preset DB failure",
			map[string]string{"elementalcloud": "18", "elastictranscoder": "18384284-0002"},
			true,

			http.StatusInternalServerError,
			map[string]interface{}{"error": "database error"},
		},
	}
	for _, test := range tests {
		srvr := server.NewSimpleServer(nil)
		fakeDB := newFakeDB(test.givenTriggerDBError)
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDB,
		})
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
			preset, err := fakeDB.GetPreset(got["presetId"].(string))
			if err != nil {
				t.Error(err)
			} else if !reflect.DeepEqual(preset.ProviderMapping, test.givenRequestData) {
				t.Errorf("%s: didn't save the preset in the database. Want %#v. Got %#v", test.givenTestCase, test.givenRequestData, preset.ProviderMapping)
			}
		}
	}
}

func TestDeletePreset(t *testing.T) {
	tests := []struct {
		givenTestCase string
		givenPresetID string
		wantCode      int
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
		fakeDB.SavePreset(&db.Preset{ID: "preset-1"})
		srvr.Register(&TranscodingService{
			config: &config.Config{},
			db:     fakeDB,
		})
		r, _ := http.NewRequest("DELETE", "/presets/"+test.givenPresetID, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("%s: wrong response code. Want %d. Got %d", test.givenTestCase, test.wantCode, w.Code)
		}
		if test.wantCode == http.StatusOK {
			_, err := fakeDB.GetPreset(test.givenPresetID)
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
					ID:              "preset-1",
					ProviderMapping: map[string]string{"elementalcloud": "abc123"},
				},
				{
					ID:              "preset-2",
					ProviderMapping: map[string]string{"elementalcloud": "abc124"},
				},
				{
					ID:              "preset-3",
					ProviderMapping: map[string]string{"elementalcloud": "abc125"},
				},
			},
			http.StatusOK,
			map[string]db.Preset{
				"preset-1": {
					ID:              "preset-1",
					ProviderMapping: map[string]string{"elementalcloud": "abc123"},
				},
				"preset-2": {
					ID:              "preset-2",
					ProviderMapping: map[string]string{"elementalcloud": "abc124"},
				},
				"preset-3": {
					ID:              "preset-3",
					ProviderMapping: map[string]string{"elementalcloud": "abc125"},
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
			fakeDB.SavePreset(&test.givenPresets[i])
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
