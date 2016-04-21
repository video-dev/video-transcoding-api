package service

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/provider"
)

// swagger:parameters newPreset
type newPresetMapInput struct {
	// in: body
	// required: true
	Payload db.PresetMap
}

type newPresetInput struct {
	Providers []string        `json:"providers"`
	Preset    provider.Preset `json:"preset"`
}

// list of the results of the attempt to create a preset
// in each provider.
//
// swagger:response newPresetOutputs
type newPresetOutputs struct {
	// in: body
	// required: true
	Results   map[string]newPresetOutput
	PresetMap string
}

type newPresetOutput struct {
	PresetID string
	Error    string
}

// list of the results of the attempt to delete a preset
// in each provider.
//
// swagger:response deletePresetOutputs
type deletePresetOutputs struct {
	// in: body
	// required: true
	Results   map[string]deletePresetOutput `json:"results"`
	PresetMap string                        `json:"presetMap"`
}

type deletePresetOutput struct {
	PresetID string `json:"presetId"`
	Error    string `json:"error,omitempty"`
}

// Preset loads the input from the request body, validates them and returns the
// preset.
func (p *newPresetMapInput) PresetMap(body io.Reader) (db.PresetMap, error) {
	err := json.NewDecoder(body).Decode(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	err = validatePresetMap(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	return p.Payload, nil
}

// swagger:parameters getPreset deletePreset deletePresetMap
type getPresetMapInput struct {
	// in: path
	// required: true
	Name string `json:"name"`
}

func (p *getPresetMapInput) loadParams(paramsMap map[string]string) {
	p.Name = paramsMap["name"]
}

// swagger:parameters updatePreset
type updatePresetMapInput struct {
	// in: path
	// required: true
	Name string `json:"name"`

	// in: body
	// required: true
	Payload db.PresetMap

	newPresetMapInput
}

func (p *updatePresetMapInput) PresetMap(paramsMap map[string]string, body io.Reader) (db.PresetMap, error) {
	p.Name = paramsMap["name"]
	err := json.NewDecoder(body).Decode(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	p.Payload.Name = p.Name
	err = validatePresetMap(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	return p.Payload, nil
}

func validatePresetMap(p *db.PresetMap) error {
	if p.Name == "" {
		return errors.New("missing field name from the request")
	}
	if len(p.ProviderMapping) == 0 {
		return errors.New("missing field providerMapping from the request")
	}
	return nil
}
