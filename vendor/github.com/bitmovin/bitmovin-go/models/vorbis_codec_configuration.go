package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type VorbisCodecConfiguration struct {
	ID            *string                     `json:"id,omitempty"`
	Name          *string                     `json:"name,omitempty"`
	Description   *string                     `json:"description,omitempty"`
	CustomData    map[string]interface{}      `json:"customData,omitempty"`
	Bitrate       *int64                      `json:"bitrate,omitempty"`
	SamplingRate  *float64                    `json:"rate,omitempty"`
	ChannelLayout bitmovintypes.ChannelLayout `json:"channelLayout,omitempty"`
}

type VorbisCodecConfigurationData struct {
	//Success fields
	Result   VorbisCodecConfiguration `json:"result,omitempty"`
	Messages []Message                `json:"messages,omitempty"`

	//Error fields
	// TODO: type all the error codes similarly to the http status codes in golang
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type VorbisCodecConfigurationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      VorbisCodecConfigurationData `json:"data,omitempty"`
}

type VorbisCodecConfigurationListResult struct {
	TotalCount *int64                     `json:"totalCount,omitempty"`
	Previous   *string                    `json:"previous,omitempty"`
	Next       *string                    `json:"next,omitempty"`
	Items      []VorbisCodecConfiguration `json:"items,omitempty"`
}

type VorbisCodecConfigurationListData struct {
	Result VorbisCodecConfigurationListResult `json:"result,omitempty"`
}

type VorbisCodecConfigurationListResponse struct {
	RequestID *string                          `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus     `json:"status,omitempty"`
	Data      VorbisCodecConfigurationListData `json:"data,omitempty"`
}
