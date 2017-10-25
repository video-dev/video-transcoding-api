package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type H265CodecConfigurationService struct {
	RestService *RestService
}

const (
	H265CodecConfigurationEndpoint string = "encoding/configurations/video/h265"
)

func NewH265CodecConfigurationService(bitmovin *bitmovin.Bitmovin) *H265CodecConfigurationService {
	r := NewRestService(bitmovin)

	return &H265CodecConfigurationService{RestService: r}
}

func (s *H265CodecConfigurationService) Create(a *models.H265CodecConfiguration) (*models.H265CodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(H265CodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.H265CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *H265CodecConfigurationService) Retrieve(id string) (*models.H265CodecConfigurationResponse, error) {
	path := H265CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.H265CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *H265CodecConfigurationService) Delete(id string) (*models.H265CodecConfigurationResponse, error) {
	path := H265CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.H265CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *H265CodecConfigurationService) List(offset int64, limit int64) (*models.H265CodecConfigurationListResponse, error) {
	o, err := s.RestService.List(H265CodecConfigurationEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.H265CodecConfigurationListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *H265CodecConfigurationService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := H265CodecConfigurationEndpoint + "/" + id
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
