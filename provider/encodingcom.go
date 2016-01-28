package provider

import (
	"errors"
	"strconv"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
)

// ErrMissingData is the error returned by the factory when required data is
// missing.
var ErrMissingData = InvalidConfigError("missing Encoding.com user id or key. Please define the environment variables ENCODINGCOM_USER_ID and ENCODINGCOM_USER_KEY")

type encodingComProvider struct {
	config *config.Config
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
		TwoPass:             encodingcom.YesNoBoolean(profile.TwoPassEncoding),
	}
	if profile.Rotate.set {
		format.Rotate = strconv.FormatUint(uint64(profile.Rotate.value), 10)
	} else {
		format.Rotate = "def"
	}
	return &format
}

func (e *encodingComProvider) JobStatus(id string) (*JobStatus, error) {
	resp, err := e.client.GetStatus([]string{id})
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, errors.New("invalid value returned by the Encoding.com API: []")
	}
	return &JobStatus{
		ProviderJobID: id,
		ProviderName:  "encoding.com",
		Status:        e.statusMap(resp[0].MediaStatus),
		ProviderStatus: map[string]interface{}{
			"progress":   resp[0].Progress,
			"sourcefile": resp[0].SourceFile,
			"timeleft":   resp[0].TimeLeft,
			"created":    resp[0].CreateDate,
			"started":    resp[0].StartDate,
			"finished":   resp[0].FinishDate,
		},
	}, nil
}

func (e *encodingComProvider) statusMap(encodingComStatus string) status {
	switch encodingComStatus {
	case "New":
		return StatusQueued
	case "Downloading", "Ready to process", "Waiting for encoder", "Processing", "Saving":
		return StatusStarted
	case "Finished":
		return StatusFinished
	default:
		return StatusFailed
	}
}

// EncodingComProvider is the factory function for the Encoding.com provider.
func EncodingComProvider(cfg *config.Config) (TranscodingProvider, error) {
	if cfg.EncodingCom.UserID == "" || cfg.EncodingCom.UserKey == "" {
		return nil, ErrMissingData
	}
	client, err := encodingcom.NewClient("https://manage.encoding.com", cfg.EncodingCom.UserID, cfg.EncodingCom.UserKey)
	if err != nil {
		return nil, err
	}
	return &encodingComProvider{client: client, config: cfg}, nil
}
