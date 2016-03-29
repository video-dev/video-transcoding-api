package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
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
func (s *TranscodingService) newPresetMap(r *http.Request) gizmoResponse {
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
		return newErrorResponse(err)
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
func (s *TranscodingService) getPresetMap(r *http.Request) gizmoResponse {
	var params getPresetMapInput
	params.loadParams(mux.Vars(r))
	preset, err := s.db.GetPresetMap(params.Name)

	switch err {
	case nil:
		return newPresetMapResponse(preset)
	case db.ErrPresetMapNotFound:
		return newPresetMapNotFoundResponse(err)
	default:
		return newErrorResponse(err)
	}
}

// swagger:route PUT /presetmaps/{name} presets updatePreset
//
// Updates a preset using its name.
//
//     Responses:
//       200: preset
//       400: invalidPreset
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) updatePresetMap(r *http.Request) gizmoResponse {
	defer r.Body.Close()
	var input updatePresetMapInput
	presetMap, err := input.PresetMap(mux.Vars(r), r.Body)
	if err != nil {
		return newInvalidPresetMapResponse(err)
	}
	err = s.db.UpdatePresetMap(&presetMap)

	switch err {
	case nil:
		return newPresetMapResponse(&presetMap)
	case db.ErrPresetMapNotFound:
		return newPresetMapNotFoundResponse(err)
	default:
		return newErrorResponse(err)
	}
}

// swagger:route DELETE /presetmaps/{name} presets deletePreset
//
// Deletes a preset by name.
//
//     Responses:
//       200: emptyResponse
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) deletePresetMap(r *http.Request) gizmoResponse {
	var params getPresetMapInput
	params.loadParams(mux.Vars(r))
	err := s.db.DeletePresetMap(&db.PresetMap{Name: params.Name})

	switch err {
	case nil:
		return emptyResponse(http.StatusOK)
	case db.ErrPresetMapNotFound:
		return newPresetMapNotFoundResponse(err)
	default:
		return newErrorResponse(err)
	}
}

// swagger:route GET /presetmaps presets listPresets
//
// List available presets on the API.
//
//     Responses:
//       200: listPresetMaps
//       500: genericError
func (s *TranscodingService) listPresetMaps(r *http.Request) gizmoResponse {
	presetsMap, err := s.db.ListPresetMaps()
	if err != nil {
		return newErrorResponse(err)
	}
	return newListPresetMapsResponse(presetsMap)
}

// swagger:route POST /presets presets Output
//
// Creates a new preset on given providers.
//     Responses:
//       200: newPresetOutputs
//       500: genericError
func (s *TranscodingService) newPreset(r *http.Request) gizmoResponse {
	defer r.Body.Close()
	var input newPresetInput
	var output = make(newPresetOutputs)

	respData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return newErrorResponse(err)
	}

	err = json.Unmarshal(respData, &input)
	if err != nil {
		return newErrorResponse(err)
	}

	for _, p := range input.Providers {
		providerFactory, err := provider.GetProviderFactory(p)
		if err != nil {
			output[p] = newPresetOutput{PresetID: "", Error: "getting factory: " + err.Error()}
			continue
		}
		providerObj, err := providerFactory(s.config)
		if err != nil {
			output[p] = newPresetOutput{PresetID: "", Error: "initializing provider: " + err.Error()}
			continue
		}
		presetID, err := providerObj.CreatePreset(input.Preset)
		if err != nil {
			output[p] = newPresetOutput{PresetID: "", Error: "creating preset: " + err.Error()}
		} else {
			output[p] = newPresetOutput{PresetID: presetID, Error: ""}
		}
	}

	return &newPresetResponse{
		baseResponse: baseResponse{
			payload: output,
			status:  http.StatusOK,
		},
	}
}
