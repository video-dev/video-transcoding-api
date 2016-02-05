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
	return http.StatusOK, map[string]string{"presetId": preset.ID}, nil
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

func (s *TranscodingService) listPresets(r *http.Request) (int, interface{}, error) {
	presets, err := s.db.ListPresets()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	presetsMap := make(map[string]db.Preset, len(presets))
	for _, preset := range presets {
		presetsMap[preset.ID] = preset
	}
	return http.StatusOK, presetsMap, nil
}
