package bitmovintypes

type MediaType string

const (
	MediaTypeAudio          MediaType = "AUDIO"
	MediaTypeVideo          MediaType = "VIDEO"
	MediaTypeSubtitles      MediaType = "SUBTITLES"
	MediaTypeClosedCaptions MediaType = "CLOSED_CAPTIONS"
)
