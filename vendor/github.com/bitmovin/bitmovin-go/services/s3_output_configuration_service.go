package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type S3OutputService struct {
	RestService *RestService
}

const (
	S3OutputEndpoint string = "encoding/outputs/s3"
)

func NewS3OutputService(bitmovin *bitmovin.Bitmovin) *S3OutputService {
	r := NewRestService(bitmovin)

	return &S3OutputService{RestService: r}
}

func (s *S3OutputService) Create(a *models.S3Output) (*models.S3OutputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(S3OutputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.S3OutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *S3OutputService) Retrieve(id string) (*models.S3OutputResponse, error) {
	path := S3OutputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.S3OutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *S3OutputService) Delete(id string) (*models.S3OutputResponse, error) {
	path := S3OutputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.S3OutputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *S3OutputService) List(offset int64, limit int64) (*models.S3OutputListResponse, error) {
	o, err := s.RestService.List(S3OutputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.S3OutputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *S3OutputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := S3OutputEndpoint + "/" + id
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
