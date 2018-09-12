package bitmovintypes

type DecodingErrorMode string

const (
	DecodingErrorModeFailOnError     DecodingErrorMode = "FAIL_ON_ERROR"
	DecodingErrorModeDuplicateFrames DecodingErrorMode = "DUPLICATE_FRAMES"
)
