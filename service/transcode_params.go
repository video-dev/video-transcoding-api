package service

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/nytm/video-transcoding-api/provider"
)

// swagger:parameters newJob
type newTranscodeJobInput struct {
	// in: body
	// required: true
	Payload struct {
		// source media for the transcode job.
		Source string

		// profiles to be used. This parameter is exclusive with the list of
		// presets. One, and only one, of both should be provided.
		Profiles []provider.Profile

		// presets to be used. this parameter is exclusive with the list of
		// profiles. One, and only one, of both should be provided.
		Presets []string

		// provider to use in this job
		Provider string
	}
}

// ProviderFactory loads and validates the parameters, and then returns the
// provider factory.
func (p *newTranscodeJobInput) ProviderFactory(body io.Reader) (provider.Factory, error) {
	err := p.loadParams(body)
	if err != nil {
		return nil, err
	}
	err = p.validate()
	if err != nil {
		return nil, err
	}
	return provider.GetProviderFactory(p.Payload.Provider)
}

func (p *newTranscodeJobInput) loadParams(body io.Reader) error {
	return json.NewDecoder(body).Decode(&p.Payload)
}

func (p *newTranscodeJobInput) validate() error {
	if p.Payload.Provider == "" {
		return errors.New("missing provider from request")
	}
	if p.Payload.Source == "" {
		return errors.New("missing source media from request")
	}
	if len(p.Payload.Profiles) == 0 && len(p.Payload.Presets) == 0 {
		return errors.New("please specify either the list of presets or the list of profiles")
	}
	if len(p.Payload.Profiles) > 0 && len(p.Payload.Presets) > 0 {
		return errors.New("presets and profiles are mutually exclusive, please specify only one of them")
	}
	return nil
}

// swagger:parameters getJob
type getTranscodeJobParams struct {
	// in: path
	// required: true
	JobID string
}

func (p *getTranscodeJobParams) loadParams(paramsMap map[string]string) {
	p.JobID = paramsMap["jobId"]
}
