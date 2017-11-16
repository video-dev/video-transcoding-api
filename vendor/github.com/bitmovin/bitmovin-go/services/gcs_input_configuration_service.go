package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type GcsInputService struct {
	RestService *RestService
}

const (
	GcsInputEndpoint string = "encoding/inputs/gcs"
)

func NewGCSInputService(bitmovin *bitmovin.Bitmovin) *GcsInputService {
	r := NewRestService(bitmovin)
	return &GcsInputService{RestService: r}
}

func (s *GcsInputService) Create(a *models.GCSInput) (*models.GCSInputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(GcsInputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.GCSInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GcsInputService) Retrieve(id string) (*models.GCSInputResponse, error) {
	path := GcsInputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.GCSInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GcsInputService) Delete(id string) (*models.GCSInputResponse, error) {
	path := GcsInputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.GCSInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GcsInputService) List(offset int64, limit int64) (*models.GCSInputListResponse, error) {
	o, err := s.RestService.List(GcsInputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.GCSInputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GcsInputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := GcsInputEndpoint + "/" + id
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
