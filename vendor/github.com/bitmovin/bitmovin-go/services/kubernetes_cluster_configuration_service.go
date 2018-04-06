package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

const (
	KubernetesConfigurationEndpoint string = "encoding/infrastructure/kubernetes"
)

type KubernetesClusterConfigurationService struct {
	RestService *RestService
}

func NewKubernetesClusterConfigurationService(bitmovin *bitmovin.Bitmovin) *KubernetesClusterConfigurationService {
	r := NewRestService(bitmovin)
	return &KubernetesClusterConfigurationService{RestService: r}
}

func (s *KubernetesClusterConfigurationService) Upsert(ID string, configRequest *models.KubernetesClusterConfigurationRequest) (*models.KubernetesClusterConfigurationResponse, error) {
	b, err := json.Marshal(*configRequest)
	if err != nil {
		return nil, err
	}

	path := KubernetesConfigurationEndpoint + "/" + ID + "/" + "configuration"
	responseBody, err := s.RestService.Update(path, b)
	if err != nil {
		return nil, err
	}

	var responseValue models.KubernetesClusterConfigurationResponse
	err = json.Unmarshal(responseBody, &responseValue)
	if err != nil {
		return nil, err
	}

	return &responseValue, nil
}

func (s *KubernetesClusterConfigurationService) Retrieve(ID string) (*models.KubernetesClusterConfigurationDetail, error) {
	path := KubernetesConfigurationEndpoint + "/" + ID + "/" + "configuration"

	responseBody, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}

	var responseValue models.KubernetesClusterConfigurationResponse
	err = json.Unmarshal(responseBody, &responseValue)
	if err != nil {
		return nil, err
	}

	return &responseValue.Data.Result, nil
}
