package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type VP8CodecConfiguration struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`

	Bitrate   *int64   `json:"bitrate,omitempty"`
	FrameRate *float64 `json:"rate,omitempty"`
	Width     *int64   `json:"width,omitempty"`
	Height    *int64   `json:"height,omitempty"`

	CRF               *float64                 `json:"crf,omitempty"`
	LagInFrames       *int64                   `json:"lagInFrames,omitempty"`
	MaxIntraRate      *int64                   `json:"maxIntraRate,omitempty"`
	QPMin             *int64                   `json:"qpMin,omitempty"`
	QPMax             *int64                   `json:"qpMax,omitempty"`
	RateUndershootPct *int64                   `json:"rateUndershootPct,omitempty"`
	RateOvershootPct  *int64                   `json:"rateUndershootPct,omitempty"`
	Sharpness         *int64                   `json:"sharpness,omitempty"`
	Quality           bitmovintypes.VPQuality  `json:"quality,omitempty"`
	StaticThresh      *int64                   `json:"staticThresh,omitempty"`
	ARNRMaxFrames     *int64                   `json:"arnrMaxFrames,omitempty"`
	ARNRStrength      *int64                   `json:"arnrStrength,omitempty"`
	ARNRType          bitmovintypes.VPARNRType `json:"arnrType,omitempty"`

	NoiseSensitivity bitmovintypes.VP8NoiseSensitivity `json:"noiseSensitivity,omitempty"`
}

type VP8CodecConfigurationData struct {
	//Success fields
	Result   VP8CodecConfiguration `json:"result,omitempty"`
	Messages []Message             `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type VP8CodecConfigurationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      VP8CodecConfigurationData    `json:"data,omitempty"`
}

type VP8CodecConfigurationListResult struct {
	TotalCount *int64                  `json:"totalCount,omitempty"`
	Previous   *string                 `json:"previous,omitempty"`
	Next       *string                 `json:"next,omitempty"`
	Items      []VP8CodecConfiguration `json:"items,omitempty"`
}

type VP8CodecConfigurationListData struct {
	Result VP8CodecConfigurationListResult `json:"result,omitempty"`
}

type VP8CodecConfigurationListResponse struct {
	RequestID *string                       `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus  `json:"status,omitempty"`
	Data      VP8CodecConfigurationListData `json:"data,omitempty"`
}

type VP9CodecConfiguration struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`

	Bitrate   *int64   `json:"bitrate,omitempty"`
	FrameRate *float64 `json:"rate,omitempty"`
	Width     *int64   `json:"width,omitempty"`
	Height    *int64   `json:"height,omitempty"`

	CRF               *float64                 `json:"crf,omitempty"`
	LagInFrames       *int64                   `json:"lagInFrames,omitempty"`
	MaxIntraRate      *int64                   `json:"maxIntraRate,omitempty"`
	QPMin             *int64                   `json:"qpMin,omitempty"`
	QPMax             *int64                   `json:"qpMax,omitempty"`
	RateUndershootPct *int64                   `json:"rateUndershootPct,omitempty"`
	RateOvershootPct  *int64                   `json:"rateUndershootPct,omitempty"`
	Sharpness         *int64                   `json:"sharpness,omitempty"`
	Quality           bitmovintypes.VPQuality  `json:"quality,omitempty"`
	StaticThresh      *int64                   `json:"staticThresh,omitempty"`
	ARNRMaxFrames     *int64                   `json:"arnrMaxFrames,omitempty"`
	ARNRStrength      *int64                   `json:"arnrStrength,omitempty"`
	ARNRType          bitmovintypes.VPARNRType `json:"arnrType,omitempty"`

	TileColumns      *int64                  `json:"tileColumns,omitempty"`
	TileRows         *int64                  `json:"tileRows,omitempty"`
	FrameParallel    *bool                   `json:"frameParallel,omitempty"`
	NoiseSensitivity *bool                   `json:"noiseSensitivity,omitempty"`
	Lossless         *bool                   `json:"lossless,omitempty"`
	AQMode           bitmovintypes.VP9AQMode `json:"aqMode,omitempty"`
}

type VP9CodecConfigurationData struct {
	//Success fields
	Result   VP9CodecConfiguration `json:"result,omitempty"`
	Messages []Message             `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type VP9CodecConfigurationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      VP9CodecConfigurationData    `json:"data,omitempty"`
}

type VP9CodecConfigurationListResult struct {
	TotalCount *int64                  `json:"totalCount,omitempty"`
	Previous   *string                 `json:"previous,omitempty"`
	Next       *string                 `json:"next,omitempty"`
	Items      []VP9CodecConfiguration `json:"items,omitempty"`
}

type VP9CodecConfigurationListData struct {
	Result VP9CodecConfigurationListResult `json:"result,omitempty"`
}

type VP9CodecConfigurationListResponse struct {
	RequestID *string                       `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus  `json:"status,omitempty"`
	Data      VP9CodecConfigurationListData `json:"data,omitempty"`
}
