package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type ThumbnailBuilder struct {
	Thumbnail *Thumbnail
}

func (t *ThumbnailBuilder) Name(name string) *ThumbnailBuilder {
	t.Thumbnail.Name = &name
	return t
}
func (t *ThumbnailBuilder) Description(desc string) *ThumbnailBuilder {
	t.Thumbnail.Description = &desc
	return t
}

func (t *ThumbnailBuilder) Height(h int) *ThumbnailBuilder {
	t.Thumbnail.Height = h
	return t
}
func (t *ThumbnailBuilder) PositionUnit(u bitmovintypes.PositionUnit) *ThumbnailBuilder {
	t.Thumbnail.PositionUnit = u
	return t
}
func (t *ThumbnailBuilder) Positions(pos []float64) *ThumbnailBuilder {
	t.Thumbnail.Positions = pos
	return t
}
func (t *ThumbnailBuilder) Pattern(p string) *ThumbnailBuilder {
	t.Thumbnail.Pattern = &p
	return t
}
func (t *ThumbnailBuilder) Outputs(o []Output) *ThumbnailBuilder {
	t.Thumbnail.Outputs = o
	return t
}
func (t *ThumbnailBuilder) Build() *Thumbnail {
	return t.Thumbnail
}
