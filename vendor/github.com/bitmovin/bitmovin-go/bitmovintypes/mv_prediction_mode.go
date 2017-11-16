package bitmovintypes

type MVPredictionMode string

const (
	MVPredictionModeNone     MVPredictionMode = "NONE"
	MVPredictionModeSpatial  MVPredictionMode = "SPATIAL"
	MVPredictionModeTemporal MVPredictionMode = "TEMPORAL"
	MVPredictionModeAuto     MVPredictionMode = "AUTO"
)
