package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type GCSOutputService struct {
	RestService *RestService
}

const (
	GCSOutputEndpoint string = "encoding/outputs/gcs"
)

func NewGCSOutputService(bitmovin *bitmovin.Bitmovin) *GCSOutputService {
	r := NewRestService(bitmovin)

	return &GCSOutputService{RestService: r}
}

func (s *GCSOutputService) Create(a *models.GCSOutput) (*models.GCSOutputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(GCSOutputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.GCSOutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GCSOutputService) Retrieve(id string) (*models.GCSOutputResponse, error) {
	path := GCSOutputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.GCSOutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GCSOutputService) Delete(id string) (*models.GCSOutputResponse, error) {
	path := GCSOutputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.GCSOutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GCSOutputService) List(offset int64, limit int64) (*models.GCSOutputListResponse, error) {
	o, err := s.RestService.List(GCSOutputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.GCSOutputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GCSOutputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := GCSOutputEndpoint + "/" + id
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
