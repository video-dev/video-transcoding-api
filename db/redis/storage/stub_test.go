package storage

import "time"

type Person struct {
	ID               string    `redis-hash:"-"`
	Name             string    `redis-hash:"name"`
	Address          Address   `redis-hash:"address,expand"`
	Age              uint      `redis-hash:"age"`
	Weight           float64   `redis-hash:"weight"`
	BirthTime        time.Time `redis-hash:"birth"`
	PreferredColors  []string  `redis-hash:"colors"`
	NonTagged        string
	unexported       string
	unexportedTagged string `redis-hash:"unexported"`
}

type Address struct {
	Data   map[string]string `redis-hash:"data,expand"`
	Number int               `redis-hash:"number"`
	Main   bool              `redis-hash:"main"`
	City   *City             `redis-hash:"city,expand"`
}

type City struct {
	Name string `redis-hash:"name"`
}

type InvalidStruct struct {
	Name string `redis-hash:"name,expand"`
}

type InvalidInnerStruct struct {
	Data map[string]int `redis-hash:"data,expand"`
}

type Job struct {
	ID              string            `redis-hash:"jobID"`
	ProviderName    string            `redis-hash:"providerName"`
	ProviderJobID   string            `redis-hash:"providerJobID"`
	StreamingParams StreamingParams   `redis-hash:"streamingparams,expand"`
	CreationTime    time.Time         `redis-hash:"creationTime"`
	SourceMedia     string            `redis-hash:"source"`
	Outputs         []TranscodeOutput `redis-hash:"-"`
}

type TranscodeOutput struct {
	Preset   PresetMap `redis-hash:"presetmap,expand"`
	FileName string    `redis-hash:"filename"`
}

type PresetMap struct {
	Name string `redis-hash:"presetmap_name"`
}

type StreamingParams struct {
	SegmentDuration  uint   `redis-hash:"segmentDuration"`
	Protocol         string `redis-hash:"protocol"`
	PlaylistFileName string `redis-hash:"playlistFileName"`
}

type LocalPreset struct {
	Name   string `redis-hash:"-"`
	Preset Preset `redis-hash:"preset,expand"`
}

type Preset struct {
	Name        string      `redis-hash:"name"`
	Description string      `redis-hash:"description,omitempty"`
	Container   string      `redis-hash:"container,omitempty"`
	RateControl string      `redis-hash:"ratecontrol,omitempty"`
	Video       VideoPreset `redis-hash:"video,expand"`
	Audio       AudioPreset `redis-hash:"audio,expand"`
}

type VideoPreset struct {
	Profile       string `redis-hash:"profile,omitempty"`
	ProfileLevel  string `redis-hash:"profilelevel,omitempty"`
	Width         string `redis-hash:"width,omitempty"`
	Height        string `redis-hash:"height,omitempty"`
	Codec         string `redis-hash:"codec,omitempty"`
	Bitrate       string `redis-hash:"bitrate,omitempty"`
	GopSize       string `redis-hash:"gopsize,omitempty"`
	GopMode       string `redis-hash:"gopmode,omitempty"`
	InterlaceMode string `redis-hash:"interlacemode,omitempty"`
}

type AudioPreset struct {
	Codec   string `redis-hash:"codec,omitempty"`
	Bitrate string `redis-hash:"bitrate,omitempty"`
}
