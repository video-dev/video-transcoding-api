package elementalconductor

import (
	"io/ioutil"
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
	req  *http.Request
	body []byte
}

func (s *S) startServer(status int, content string) (*httptest.Server, chan fakeServerRequest) {
	requests := make(chan fakeServerRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		fakeRequest := fakeServerRequest{req: r, body: data}
		requests <- fakeRequest
		w.WriteHeader(status)
		w.Write([]byte(content))
	}))
	return server, requests
}
