package service

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/gziphandler"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/db/redis"
)

// TranscodingService will implement server.JSONService and handle all requests
// to the server.
type TranscodingService struct {
	config *config.Config
	db     db.Repository
}

// NewTranscodingService will instantiate a JSONService
// with the given configuration.
func NewTranscodingService(cfg *config.Config) (*TranscodingService, error) {
	dbRepo, err := redis.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error initializing Redis client: %s", err)
	}
	return &TranscodingService{config: cfg, db: dbRepo}, nil
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
	return gziphandler.GzipHandler(server.CORSHandler(h, ""))
}

// JSONMiddleware provides a JSONEndpoint hook wrapped around all requests.
func (s *TranscodingService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {
		status, res, err := j(r)
		if err != nil {
			return newErrorResponse(err).withStatus(status).Result()
		}
		return status, res, nil
	}
}

// JSONEndpoints is a listing of all endpoints available in the JSONService.
func (s *TranscodingService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{
		"/jobs": {
			"POST": handlerToEndpoint(s.newTranscodeJob),
		},
		"/jobs/{jobId:[^/]+}": {
			"GET": handlerToEndpoint(s.getTranscodeJob),
		},
		"/presets": {
			"POST": handlerToEndpoint(s.newPreset),
			"GET":  handlerToEndpoint(s.listPresets),
		},
		"/presets2": {
			"POST": handlerToEndpoint(s.newPreset2),
		},
		"/presets/{name:[^/]+}": {
			"GET":    handlerToEndpoint(s.getPreset),
			"PUT":    handlerToEndpoint(s.updatePreset),
			"DELETE": handlerToEndpoint(s.deletePreset),
		},
		"/providers": {
			"GET": handlerToEndpoint(s.listProviders),
		},
		"/providers/{name:[^/]+}": {
			"GET": handlerToEndpoint(s.getProvider),
		},
	}
}

// Endpoints is a list of all non-json endpoints.
func (s *TranscodingService) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"/swagger.json": {
			"GET": s.swaggerManifest,
		},
	}
}
