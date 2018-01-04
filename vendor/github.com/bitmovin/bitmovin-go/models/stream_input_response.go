package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type StreamInputData struct {
	Result StreamInputResult `json:"result,omitempty"`
}

type StreamInputResult struct {
	FormatName   *string                `json:"formatName,omitempty"`
	StartTime    *float64               `json:"startTime,omitempty"`
	Duration     *float64               `json:"duration,omitempty"`
	Size         *int64                 `json:"size,omitempty"`
	Bitrate      *int64                 `json:"bitrate,omitempty"`
	AudioStreams []StreamInputAudio     `json:"audioStreams,omitempty"`
	VideoStreams []StreamInputVideo     `json:"videoStreams,omitempty"`
	Tags         map[string]interface{} `json:"tags,omitempty"`
}

type StreamInputAudio struct {
	ID              *string  `json:"id,omitempty"`
	Position        *int64   `json:"position,omitempty"`
	Duration        *float64 `json:"duration,omitempty"`
	Codec           *string  `json:"codec,omitempty"`
	SampleRate      *int64   `json:"sampleRate,omitempty"`
	Bitrate         *int64   `json:"bitrate,string,omitempty"`
	ChannelFormat   *string  `json:"channelFormat,omitempty"`
	Language        *string  `json:"language,omitempty"`
	HearingImpaired *bool    `json:"hearingImpaired,omitempty"`
}

type StreamInputVideo struct {
	ID       *string  `json:"id,omitempty"`
	Position *int64   `json:"position,omitempty"`
	Duration *float64 `json:"duration,omitempty"`
	Codec    *string  `json:"codec,omitempty"`
	FPS      *string  `json:"fps,omitempty"`
	Bitrate  *int64   `json:"bitrate,string,omitempty"`
	Width    *int64   `json:"width,omitempty"`
	Height   *int64   `json:"height,omitempty"`
}

type StreamInputResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      StreamInputData              `json:"data,omitempty"`
}
