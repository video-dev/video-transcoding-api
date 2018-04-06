package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type GenericS3OutputService struct {
	RestService *RestService
}

const (
	GenericS3OutputEndpoint string = "encoding/outputs/generic-s3"
)

func NewGenericS3OutputService(bitmovin *bitmovin.Bitmovin) *GenericS3OutputService {
	r := NewRestService(bitmovin)

	return &GenericS3OutputService{RestService: r}
}

func (s *GenericS3OutputService) Create(a *models.GenericS3Output) (*models.GenericS3OutputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(GenericS3OutputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.GenericS3OutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GenericS3OutputService) Retrieve(id string) (*models.GenericS3OutputResponse, error) {
	path := GenericS3OutputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.GenericS3OutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GenericS3OutputService) Delete(id string) (*models.GenericS3OutputResponse, error) {
	path := GenericS3OutputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.GenericS3OutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GenericS3OutputService) List(offset int64, limit int64) (*models.GenericS3OutputListResponse, error) {
	o, err := s.RestService.List(GenericS3OutputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.GenericS3OutputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *GenericS3OutputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := GenericS3OutputEndpoint + "/" + id
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
