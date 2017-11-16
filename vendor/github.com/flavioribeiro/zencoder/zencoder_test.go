package zencoder

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenericCall(t *testing.T) {
	var headers http.Header

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		headers = r.Header
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	resp, err := zc.call("GET", "test", nil, []int{http.StatusOK})
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if resp == nil {
		t.Fatal("Expected a response")
	}

	if len(headers) == 0 {
		t.Fatal("Expected headers")
	}

	if len(headers["User-Agent"]) == 0 {
		t.Fatal("Expected User-Agent")
	}

	if headers["User-Agent"][0] != "gozencoder v1" {
		t.Fatal("Expected User-Agent=gozencoder v1", headers["User-Agent"])
	}

	if len(headers["Accept"]) == 0 {
		t.Fatal("Expected Accept")
	}

	if headers["Accept"][0] != "application/json" {
		t.Fatal("Expected Accept=application/json", headers["Accept"])
	}

	if len(headers["Content-Type"]) == 0 {
		t.Fatal("Expected Content-Type")
	}

	if headers["Content-Type"][0] != "application/json" {
		t.Fatal("Expected Content-Type=application/json", headers["Content-Type"])
	}

	if len(headers["Zencoder-Api-Key"]) == 0 {
		t.Fatal("Expected Zencoder-Api-Key")
	}

	if headers["Zencoder-Api-Key"][0] != "abc" {
		t.Fatal("Expected Zencoder-Api-Key=abc", headers["Zencoder-Api-Key"])
	}
}
