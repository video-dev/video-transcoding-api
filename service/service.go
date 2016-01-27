package service

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/gziphandler"
	"github.com/Sirupsen/logrus"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// TranscodingService will implement server.JSONService and handle all requests
// to the server.
type TranscodingService struct {
	config    *config.Config
	db        db.JobRepository
	providers map[string]provider.Factory
}

// NewTranscodingService will instantiate a JSONService
// with the given configuration.
func NewTranscodingService(cfg *config.Config) (*TranscodingService, error) {
	dbRepo, err := db.NewRedisJobRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error initializing Redis client: %s", err)
	}
	return &TranscodingService{
		config: cfg,
		db:     dbRepo,
		providers: map[string]provider.Factory{
			"encoding.com": provider.EncodingComProvider,
		},
	}, nil
}

// Prefix returns the string prefix used for all endpoints within
// this service.
func (s *TranscodingService) Prefix() string {
	return "/"
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s *TranscodingService) Middleware(h http.Handler) http.Handler {
	return gziphandler.GzipHandler(h)
}

// JSONMiddleware provides a JSONEndpoint hook wrapped around all requests.
// In this implementation, we're using it to provide application logging and to check errors
// and provide generic responses.
func (s *TranscodingService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {

		status, res, err := j(r)
		if err != nil {
			if status == http.StatusServiceUnavailable {
				server.LogWithFields(r).WithFields(logrus.Fields{
					"error": err,
				}).Error("problems with serving request")
				return http.StatusServiceUnavailable, nil, &jsonErr{"sorry, this service is unavailable"}
			}
			return status, nil, &jsonErr{err.Error()}
		}

		server.LogWithFields(r).Info("success!")
		return status, res, nil
	}
}

// JSONEndpoints is a listing of all endpoints available in the JSONService.
func (s *TranscodingService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{
		"/jobs": {
			"POST": s.newTranscodeJob,
		},
		"/jobs/{jobId:[^/]+}": {
			"GET": s.getTranscodeJob,
		},
	}
}

type jsonErr struct {
	Err string `json:"error"`
}

func (e *jsonErr) Error() string {
	return e.Err
}
