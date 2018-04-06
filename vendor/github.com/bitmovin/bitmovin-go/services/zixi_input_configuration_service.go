package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type ZixiInputService struct {
	RestService *RestService
}

const (
	ZixiInputEndpoint string = "encoding/inputs/zixi"
)

func NewZixiInputService(bitmovin *bitmovin.Bitmovin) *ZixiInputService {
	r := NewRestService(bitmovin)
	return &ZixiInputService{RestService: r}
}

func (s *ZixiInputService) Create(a *models.ZixiInput) (*models.ZixiInputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(ZixiInputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.ZixiInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *ZixiInputService) Retrieve(id string) (*models.ZixiInputResponse, error) {
	path := ZixiInputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.ZixiInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *ZixiInputService) Delete(id string) (*models.ZixiInputResponse, error) {
	path := ZixiInputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.ZixiInputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *ZixiInputService) List(offset int64, limit int64) (*models.ZixiInputListResponse, error) {
	o, err := s.RestService.List(ZixiInputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.ZixiInputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *ZixiInputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := ZixiInputEndpoint + "/" + id
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
