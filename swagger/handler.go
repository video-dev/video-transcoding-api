package swagger

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
)

// GizmoJSONResponse represents a response type that can be converted to
// Gizmo's JSONEndpoint result format.
type GizmoJSONResponse interface {
	Result() (int, interface{}, error)
}

// Handler represents a function that receives an HTTP request and returns a
// GizmoJSONResponse. It's a Gizmo JSONEndpoint structured for goswagger's
// annotation.
type Handler func(*http.Request) GizmoJSONResponse

// HandlerToJSONEndpoint converts a handler to a proper Gizmo JSONEndpoint.
func HandlerToJSONEndpoint(h Handler) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {
		return h(r).Result()
	}
}
