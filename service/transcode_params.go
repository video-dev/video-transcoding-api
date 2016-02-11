package service

import (
	"errors"

	"github.com/nytm/video-transcoding-api/provider"
)

// swagger:parameters newJob
type newTranscodeParams struct {
	// source media for the transcode job.
	//
	// required: true
	// in: body
	Source string

	// profiles to be used. This parameter is exclusive with the list of
	// presets. One, and only one, of both should be provided.
	//
	// in: body
	Profiles []provider.Profile

	// presets to be used. this parameter is exclusive with the list of
	// profiles. One, and only one, of both should be provided.
	//
	// in: body
	Presets []string

	// provider to use in this job
	//
	// required: true
	// in: body
	Provider string
}

// ProviderFactory gets the factory of the provider after validating all
// parameters.
func (p *newTranscodeParams) ProviderFactory() (provider.Factory, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}
	return provider.GetProviderFactory(p.Provider)
}

func (p *newTranscodeParams) validate() error {
	if p.Provider == "" {
		return errors.New("missing provider from request")
	}
	if p.Source == "" {
		return errors.New("missing source media from request")
	}
	if len(p.Profiles) == 0 && len(p.Presets) == 0 {
		return errors.New("please specify either the list of presets or the list of profiles")
	}
	if len(p.Profiles) > 0 && len(p.Presets) > 0 {
		return errors.New("presets and profiles are mutually exclusive, please specify only one of them")
	}
	return nil
}
