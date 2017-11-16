package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type HTTPSInputService struct {
	RestService *RestService
}

const (
	HTTPSInputEndpoint string = "encoding/inputs/https"
)

func NewHTTPSInputService(bitmovin *bitmovin.Bitmovin) *HTTPSInputService {
	r := NewRestService(bitmovin)

	return &HTTPSInputService{RestService: r}
}

func (s *HTTPSInputService) Create(a *models.HTTPSInput) (*models.HTTPSInputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(b))
	o, err := s.RestService.Create(HTTPSInputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.HTTPSInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *HTTPSInputService) Retrieve(id string) (*models.HTTPSInputResponse, error) {
	path := HTTPSInputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.HTTPSInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *HTTPSInputService) Delete(id string) (*models.HTTPSInputResponse, error) {
	path := HTTPSInputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.HTTPSInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *HTTPSInputService) List(offset int64, limit int64) (*models.HTTPSInputListResponse, error) {
	o, err := s.RestService.List(HTTPSInputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.HTTPSInputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *HTTPSInputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := HTTPSInputEndpoint + "/" + id
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
