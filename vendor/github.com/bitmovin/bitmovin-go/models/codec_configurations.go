package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type CodecConfigurationItem struct {
	ID *string `json:"id,omitempty"`
	// TODO: Codec typing
	Type *string `json:"type,omitempty"`
}

type CodecConfigurationListResult struct {
	TotalCount *int64                   `json:"totalCount,omitempty"`
	Previous   *string                  `json:"previous,omitempty"`
	Next       *string                  `json:"next,omitempty"`
	Items      []CodecConfigurationItem `json:"items,omitempty"`
}

type CodecConfigurationListData struct {
	Result CodecConfigurationListResult `json:"result,omitempty"`
}

type CodecConfigurationListResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      CodecConfigurationListData   `json:"data,omitempty"`
}

type AACCodecConfiguration struct {
	ID            *string                        `json:"id,omitempty"`
	Name          *string                        `json:"name,omitempty"`
	Description   *string                        `json:"description,omitempty"`
	CustomData    map[string]interface{}         `json:"customData,omitempty"`
	Bitrate       *int64                         `json:"bitrate,omitempty"`
	SamplingRate  *float64                       `json:"rate,omitempty"`
	ChannelLayout bitmovintypes.AACChannelLayout `json:"channelLayout,omitempty"`
	VolumeAdjust  *int64                         `json:"volumeAdjust,omitempty"`
	Normalize     *bool                          `json:"normalize,omitempty"`
}

type AACCodecConfigurationData struct {
	//Success fields
	Result   AACCodecConfiguration `json:"result,omitempty"`
	Messages []Message             `json:"messages,omitempty"`

	//Error fields
	// TODO: type all the error codes similarly to the http status codes in golang
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type AACCodecConfigurationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      AACCodecConfigurationData    `json:"data,omitempty"`
}

type AACCodecConfigurationListResult struct {
	TotalCount *int64                  `json:"totalCount,omitempty"`
	Previous   *string                 `json:"previous,omitempty"`
	Next       *string                 `json:"next,omitempty"`
	Items      []AACCodecConfiguration `json:"items,omitempty"`
}

type AACCodecConfigurationListData struct {
	Result AACCodecConfigurationListResult `json:"result,omitempty"`
}

type AACCodecConfigurationListResponse struct {
	RequestID *string                       `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus  `json:"status,omitempty"`
	Data      AACCodecConfigurationListData `json:"data,omitempty"`
}

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

type H265CodecConfiguration struct {
	ID                       *string                    `json:"id,omitempty"`
	Name                     *string                    `json:"name,omitempty"`
	Description              *string                    `json:"description,omitempty"`
	CustomData               map[string]interface{}     `json:"customData,omitempty"`
	Bitrate                  *int64                     `json:"bitrate,omitempty"`
	FrameRate                *float64                   `json:"rate,omitempty"`
	Width                    *int64                     `json:"width,omitempty"`
	Height                   *int64                     `json:"height,omitempty"`
	Profile                  bitmovintypes.H265Profile  `json:"profile,omitempty"`
	BFrames                  *int64                     `json:"bFrames,omitempty"`
	RefFrames                *int64                     `json:"refFrames,omitempty"`
	QP                       *int64                     `json:"qp,omitempty"`
	MaxBitrate               *int64                     `json:"maxBitrate,omitempty"`
	MinBitrate               *int64                     `json:"minBitrate,omitempty"`
	BufSize                  *int64                     `json:"bufsize,omitempty"`
	MinGOP                   *int64                     `json:"minGop,omitempty"`
	MaxGOP                   *int64                     `json:"maxGop,omitempty"`
	Level                    bitmovintypes.H265Level    `json:"level,omitempty"`
	RCLookahead              *int64                     `json:"rcLookahead,omitempty"`
	BAdapt                   bitmovintypes.BAdapt       `json:"bAdapt,omitempty"`
	MaxCTUSize               bitmovintypes.MaxCTUSize   `json:"maxCTUSize,omitempty"`
	TUIntraDepth             bitmovintypes.TUIntraDepth `json:"tuIntraDepth,omitempty"`
	TUInterDepth             bitmovintypes.TUInterDepth `json:"tuInterDepth,omitempty"`
	MotionSearch             bitmovintypes.MotionSearch `json:"motionSearch,omitempty"`
	SubMe                    *int64                     `json:"subMe,omitempty"`
	MotionSearchRange        *int64                     `json:"motionSearchRange,omitempty"`
	WeightPredictionOnPSlice *bool                      `json:"weightPredictionOnPSlice,omitempty"`
	WeightPredictionOnBSlice *bool                      `json:"weightPredictionOnBSlice,omitempty"`
	SAO                      *bool                      `json:"sao,omitempty"`
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
