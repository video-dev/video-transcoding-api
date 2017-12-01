package bitmovintypes

type VP8NoiseSensitivity string

const (
	VP8NoiseSensitivityOff             VP8NoiseSensitivity = "OFF"
	VP8NoiseSensitivityOnYOnly         VP8NoiseSensitivity = "ON_Y_ONLY"
	VP8NoiseSensitivityOnYUV           VP8NoiseSensitivity = "ON_YUV"
	VP8NoiseSensitivityOnYUVAggressive VP8NoiseSensitivity = "ON_YUV_AGRESSIVE"
	VP8NoiseSensitivityAdaptive        VP8NoiseSensitivity = "ADAPTIVE"
)
