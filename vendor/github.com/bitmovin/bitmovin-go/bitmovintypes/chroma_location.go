package bitmovintypes

type ChromaLocation string

const (
	ChromaLocationUnspecified ChromaLocation = "UNSPECIFIED"
	ChromaLocationLeft        ChromaLocation = "LEFT"
	ChromaLocationCenter      ChromaLocation = "CENTER"
	ChromaLocationTopLeft     ChromaLocation = "TOPLEFT"
	ChromaLocationTop         ChromaLocation = "TOP"
	ChromaLocationBottomLeft  ChromaLocation = "BOTTOMLEFT"
	ChromaLocationBottom      ChromaLocation = "BOTTOM"
)
