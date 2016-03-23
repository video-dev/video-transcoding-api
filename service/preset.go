package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// swagger:route POST /presetsmap presets newPreset
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

// swagger:route GET /presetsmap/{name} presets getPreset
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

// swagger:route PUT /presets/{name} presets updatePreset
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

// swagger:route DELETE /presets/{name} presets deletePreset
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

// swagger:route GET /presetsmap presets listPresets
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

// swagger:route POST /presets2 presets Output
//
// Creates a new preset on given providers.
//     Responses:
//       200: newPresetOutputs
//       500: genericError
func (s *TranscodingService) newPreset2(r *http.Request) gizmoResponse {
	defer r.Body.Close()
	var input newPresetInput2
	var result interface{}
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
			output[p] = newPresetOutput{Output: "", Error: "getting factory: " + err.Error()}
			continue
		}
		providerObj, err := providerFactory(s.config)
		if err != nil {
			output[p] = newPresetOutput{Output: "", Error: "initializing provider: " + err.Error()}
			continue
		}
		result, err = providerObj.CreatePreset(input.Preset)
		if err != nil {
			output[p] = newPresetOutput{Output: "", Error: "creating preset: " + err.Error()}
		} else {
			output[p] = newPresetOutput{Output: result, Error: ""}
		}
	}

	return &newPresetResponse2{
		baseResponse: baseResponse{
			payload: output,
			status:  http.StatusOK,
		},
	}
}
