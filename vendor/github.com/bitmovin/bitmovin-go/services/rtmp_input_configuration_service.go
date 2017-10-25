package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type RTMPInputService struct {
	RestService *RestService
}

const (
	RTMPInputEndpoint string = "encoding/inputs/rtmp"
)

func NewRTMPInputService(bitmovin *bitmovin.Bitmovin) *RTMPInputService {
	r := NewRestService(bitmovin)

	return &RTMPInputService{RestService: r}
}

func (s *RTMPInputService) Retrieve(id string) (*models.RTMPInputResponse, error) {
	path := RTMPInputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.RTMPInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *RTMPInputService) List(offset int64, limit int64) (*models.RTMPInputListResponse, error) {
	o, err := s.RestService.List(RTMPInputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.RTMPInputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
