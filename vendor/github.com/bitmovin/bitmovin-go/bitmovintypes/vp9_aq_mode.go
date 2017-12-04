package bitmovintypes

type VP9AQMode string

const (
	VP9AQModeNone       VP9AQMode = "NONE"
	VP9AQModeVariance   VP9AQMode = "VARIANCE"
	VP9AQModeComplexity VP9AQMode = "COMPLEXITY"
	VP9AQModeCyclic     VP9AQMode = "CYCLIC"
)
