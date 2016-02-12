package service

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nytm/video-transcoding-api/db"
)

// swagger:parameters newPreset
type newPresetInput struct {
	// in: body
	// required: true
	Payload db.Preset
}

// Preset loads the input from the request body, validates them and returns the
// preset.
func (p *newPresetInput) Preset(body io.Reader) (db.Preset, error) {
	err := json.NewDecoder(body).Decode(&p.Payload)
	if err != nil {
		return p.Payload, err
	}
	if field, valid := p.validate(); !valid {
		return p.Payload, fmt.Errorf("missing field %s from the request", field)
	}
	return p.Payload, nil
}

func (p *newPresetInput) validate() (fieldName string, valid bool) {
	if p.Payload.Name == "" {
		return "name", false
	}
	if len(p.Payload.ProviderMapping) == 0 {
		return "providerMapping", false
	}
	return "", true
}

// swagger:parameters getPreset deletePreset
type getPresetParams struct {
	// in: path
	// required: true
	Name string
}

func (p *getPresetParams) loadParams(paramsMap map[string]string) {
	p.Name = paramsMap["name"]
}
