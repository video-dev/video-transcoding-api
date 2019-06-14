package configuration

import (
	"github.com/NYTimes/video-transcoding-api/provider/bitmovinnewsdk/configuration/codec"

	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/pkg/errors"
)

// VP8Vorbis is a configuration service for content in this codec pair
type VP8Vorbis struct {
	api *bitmovin.BitmovinApi
}

// NewVP8Vorbis returns a service for managing VP8 / Vorbis configurations
func NewVP8Vorbis(api *bitmovin.BitmovinApi) *VP8Vorbis {
	return &VP8Vorbis{api: api}
}

// Create will create a new VP8 configuration based on a preset
func (c *VP8Vorbis) Create(preset db.Preset) (string, error) {
	audCfgID, err := codec.NewVorbis(c.api, preset.Audio.Bitrate)
	if err != nil {
		return "", err
	}

	vidCfgID, err := codec.NewVP8(c.api, preset, customDataWith(audCfgID, preset.Container))
	if err != nil {
		return "", err
	}

	return vidCfgID, nil
}

// Get retrieves audio / video configuration with a presetID
func (c *VP8Vorbis) Get(presetID string) (bool, Details, error) {
	vidCfg, err := c.api.Encoding.Configurations.Video.Vp8.Get(presetID)
	if err != nil {
		return false, Details{}, errors.Wrap(err, "retrieving configuration with presetID")
	}

	dataResp, err := c.api.Encoding.Configurations.Video.Vp8.Customdata.Get(vidCfg.Id)
	if err != nil {
		return false, Details{}, errors.Wrap(err, "retrieving custom data with config ID")
	}

	audCfgID, err := AudCfgIDFrom(dataResp.CustomData)
	if err != nil {
		return false, Details{}, err
	}

	audCfg, err := c.api.Encoding.Configurations.Audio.Vorbis.Get(audCfgID)
	if err != nil {
		return false, Details{}, errors.Wrapf(err, "getting the audio config with ID %q", audCfgID)
	}

	return true, Details{vidCfg, audCfg, dataResp.CustomData}, nil
}

// Delete removes the audio / video configurations
func (c *VP8Vorbis) Delete(presetID string) (found bool, e error) {
	vidCfg, err := c.api.Encoding.Configurations.Video.Vp8.Get(presetID)
	if err != nil {
		return found, errors.Wrap(err, "retrieving video configuration with presetID")
	}

	audCfgID, err := AudCfgIDFrom(vidCfg.CustomData)
	if err != nil {
		return found, err
	}

	audCfg, err := c.api.Encoding.Configurations.Audio.Vorbis.Get(audCfgID)
	if err != nil {
		return found, errors.Wrap(err, "retrieving audio configuration")
	}
	found = true

	_, err = c.api.Encoding.Configurations.Audio.Vorbis.Delete(audCfg.Id)
	if err != nil {
		return found, errors.Wrap(err, "removing the audio config")
	}

	_, err = c.api.Encoding.Configurations.Video.Vp8.Delete(vidCfg.Id)
	if err != nil {
		return found, errors.Wrap(err, "removing the video config")
	}

	return found, nil
}
