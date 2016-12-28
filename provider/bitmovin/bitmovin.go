package bitmovin

import "github.com/NYTimes/video-transcoding-api/provider"

// Name is the name used for registering the bitmovin provider in the
// registry of providers.
const Name = "bitmovin"

type bitmovinProvider struct {
}

// func (p *bitmovinProvider) CreatePreset(db.Preset) (string, error) {

// }

// func (p *bitmovinProvider) DeletePreset(presetID string) error {

// }

// func (p *bitmovinProvider) GetPreset(presetID string) (interface{}, error) {

// }

// func (p *bitmovinProvider) Transcode(db.Job) (*provider.JobStatus, error) {

// }

// func (p *bitmovinProvider) JobStatus(db.Job) (*provider.JobStatus, error) {

// }

// func (p *bitmovinProvider) CancelJob(jobID string) error {

// }

// func (p *bitmovinProvider) Healthcheck() error {

// }

func (p *bitmovinProvider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls"},
		Destinations:  []string{"s3"},
	}
}
