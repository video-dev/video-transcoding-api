package bitmovintypes

type ChannelLayout string

// Constants in Golang are not supposed to be all CAPS
const (
	ChannelLayoutNone          ChannelLayout = "NONE"
	ChannelLayoutMono          ChannelLayout = "MONO"
	ChannelLayoutStereo        ChannelLayout = "STEREO"
	ChannelLayoutSurround      ChannelLayout = "SURROUND"
	ChannelLayoutQuad          ChannelLayout = "QUAD"
	ChannelLayoutHexagonal     ChannelLayout = "HEXAGONAL"
	ChannelLayoutOctogonal     ChannelLayout = "OCTOGONAL"
	ChannelLayoutStereoDownmix ChannelLayout = "STEREO_DOWNMIX"
	ChannelLayout21            ChannelLayout = "2.1"
	ChannelLayout22            ChannelLayout = "2.2"
	ChannelLayout31            ChannelLayout = "3.1"
	ChannelLayout40            ChannelLayout = "4.0"
	ChannelLayout41            ChannelLayout = "4.1"
	ChannelLayout50            ChannelLayout = "5.0"
	ChannelLayout50Back        ChannelLayout = "5.0_BACK"
	ChannelLayout51            ChannelLayout = "5.1"
	ChannelLayout51Back        ChannelLayout = "5.1_BACK"
	ChannelLayout60            ChannelLayout = "6.0"
	ChannelLayout60Front       ChannelLayout = "6.0_FRONT"
	ChannelLayout61            ChannelLayout = "6.1"
	ChannelLayout61Back        ChannelLayout = "6.1_BACK"
	ChannelLayout61Front       ChannelLayout = "6.1_FRONT"
	ChannelLayout70            ChannelLayout = "7.0"
	ChannelLayout70Front       ChannelLayout = "7.0_FRONT"
	ChannelLayout71            ChannelLayout = "7.1"
	ChannelLayout70Wide        ChannelLayout = "7.0_WIDE"
	ChannelLayout70WideBack    ChannelLayout = "7.0_WIDE_BACK"
)
