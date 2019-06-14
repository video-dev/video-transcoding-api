package configuration

import (
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider/bitmovinnewsdk/types"
	"github.com/bitmovin/bitmovin-api-sdk-go/model"
	"github.com/pkg/errors"
)

// Store is the interface for any underlying codec config services
type Store interface {
	Create(preset db.Preset) (string, error)
	Get(presetID string) (bool, Details, error)
	Delete(presetID string) (found bool, e error)
}

// Details hold details for video / audio configurations
type Details struct {
	Video      model.CodecConfiguration
	Audio      model.CodecConfiguration
	CustomData types.CustomData
}

const (
	customDataKeyAudio         = "audio"
	customDataKeyAudioID       = "id"
	customDataKeyContainer     = "container"
	customDataKeyContainerName = "name"
)

// ContainerFrom is a helper for extracting the container value from customData
func ContainerFrom(data types.CustomData) (string, error) {
	contnr, err := types.CustomDataStringValAtKeys(data, customDataKeyContainer, customDataKeyContainerName)
	if err != nil {
		return "", errors.New("extracting container from custom data")
	}

	return contnr, nil
}

// AudCfgIDFrom extracts the audio configuration id from customData
func AudCfgIDFrom(data types.CustomData) (string, error) {
	audCfgID, err := types.CustomDataStringValAtKeys(data, customDataKeyAudio, customDataKeyAudioID)
	if err != nil {
		return "", errors.New("extracting audio config ID from custom data")
	}

	return audCfgID, nil
}

func customDataWith(audCfgID, container string) types.CustomData {
	return &map[string]map[string]interface{}{
		customDataKeyAudio:     {customDataKeyAudioID: audCfgID},
		customDataKeyContainer: {customDataKeyContainerName: container},
	}
}
