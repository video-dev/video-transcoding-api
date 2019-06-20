package codec

import (
	"fmt"
	"strconv"

	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/bitmovin/bitmovin-api-sdk-go/model"
	"github.com/pkg/errors"
)

const defaultVorbisSampleRate = 48000

// NewVorbis creates a Vorbis codec configuration and returns its ID
func NewVorbis(api *bitmovin.BitmovinApi, bitrate string) (string, error) {
	createCfg, err := vorbisConfigFrom(bitrate)
	if err != nil {
		return "", err
	}

	cfg, err := api.Encoding.Configurations.Audio.Vorbis.Create(createCfg)
	if err != nil {
		return "", errors.Wrap(err, "creating audio config")
	}

	return cfg.Id, nil
}

func vorbisConfigFrom(bitrate string) (model.VorbisAudioConfiguration, error) {
	convertedBitrate, err := strconv.ParseInt(bitrate, 10, 64)
	if err != nil {
		return model.VorbisAudioConfiguration{}, errors.Wrapf(err, "parsing audio bitrate %q to int64", bitrate)
	}

	return model.VorbisAudioConfiguration{
		Name:    fmt.Sprintf("vorbis_%s_%d", bitrate, defaultVorbisSampleRate),
		Bitrate: &convertedBitrate,
		Rate:    floatToPtr(defaultVorbisSampleRate),
	}, nil
}
