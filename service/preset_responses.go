package service

import (
	"net/http"

	"github.com/video-dev/video-transcoding-api/swagger"
)

type newPresetResponse struct {
	baseResponse
}

type deletePresetResponse struct {
	baseResponse
}

// error returned when the given preset data is not valid.
//
// swagger:response invalidPreset
type invalidPresetResponse struct {
	// in: body
	Error *swagger.ErrorResponse
}

func newInvalidPresetResponse(err error) *invalidPresetResponse {
	return &invalidPresetResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusBadRequest)}
}

func (r *invalidPresetResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}
