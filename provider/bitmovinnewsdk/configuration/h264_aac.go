package configuration

import (
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider/bitmovinnewsdk/configuration/codec"
	"github.com/NYTimes/video-transcoding-api/provider/bitmovinnewsdk/types"
	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/bitmovin/bitmovin-api-sdk-go/model"
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
func (c *H264AAC) Get(cfgID string) (bool, Details, error) {
	vidCfg, customData, err := c.vidConfigWithCustomDataFrom(cfgID)
	if err != nil {
		return false, Details{}, err
	}

	audCfgID, err := AudCfgIDFrom(customData)
	if err != nil {
		return false, Details{}, err
	}

	audCfg, err := c.api.Encoding.Configurations.Audio.Aac.Get(audCfgID)
	if err != nil {
		return false, Details{}, errors.Wrapf(err, "getting the audio configuration with ID %q", audCfgID)
	}

	return true, Details{vidCfg, audCfg, customData}, nil
}

// Delete removes the audio / video configurations
func (c *H264AAC) Delete(cfgID string) (found bool, e error) {
	customData, err := c.vidCustomDataFrom(cfgID)
	if err != nil {
		return found, err
	}

	audCfgID, err := AudCfgIDFrom(customData)
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

	_, err = c.api.Encoding.Configurations.Video.H264.Delete(cfgID)
	if err != nil {
		return found, errors.Wrap(err, "removing the video config")
	}

	return found, nil
}

func (c *H264AAC) vidConfigWithCustomDataFrom(cfgID string) (*model.H264VideoConfiguration, types.CustomData, error) {
	vidCfg, err := c.api.Encoding.Configurations.Video.H264.Get(cfgID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "retrieving configuration with config ID")
	}

	data, err := c.vidCustomDataFrom(vidCfg.Id)
	if err != nil {
		return nil, nil, err
	}

	return vidCfg, data, nil
}

func (c *H264AAC) vidCustomDataFrom(cfgID string) (types.CustomData, error) {
	data, err := c.api.Encoding.Configurations.Video.H264.Customdata.Get(cfgID)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving custom data with config ID")
	}

	return data.CustomData, nil
}
