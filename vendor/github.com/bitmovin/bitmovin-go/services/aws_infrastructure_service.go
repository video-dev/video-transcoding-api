package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type AWSInfrastructureService struct {
	RestService *RestService
}

const (
	AWSInfrastructureEndpoint string = "encoding/infrastructure/aws"
)

func NewAWSInfrastructureService(bitmovin *bitmovin.Bitmovin) *AWSInfrastructureService {
	r := NewRestService(bitmovin)
	return &AWSInfrastructureService{RestService: r}
}

func (s *AWSInfrastructureService) Create(i *models.CreateAWSInfrastructureRequest) (*models.AWSInfrastructureDetail, error) {
	b, err := json.Marshal(*i)
	if err != nil {
		return nil, err
	}
	responseBody, err := s.RestService.Create(AWSInfrastructureEndpoint, b)
	if err != nil {
		return nil, err
	}

	return MarshalSingleResponseAWS(responseBody)
}

func (s *AWSInfrastructureService) Retrieve(id string) (*models.AWSInfrastructureDetail, error) {
	path := AWSInfrastructureEndpoint + "/" + id

	responseBody, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	return MarshalSingleResponseAWS(responseBody)
}

func (s *AWSInfrastructureService) List(offset int64, limit int64) (*[]models.AWSInfrastructureDetail, error) {
	o, err := s.RestService.List(AWSInfrastructureEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.AWSInfrastructureListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r.Data.Result.Items, nil
}

func MarshalSingleResponseAWS(responseString []byte) (*models.AWSInfrastructureDetail, error) {
	var responseValue models.AWSInfrastructureResponse
	err := json.Unmarshal(responseString, &responseValue)
	if err != nil {
		return nil, err
	}
	return &responseValue.Data.Result, nil
}
