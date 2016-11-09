package provider

import (
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
)

type fakeProvider struct {
	cap       Capabilities
	healthErr error
}

func (*fakeProvider) Transcode(*db.Job) (*JobStatus, error) {
	return nil, nil
}

func (*fakeProvider) JobStatus(*db.Job) (*JobStatus, error) {
	return nil, nil
}

func (*fakeProvider) CreatePreset(db.Preset) (string, error) {
	return "", nil
}

func (*fakeProvider) GetPreset(string) (interface{}, error) {
	return "", nil
}

func (*fakeProvider) DeletePreset(string) error {
	return nil
}

func (*fakeProvider) CancelJob(string) error {
	return nil
}

func (f *fakeProvider) Healthcheck() error {
	return f.healthErr
}

func (f *fakeProvider) Capabilities() Capabilities {
	return f.cap
}

func getFactory(fErr error, healthErr error, capabilities Capabilities) Factory {
	return func(*config.Config) (TranscodingProvider, error) {
		if fErr != nil {
			return nil, fErr
		}
		return &fakeProvider{healthErr: healthErr, cap: capabilities}, nil
	}
}
