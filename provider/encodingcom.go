package provider

import (
	"errors"
	"strconv"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
)

// ErrMissingData is the error returned by the factory when required data is
// missing.
var ErrMissingData = errors.New("missing Encoding.com user id or key. Please define the environment variables ENCODINGCOM_USER_ID and ENCODINGCOM_USER_KEY")

type encodingComProvider struct {
	client *encodingcom.Client
}

func (e *encodingComProvider) Transcode(sourceMedia, destination string, profile Profile) (*JobStatus, error) {
	format := e.profileToFormat(profile)
	format.Destination = []string{destination}
	resp, err := e.client.AddMedia([]string{sourceMedia}, format)
	if err != nil {
		return nil, err
	}
	return &JobStatus{ProviderJobID: resp.MediaID, StatusMessage: resp.Message}, nil
}

func (e *encodingComProvider) profileToFormat(profile Profile) *encodingcom.Format {
	format := encodingcom.Format{
		Output:              []string{profile.Output},
		Size:                profile.Size.String(),
		AudioCodec:          profile.AudioCodec,
		AudioBitrate:        profile.AudioBitRate,
		AudioChannelsNumber: profile.AudioChannelsNumber,
		AudioSampleRate:     profile.AudioSampleRate,
		Bitrate:             profile.BitRate,
		Framerate:           profile.FrameRate,
		KeepAspectRatio:     encodingcom.YesNoBoolean(profile.KeepAspectRatio),
		VideoCodec:          profile.VideoCodec,
		Keyframe:            []string{profile.KeyFrame},
		AudioVolume:         profile.AudioVolume,
	}
	if profile.Rotate.set {
		format.Rotate = strconv.FormatUint(uint64(profile.Rotate.value), 10)
	} else {
		format.Rotate = "def"
	}
	return &format
}

func (e *encodingComProvider) JobStatus(id string) (*JobStatus, error) {
	return nil, nil
}

// EncodingComProvider is the factory function for the Encoding.com provider.
func EncodingComProvider(cfg *config.Config) (TranscodingProvider, error) {
	if cfg.EncodingComUserID == "" || cfg.EncodingComUserKey == "" {
		return nil, ErrMissingData
	}
	client, err := encodingcom.NewClient("https://manage.encoding.com", cfg.EncodingComUserID, cfg.EncodingComUserKey)
	if err != nil {
		return nil, err
	}
	return &encodingComProvider{client: client}, nil
}
