package codec

import (
	"strconv"

	"github.com/pkg/errors"
)

func dimensionFrom(presetDimension string) (*int32, error) {
	dim, err := strconv.ParseInt(presetDimension, 10, 32)
	if err != nil {
		return nil, err
	}

	return int32ToPtr(int32(dim)), nil
}

func bitrateFrom(presetBitrate string) (*int64, error) {
	if presetBitrate == "" {
		return nil, errors.New("video bitrate must be set")
	}

	bitrate, err := strconv.ParseInt(presetBitrate, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "parsing bitrate to int64")
	}

	return intToPtr(bitrate), nil
}
