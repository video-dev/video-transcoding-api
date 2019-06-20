package codec

import (
	"fmt"
	"strconv"

	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/bitmovin/bitmovin-api-sdk-go/model"
	"github.com/pkg/errors"
)

const defaultAACSampleRate = 48000

// NewAAC creates a AAC codec configuration and returns its ID
func NewAAC(api *bitmovin.BitmovinApi, bitrate string) (string, error) {
	createCfg, err := aacConfigFrom(bitrate)
	if err != nil {
		return "", err
	}

	cfg, err := api.Encoding.Configurations.Audio.Aac.Create(createCfg)
	if err != nil {
		return "", errors.Wrap(err, "creating audio cfg")
	}

	return cfg.Id, nil
}

func aacConfigFrom(bitrate string) (model.AacAudioConfiguration, error) {
	convertedBitrate, err := strconv.ParseInt(bitrate, 10, 64)
	if err != nil {
		return model.AacAudioConfiguration{}, errors.Wrapf(err, "parsing audio bitrate %q to int64", bitrate)
	}

	return model.AacAudioConfiguration{
		Name:    fmt.Sprintf("aac_%s_%d", bitrate, defaultAACSampleRate),
		Bitrate: &convertedBitrate,
		Rate:    floatToPtr(defaultAACSampleRate),
	}, nil
}
