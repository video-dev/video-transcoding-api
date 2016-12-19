package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NYTimes/gizmo/web"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/NYTimes/video-transcoding-api/swagger"
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
			providerFactory, ierr := provider.GetProviderFactory(p)
			if ierr != nil {
				output.Results[p] = deletePresetOutput{PresetID: "", Error: "getting factory: " + ierr.Error()}
				continue
			}
			providerObj, ierr := providerFactory(s.config)
			if ierr != nil {
				output.Results[p] = deletePresetOutput{PresetID: "", Error: "initializing provider: " + ierr.Error()}
				continue
			}
			ierr = providerObj.DeletePreset(presetID)
			if ierr != nil {
				output.Results[p] = deletePresetOutput{PresetID: "", Error: "deleting preset: " + ierr.Error()}
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
//       400: invalidPreset
//       500: genericError
func (s *TranscodingService) newPreset(r *http.Request) swagger.GizmoJSONResponse {
	defer r.Body.Close()
	var input newPresetInput
	var output newPresetOutputs
	var presetMap *db.PresetMap
	var providers []string
	var shouldCreatePresetMap bool

	respData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return swagger.NewErrorResponse(err)
	}

	err = json.Unmarshal(respData, &input)
	if err != nil {
		return swagger.NewErrorResponse(err)
	}

	output.Results = make(map[string]newPresetOutput)

	// Sometimes we try to create a new preset in a new provider but we already
	// have the PresetMap stored. We want to update the PresetMap in such cases.
	presetMap, err = s.db.GetPresetMap(input.Preset.Name)
	if err == db.ErrPresetMapNotFound {
		presetMap = &db.PresetMap{Name: input.Preset.Name}
		presetMap.OutputOpts = input.OutputOptions
		presetMap.OutputOpts.Extension = input.Preset.Container
		presetMap.ProviderMapping = make(map[string]string)
		if err = presetMap.OutputOpts.Validate(); err != nil {
			return newInvalidPresetResponse(fmt.Errorf("invalid outputOptions: %s", err))
		}
		shouldCreatePresetMap = true
		providers = input.Providers
	} else if err != nil {
		return swagger.NewErrorResponse(err)
	} else {
		// If we already have a PresetMap for this preset, we just need to create the
		// preset on the providers that are not mapped yet.
		providers = s.getMissingProviders(input.Providers, presetMap.ProviderMapping)

		// We also want to add the existent presets on the result.
		for provider, presetID := range presetMap.ProviderMapping {
			output.Results[provider] = newPresetOutput{PresetID: presetID, Error: ""}
		}
	}

	for _, p := range providers {
		providerFactory, ierr := provider.GetProviderFactory(p)
		if ierr != nil {
			output.Results[p] = newPresetOutput{PresetID: "", Error: "getting factory: " + ierr.Error()}
			continue
		}
		providerObj, ierr := providerFactory(s.config)
		if ierr != nil {
			output.Results[p] = newPresetOutput{PresetID: "", Error: "initializing provider: " + ierr.Error()}
			continue
		}
		presetID, ierr := providerObj.CreatePreset(input.Preset)
		if ierr != nil {
			output.Results[p] = newPresetOutput{PresetID: "", Error: "creating preset: " + ierr.Error()}
			continue
		}
		presetMap.ProviderMapping[p] = presetID
		output.Results[p] = newPresetOutput{PresetID: presetID, Error: ""}
	}

	status := http.StatusOK
	if len(presetMap.ProviderMapping) > 0 {
		if shouldCreatePresetMap {
			err = s.db.CreatePresetMap(presetMap)
		} else {
			err = s.db.UpdatePresetMap(presetMap)
		}
		if err != nil {
			return newInvalidPresetResponse(fmt.Errorf("failed creating/updating presetmap after creating presets: %s", err))
		}
		output.PresetMap = presetMap.Name
	} else {
		status = http.StatusInternalServerError
	}

	return &newPresetResponse{
		baseResponse: baseResponse{
			payload: output,
			status:  status,
		},
	}
}

// getMissingProviders will check what providers already have a preset associated to it
// and return the missing ones. This method is used when a request to create a new preset
// is done but we already have a PresetMap stored locally.
func (s *TranscodingService) getMissingProviders(inputProviders []string, providerMapping map[string]string) []string {
	var missingProviders []string
	for _, provider := range inputProviders {
		if _, ok := providerMapping[provider]; !ok {
			missingProviders = append(missingProviders, provider)
		}
	}
	return missingProviders
}
