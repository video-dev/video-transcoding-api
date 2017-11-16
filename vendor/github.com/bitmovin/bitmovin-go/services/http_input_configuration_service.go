package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type HTTPInputService struct {
	RestService *RestService
}

const (
	HTTPInputEndpoint string = "encoding/inputs/http"
)

func NewHTTPInputService(bitmovin *bitmovin.Bitmovin) *HTTPInputService {
	r := NewRestService(bitmovin)

	return &HTTPInputService{RestService: r}
}

func (s *HTTPInputService) Create(a *models.HTTPInput) (*models.HTTPInputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(HTTPInputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.HTTPInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *HTTPInputService) Retrieve(id string) (*models.HTTPInputResponse, error) {
	path := HTTPInputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.HTTPInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *HTTPInputService) Delete(id string) (*models.HTTPInputResponse, error) {
	path := HTTPInputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.HTTPInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *HTTPInputService) List(offset int64, limit int64) (*models.HTTPInputListResponse, error) {
	o, err := s.RestService.List(HTTPInputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.HTTPInputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *HTTPInputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := HTTPInputEndpoint + "/" + id
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
