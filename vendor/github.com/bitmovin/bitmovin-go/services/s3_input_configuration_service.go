package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type S3InputService struct {
	RestService *RestService
}

const (
	S3InputEndpoint string = "encoding/inputs/s3"
)

func NewS3InputService(bitmovin *bitmovin.Bitmovin) *S3InputService {
	r := NewRestService(bitmovin)

	return &S3InputService{RestService: r}
}

func (s *S3InputService) Create(a *models.S3Input) (*models.S3InputResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	o, err := s.RestService.Create(S3InputEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.S3InputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *S3InputService) Retrieve(id string) (*models.S3InputResponse, error) {
	path := S3InputEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.S3InputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *S3InputService) Delete(id string) (*models.S3InputResponse, error) {
	path := S3InputEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.S3InputResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *S3InputService) List(offset int64, limit int64) (*models.S3InputListResponse, error) {
	o, err := s.RestService.List(S3InputEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.S3InputListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *S3InputService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := S3InputEndpoint + "/" + id
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
