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
//       409: presetAlreadyExists
//       500: genericError
func (s *TranscodingService) newPreset(r *http.Request) gizmoResponse {
	var params newPresetParams
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return newErrorResponse(err)
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

// swagger:route GET /presets/{presetId} getPreset
//
// Finds a preset using its id.
//
//     Responses:
//       200: preset
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) getPreset(r *http.Request) gizmoResponse {
	presetID := mux.Vars(r)["presetId"]
	preset, err := s.db.GetPreset(presetID)

	switch err {
	case nil:
		return newPresetResponse(preset)
	case db.ErrPresetNotFound:
		return newPresetNotFoundResponse(err)
	default:
		return newErrorResponse(err)
	}
}

// swagger:route DELETE /presets/{presetId} deletePreset
//
// Deletes a preset by id.
//
//     Responses:
//       200: emptyResponse
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) deletePreset(r *http.Request) gizmoResponse {
	presetID := mux.Vars(r)["presetId"]
	err := s.db.DeletePreset(&db.Preset{ID: presetID})

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
