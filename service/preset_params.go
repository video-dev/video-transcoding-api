package service

import (
	"github.com/video-dev/video-transcoding-api/db"
)

type newPresetInput struct {
	Providers     []string         `json:"providers"`
	Preset        db.Preset        `json:"preset"`
	OutputOptions db.OutputOptions `json:"outputOptions"`
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
