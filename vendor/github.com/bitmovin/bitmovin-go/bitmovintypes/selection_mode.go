package bitmovintypes

type SelectionMode string

const (
	SelectionModeAuto             SelectionMode = "AUTO"
	SelectionModePositionAbsolute SelectionMode = "POSITION_ABSOLUTE"
	SelectionModeVideoRelative    SelectionMode = "VIDEO_RELATIVE"
	SelectionModeAudioRelative    SelectionMode = "AUDIO_RELATIVE"
)
