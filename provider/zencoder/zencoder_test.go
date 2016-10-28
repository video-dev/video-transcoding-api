package zencoder

import (
	"reflect"
	"testing"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/brandscreen/zencoder"
)

func TestFactoryIsRegistered(t *testing.T) {
	_, err := provider.GetProviderFactory(Name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestZencoderFactory(t *testing.T) {
	cfg := config.Config{
		Zencoder: &config.Zencoder{
			APIKey: "api-key-here",
		},
	}
	provider, err := zencoderFactory(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	zencoderProvider, ok := provider.(*zencoderProvider)
	if !ok {
		t.Fatalf("Wrong provider returned. Want zencoderProvider instance. Got %#v.", provider)
	}
	expected := zencoder.NewZencoder("api-key-here")
	if !reflect.DeepEqual(zencoderProvider.client, expected) {
		t.Errorf("Factory: wrong client returned. Want %#v. Got %#v.", expected, zencoderProvider.client)
	}
	if !reflect.DeepEqual(zencoderProvider.config, &cfg) {
		t.Errorf("Factory: wrong config returned. Want %#v. Got %#v.", &cfg, zencoderProvider.config)
	}
}

func TestZencoderFactoryValidation(t *testing.T) {
	cfg := config.Config{Zencoder: &config.Zencoder{APIKey: "api-key"}}
	provider, err := zencoderFactory(&cfg)
	if provider == nil {
		t.Errorf("Unexpected nil provider: %#v", provider)
	}
	if err != nil {
		t.Errorf("Unexpected Error returned. Got %#v", err)
	}

	cfg = config.Config{Zencoder: &config.Zencoder{APIKey: ""}}
	provider, err = zencoderFactory(&cfg)
	if provider != nil {
		t.Errorf("Unexpected non-nil provider: %#v", provider)
	}
	if err != errZencoderInvalidConfig {
		t.Errorf("Wrong error returned. Want errZencoderInvalidConfig. Got %#v", err)
	}
}
