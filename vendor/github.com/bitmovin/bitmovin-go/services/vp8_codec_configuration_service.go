package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type VP8CodecConfigurationService struct {
	RestService *RestService
}

const (
	VP8CodecConfigurationEndpoint string = "encoding/configurations/video/vp8"
)

func NewVP8CodecConfigurationService(bitmovin *bitmovin.Bitmovin) *VP8CodecConfigurationService {
	r := NewRestService(bitmovin)

	return &VP8CodecConfigurationService{RestService: r}
}

func (s *VP8CodecConfigurationService) Create(a *models.VP8CodecConfiguration) (*models.VP8CodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(VP8CodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.VP8CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VP8CodecConfigurationService) Retrieve(id string) (*models.VP8CodecConfigurationResponse, error) {
	path := VP8CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.VP8CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VP8CodecConfigurationService) Delete(id string) (*models.VP8CodecConfigurationResponse, error) {
	path := VP8CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.VP8CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VP8CodecConfigurationService) List(offset int64, limit int64) (*models.VP8CodecConfigurationListResponse, error) {
	o, err := s.RestService.List(VP8CodecConfigurationEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.VP8CodecConfigurationListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VP8CodecConfigurationService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := VP8CodecConfigurationEndpoint + "/" + id
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
