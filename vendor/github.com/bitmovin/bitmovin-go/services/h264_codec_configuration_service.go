package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type H264CodecConfigurationService struct {
	RestService *RestService
}

const (
	H264CodecConfigurationEndpoint string = "encoding/configurations/video/h264"
)

func NewH264CodecConfigurationService(bitmovin *bitmovin.Bitmovin) *H264CodecConfigurationService {
	r := NewRestService(bitmovin)

	return &H264CodecConfigurationService{RestService: r}
}

func (s *H264CodecConfigurationService) Create(a *models.H264CodecConfiguration) (*models.H264CodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(H264CodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.H264CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *H264CodecConfigurationService) Retrieve(id string) (*models.H264CodecConfigurationResponse, error) {
	path := H264CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.H264CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *H264CodecConfigurationService) Delete(id string) (*models.H264CodecConfigurationResponse, error) {
	path := H264CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.H264CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *H264CodecConfigurationService) List(offset int64, limit int64) (*models.H264CodecConfigurationListResponse, error) {
	o, err := s.RestService.List(H264CodecConfigurationEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.H264CodecConfigurationListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *H264CodecConfigurationService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := H264CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.RetrieveCustomData(path)
	if err != nil {
		return nil, err
	}
	var r models.CustomDataResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
