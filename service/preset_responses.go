package service

import (
	"net/http"

	"github.com/nytm/video-transcoding-api/db"
)

// JSON-encoded preset returned on the newPreset and getPreset operations.
//
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

// error returned when the given preset name is not found on the API (either on
// getPreset or deletePreset operations).
//
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

// error returned when the given preset data is not valid.
//
// swagger:response invalidPreset
type invalidPresetResponse struct {
	// in: body
	Error *errorResponse
}

func newInvalidPresetResponse(err error) *invalidPresetResponse {
	return &invalidPresetResponse{Error: newErrorResponse(err).withStatus(http.StatusBadRequest)}
}

func (r *invalidPresetResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// error returned when trying to create a new preset using a name that is
// already in-use.
//
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

// response for the listPresets operation. It's actually a JSON-encoded object
// instead of an array, in the format `presetName: presetObject`
//
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
