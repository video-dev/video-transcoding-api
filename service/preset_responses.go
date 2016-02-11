package service

import (
	"fmt"
	"net/http"

	"github.com/nytm/video-transcoding-api/db"
)

// swagger:response preset
type presetResponse struct {
	// in: body
	Payload *db.Preset

	baseResponse
}

func newPresetResponse(preset *db.Preset) *presetResponse {
	return &presetResponse{
		baseResponse: baseResponse{
			payload: preset,
			status:  http.StatusOK,
		},
	}
}

// swagger:response presetNotFound
type presetNotFoundResponse struct {
	// in: body
	Error *errorResponse
}

func newPresetNotFoundResponse(err error) *presetNotFoundResponse {
	return &presetNotFoundResponse{Error: newErrorResponse(err).withStatus(http.StatusNotFound)}
}

func (r *presetNotFoundResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// swagger:response invalidPreset
type invalidPresetResponse struct {
	// in: body
	Error *errorResponse
}

func newInvalidPresetResponse(field string) *invalidPresetResponse {
	err := fmt.Errorf("missing field %s from the request", field)
	return &invalidPresetResponse{Error: newErrorResponse(err).withStatus(http.StatusBadRequest)}
}

func (r *invalidPresetResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// swagger:response presetAlreadyExists
type presetAlreadyExistsResponse struct {
	// in: body
	Error *errorResponse
}

func newPresetAlreadyExistsResponse(err error) *presetAlreadyExistsResponse {
	return &presetAlreadyExistsResponse{Error: newErrorResponse(err).withStatus(http.StatusConflict)}
}

func (r *presetAlreadyExistsResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// swagger:response listPresets
type listPresetsResponse struct {
	// in: body
	PresetsMap map[string]db.Preset

	baseResponse
}

func newListPresetsResponse(presets []db.Preset) *listPresetsResponse {
	presetsMap := make(map[string]db.Preset, len(presets))
	for _, preset := range presets {
		presetsMap[preset.Name] = preset
	}
	return &listPresetsResponse{
		baseResponse: baseResponse{
			status:  http.StatusOK,
			payload: presetsMap,
		},
	}
}
