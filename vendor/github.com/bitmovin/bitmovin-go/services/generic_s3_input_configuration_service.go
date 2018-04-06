package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type GenericS3InputService struct {
	RestService *RestService
}

const (
	GenericS3InputEndpoint string = "encoding/inputs/generic-s3"
)

func NewGenericS3InputService(bitmovin *bitmovin.Bitmovin) *GenericS3InputService {
	r := NewRestService(bitmovin)

	return &GenericS3InputService{RestService: r}
}

func (s *GenericS3InputService) Create(a *models.GenericS3Input) (*models.GenericS3InputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(GenericS3InputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.GenericS3InputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GenericS3InputService) Retrieve(id string) (*models.GenericS3InputResponse, error) {
	path := GenericS3InputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.GenericS3InputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GenericS3InputService) Delete(id string) (*models.GenericS3InputResponse, error) {
	path := GenericS3InputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.GenericS3InputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GenericS3InputService) List(offset int64, limit int64) (*models.GenericS3InputListResponse, error) {
	o, err := s.RestService.List(GenericS3InputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.GenericS3InputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GenericS3InputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := GenericS3InputEndpoint + "/" + id
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
