package service

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/sirupsen/logrus"
)

func TestSwaggerManifest(t *testing.T) {
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
		"Access-Control-Allow-Headers": "Content-Type",
		"Content-Type":                 "application/json",
	}
	srvr := server.NewSimpleServer(nil)
	srvr.Register(&TranscodingService{
		config: &config.Config{
			SwaggerManifest: "testdata/swagger.json",
			Server:          &server.Config{},
		},
		logger: logrus.New(),
	})
	r, _ := http.NewRequest("GET", "/swagger.json", nil)
	w := httptest.NewRecorder()
	srvr.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Swagger manifest: wrong status code. Want %d. Got %d", http.StatusOK, w.Code)
	}
	for k, v := range expectedHeaders {
		got := w.Header().Get(k)
		if got != v {
			t.Errorf("Swagger manifest: wrong header value for key=%q. Want %q. Got %q", k, v, got)
		}
	}
	expectedData, err := ioutil.ReadFile("testdata/swagger.json")
	if err != nil {
		t.Fatal(err)
	}
	if string(expectedData) != w.Body.String() {
		t.Errorf("Swagger manifest: wrong body\nWant: %s\nGot:  %s", expectedData, w.Body)
	}
}
