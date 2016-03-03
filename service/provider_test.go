package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

func TestListProviders(t *testing.T) {
	srvr := server.NewSimpleServer(nil)
	srvr.Register(&TranscodingService{config: &config.Config{}})
	r, _ := http.NewRequest("GET", "/providers", nil)
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("listProviders: wrong status code. Want %d. Got %d", http.StatusOK, w.Code)
	}
	var providers []provider.Descriptor
	err := json.NewDecoder(w.Body).Decode(&providers)
	if err != nil {
		t.Fatal(err)
	}
	expected := []provider.Descriptor{
		{
			Name:   "fake",
			Health: provider.Health{OK: true},
			Capabilities: provider.Capabilities{
				InputFormats:  []string{"prores", "h264"},
				OutputFormats: []string{"mp4", "webm", "hls"},
				Destinations:  []string{"akamai", "s3"},
			},
		},
	}
	if !reflect.DeepEqual(providers, expected) {
		t.Errorf("listProviders: wrong body. Want %#v. Got %#v", expected, providers)
	}
}
