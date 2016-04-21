package service

import (
	"net/http"

	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/swagger"
)

// JSON-encoded preset returned on the newPreset and getPreset operations.
//
// swagger:response preset
type presetMapResponse struct {
	// in: body
	Payload *db.PresetMap

	baseResponse
}

type newPresetResponse struct {
	baseResponse
}

type deletePresetResponse struct {
	baseResponse
}

func newPresetMapResponse(preset *db.PresetMap) *presetMapResponse {
	return &presetMapResponse{
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
type presetMapNotFoundResponse struct {
	// in: body
	Error *swagger.ErrorResponse
}

func newPresetMapNotFoundResponse(err error) *presetMapNotFoundResponse {
	return &presetMapNotFoundResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusNotFound)}
}

func (r *presetMapNotFoundResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// error returned when the given preset data is not valid.
//
// swagger:response invalidPreset
type invalidPresetMapResponse struct {
	// in: body
	Error *swagger.ErrorResponse
}

func newInvalidPresetMapResponse(err error) *invalidPresetMapResponse {
	return &invalidPresetMapResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusBadRequest)}
}

func (r *invalidPresetMapResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// error returned when trying to create a new preset using a name that is
// already in-use.
//
// swagger:response presetAlreadyExists
type presetMapAlreadyExistsResponse struct {
	// in: body
	Error *swagger.ErrorResponse
}

func newPresetMapAlreadyExistsResponse(err error) *presetMapAlreadyExistsResponse {
	return &presetMapAlreadyExistsResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusConflict)}
}

func (r *presetMapAlreadyExistsResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

// response for the listPresetMaps operation. It's actually a JSON-encoded object
// instead of an array, in the format `presetName: presetObject`
//
// swagger:response listPresetMaps
type listPresetMapsResponse struct {
	// in: body
	PresetMaps map[string]db.PresetMap

	baseResponse
}

func newListPresetMapsResponse(presetsMap []db.PresetMap) *listPresetMapsResponse {
	Map := make(map[string]db.PresetMap, len(presetsMap))
	for _, presetMap := range presetsMap {
		Map[presetMap.Name] = presetMap
	}
	return &listPresetMapsResponse{
		baseResponse: baseResponse{
			status:  http.StatusOK,
			payload: Map,
		},
	}
}
