package provider

import (
	"errors"
	"reflect"
	"testing"

	"github.com/video-dev/video-transcoding-api/v2/config"
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

func TestListProviders(t *testing.T) {
	cap := Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls"},
		Destinations:  []string{"s3", "akamai"},
	}
	providers = map[string]Factory{
		"cap-and-unhealthy": getFactory(nil, errors.New("api is down"), cap),
		"factory-err":       getFactory(errors.New("invalid config"), nil, cap),
		"cap-and-healthy":   getFactory(nil, nil, cap),
	}
	expected := []string{"cap-and-healthy", "cap-and-unhealthy"}
	got := ListProviders(&config.Config{})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("DescribeProviders: want %#v. Got %#v", expected, got)
	}
}

func TestListProvidersEmpty(t *testing.T) {
	providers = nil
	providerNames := ListProviders(&config.Config{})
	if len(providerNames) != 0 {
		t.Errorf("Unexpected non-empty provider list: %#v", providerNames)
	}
}

func TestDescribeProvider(t *testing.T) {
	cap := Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls"},
		Destinations:  []string{"s3", "akamai"},
	}
	providers = map[string]Factory{
		"cap-and-unhealthy": getFactory(nil, errors.New("api is down"), cap),
		"factory-err":       getFactory(errors.New("invalid config"), nil, cap),
		"cap-and-healthy":   getFactory(nil, nil, cap),
	}
	tests := []struct {
		input    string
		expected Description
	}{
		{
			"factory-err",
			Description{Name: "factory-err", Enabled: false},
		},
		{
			"cap-and-healthy",
			Description{
				Name:         "cap-and-healthy",
				Capabilities: cap,
				Health:       Health{OK: true},
				Enabled:      true,
			},
		},
		{
			"cap-and-unhealthy",
			Description{
				Name:         "cap-and-unhealthy",
				Capabilities: cap,
				Health:       Health{OK: false, Message: "api is down"},
				Enabled:      true,
			},
		},
	}
	for _, test := range tests {
		description, err := DescribeProvider(test.input, &config.Config{})
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(*description, test.expected) {
			t.Errorf("DescribeProvider(%q): want %#v. Got %#v", test.input, test.expected, *description)
		}
	}
}

func TestDescribeProviderNotFound(t *testing.T) {
	providers = nil
	description, err := DescribeProvider("anything", nil)
	if err != ErrProviderNotFound {
		t.Errorf("Wrong error. Want %#v. Got %#v", ErrProviderNotFound, err)
	}
	if description != nil {
		t.Errorf("Unexpected non-nil description: %#v", description)
	}
}
