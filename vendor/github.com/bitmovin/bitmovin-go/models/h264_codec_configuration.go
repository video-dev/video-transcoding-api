package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type H264CodecConfiguration struct {
	ID               *string                        `json:"id,omitempty"`
	Name             *string                        `json:"name,omitempty"`
	Description      *string                        `json:"description,omitempty"`
	CustomData       map[string]interface{}         `json:"customData,omitempty"`
	Bitrate          *int64                         `json:"bitrate,omitempty"`
	FrameRate        *float64                       `json:"rate,omitempty"`
	Width            *int64                         `json:"width,omitempty"`
	Height           *int64                         `json:"height,omitempty"`
	Profile          bitmovintypes.H264Profile      `json:"profile,omitempty"`
	BFrames          *int64                         `json:"bFrames,omitempty"`
	RefFrames        *int64                         `json:"refFrames,omitempty"`
	QPMin            *int64                         `json:"qpMin,omitempty"`
	QPMax            *int64                         `json:"qpMax,omitempty"`
	MVPredictionMode bitmovintypes.MVPredictionMode `json:"mvPredictionMode,omitempty"`
	MVSearchRangeMax *int64                         `json:"mvSearchRangeMax,omitempty"`
	CABAC            *bool                          `json:"cabac,omitempty"`
	MaxBitrate       *int64                         `json:"maxBitrate,omitempty"`
	MinBitrate       *int64                         `json:"minBitrate,omitempty"`
	BufSize          *int64                         `json:"bufsize,omitempty"`
	MinGOP           *int64                         `json:"minGop,omitempty"`
	MaxGOP           *int64                         `json:"maxGop,omitempty"`
	Level            bitmovintypes.H264Level        `json:"level,omitempty"`
	Trellis          bitmovintypes.Trellis          `json:"trellis,omitempty"`
	RcLookahead      *int64                         `json:"rcLookahead,omitempty"`
	Partitions       []bitmovintypes.Partition      `json:"partitions,omitempty"`
	CRF              *float64                       `json:"crf,omitempty"`
}

type H264CodecConfigurationData struct {
	//Success fields
	Result   H264CodecConfiguration `json:"result,omitempty"`
	Messages []Message              `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type H264CodecConfigurationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      H264CodecConfigurationData   `json:"data,omitempty"`
}

type H264CodecConfigurationListResult struct {
	TotalCount *int64                   `json:"totalCount,omitempty"`
	Previous   *string                  `json:"previous,omitempty"`
	Next       *string                  `json:"next,omitempty"`
	Items      []H264CodecConfiguration `json:"items,omitempty"`
}

type H264CodecConfigurationListData struct {
	Result H264CodecConfigurationListResult `json:"result,omitempty"`
}

type H264CodecConfigurationListResponse struct {
	RequestID *string                        `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus   `json:"status,omitempty"`
	Data      H264CodecConfigurationListData `json:"data,omitempty"`
}
