package configuration

import (
	"github.com/NYTimes/video-transcoding-api/provider/bitmovinnewsdk/configuration/codec"

	"github.com/bitmovin/bitmovin-api-sdk-go"

	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/pkg/errors"
)

// H264AAC is a configuration service for content in this codec pair
type H264AAC struct {
	api *bitmovin.BitmovinApi
}

// NewH264AAC returns a service for managing H264 / AAC configurations
func NewH264AAC(api *bitmovin.BitmovinApi) *H264AAC {
	return &H264AAC{api: api}
}

// Create will create a new H264AAC configuration based on a preset
func (c *H264AAC) Create(preset db.Preset) (string, error) {
	audCfgID, err := codec.NewAAC(c.api, preset.Audio.Bitrate)
	if err != nil {
		return "", err
	}

	vidCfgID, err := codec.NewH264(c.api, preset, customDataWith(audCfgID, preset.Container))
	if err != nil {
		return "", err
	}

	return vidCfgID, nil
}

// Get retrieves audio / video configuration with a presetID
// the function will return a boolean indicating whether the video
// configuration was found, a config object and an optional error
func (c *H264AAC) Get(presetID string) (bool, Details, error) {
	vidCfg, err := c.api.Encoding.Configurations.Video.H264.Get(presetID)
	if err != nil {
		return false, Details{}, errors.Wrap(err, "retrieving configuration with config ID")
	}

	dataResp, err := c.api.Encoding.Configurations.Video.H264.Customdata.Get(vidCfg.Id)
	if err != nil {
		return false, Details{}, errors.Wrap(err, "retrieving custom data with config ID")
	}

	audCfgID, err := AudCfgIDFrom(dataResp.CustomData)
	if err != nil {
		return false, Details{}, err
	}

	audCfg, err := c.api.Encoding.Configurations.Audio.Aac.Get(audCfgID)
	if err != nil {
		return false, Details{}, errors.Wrapf(err, "getting the audio configuration with ID %q", audCfgID)
	}

	return true, Details{vidCfg, audCfg, dataResp.CustomData}, nil
}

// Delete removes the audio / video configurations
func (c *H264AAC) Delete(presetID string) (found bool, e error) {
	vidCfg, err := c.api.Encoding.Configurations.Video.H264.Get(presetID)
	if err != nil {
		return found, errors.Wrap(err, "retrieving video configuration with presetID")
	}

	audCfgID, err := AudCfgIDFrom(vidCfg.CustomData)
	if err != nil {
		return found, err
	}

	audCfg, err := c.api.Encoding.Configurations.Audio.Aac.Get(audCfgID)
	if err != nil {
		return found, errors.Wrap(err, "retrieving audio configuration")
	}
	found = true

	_, err = c.api.Encoding.Configurations.Audio.Aac.Delete(audCfg.Id)
	if err != nil {
		return found, errors.Wrap(err, "removing the audio config")
	}

	_, err = c.api.Encoding.Configurations.Video.H264.Delete(vidCfg.Id)
	if err != nil {
		return found, errors.Wrap(err, "removing the video config")
	}

	return found, nil
}
