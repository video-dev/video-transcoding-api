package service

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/gziphandler"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis"
	"github.com/NYTimes/video-transcoding-api/swagger"
	"github.com/fsouza/ctxlogger"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
)

// TranscodingService will implement server.JSONService and handle all requests
// to the server.
type TranscodingService struct {
	config *config.Config
	db     db.Repository
	logger *logrus.Logger
}

// NewTranscodingService will instantiate a JSONService
// with the given configuration.
func NewTranscodingService(cfg *config.Config, logger *logrus.Logger) (*TranscodingService, error) {
	dbRepo, err := redis.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error initializing Redis client: %s", err)
	}
	return &TranscodingService{config: cfg, db: dbRepo, logger: logger}, nil
}

// Prefix returns the string prefix used for all endpoints within
// this service.
func (s *TranscodingService) Prefix() string {
	return ""
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s *TranscodingService) Middleware(h http.Handler) http.Handler {
	logMiddleware := ctxlogger.ContextLogger(s.logger)
	h = logMiddleware(h)
	if s.config.Server.HTTPAccessLog == nil {
		h = handlers.LoggingHandler(s.logger.Writer(), h)
	}
	return gziphandler.GzipHandler(server.CORSHandler(h, ""))
}

// JSONMiddleware provides a JSONEndpoint hook wrapped around all requests.
func (s *TranscodingService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {
		status, res, err := j(r)
		if err != nil {
			return swagger.NewErrorResponse(err).WithStatus(status).Result()
		}
		return status, res, nil
	}
}

// JSONEndpoints is a listing of all endpoints available in the JSONService.
func (s *TranscodingService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{
		"/jobs": {
			"POST": swagger.HandlerToJSONEndpoint(s.newTranscodeJob),
		},
		"/jobs/:jobId": {
			"GET": swagger.HandlerToJSONEndpoint(s.getTranscodeJob),
		},
		"/jobs/:jobId/cancel": {
			"POST": swagger.HandlerToJSONEndpoint(s.cancelTranscodeJob),
		},
		"/presets": {
			"POST": swagger.HandlerToJSONEndpoint(s.newPreset),
		},
		"/presets/:name": {
			"DELETE": swagger.HandlerToJSONEndpoint(s.deletePreset),
		},
		"/presetmaps": {
			"POST": swagger.HandlerToJSONEndpoint(s.newPresetMap),
			"GET":  swagger.HandlerToJSONEndpoint(s.listPresetMaps),
		},
		"/presetmaps/:name": {
			"GET":    swagger.HandlerToJSONEndpoint(s.getPresetMap),
			"PUT":    swagger.HandlerToJSONEndpoint(s.updatePresetMap),
			"DELETE": swagger.HandlerToJSONEndpoint(s.deletePresetMap),
		},
		"/providers": {
			"GET": swagger.HandlerToJSONEndpoint(s.listProviders),
		},
		"/providers/:name": {
			"GET": swagger.HandlerToJSONEndpoint(s.getProvider),
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
