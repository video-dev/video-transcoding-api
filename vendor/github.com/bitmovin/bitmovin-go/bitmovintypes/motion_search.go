package bitmovintypes

type MotionSearch string

const (
	MotionSearchDiamond            MotionSearch = "DIA"
	MotionSearchHexagon            MotionSearch = "HEX"
	MotionSearchUnevenMultiHexagon MotionSearch = "UMH"
	MotionSearchStar               MotionSearch = "STAR"
	MotionSearchFull               MotionSearch = "FULL"
)
