package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type SmoothStreamingService struct {
	RestService *RestService
}

const (
	SmoothStreamingManifestEndpoint = "encoding/manifests/smooth"
)

func NewSmoothStreamingService(bitmovin *bitmovin.Bitmovin) *SmoothStreamingService {
	return &SmoothStreamingService{RestService: NewRestService(bitmovin)}
}

func (s *SmoothStreamingService) Create(a *models.SmoothStreamingManifest) (*models.SmoothStreamingManifestResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}

	o, err := s.RestService.Create(SmoothStreamingManifestEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.SmoothStreamingManifestResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *SmoothStreamingService) Retrieve(id string) (*models.SmoothStreamingManifestResponse, error) {
	path := SmoothStreamingManifestEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.SmoothStreamingManifestResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *SmoothStreamingService) Delete(id string) (*models.SmoothStreamingManifestResponse, error) {
	path := SmoothStreamingManifestEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.SmoothStreamingManifestResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *SmoothStreamingService) AddMp4Representation(manifestID string, a *models.SmoothStreamingMp4Representation) (*models.SmoothStreamingMp4RepresentationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := SmoothStreamingManifestEndpoint + "/" + manifestID + "/representations/mp4"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.SmoothStreamingMp4RepresentationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *SmoothStreamingService) AddContentProtection(manifestID string, a *models.SmoothStreamingContentProtection) (*models.SmoothStreamingContentProtectionResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := SmoothStreamingManifestEndpoint + "/" + manifestID + "/contentprotection"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.SmoothStreamingContentProtectionResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *SmoothStreamingService) Start(manifestID string) (*models.StartStopResponse, error) {
	path := SmoothStreamingManifestEndpoint + "/" + manifestID + "/start"
	o, err := s.RestService.Create(path, nil)
	if err != nil {
		return nil, err
	}
	var r models.StartStopResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *SmoothStreamingService) RetrieveStatus(manifestID string) (*models.StatusResponse, error) {
	path := SmoothStreamingManifestEndpoint + "/" + manifestID + "/status"
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.StatusResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
