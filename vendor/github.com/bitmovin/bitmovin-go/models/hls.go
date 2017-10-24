package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type HLSManifest struct {
	ID           *string  `json:"id,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Outputs      []Output `json:"outputs,omitempty"`
	ManifestName *string  `json:"manifestName,omitempty"`
}

type HLSManifestData struct {
	//Success fields
	Result   HLSManifest `json:"result,omitempty"`
	Messages []Message   `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type HLSManifestResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      HLSManifestData              `json:"data,omitempty"`
}

type LiveHLSManifest struct {
	ManifestID *string  `json:"manifestId,omitempty"`
	TimeShift  *float64 `json:"timeShift,omitempty"`
}

type StreamInfo struct {
	ID                 *string `json:"id,omitempty"`
	URI                *string `json:"uri,omitempty"`
	Audio              *string `json:"audio,omitempty"`
	Video              *string `json:"video,omitempty"`
	Subtitles          *string `json:"subtitles,omitempty"`
	ClosedCaptions     *string `json:"closedCaptions,omitempty"`
	SegmentPath        *string `json:"segmentPath,omitempty"`
	EncodingID         *string `json:"encodingId,omitempty"`
	StreamID           *string `json:"streamId,omitempty"`
	MuxingID           *string `json:"muxingId,omitempty"`
	DRMID              *string `json:"drmId,omitempty"`
	StartSegmentNumber *int64  `json:"startSegmentNumber,omitempty"`
	EndSegmentNumber   *int64  `json:"endSegmentNumber,omitempty"`
}

type StreamInfoData struct {
	//Success fields
	Result   StreamInfo `json:"result,omitempty"`
	Messages []Message  `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type StreamInfoResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      StreamInfoData               `json:"data,omitempty"`
}

type MediaInfo struct {
	ID                 *string                 `json:"id,omitempty"`
	Type               bitmovintypes.MediaType `json:"type,omitempty"`
	URI                *string                 `json:"uri,omitempty"`
	GroupID            *string                 `json:"groupId,omitempty"`
	Language           *string                 `json:"language,omitempty"`
	AssociatedLanguage *string                 `json:"assocLanguage,omitempty"`
	Name               *string                 `json:"name,omitempty"`
	IsDefault          *bool                   `json:"isDefault,omitempty"`
	Autoselect         *bool                   `json:"autoselect,omitempty"`
	Forced             *bool                   `json:"forced,omitempty"`
	InstreamID         *string                 `json:"instreamId,omitempty"`
	Characteristics    []string                `json:"characteristics,omitempty"`
	SegmentPath        *string                 `json:"segmentPath,omitempty"`
	EncodingID         *string                 `json:"encodingId,omitempty"`
	StreamID           *string                 `json:"streamId,omitempty"`
	MuxingID           *string                 `json:"muxingId,omitempty"`
	DRMID              *string                 `json:"drmId,omitempty"`
	StartSegmentNumber *int64                  `json:"startSegmentNumber,omitempty"`
	EndSegmentNumber   *int64                  `json:"endSegmentNumber,omitempty"`
}

type MediaInfoData struct {
	//Success fields
	Result   MediaInfo `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type MediaInfoResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      MediaInfoData                `json:"data,omitempty"`
}
