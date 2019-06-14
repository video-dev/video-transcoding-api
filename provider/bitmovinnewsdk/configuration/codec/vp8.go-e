package codec

import (
	"strings"

	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/bitmovin/bitmovin-api-sdk-go/model"
	"github.com/cbsinteractive/video-transcoding-api/db"
	"github.com/pkg/errors"
)

// NewVP8 creates a VP8 codec configuration and returns its ID
func NewVP8(api *bitmovin.BitmovinApi, preset db.Preset, customData *map[string]map[string]interface{}) (string, error) {
	newVidCfg, err := vp8ConfigFrom(preset, customData)
	if err != nil {
		return "", errors.Wrap(err, "creating vp8 config object")
	}

	vidCfg, err := api.Encoding.Configurations.Video.Vp8.Create(newVidCfg)
	if err != nil {
		return "", errors.Wrap(err, "creating vp8 config with the API")
	}

	return vidCfg.Id, nil
}

func vp8ConfigFrom(preset db.Preset, customData *map[string]map[string]interface{}) (model.Vp8VideoConfiguration, error) {
	cfg := model.Vp8VideoConfiguration{}

	cfg.CustomData = customData

	cfg.Name = strings.ToLower(preset.Name)

	presetWidth := preset.Video.Width
	if presetWidth != "" {
		width, err := dimensionFrom(presetWidth)
		if err != nil {
			return model.Vp8VideoConfiguration{}, err
		}
		cfg.Width = width
	}

	presetHeight := preset.Video.Height
	if presetHeight != "" {
		height, err := dimensionFrom(presetHeight)
		if err != nil {
			return model.Vp8VideoConfiguration{}, err
		}
		cfg.Height = height
	}

	bitrate, err := bitrateFrom(preset.Video.Bitrate)
	if err != nil {
		return model.Vp8VideoConfiguration{}, err
	}
	cfg.Bitrate = bitrate

	return cfg, nil
}
