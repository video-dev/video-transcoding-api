package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type Sprite struct {
	Name        *string  `json:"name,omitempty"`
	SpriteName  *string  `json:"spriteName,omitempty"`
	Description *string  `json:"description,omitempty"`
	VTTName     *string  `json:"vttName,omitempty"`
	Height      int      `json:"height,omitempty"`
	Width       int      `json:"width,omitempty"`
	Distance    int      `json:"distance,omitempty"`
	Outputs     []Output `json:"outputs,omitempty"`
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

func NewSprite(name, spriteName, description, vttName *string, height, width, distance int, outputs []Output) *Sprite {
	return &Sprite{
		Name:        name,
		SpriteName:  spriteName,
		Description: description,
		VTTName:     vttName,
		Height:      height,
		Width:       width,
		Distance:    distance,
		Outputs:     outputs,
	}
}
