package bitmovintypes

type StreamMode string

const (
	StreamModeStandard         StreamMode = "STANDARD"
	StreamModePerTitleTemplate StreamMode = "PER_TITLE_TEMPLATE"
	StreamModePerTitleResult   StreamMode = "PER_TITLE_RESULT"
)
