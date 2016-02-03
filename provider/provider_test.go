package provider

import (
	"testing"

	"github.com/nytm/video-transcoding-api/config"
)

func noopFactory(*config.Config) (TranscodingProvider, error) {
	return nil, nil
}

func TestRegister(t *testing.T) {
	providers = nil
	err := Register("noop", noopFactory)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := providers["noop"]; !ok {
		t.Errorf("expected to get the noop factory register. Got map %#v", providers)
	}
}

func TestRegisterMultiple(t *testing.T) {
	providers = nil
	err := Register("noop", noopFactory)
	if err != nil {
		t.Fatal(err)
	}
	err = Register("noope", noopFactory)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := providers["noop"]; !ok {
		t.Errorf("expected to get the noop factory register. Got map %#v", providers)
	}
	if _, ok := providers["noope"]; !ok {
		t.Errorf("expected to get the noope factory register. Got map %#v", providers)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	providers = nil
	err := Register("noop", noopFactory)
	if err != nil {
		t.Fatal(err)
	}
	err = Register("noop", noopFactory)
	if err != ErrProviderAlreadyRegistered {
		t.Errorf("Got wrong error when registering provider twice. Want %#v. Got %#v", ErrProviderAlreadyRegistered, err)
	}
}

func TestGetProviderFactory(t *testing.T) {
	providers = nil
	var called bool
	err := Register("noop", func(*config.Config) (TranscodingProvider, error) {
		called = true
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	factory, err := GetProviderFactory("noop")
	if err != nil {
		t.Fatal(err)
	}
	factory(nil)
	if !called {
		t.Errorf("Did not call the expected factory. Got %#v", factory)
	}
}

func TestGetProviderFactoryNotRegistered(t *testing.T) {
	providers = nil
	factory, err := GetProviderFactory("noop")
	if factory != nil {
		t.Errorf("Got unexpected non-nil factory: %#v", factory)
	}
	if err != ErrProviderNotFound {
		t.Errorf("Got wrong error when getting an unregistered provider. Want %#v. Got %#v", ErrProviderNotFound, err)
	}
}
