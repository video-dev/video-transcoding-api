package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type SmoothStreamingManifest struct {
	ID                 *string                `json:"id,omitempty"`
	Name               *string                `json:"name,omitempty"`
	Description        *string                `json:"description,omitempty"`
	CustomData         map[string]interface{} `json:"customData,omitempty"`
	Outputs            []Output               `json:"outputs,omitempty"`
	ServerManifestName *string                `json:"serverManifestName,omitempty"`
	ClientManifestName *string                `json:"clientManifestName,omitempty"`
}

type SmoothStreamingManifestData struct {
	Result           SmoothStreamingManifest `json:"result,omitempty"`
	Messages         []Message               `json:"messages,omitempty"`
	Code             *int64                  `json:"code,omitempty"`
	Message          *string                 `json:"message,omitempty"`
	DeveloperMessage *string                 `json:"developerMessage,omitempty"`
	Links            []Link                  `json:"links,omitempty"`
	Details          []Detail                `json:"details,omitempty"`
}

type SmoothStreamingManifestResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      SmoothStreamingManifestData  `json:"data,omitempty"`
}

type SmoothStreamingMp4Representation struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	EncodingID  *string `json:"encodingId,omitempty"`
	MuxingID    *string `json:"muxingId,omitempty"`
	MediaFile   *string `json:"mediaFile,omitempty"`
	Language    *string `json:"language,omitempty"`
	TrackName   *string `json:"trackName,omitempty"`
}

type SmoothStreamingMp4RepresentationData struct {
	Result           SmoothStreamingManifestResponse `json:"result,omitempty"`
	Messages         []Message                       `json:"messages,omitempty"`
	Code             *int64                          `json:"code,omitempty"`
	Message          *string                         `json:"message,omitempty"`
	DeveloperMessage *string                         `json:"developerMessage,omitempty"`
	Links            []Link                          `json:"links,omitempty"`
	Details          []Detail                        `json:"details,omitempty"`
}

type SmoothStreamingMp4RepresentationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      SmoothStreamingManifestData  `json:"data,omitempty"`
}

type SmoothStreamingContentProtection struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	EncodingID  *string `json:"encodingId,omitempty"`
	MuxingID    *string `json:"muxingId,omitempty"`
	DrmID       *string `json:"drmId,omitempty"`
}

type SmoothStreamingContentProtectionData struct {
	Result           SmoothStreamingContentProtection `json:"result,omitempty"`
	Messages         []Message                        `json:"messages,omitempty"`
	Code             *int64                           `json:"code,omitempty"`
	Message          *string                          `json:"message,omitempty"`
	DeveloperMessage *string                          `json:"developerMessage,omitempty"`
	Links            []Link                           `json:"links,omitempty"`
	Details          []Detail                         `json:"details,omitempty"`
}

type SmoothStreamingContentProtectionResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      SmoothStreamingManifestData  `json:"data,omitempty"`
}
