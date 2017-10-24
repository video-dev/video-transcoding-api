package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type H264CodecConfigBuilder struct {
	Config *H264CodecConfiguration
}

func NewH264CodecConfigBuilder(name string) *H264CodecConfigBuilder {
	return &H264CodecConfigBuilder{
		Config: &H264CodecConfiguration{
			Name: &name,
		},
	}
}

func (b *H264CodecConfigBuilder) Width(w int64) *H264CodecConfigBuilder {
	b.Config.Width = &w
	return b
}
func (b *H264CodecConfigBuilder) Height(h int64) *H264CodecConfigBuilder {
	b.Config.Height = &h
	return b
}
func (b *H264CodecConfigBuilder) Bitrate(br int64) *H264CodecConfigBuilder {
	b.Config.Bitrate = &br
	return b
}
func (b *H264CodecConfigBuilder) Framerate(f float64) *H264CodecConfigBuilder {
	b.Config.FrameRate = &f
	return b
}
func (b *H264CodecConfigBuilder) Profile(p bitmovintypes.H264Profile) *H264CodecConfigBuilder {
	b.Config.Profile = p
	return b
}
func (b *H264CodecConfigBuilder) BFrames(bf int64) *H264CodecConfigBuilder {
	b.Config.BFrames = &bf
	return b
}
func (b *H264CodecConfigBuilder) RefFrames(r int64) *H264CodecConfigBuilder {
	b.Config.RefFrames = &r
	return b
}
func (b *H264CodecConfigBuilder) MVPredictionMode(m bitmovintypes.MVPredictionMode) *H264CodecConfigBuilder {
	b.Config.MVPredictionMode = m
	return b
}
func (b *H264CodecConfigBuilder) MVSearchRangeMax(r int64) *H264CodecConfigBuilder {
	b.Config.MVSearchRangeMax = &r
	return b
}
func (b *H264CodecConfigBuilder) CABAC(r bool) *H264CodecConfigBuilder {
	b.Config.CABAC = &r
	return b
}
func (b *H264CodecConfigBuilder) Trellis(r bitmovintypes.Trellis) *H264CodecConfigBuilder {
	b.Config.Trellis = r
	return b
}
func (b *H264CodecConfigBuilder) RcLookahead(r int64) *H264CodecConfigBuilder {
	b.Config.RcLookahead = &r
	return b
}
func (b *H264CodecConfigBuilder) Partitions(r []bitmovintypes.Partition) *H264CodecConfigBuilder {
	b.Config.Partitions = r
	return b
}
func (b *H264CodecConfigBuilder) Build() *H264CodecConfiguration {
	return b.Config
}
