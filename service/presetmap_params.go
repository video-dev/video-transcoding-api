package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/video-dev/video-transcoding-api/v2/db"
	"github.com/video-dev/video-transcoding-api/v2/swagger"
)

// JSON-encoded preset returned on the newPreset and getPreset operations.
//
// swagger:response preset
type presetMapResponse struct {
	// in: body
	Payload *db.PresetMap

	baseResponse
}

// swagger:parameters getPreset deletePreset deletePresetMap
type getPresetMapInput struct {
	// in: path
	// required: true
	Name string `json:"name"`
}

// swagger:parameters updatePreset
type updatePresetMapInput struct {
	// in: path
	// required: true
	Name string `json:"name"`

	// in: body
	// required: true
	Payload db.PresetMap

	newPresetMapInput
}

// swagger:parameters newPreset
type newPresetMapInput struct {
	// in: body
	// required: true
	Payload db.PresetMap
}

// error returned when the given preset name is not found on the API (either on
// getPreset or deletePreset operations).
//
// swagger:response presetNotFound
type presetMapNotFoundResponse struct {
	// in: body
	Error *swagger.ErrorResponse
}

// error returned when the given preset data is not valid.
//
// swagger:response invalidPreset
type invalidPresetMapResponse struct {
	// in: body
	Error *swagger.ErrorResponse
}

// error returned when trying to create a new preset using a name that is
// already in-use.
//
// swagger:response presetAlreadyExists
type presetMapAlreadyExistsResponse struct {
	// in: body
	Error *swagger.ErrorResponse
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

func newPresetMapResponse(preset *db.PresetMap) *presetMapResponse {
	return &presetMapResponse{
		baseResponse: baseResponse{
			payload: preset,
			status:  http.StatusOK,
		},
	}
}

func newPresetMapNotFoundResponse(err error) *presetMapNotFoundResponse {
	return &presetMapNotFoundResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusNotFound)}
}

func (r *presetMapNotFoundResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

func newInvalidPresetMapResponse(err error) *invalidPresetMapResponse {
	return &invalidPresetMapResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusBadRequest)}
}

func (r *invalidPresetMapResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
}

func newPresetMapAlreadyExistsResponse(err error) *presetMapAlreadyExistsResponse {
	return &presetMapAlreadyExistsResponse{Error: swagger.NewErrorResponse(err).WithStatus(http.StatusConflict)}
}

func (r *presetMapAlreadyExistsResponse) Result() (int, interface{}, error) {
	return r.Error.Result()
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

// Preset loads the input from the request body, validates them and returns the
// preset.
func (p *newPresetMapInput) PresetMap(body io.Reader) (db.PresetMap, error) {
	err := json.NewDecoder(body).Decode(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	err = validatePresetMap(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	err = p.Payload.OutputOpts.Validate()
	if err != nil {
		return p.Payload, fmt.Errorf("invalid output: %s", err)
	}
	return p.Payload, nil
}

func (p *getPresetMapInput) loadParams(paramsMap map[string]string) {
	p.Name = paramsMap["name"]
}

func (p *updatePresetMapInput) PresetMap(paramsMap map[string]string, body io.Reader) (db.PresetMap, error) {
	p.Name = paramsMap["name"]
	err := json.NewDecoder(body).Decode(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	p.Payload.Name = p.Name
	err = validatePresetMap(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	return p.Payload, nil
}

func validatePresetMap(p *db.PresetMap) error {
	if p.Name == "" {
		return errors.New("missing field name from the request")
	}
	if len(p.ProviderMapping) == 0 {
		return errors.New("missing field providerMapping from the request")
	}
	return nil
}
