package provider

import "github.com/nytm/video-transcoding-api/config"

type fakeProvider struct {
	cap       Capabilities
	healthErr error
}

func (*fakeProvider) Transcode(TranscodeProfile) (*JobStatus, error) {
	return nil, nil
}

func (*fakeProvider) JobStatus(string) (*JobStatus, error) {
	return nil, nil
}

func (*fakeProvider) CreatePreset(Preset) (string, error) {
	return "", nil
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
