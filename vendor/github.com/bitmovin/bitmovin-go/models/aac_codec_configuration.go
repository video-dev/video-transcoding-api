package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

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
