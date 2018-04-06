package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/sirupsen/logrus"
)

func TestListProviders(t *testing.T) {
	srvr := server.NewSimpleServer(&server.Config{RouterType: "fast"})
	service, err := NewTranscodingService(&config.Config{Server: &server.Config{}}, logrus.New())
	if err != nil {
		t.Fatal(err)
	}
	srvr.Register(service)
	r, _ := http.NewRequest("GET", "/providers", nil)
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("listProviders: wrong status code. Want %d. Got %d", http.StatusOK, w.Code)
	}
	var providers []string
	err = json.NewDecoder(w.Body).Decode(&providers)
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"fake", "zencoder"}
	if !reflect.DeepEqual(providers, expected) {
		t.Errorf("listProviders: wrong body. Want %#v. Got %#v", expected, providers)
	}
}

func TestGetProvider(t *testing.T) {
	var tests = []struct {
		testCase string
		name     string

		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			"get provider",
			"fake",
			http.StatusOK,
			map[string]interface{}{
				"name":   "fake",
				"health": map[string]interface{}{"ok": true},
				"capabilities": map[string]interface{}{
					"input":        []interface{}{"prores", "h264"},
					"output":       []interface{}{"mp4", "webm", "hls"},
					"destinations": []interface{}{"akamai", "s3"},
				},
				"enabled": true,
			},
		},
		{
			"provider not found",
			"whatever",
			http.StatusNotFound,
			map[string]interface{}{"error": "provider not found"},
		},
	}
	for _, test := range tests {
		srvr := server.NewSimpleServer(&server.Config{RouterType: "fast"})
		service, err := NewTranscodingService(&config.Config{Server: &server.Config{}}, logrus.New())
		if err != nil {
			t.Fatal(err)
		}
		srvr.Register(service)
		r, _ := http.NewRequest("GET", "/providers/"+test.name, nil)
		w := httptest.NewRecorder()
		srvr.ServeHTTP(w, r)
		if w.Code != test.expectedStatus {
			t.Errorf("%s: wrong status code. Want %d. Got %d", test.testCase, test.expectedStatus, w.Code)
		}
		var gotBody map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&gotBody)
		if err != nil {
			t.Errorf("%s: %s", test.testCase, err)
		}
		if !reflect.DeepEqual(gotBody, test.expectedBody) {
			t.Errorf("%s: wrong body.\nWant %#v.\nGot  %#v", test.testCase, test.expectedBody, gotBody)
		}
	}
}
