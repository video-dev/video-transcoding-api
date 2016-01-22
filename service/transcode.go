package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nytm/video-transcoding-api/provider"
)

type newTranscodeRequest struct {
	Source      string
	Destination string
	Profile     string
	Provider    string
}

func (s *TranscodingService) newTranscodeJob(r *http.Request) (int, interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	var reqObject newTranscodeRequest
	err := decoder.Decode(&reqObject)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("Error while parsing request: %s", err)
	}
	if reqObject.Provider == "" {
		return http.StatusBadRequest, nil, fmt.Errorf("Missing provider from request")
	}
	if reqObject.Source == "" {
		return http.StatusBadRequest, nil, fmt.Errorf("Missing source from request")
	}
	if reqObject.Destination == "" {
		return http.StatusBadRequest, nil, fmt.Errorf("Missing destination from request")
	}
	if reqObject.Profile == "" {
		return http.StatusBadRequest, nil, fmt.Errorf("Missing profile from request")
	}
	decoder = json.NewDecoder(strings.NewReader(reqObject.Profile))
	var reqProfile provider.Profile
	err = decoder.Decode(&reqProfile)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("Error while parsing profile in request: %s", err)
	}
	providerFactory := s.providers[reqObject.Provider]
	if providerFactory == nil {
		return http.StatusBadRequest, nil, fmt.Errorf("Unknown provider found in request: %s", reqObject.Provider)
	}
	providerObj, err := providerFactory(s.config)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(provider.InvalidConfigError); ok {
			statusCode = http.StatusBadRequest
		}
		return statusCode, nil, fmt.Errorf("Error initializing provider %s: %s", providerObj, err)
	}

	jobStatus, err := providerObj.Transcode(reqObject.Source, reqObject.Destination, reqProfile)
	jobStatus.ProviderName = reqObject.Provider
	return 200, jobStatus, nil
}

func (s *TranscodingService) getTranscodeJob(r *http.Request) (int, interface{}, error) {
	return 0, nil, nil
}
