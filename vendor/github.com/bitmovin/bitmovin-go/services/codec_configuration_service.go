package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type CodecConfigurationService struct {
	RestService *RestService
}

const (
	CodecConfigurationEndpoint string = "configurations"
)

func NewCodecConfigurationService(bitmovin *bitmovin.Bitmovin) *CodecConfigurationService {
	// FIXME correct endpoint?
	r := NewRestService(bitmovin)

	return &CodecConfigurationService{RestService: r}
}

func (s *CodecConfigurationService) List(offset int64, limit int64) (*models.CodecConfigurationListResponse, error) {
	o, err := s.RestService.List(CodecConfigurationEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.CodecConfigurationListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
