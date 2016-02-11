package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nytm/video-transcoding-api/db"
)

// swagger:route POST /presets newPreset
//
// Creates a new preset in the API.
//
//     Responses:
//       200: preset
//       400: validationError
//       409: presetAlreadyExists
//       500: genericError
func (s *TranscodingService) newPreset(r *http.Request) gizmoResponse {
	var params newPresetParams
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return newErrorResponse(err)
	}
	if fieldName, valid := params.Validate(); !valid {
		return newInvalidPresetResponse(fieldName)
	}
	preset := params.Preset()
	err = s.db.SavePreset(&preset)
	switch err {
	case nil:
		return newPresetResponse(&preset)
	case db.ErrPresetAlreadyExists:
		return newPresetAlreadyExistsResponse(err)
	default:
		return newErrorResponse(err)
	}
}

// swagger:route GET /presets/{name} getPreset
//
// Finds a preset using its name.
//
//     Responses:
//       200: preset
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) getPreset(r *http.Request) gizmoResponse {
	name := mux.Vars(r)["name"]
	preset, err := s.db.GetPreset(name)

	switch err {
	case nil:
		return newPresetResponse(preset)
	case db.ErrPresetNotFound:
		return newPresetNotFoundResponse(err)
	default:
		return newErrorResponse(err)
	}
}

// swagger:route DELETE /presets/{name} deletePreset
//
// Deletes a preset by name.
//
//     Responses:
//       200: emptyResponse
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) deletePreset(r *http.Request) gizmoResponse {
	name := mux.Vars(r)["name"]
	err := s.db.DeletePreset(&db.Preset{Name: name})

	switch err {
	case nil:
		return emptyResponse(http.StatusOK)
	case db.ErrPresetNotFound:
		return newPresetNotFoundResponse(err)
	default:
		return newErrorResponse(err)
	}
}

// swagger:route GET /presets listPresets presets
//
// List available presets on the API.
//
//     Responses:
//       200: listPresets
//       500: genericError
func (s *TranscodingService) listPresets(r *http.Request) gizmoResponse {
	presets, err := s.db.ListPresets()
	if err != nil {
		return newErrorResponse(err)
	}
	return newListPresetsResponse(presets)
}
