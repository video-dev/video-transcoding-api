package bitmovintypes

type AACChannelLayout string

// Constants in Golang are not supposed to be all CAPS
const (
	AACChannelLayoutNone          AACChannelLayout = "NONE"
	AACChannelLayoutMono          AACChannelLayout = "MONO"
	AACChannelLayoutStereo        AACChannelLayout = "STEREO"
	AACChannelLayoutSurround      AACChannelLayout = "SURROUND"
	AACChannelLayoutQuad          AACChannelLayout = "QUAD"
	AACChannelLayoutHexagonal     AACChannelLayout = "HEXAGONAL"
	AACChannelLayoutOctogonal     AACChannelLayout = "OCTOGONAL"
	AACChannelLayoutStereoDownmix AACChannelLayout = "STEREO_DOWNMIX"
	AACChannelLayout21            AACChannelLayout = "2.1"
	AACChannelLayout22            AACChannelLayout = "2.2"
	AACChannelLayout31            AACChannelLayout = "3.1"
	AACChannelLayout40            AACChannelLayout = "4.0"
	AACChannelLayout41            AACChannelLayout = "4.1"
	AACChannelLayout50            AACChannelLayout = "5.0"
	AACChannelLayout50Back        AACChannelLayout = "5.0_BACK"
	AACChannelLayout51            AACChannelLayout = "5.1"
	AACChannelLayout51Back        AACChannelLayout = "5.1_BACK"
	AACChannelLayout60            AACChannelLayout = "6.0"
	AACChannelLayout60Front       AACChannelLayout = "6.0_FRONT"
	AACChannelLayout61            AACChannelLayout = "6.1"
	AACChannelLayout61Back        AACChannelLayout = "6.1_BACK"
	AACChannelLayout61Front       AACChannelLayout = "6.1_FRONT"
	AACChannelLayout70            AACChannelLayout = "7.0"
	AACChannelLayout70Front       AACChannelLayout = "7.0_FRONT"
	AACChannelLayout71            AACChannelLayout = "7.1"
	AACChannelLayout70Wide        AACChannelLayout = "7.0_WIDE"
	AACChannelLayout70WideBack    AACChannelLayout = "7.0_WIDE_BACK"
)
