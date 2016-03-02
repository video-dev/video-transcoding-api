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
		// source media for the transcoding job.
		Source string

		// presets to use in the transcoding job.
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
	if len(p.Payload.Presets) == 0 {
		return errors.New("missing preset list from request")
	}
	return nil
}

// swagger:parameters getJob
type getTranscodeJobInput struct {
	// in: path
	// required: true
	JobID string `json:"jobId"`
}

func (p *getTranscodeJobInput) loadParams(paramsMap map[string]string) {
	p.JobID = paramsMap["jobId"]
}
