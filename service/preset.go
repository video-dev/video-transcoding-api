package service

import (
	"encoding/json"
	"net/http"

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
