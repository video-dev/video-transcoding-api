package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type InfrastructureService struct {
	RestService *RestService
}

const (
	InfrastructureEndpoint string = "encoding/infrastructure/kubernetes"
)

func NewInfrastructureService(bitmovin *bitmovin.Bitmovin) *InfrastructureService {
	r := NewRestService(bitmovin)
	return &InfrastructureService{RestService: r}
}

func (s *InfrastructureService) Create(i *models.CreateInfrastructureRequest) (*models.InfrastructureDetail, error) {
	b, err := json.Marshal(*i)
	if err != nil {
		return nil, err
	}
	responseBody, err := s.RestService.Create(InfrastructureEndpoint, b)
	if err != nil {
		return nil, err
	}

	return MarshalSingleResponse(responseBody)
}

func (s *InfrastructureService) Retrieve(id string) (*models.InfrastructureDetail, error) {
	path := InfrastructureEndpoint + "/" + id

	responseBody, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	return MarshalSingleResponse(responseBody)
}

func (s *InfrastructureService) List(offset int64, limit int64) (*[]models.InfrastructureDetail, error) {
	o, err := s.RestService.List(InfrastructureEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.InfrastructureListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r.Data.Result.Items, nil
}

func MarshalSingleResponse(responseString []byte) (*models.InfrastructureDetail, error) {
	var responseValue models.InfrastructureResponse
	err := json.Unmarshal(responseString, &responseValue)
	if err != nil {
		return nil, err
	}
	return &responseValue.Data.Result, nil
}
