package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type FilterService struct {
	RestService *RestService
}

const (
	DeinterlaceEndpoint string = "encoding/filters/deinterlace"
)

func NewFilterService(client *bitmovin.Bitmovin) *FilterService {
	return &FilterService{RestService: NewRestService(client)}
}

func (f *FilterService) CreateDeinterlacingFilter(filter *models.DeinterlacingFilter) (*models.DeinterlacingFilterResponse, error) {
	b, err := json.Marshal(filter)
	if err != nil {
		return nil, err
	}
	response, err := f.RestService.Create(DeinterlaceEndpoint, b)
	if err != nil {
		return nil, err
	}

	var result models.DeinterlacingFilterResponse
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
