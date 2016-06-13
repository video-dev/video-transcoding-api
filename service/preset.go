package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/NYTimes/gizmo/web"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
	"github.com/nytm/video-transcoding-api/swagger"
)

// swagger:route DELETE /presets/{name} presets deletePreset
//
// Deletes a preset by name.
//
//     Responses:
//       200: deletePresetOutputs
//       404: presetNotFound
//       500: genericError
func (s *TranscodingService) deletePreset(r *http.Request) swagger.GizmoJSONResponse {
	var output deletePresetOutputs
	var params getPresetMapInput
	params.loadParams(web.Vars(r))

	output.Results = make(map[string]deletePresetOutput)

	presetmap, err := s.db.GetPresetMap(params.Name)
	if err != nil {
		output.PresetMap = "couldn't retrieve: " + err.Error()
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
			output.PresetMap = "error: " + err.Error()
		} else {
			output.PresetMap = "removed successfully"
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
func (s *TranscodingService) newPreset(r *http.Request) swagger.GizmoJSONResponse {
	defer r.Body.Close()
	var input newPresetInput
	var output newPresetOutputs
	var presetMap db.PresetMap

	respData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return swagger.NewErrorResponse(err)
	}

	err = json.Unmarshal(respData, &input)
	if err != nil {
		return swagger.NewErrorResponse(err)
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
