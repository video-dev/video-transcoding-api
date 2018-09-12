package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
	"github.com/bitmovin/bitmovin-go/models"
)

const (
	AWSInfrastructureRegionSettingsConfigurationEndpoint string = "encoding/infrastructure/aws"
)

type AWSInfrastructureRegionSettingsService struct {
	RestService *RestService
}

func NewAWSInfrastructureRegionSettingsService(bitmovin *bitmovin.Bitmovin) *AWSInfrastructureRegionSettingsService {
	r := NewRestService(bitmovin)
	return &AWSInfrastructureRegionSettingsService{RestService: r}
}

func (s *AWSInfrastructureRegionSettingsService) Create(infrastructureID string, region bitmovintypes.AWSCloudRegion, i *models.CreateAWSInfrastructureRegionSettingsRequest) (*models.AWSInfrastructureRegionSettingsDetail, error) {
	path := AWSInfrastructureRegionSettingsConfigurationEndpoint + "/" + infrastructureID + "/" + "regions/" + string(region)

	b, err := json.Marshal(*i)
	if err != nil {
		return nil, err
	}
	responseBody, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}

	return MarshalSingleResponseAWSRegionSetting(responseBody)
}

func (s *AWSInfrastructureRegionSettingsService) Retrieve(infrastructureID string, region bitmovintypes.AWSCloudRegion) (*models.AWSInfrastructureRegionSettingsDetail, error) {
	path := AWSInfrastructureRegionSettingsConfigurationEndpoint + "/" + infrastructureID + "/" + "regions/" + string(region)

	responseBody, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}

	var responseValue models.AWSInfrastructureRegionSettingsResponse
	err = json.Unmarshal(responseBody, &responseValue)
	if err != nil {
		return nil, err
	}

	return &responseValue.Data.Result, nil
}

func (s *AWSInfrastructureRegionSettingsService) List(infrastructureID string, offset int64, limit int64) (*[]models.AWSInfrastructureRegionSettingsDetail, error) {
	path := AWSInfrastructureRegionSettingsConfigurationEndpoint + "/" + infrastructureID + "/" + "regions"

	o, err := s.RestService.List(path, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.AWSInfrastructureRegionSettingsListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r.Data.Result.Items, nil
}

func (s *AWSInfrastructureRegionSettingsService) Delete(infrastructureID string, region bitmovintypes.AWSCloudRegion) (*models.AWSInfrastructureRegionSettingsResponse, error) {
	path := AWSInfrastructureRegionSettingsConfigurationEndpoint + "/" + infrastructureID + "/" + "regions/" + string(region)
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.AWSInfrastructureRegionSettingsResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func MarshalSingleResponseAWSRegionSetting(responseString []byte) (*models.AWSInfrastructureRegionSettingsDetail, error) {
	var responseValue models.AWSInfrastructureRegionSettingsResponse
	err := json.Unmarshal(responseString, &responseValue)
	if err != nil {
		return nil, err
	}
	return &responseValue.Data.Result, nil
}
