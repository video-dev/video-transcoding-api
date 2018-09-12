package services

import (
	"encoding/json"
	"fmt"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

type AnalyticsService struct {
	RestService *RestService
}

const (
	path string = "analytics/queries/"
)

func NewAnalyticsService(bitmovin *bitmovin.Bitmovin) *AnalyticsService {
	r := NewRestService(bitmovin)

	return &AnalyticsService{RestService: r}
}

func (s *AnalyticsService) Count(a *models.Query) (*models.QueryResponse, error) {
	return s.doAnalytics(a, "count")
}

func (s *AnalyticsService) Sum(a *models.Query) (*models.QueryResponse, error) {
	return s.doAnalytics(a, "sum")
}

func (s *AnalyticsService) Avg(a *models.Query) (*models.QueryResponse, error) {
	return s.doAnalytics(a, "avg")
}

func (s *AnalyticsService) Min(a *models.Query) (*models.QueryResponse, error) {
	return s.doAnalytics(a, "min")
}

func (s *AnalyticsService) Max(a *models.Query) (*models.QueryResponse, error) {
	return s.doAnalytics(a, "max")
}

func (s *AnalyticsService) Stddev(a *models.Query) (*models.QueryResponse, error) {
	return s.doAnalytics(a, "stddev")
}

func (s *AnalyticsService) Percentile(a *models.PercentileQuery) (*models.QueryResponse, error) {
	marshaledQuery, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	p := fmt.Sprintf("%spercentile", path)
	output, err := s.RestService.Create(p, marshaledQuery)
	if err != nil {
		return nil, err
	}
	var response models.QueryResponse
	err = json.Unmarshal(output, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *AnalyticsService) Variance(a *models.Query) (*models.QueryResponse, error) {
	return s.doAnalytics(a, "variance")
}

func (s *AnalyticsService) Median(a *models.Query) (*models.QueryResponse, error) {
	return s.doAnalytics(a, "median")
}

func (s *AnalyticsService) doAnalytics(a *models.Query, method string) (*models.QueryResponse, error) {
	marshaledQuery, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	p := fmt.Sprintf("%s%s", path, method)
	output, err := s.RestService.Create(p, marshaledQuery)
	if err != nil {
		return nil, err
	}
	var response models.QueryResponse
	err = json.Unmarshal(output, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
