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
