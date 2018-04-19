package bitmovintypes

type ConditionAttribute string

const (
	ConditionAttributeHeight      ConditionAttribute = "HEIGHT"
	ConditionAttributeWidth       ConditionAttribute = "WIDTH"
	ConditionAttributeFPS         ConditionAttribute = "FPS"
	ConditionAttributeBitrate     ConditionAttribute = "BITRATE"
	ConditionAttributeAspectRatio ConditionAttribute = "ASPECTRATIO"
	ConditionAttributeInputStream ConditionAttribute = "INPUTSTREAM"
)

type ConditionType string

const (
	ConditionTypeAnd       ConditionType = "AND"
	ConditionTypeOr        ConditionType = "OR"
	ConditionTypeCondition ConditionType = "CONDITION"
)

type ConditionMode string

const (
	ConditionModeDropMuxing ConditionMode = "DROP_MUXING"
	ConditionModeDropStream ConditionMode = "DROP_STREAM"
)
