package bitmovintypes

type Trellis string

const (
	TrellisDisabled       = `DISABLED`
	TrellisEnabledFinalMB = `ENABLED_FINAL_MB`
	TrellisEnabledAll     = `ENABLED_ALL`
	TrellisDefault        = TrellisEnabledFinalMB
)
