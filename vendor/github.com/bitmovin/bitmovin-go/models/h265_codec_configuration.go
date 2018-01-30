package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type H265CodecConfiguration struct {
	ID                          *string                    `json:"id,omitempty"`
	Name                        *string                    `json:"name,omitempty"`
	Description                 *string                    `json:"description,omitempty"`
	CustomData                  map[string]interface{}     `json:"customData,omitempty"`
	Bitrate                     *int64                     `json:"bitrate,omitempty"`
	FrameRate                   *float64                   `json:"rate,omitempty"`
	Width                       *int64                     `json:"width,omitempty"`
	Height                      *int64                     `json:"height,omitempty"`
	Profile                     bitmovintypes.H265Profile  `json:"profile,omitempty"`
	BFrames                     *int64                     `json:"bFrames,omitempty"`
	RefFrames                   *int64                     `json:"refFrames,omitempty"`
	QP                          *int64                     `json:"qp,omitempty"`
	MaxBitrate                  *int64                     `json:"maxBitrate,omitempty"`
	MinBitrate                  *int64                     `json:"minBitrate,omitempty"`
	BufSize                     *int64                     `json:"bufsize,omitempty"`
	MinGOP                      *int64                     `json:"minGop,omitempty"`
	MaxGOP                      *int64                     `json:"maxGop,omitempty"`
	Level                       bitmovintypes.H265Level    `json:"level,omitempty"`
	RCLookahead                 *int64                     `json:"rcLookahead,omitempty"`
	BAdapt                      bitmovintypes.BAdapt       `json:"bAdapt,omitempty"`
	MaxCTUSize                  bitmovintypes.MaxCTUSize   `json:"maxCTUSize,omitempty"`
	TUIntraDepth                bitmovintypes.TUIntraDepth `json:"tuIntraDepth,omitempty"`
	TUInterDepth                bitmovintypes.TUInterDepth `json:"tuInterDepth,omitempty"`
	MotionSearch                bitmovintypes.MotionSearch `json:"motionSearch,omitempty"`
	SubMe                       *int64                     `json:"subMe,omitempty"`
	MotionSearchRange           *int64                     `json:"motionSearchRange,omitempty"`
	WeightPredictionOnPSlice    *bool                      `json:"weightPredictionOnPSlice,omitempty"`
	WeightPredictionOnBSlice    *bool                      `json:"weightPredictionOnBSlice,omitempty"`
	SAO                         *bool                      `json:"sao,omitempty"`
	CRF                         *float64                   `json:"crf,omitempty"`
	ColorConfig                 ColorConfig                `json:"colorConfig,omitempty"`
	MasterDisplay               *string                    `json:"masterDisplay,omitempty"`
	MaxContentLightLevel        *int64                     `json:"maxContentLightLevel,omitempty"`
	MaxPictureAverageLightLevel *int64                     `json:"maxPictureAverageLightLevel,omitempty"`
	HDR                         *bool                      `json:"hdr,omitempty"`
	PixelFormat                 bitmovintypes.PixelFormat  `json:"pixelFormat,omitempty"`
}

type H265CodecConfigurationData struct {
	//Success fields
	Result   H265CodecConfiguration `json:"result,omitempty"`
	Messages []Message              `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type H265CodecConfigurationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      H265CodecConfigurationData   `json:"data,omitempty"`
}

type H265CodecConfigurationListResult struct {
	TotalCount *int64                   `json:"totalCount,omitempty"`
	Previous   *string                  `json:"previous,omitempty"`
	Next       *string                  `json:"next,omitempty"`
	Items      []H265CodecConfiguration `json:"items,omitempty"`
}

type H265CodecConfigurationListData struct {
	Result H265CodecConfigurationListResult `json:"result,omitempty"`
}

type H265CodecConfigurationListResponse struct {
	RequestID *string                        `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus   `json:"status,omitempty"`
	Data      H265CodecConfigurationListData `json:"data,omitempty"`
}
