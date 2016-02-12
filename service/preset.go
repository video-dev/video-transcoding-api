package service

import (
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
//       400: invalidPreset
//       409: presetAlreadyExists
//       500: genericError
func (s *TranscodingService) newPreset(r *http.Request) gizmoResponse {
	var input newPresetInput
	defer r.Body.Close()
	preset, err := input.Preset(r.Body)
	if err != nil {
		return newInvalidPresetResponse(err)
	}
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

// swagger:route GET /presets/{Name} getPreset
//
// Finds a preset using its name.
//
//     Responses:
//       200: preset
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) getPreset(r *http.Request) gizmoResponse {
	var params getPresetParams
	params.loadParams(mux.Vars(r))
	preset, err := s.db.GetPreset(params.Name)

	switch err {
	case nil:
		return newPresetResponse(preset)
	case db.ErrPresetNotFound:
		return newPresetNotFoundResponse(err)
	default:
		return newErrorResponse(err)
	}
}

// swagger:route DELETE /presets/{Name} deletePreset
//
// Deletes a preset by name.
//
//     Responses:
//       200: emptyResponse
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) deletePreset(r *http.Request) gizmoResponse {
	var params getPresetParams
	params.loadParams(mux.Vars(r))
	err := s.db.DeletePreset(&db.Preset{Name: params.Name})

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
