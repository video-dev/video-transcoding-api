package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type DashManifest struct {
	ID           *string  `json:"id"`
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Outputs      []Output `json:"outputs"`
	ManifestName *string  `json:"manifestName"`
}

func (h *DashManifest) AddOutput(output *Output) {
	h.Outputs = append(h.Outputs, *output)
}

type DashManifestData struct {
	//Success fields
	Result   DashManifest `json:"result,omitempty"`
	Messages []Message    `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type DashManifestResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      DashManifestData             `json:"data,omitempty"`
}

type LiveDashManifest struct {
	ManifestID     *string  `json:"manifestId,omitempty"`
	TimeShift      *float64 `json:"timeShift,omitempty"`
	LiveEdgeOffset *float64 `json:"liveEdgeOffset,omitempty"`
}

type Period struct {
	ID       *string  `json:"id"`
	Start    *float64 `json:"start"`
	Duration *float64 `json:"duration"`
}

type PeriodData struct {
	//Success fields
	Result   Period    `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type PeriodResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      PeriodData                   `json:"data,omitempty"`
}

type CustomAttribute struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type AudioAdaptationSet struct {
	ID               *string           `json:"id,omitempty"`
	CustomAttributes []CustomAttribute `json:"customAttributes,omitempty"`
	Language         *string           `json:"lang,omitempty"`
}

type AudioAdaptationSetData struct {
	//Success fields
	Result   AudioAdaptationSet `json:"result,omitempty"`
	Messages []Message          `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type AudioAdaptationSetResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      AudioAdaptationSetData       `json:"data,omitempty"`
}

type VideoAdaptationSet struct {
	ID               *string           `json:"id,omitempty"`
	CustomAttributes []CustomAttribute `json:"customAttributes,omitempty"`
}

type VideoAdaptationSetData struct {
	//Success fields
	Result   VideoAdaptationSet `json:"result,omitempty"`
	Messages []Message          `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type VideoAdaptationSetResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      VideoAdaptationSetData       `json:"data,omitempty"`
}

type FMP4Representation struct {
	ID                 *string                              `json:"id,omitempty"`
	Type               bitmovintypes.FMP4RepresentationType `json:"type,omitempty"`
	MuxingID           *string                              `json:"muxingId,omitempty"`
	EncodingID         *string                              `json:"encodingId,omitempty"`
	StartSegmentNumber *int64                               `json:"startSegmentNumber,omitempty"`
	SegmentPath        *string                              `json:"segmentPath,omitempty"`
}

type FMP4RepresentationData struct {
	//Success fields
	Result   FMP4Representation `json:"result,omitempty"`
	Messages []Message          `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type FMP4RepresentationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      FMP4RepresentationData       `json:"data,omitempty"`
}

type DrmFMP4Representation struct {
	ID                 *string                              `json:"id,omitempty"`
	Type               bitmovintypes.FMP4RepresentationType `json:"type,omitempty"`
	MuxingID           *string                              `json:"muxingId,omitempty"`
	EncodingID         *string                              `json:"encodingId,omitempty"`
	StartSegmentNumber *string                              `json:"startSegmentNumber,omitempty"`
	EndSegmentNumber   *string                              `json:"endSegmentNumber,omitempty"`
	SegmentPath        *string                              `json:"segmentPath,omitempty"`
	DrmID              *string                              `json:"drmId,omitempty"`
}

type DrmFMP4RepresentationData struct {
	//Success fields
	Result   DrmFMP4Representation `json:"result,omitempty"`
	Messages []Message             `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type DrmFMP4RepresentationResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      DrmFMP4RepresentationData    `json:"data,omitempty"`
}

type AdaptationSetContentProtection struct {
	EncodingId *string `json:"encodingId,omitempty"`
	MuxingId   *string `json:"muxingId,omitempty"`
	DrmId      *string `json:"drmId"`
}

type AdaptationSetContentProtectionData struct {
	//Success fields
	Result   AdaptationSetContentProtection `json:"result,omitempty"`
	Messages []Message                      `json:"messages,omitempty"`

	//Error fields
	Code             *int64   `json:"code,omitempty"`
	Message          *string  `json:"message,omitempty"`
	DeveloperMessage *string  `json:"developerMessage,omitempty"`
	Links            []Link   `json:"links,omitempty"`
	Details          []Detail `json:"details,omitempty"`
}

type AdaptationSetContentProtectionResponse struct {
	RequestID *string                            `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus       `json:"status,omitempty"`
	Data      AdaptationSetContentProtectionData `json:"data,omitempty"`
}
