package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// swagger:route POST /presets presets newPreset
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
	err = s.db.CreatePreset(&preset)
	switch err {
	case nil:
		return newPresetResponse(&preset)
	case db.ErrPresetAlreadyExists:
		return newPresetAlreadyExistsResponse(err)
	default:
		return newErrorResponse(err)
	}
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

// swagger:route GET /presets/{name} presets getPreset
//
// Finds a preset using its name.
//
//     Responses:
//       200: preset
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) getPreset(r *http.Request) gizmoResponse {
	var params getPresetInput
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

// swagger:route PUT /presets/{name} presets updatePreset
//
// Updates a preset using its name.
//
//     Responses:
//       200: preset
//       400: invalidPreset
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) updatePreset(r *http.Request) gizmoResponse {
	defer r.Body.Close()
	var input updatePresetInput
	preset, err := input.Preset(mux.Vars(r), r.Body)
	if err != nil {
		return newInvalidPresetResponse(err)
	}
	err = s.db.UpdatePreset(&preset)

	switch err {
	case nil:
		return newPresetResponse(&preset)
	case db.ErrPresetNotFound:
		return newPresetNotFoundResponse(err)
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
func (s *TranscodingService) deletePreset(r *http.Request) gizmoResponse {
	var params getPresetInput
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

// swagger:route GET /presets presets listPresets
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
