package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type Sprite struct {
	Name        *string                    `json:"name,omitempty"`
	SpriteName  *string                    `json:"spriteName,omitempty"`
	Description *string                    `json:"description,omitempty"`
	VTTName     *string                    `json:"vttName,omitempty"`
	Unit        bitmovintypes.DistanceUnit `json:"unit",omitempty"`
	Height      int64                      `json:"height,omitempty"`
	Width       int64                      `json:"width,omitempty"`
	Distance    float64                    `json:"distance,omitempty"`
	Outputs     []Output                   `json:"outputs,omitempty"`
}

type SpriteData struct {
	Result   Sprite    `json:"result,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type SpriteResponse struct {
	RequestID *string                      `json:"requestId,omitempty"`
	Status    bitmovintypes.ResponseStatus `json:"status,omitempty"`
	Data      SpriteData                   `json:"data,omitempty"`
}

type SpriteListResponse struct {
	TotalCount *int64   `json:"totalCount,omitempty"`
	Previous   *string  `json:"previous,omitempty"`
	Next       *string  `json:"next,omitempty"`
	Items      []Sprite `json:"items,omitempty"`
}

func NewSprite(name, spriteName, description, vttName *string, height, width int64, distance float64, unit bitmovintypes.DistanceUnit, outputs []Output) *Sprite {
	return &Sprite{
		Name:        name,
		SpriteName:  spriteName,
		Description: description,
		VTTName:     vttName,
		Unit:        unit,
		Height:      height,
		Width:       width,
		Distance:    distance,
		Outputs:     outputs,
	}
}
