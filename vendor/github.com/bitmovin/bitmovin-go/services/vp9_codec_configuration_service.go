package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type VP9CodecConfigurationService struct {
	RestService *RestService
}

const (
	VP9CodecConfigurationEndpoint string = "encoding/configurations/video/vp9"
)

func NewVP9CodecConfigurationService(bitmovin *bitmovin.Bitmovin) *VP9CodecConfigurationService {
	r := NewRestService(bitmovin)

	return &VP9CodecConfigurationService{RestService: r}
}

func (s *VP9CodecConfigurationService) Create(a *models.VP9CodecConfiguration) (*models.VP9CodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(VP9CodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.VP9CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VP9CodecConfigurationService) Retrieve(id string) (*models.VP9CodecConfigurationResponse, error) {
	path := VP9CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.VP9CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VP9CodecConfigurationService) Delete(id string) (*models.VP9CodecConfigurationResponse, error) {
	path := VP9CodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.VP9CodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VP9CodecConfigurationService) List(offset int64, limit int64) (*models.VP9CodecConfigurationListResponse, error) {
	o, err := s.RestService.List(VP9CodecConfigurationEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.VP9CodecConfigurationListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VP9CodecConfigurationService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := VP9CodecConfigurationEndpoint + "/" + id
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
