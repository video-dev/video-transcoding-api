package provider

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
)

var errMediaNotFound = errors.New("media not found")

type request struct {
	Action  string              `json:"action"`
	MediaID string              `json:"mediaid"`
	Source  []string            `json:"source"`
	Format  *encodingcom.Format `json:"format"`
}

type errorResponse struct {
	Message string    `json:"message"`
	Errors  errorList `json:"errors"`
}

type errorList struct {
	Error []string `json:"error"`
}

type fakeMedia struct {
	ID      string
	Request request
}

// encodingComFakeServer is a fake version of the Encoding.com API.
type encodingComFakeServer struct {
	*httptest.Server
	medias map[string]fakeMedia
}

func newEncodingComFakeServer() *encodingComFakeServer {
	server := encodingComFakeServer{medias: make(map[string]fakeMedia)}
	server.Server = httptest.NewServer(&server)
	return &server
}

func (s *encodingComFakeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestData := r.FormValue("json")
	if requestData == "" {
		s.Error(w, "json is required")
		return
	}
	var m map[string]request
	err := json.Unmarshal([]byte(requestData), &m)
	if err != nil {
		s.Error(w, err.Error())
		return
	}
	req := m["query"]
	switch req.Action {
	case "AddMedia":
		id := generateID()
		s.medias[id] = fakeMedia{ID: id, Request: req}
		resp := map[string]encodingcom.AddMediaResponse{
			"response": {MediaID: id, Message: "it worked"},
		}
		json.NewEncoder(w).Encode(resp)
	default:
		s.Error(w, "invalid action")
	}
}

func (s *encodingComFakeServer) Error(w http.ResponseWriter, message string) {
	m := map[string]errorResponse{"response": {
		Errors: errorList{Error: []string{message}},
	}}
	json.NewEncoder(w).Encode(m)
}

func (s *encodingComFakeServer) getMedia(id string) (fakeMedia, error) {
	media, ok := s.medias[id]
	if !ok {
		return fakeMedia{}, errMediaNotFound
	}
	return media, nil
}

func generateID() string {
	var id [16]byte
	rand.Read(id[:])
	return fmt.Sprintf("%x", id[:])
}
