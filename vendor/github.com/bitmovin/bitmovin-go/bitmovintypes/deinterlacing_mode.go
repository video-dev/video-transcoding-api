package bitmovintypes

type DeinterlacingMode string

const (
	DeinterlacingModeFrame          DeinterlacingMode = "FRAME"
	DeinterlacingModeField          DeinterlacingMode = "FIELD"
	DeinterlacingModeFrameNoSpatial DeinterlacingMode = "FRAME_NOSPATIAL"
	DeinterlacingModeFieldNoSpatial DeinterlacingMode = "FIELD_NOSPATIAL"
)
