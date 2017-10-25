package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type AACCodecConfigurationService struct {
	RestService *RestService
}

const (
	AACCodecConfigurationEndpoint string = "encoding/configurations/audio/aac"
)

func NewAACCodecConfigurationService(bitmovin *bitmovin.Bitmovin) *AACCodecConfigurationService {
	r := NewRestService(bitmovin)

	return &AACCodecConfigurationService{RestService: r}
}

func (s *AACCodecConfigurationService) Create(a *models.AACCodecConfiguration) (*models.AACCodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}

	o, err := s.RestService.Create(AACCodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.AACCodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AACCodecConfigurationService) Retrieve(id string) (*models.AACCodecConfigurationResponse, error) {
	path := AACCodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.AACCodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AACCodecConfigurationService) Delete(id string) (*models.AACCodecConfigurationResponse, error) {
	path := AACCodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.AACCodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AACCodecConfigurationService) List(offset int64, limit int64) (*models.AACCodecConfigurationListResponse, error) {
	o, err := s.RestService.List(AACCodecConfigurationEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.AACCodecConfigurationListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AACCodecConfigurationService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := AACCodecConfigurationEndpoint + "/" + id
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
