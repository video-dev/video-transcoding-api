package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nytm/video-transcoding-api/db"
)

func (s *TranscodingService) newPreset(r *http.Request) (int, interface{}, error) {
	var input map[string]string
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	preset := db.Preset{ProviderMapping: input}
	err = s.db.SavePreset(&preset)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, preset, nil
}

func (s *TranscodingService) getPreset(r *http.Request) (int, interface{}, error) {
	presetID := mux.Vars(r)["presetId"]
	preset, err := s.db.GetPreset(presetID)
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusInternalServerError
		if err == db.ErrPresetNotFound {
			statusCode = http.StatusNotFound
		}
	}
	return statusCode, preset, err
}

func (s *TranscodingService) deletePreset(r *http.Request) (int, interface{}, error) {
	presetID := mux.Vars(r)["presetId"]
	err := s.db.DeletePreset(&db.Preset{ID: presetID})
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusInternalServerError
		if err == db.ErrPresetNotFound {
			statusCode = http.StatusNotFound
		}
	}
	return statusCode, nil, err
}

// swagger:route GET /presets listPresets presets
//
// List available presets on the API.
//
//     Responses:
//       200: listPresets
func (s *TranscodingService) listPresets(r *http.Request) (int, interface{}, error) {
	presets, err := s.db.ListPresets()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	resp := listPresetsResponse{PresetsMap: make(map[string]db.Preset, len(presets))}
	for _, preset := range presets {
		resp.PresetsMap[preset.ID] = preset
	}
	return http.StatusOK, resp.PresetsMap, nil
}

// swagger:response listPresets
type listPresetsResponse struct {
	// in: body
	PresetsMap map[string]db.Preset
}
