package encodingcom

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/check.v1"
)

type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&S{})

type fakeServerRequest struct {
	req   *http.Request
	query map[string]interface{}
}

func (s *S) startServer(content string) (*httptest.Server, chan fakeServerRequest) {
	requests := make(chan fakeServerRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := r.FormValue("json")
		var m map[string]interface{}
		json.Unmarshal([]byte(data), &m)
		fakeRequest := fakeServerRequest{req: r, query: m["query"].(map[string]interface{})}
		requests <- fakeRequest
		w.Write([]byte(content))
	}))
	return server, requests
}
