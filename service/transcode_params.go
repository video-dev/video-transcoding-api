package service

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/nytm/video-transcoding-api/provider"
)

const defaultStatusCallbackInterval = 5

// swagger:parameters newJob
type newTranscodeJobInput struct {
	// in: body
	// required: true
	Payload struct {
		// source media for the transcoding job.
		Source string `json:"source"`

		// presets to use in the transcoding job.
		Presets []string `json:"presets"`

		// provider to use in this job
		Provider string `json:"provider"`

		// provider Adaptive Streaming parameters
		StreamingParams provider.StreamingParams `json:"streamingParams,omitempty"`

		// if StatusCallbackURL is defined, this service will make a POST
		// request to it in the interval defined by StatusCallbackInterval
		// until the job is finished. The payload will be the same as the one
		// returned by a GET call to /jobs/<jobId>
		StatusCallbackURL string `json:"statusCallbackURL"`

		// defines the interval in seconds by which StatusCallbackURL is
		// called. If not defined, it's set to defaultStatusCallbackInterval.
		StatusCallbackInterval uint `json:"statusCallbackInterval"`

		// if CompletionCallbackURL is defined, this service will make a POST
		// request to it when the job is finished. The payload will be the same
		// as the one returned by a GET call to /jobs/<jobId> after the job is
		// done.
		CompletionCallbackURL string `json:"completionCallbackURL"`
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
	if p.Payload.StatusCallbackURL != "" && p.Payload.StatusCallbackInterval == 0 {
		p.Payload.StatusCallbackInterval = defaultStatusCallbackInterval
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
