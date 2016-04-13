package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/NYTimes/gizmo/web"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// swagger:route POST /presetmaps presets newPreset
//
// Creates a new preset in the API.
//
//     Responses:
//       200: presetmap
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
//       200: presetmap
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) getPresetMap(r *http.Request) gizmoResponse {
	var params getPresetMapInput
	params.loadParams(web.Vars(r))
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
// Updates a presetmap using its name.
//
//     Responses:
//       200: preset
//       400: invalidPreset
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) updatePresetMap(r *http.Request) gizmoResponse {
	defer r.Body.Close()
	var input updatePresetMapInput
	presetMap, err := input.PresetMap(web.Vars(r), r.Body)
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
		return newErrorResponse(err)
	}
}

// swagger:route DELETE /presetmaps/{name} presets deletePreset
//
// Deletes a presetmap by name.
//
//     Responses:
//       200: emptyResponse
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) deletePresetMap(r *http.Request) gizmoResponse {
	var params getPresetMapInput
	params.loadParams(web.Vars(r))
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

// swagger:route DELETE /presets/{name} presets deletePreset
//
// Deletes a preset by name.
//
//     Responses:
//       200: emptyResponse
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) deletePreset(r *http.Request) gizmoResponse {
	var output deletePresetOutputs
	var params getPresetMapInput
	params.loadParams(web.Vars(r))

	output.Results = make(map[string]deletePresetOutput)

	presetmap, err := s.db.GetPresetMap(params.Name)
	if err != nil {
		output.Status = "couldn't retrieve preset map: " + err.Error()
	} else {
		for p, presetID := range presetmap.ProviderMapping {
			providerFactory, err := provider.GetProviderFactory(p)
			if err != nil {
				output.Results[p] = deletePresetOutput{PresetID: "", Error: "getting factory: " + err.Error()}
				continue
			}
			providerObj, err := providerFactory(s.config)
			if err != nil {
				output.Results[p] = deletePresetOutput{PresetID: "", Error: "initializing provider: " + err.Error()}
				continue
			}
			err = providerObj.DeletePreset(presetID)
			if err != nil {
				output.Results[p] = deletePresetOutput{PresetID: "", Error: "deleting preset: " + err.Error()}
				continue
			}
			output.Results[p] = deletePresetOutput{PresetID: presetID, Error: ""}
		}
		err = s.db.DeletePresetMap(&db.PresetMap{Name: params.Name})
		if err != nil {
			output.Status = "error deleting presetmap: " + err.Error()
		} else {
			output.Status = "presetmap removed successfully"
		}
	}
	return &deletePresetResponse{
		baseResponse: baseResponse{
			payload: output,
			status:  http.StatusOK,
		},
	}
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
	var output newPresetOutputs
	var presetMap db.PresetMap

	respData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return newErrorResponse(err)
	}

	err = json.Unmarshal(respData, &input)
	if err != nil {
		return newErrorResponse(err)
	}

	presetMap.ProviderMapping = make(map[string]string)
	output.Results = make(map[string]newPresetOutput)
	for _, p := range input.Providers {
		providerFactory, err := provider.GetProviderFactory(p)
		if err != nil {
			output.Results[p] = newPresetOutput{PresetID: "", Error: "getting factory: " + err.Error()}
			continue
		}
		providerObj, err := providerFactory(s.config)
		if err != nil {
			output.Results[p] = newPresetOutput{PresetID: "", Error: "initializing provider: " + err.Error()}
			continue
		}
		presetID, err := providerObj.CreatePreset(input.Preset)
		if err != nil {
			output.Results[p] = newPresetOutput{PresetID: "", Error: "creating preset: " + err.Error()}
			continue
		}
		presetMap.ProviderMapping[p] = presetID
		output.Results[p] = newPresetOutput{PresetID: presetID, Error: ""}
	}

	output.PresetMap = ""
	if len(presetMap.ProviderMapping) > 0 {
		presetMap.Name = input.Preset.Name
		presetMap.OutputOpts.Extension = input.Preset.Container

		err = s.db.CreatePresetMap(&presetMap)
		if err == nil {
			output.PresetMap = presetMap.Name
		}
	}

	return &newPresetResponse{
		baseResponse: baseResponse{
			payload: output,
			status:  http.StatusOK,
		},
	}
}
