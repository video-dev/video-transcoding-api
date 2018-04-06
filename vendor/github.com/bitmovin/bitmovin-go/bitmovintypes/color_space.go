package bitmovintypes

type ColorSpace string

const (
	ColorSpaceUnspecified ColorSpace = "UNSPECIFIED"
	ColorSpaceRGB         ColorSpace = "RGB"
	ColorSpaceBT709       ColorSpace = "BT709"
	ColorSpaceFCC         ColorSpace = "FCC"
	ColorSpaceBT470BG     ColorSpace = "BT470BG"
	ColorSpaceSMPTE170M   ColorSpace = "SMPTE170M"
	ColorSpaceSMPTE240M   ColorSpace = "SMPTE240M"
	ColorSpaceYCGCO       ColorSpace = "YCGCO"
	ColorSpaceYCOCG       ColorSpace = "YCOCG"
	ColorSpaceBT2020_NCL  ColorSpace = "BT2020_NCL"
	ColorSpaceBT2020_CL   ColorSpace = "BT2020_CL"
	ColorSpaceSMPTE2085   ColorSpace = "SMPTE2085"
)
