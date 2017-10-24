package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type Thumbnail struct {
	Name         *string                    `json:"name,omitempty"`
	Description  *string                    `json:"description,omitempty"`
	Height       int                        `json:"height,omitempty"`
	PositionUnit bitmovintypes.PositionUnit `json:"unit,omitempty"`
	Positions    []float64                  `json:"positions,omitempty"`
	Pattern      *string                    `json:"pattern,omitempty"`
	Outputs      []Output                   `json:"outputs,omitempty"`
}

type ThumbnailData struct {
	Result   Thumbnail `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type ThumbnailResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      ThumbnailData                `json:"data,omitempty"`
}

type ThumbnailListResponse struct {
	TotalCount *int64      `json:"totalCount,omitempty"`
	Previous   *string     `json:"previous,omitempty"`
	Next       *string     `json:"next,omitempty"`
	Items      []Thumbnail `json:"items,omitempty"`
}

func NewThumbnail(height int, positions []float64, outputs []Output) *Thumbnail {
	return &Thumbnail{
		Height:    height,
		Positions: positions,
		Outputs:   outputs,
	}
}

func (t *Thumbnail) Builder() *ThumbnailBuilder {
	return &ThumbnailBuilder{
		Thumbnail: t,
	}
}
