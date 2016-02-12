package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
)

type gizmoResponse interface {
	Result() (int, interface{}, error)
}

func handlerToEndpoint(h func(r *http.Request) gizmoResponse) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {
		return h(r).Result()
	}
}
