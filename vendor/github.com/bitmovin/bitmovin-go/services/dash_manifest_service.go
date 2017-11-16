package services

import (
	"encoding/json"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type DashManifestService struct {
	RestService *RestService
}

const (
	DashManifestEndpoint string = "encoding/manifests/dash"
)

func NewDashManifestService(bitmovin *bitmovin.Bitmovin) *DashManifestService {
	r := NewRestService(bitmovin)

	return &DashManifestService{RestService: r}
}

func (s *DashManifestService) Create(a *models.DashManifest) (*models.DashManifestResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}

	o, err := s.RestService.Create(DashManifestEndpoint, b)
	if err != nil {
		return nil, err
	}
	var r models.DashManifestResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) Retrieve(id string) (*models.DashManifestResponse, error) {
	path := DashManifestEndpoint + "/" + id
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.DashManifestResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) Delete(id string) (*models.DashManifestResponse, error) {
	path := DashManifestEndpoint + "/" + id
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.DashManifestResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) AddPeriod(manifestID string, a *models.Period) (*models.PeriodResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.PeriodResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) RetrievePeriod(manifestID string, streamID string) (*models.PeriodResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + streamID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.PeriodResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) DeletePeriod(manifestID string, streamID string) (*models.PeriodResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + streamID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.PeriodResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) AddAudioAdaptationSet(manifestID string, periodID string, a *models.AudioAdaptationSet) (*models.AudioAdaptationSetResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/audio"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.AudioAdaptationSetResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) RetrieveAudioAdaptationSet(manifestID string, periodID string, adaptationSetID string, a *models.AudioAdaptationSet) (*models.AudioAdaptationSetResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/audio/" + adaptationSetID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.AudioAdaptationSetResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) DeleteAudioAdaptationSet(manifestID string, periodID string, adaptationSetID string, a *models.AudioAdaptationSet) (*models.AudioAdaptationSetResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/audio/" + adaptationSetID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.AudioAdaptationSetResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) AddVideoAdaptationSet(manifestID string, periodID string, a *models.VideoAdaptationSet) (*models.VideoAdaptationSetResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/video"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.VideoAdaptationSetResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) RetrieveVideoAdaptationSet(manifestID string, periodID string, adaptationSetID string, a *models.VideoAdaptationSet) (*models.VideoAdaptationSetResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/video/" + adaptationSetID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.VideoAdaptationSetResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) DeleteVideoAdaptationSet(manifestID string, periodID string, adaptationSetID string, a *models.VideoAdaptationSet) (*models.VideoAdaptationSetResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/video/" + adaptationSetID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.VideoAdaptationSetResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) AddFMP4Representation(manifestID string, periodID string, adaptationSetID string, a *models.FMP4Representation) (*models.FMP4RepresentationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/" + adaptationSetID + "/representations/fmp4"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.FMP4RepresentationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) RetrieveFMP4Representation(manifestID string, periodID string, adaptationSetID string, representationID string) (*models.FMP4RepresentationResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/" + adaptationSetID + "/representations/fmp4/" + representationID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.FMP4RepresentationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) DeleteFMP4Representation(manifestID string, periodID string, adaptationSetID string, representationID string) (*models.FMP4RepresentationResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/" + adaptationSetID + "/representations/fmp4/" + representationID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.FMP4RepresentationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) AddDrmFMP4Representation(manifestID string, periodID string, adaptationSetID string, a *models.DrmFMP4Representation) (*models.DrmFMP4RepresentationResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/" + adaptationSetID + "/representations/fmp4/drm"
	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.DrmFMP4RepresentationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) RetrieveDrmFMP4Representation(manifestID string, periodID string, adaptationSetID string, representationID string) (*models.DrmFMP4RepresentationResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/" + adaptationSetID + "/representations/fmp4/drm/" + representationID
	o, err := s.RestService.Retrieve(path)
	if err != nil {
		return nil, err
	}
	var r models.DrmFMP4RepresentationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) DeleteDrmFMP4Representation(manifestID string, periodID string, adaptationSetID string, representationID string) (*models.DrmFMP4RepresentationResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/" + "periods" + "/" + periodID + "/adaptationsets/" + adaptationSetID + "/representations/fmp4/drm/" + representationID
	o, err := s.RestService.Delete(path)
	if err != nil {
		return nil, err
	}
	var r models.DrmFMP4RepresentationResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *DashManifestService) Start(manifestID string) (*models.StartStopResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/start"
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

func (s *DashManifestService) RetrieveStatus(manifestID string) (*models.StatusResponse, error) {
	path := DashManifestEndpoint + "/" + manifestID + "/status"
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

func (s *DashManifestService) AddContentProtectionToAdaptationSet(manifestID string, periodID string, adaptationSetID string, a *models.AdaptationSetContentProtection) (*models.AdaptationSetContentProtectionResponse, error) {
	b, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}

	path := DashManifestEndpoint + "/" + manifestID + "/periods/" + periodID + "/adaptationsets/" + adaptationSetID + "/contentprotection"

	o, err := s.RestService.Create(path, b)
	if err != nil {
		return nil, err
	}
	var r models.AdaptationSetContentProtectionResponse
	err = json.Unmarshal(o, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
