package models

import (
	"github.com/bitmovin/bitmovin-go/bitmovintypes"
)

type Filter struct {
	ID          *string                `json:"id"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	CustomData  map[string]interface{} `json:"customData"`
}

type WatermarkFilter struct {
	Filter
	Image  *string `json:"image"`
	Left   *int64  `json:"left"`
	Right  *int64  `json:"right"`
	Top    *int64  `json:"top"`
	Bottom *int64  `json:"bottom"`
}

type CropFilter struct {
	Filter
	Left   *int64 `json:"left"`
	Right  *int64 `json:"right"`
	Top    *int64 `json:"top"`
	Bottom *int64 `json:"bottom"`
}

type RotationFilter struct {
	Filter
	Rotation *int64 `json:"rotation"`
}

type DeinterlacingFilter struct {
	Filter
	Mode   bitmovintypes.DeinterlacingMode  `json:"mode"`
	Parity bitmovintypes.PictureFieldParity `json:"parity"`
}
