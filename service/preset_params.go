package service

import (
	"github.com/nytm/video-transcoding-api/db"
)

// swagger:parameters newPreset
type newPresetParams struct {
	// name of the preset.
	//
	// in: body
	// required: true
	Name string `json:"name"`

	// the mapping of the provider name to the id of the preset in the
	// provider.
	//
	// in: body
	// required: true
	ProviderMapping map[string]string `json:"providerMapping"`
}

func (p *newPresetParams) Preset() db.Preset {
	return db.Preset{Name: p.Name, ProviderMapping: p.ProviderMapping}
}

func (p *newPresetParams) Validate() (fieldName string, valid bool) {
	if p.Name == "" {
		return "name", false
	}
	if len(p.ProviderMapping) == 0 {
		return "providerMapping", false
	}
	return "", true
}
