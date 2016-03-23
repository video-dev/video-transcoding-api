package service

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// swagger:parameters newPreset
type newPresetInput struct {
	// in: body
	// required: true
	Payload db.Preset
}

type newPresetInput2 struct {
	Providers []string        `json:"providers"`
	Preset    provider.Preset `json:"preset"`
}

// Preset loads the input from the request body, validates them and returns the
// preset.
func (p *newPresetInput) Preset(body io.Reader) (db.Preset, error) {
	err := json.NewDecoder(body).Decode(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	err = validatePreset(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	return p.Payload, nil
}

// swagger:parameters getPreset deletePreset
type getPresetInput struct {
	// in: path
	// required: true
	Name string `json:"name"`
}

func (p *getPresetInput) loadParams(paramsMap map[string]string) {
	p.Name = paramsMap["name"]
}

// swagger:parameters updatePreset
type updatePresetInput struct {
	// in: path
	// required: true
	Name string `json:"name"`

	// in: body
	// required: true
	Payload db.Preset

	newPresetInput
}

func (p *updatePresetInput) Preset(paramsMap map[string]string, body io.Reader) (db.Preset, error) {
	p.Name = paramsMap["name"]
	err := json.NewDecoder(body).Decode(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	p.Payload.Name = p.Name
	err = validatePreset(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	return p.Payload, nil
}

func validatePreset(p *db.Preset) error {
	if p.Name == "" {
		return errors.New("missing field name from the request")
	}
	if len(p.ProviderMapping) == 0 {
		return errors.New("missing field providerMapping from the request")
	}
	return nil
}
