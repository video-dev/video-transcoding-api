package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/video-dev/video-transcoding-api/db"
	"github.com/video-dev/video-transcoding-api/swagger"
)

// swagger:route POST /presetmaps presets newPreset
//
// Creates a new preset in the API.
//
//     Responses:
//       200: preset
//       400: invalidPreset
//       409: presetAlreadyExists
//       500: genericError
func (s *TranscodingService) newPresetMap(r *http.Request) swagger.GizmoJSONResponse {
	var input newPresetMapInput
	defer r.Body.Close()
	preset, err := input.PresetMap(r.Body)
	if err != nil {
		return newInvalidPresetMapResponse(err)
	}
	err = s.db.CreatePresetMap(&preset)
	switch err {
	case nil:
		return newPresetMapResponse(&preset)
	case db.ErrPresetMapAlreadyExists:
		return newPresetMapAlreadyExistsResponse(err)
	default:
		return swagger.NewErrorResponse(err)
	}
}

// swagger:route GET /presetmaps/{name} presets getPreset
//
// Finds a preset using its name.
//
//     Responses:
//       200: preset
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) getPresetMap(r *http.Request) swagger.GizmoJSONResponse {
	var params getPresetMapInput
	params.loadParams(server.Vars(r))
	preset, err := s.db.GetPresetMap(params.Name)

	switch err {
	case nil:
		return newPresetMapResponse(preset)
	case db.ErrPresetMapNotFound:
		return newPresetMapNotFoundResponse(err)
	default:
		return swagger.NewErrorResponse(err)
	}
}

// swagger:route PUT /presetmaps/{name} presets updatePreset
//
// Updates a presetmap using its name.
//
//     Responses:
//       200: preset
//       400: invalidPreset
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) updatePresetMap(r *http.Request) swagger.GizmoJSONResponse {
	defer r.Body.Close()
	var input updatePresetMapInput
	presetMap, err := input.PresetMap(server.Vars(r), r.Body)
	if err != nil {
		return newInvalidPresetMapResponse(err)
	}
	err = s.db.UpdatePresetMap(&presetMap)

	switch err {
	case nil:
		updatedPresetMap, _ := s.db.GetPresetMap(presetMap.Name)
		return newPresetMapResponse(updatedPresetMap)
	case db.ErrPresetMapNotFound:
		return newPresetMapNotFoundResponse(err)
	default:
		return swagger.NewErrorResponse(err)
	}
}

// swagger:route DELETE /presetmaps/{name} presets deletePresetMap
//
// Deletes a presetmap by name.
//
//     Responses:
//       200: emptyResponse
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) deletePresetMap(r *http.Request) swagger.GizmoJSONResponse {
	var params getPresetMapInput
	params.loadParams(server.Vars(r))
	err := s.db.DeletePresetMap(&db.PresetMap{Name: params.Name})

	switch err {
	case nil:
		return emptyResponse(http.StatusOK)
	case db.ErrPresetMapNotFound:
		return newPresetMapNotFoundResponse(err)
	default:
		return swagger.NewErrorResponse(err)
	}
}

// swagger:route GET /presetmaps presets listPresetMaps
//
// List available presets on the API.
//
//     Responses:
//       200: listPresetMaps
//       500: genericError
func (s *TranscodingService) listPresetMaps(*http.Request) swagger.GizmoJSONResponse {
	presetsMap, err := s.db.ListPresetMaps()
	if err != nil {
		return swagger.NewErrorResponse(err)
	}
	return newListPresetMapsResponse(presetsMap)
}
