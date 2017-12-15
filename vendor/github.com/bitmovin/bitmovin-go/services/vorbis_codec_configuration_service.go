package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type VorbisCodecConfigurationService struct {
	RestService *RestService
}

const (
	// FIXME Double check endpoint
	VorbisCodecConfigurationEndpoint string = "encoding/configurations/audio/vorbis"
)

func NewVorbisCodecConfigurationService(bitmovin *bitmovin.Bitmovin) *VorbisCodecConfigurationService {
	r := NewRestService(bitmovin)

	return &VorbisCodecConfigurationService{RestService: r}
}

func (s *VorbisCodecConfigurationService) Create(a *models.VorbisCodecConfiguration) (*models.VorbisCodecConfigurationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}

	o, err := s.RestService.Create(VorbisCodecConfigurationEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.VorbisCodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VorbisCodecConfigurationService) Retrieve(id string) (*models.VorbisCodecConfigurationResponse, error) {
	path := VorbisCodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.VorbisCodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VorbisCodecConfigurationService) Delete(id string) (*models.VorbisCodecConfigurationResponse, error) {
	path := VorbisCodecConfigurationEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.VorbisCodecConfigurationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VorbisCodecConfigurationService) List(offset int64, limit int64) (*models.VorbisCodecConfigurationListResponse, error) {
	o, err := s.RestService.List(VorbisCodecConfigurationEndpoint, offset, limit)
	if err != nil {
		return nil, err
	}
	var r models.VorbisCodecConfigurationListResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *VorbisCodecConfigurationService) RetrieveCustomData(id string) (*models.CustomDataResponse, error) {
	path := VorbisCodecConfigurationEndpoint + "/" + id
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
