package elementalconductor

import (
	"reflect"
	"testing"

	"github.com/NYTimes/encoding-wrapper/elementalconductor"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

func TestFactoryIsRegistered(t *testing.T) {
	_, err := provider.GetProviderFactory(Name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestElementalConductorFactory(t *testing.T) {
	cfg := config.Config{
		ElementalConductor: &config.ElementalConductor{
			Host:        "elemental-server",
			UserLogin:   "myuser",
			APIKey:      "secret-key",
			AuthExpires: 30,
		},
	}
	provider, err := elementalConductorFactory(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	econductorProvider, ok := provider.(*elementalConductorProvider)
	if !ok {
		t.Fatalf("Wrong provider returned. Want elementalConductorProvider instance. Got %#v.", provider)
	}
	expected := elementalconductor.Client{
		Host:        "elemental-server",
		UserLogin:   "myuser",
		APIKey:      "secret-key",
		AuthExpires: 30,
	}
	if !reflect.DeepEqual(*econductorProvider.client, expected) {
		t.Errorf("Factory: wrong client returned. Want %#v. Got %#v.", expected, *econductorProvider.client)
	}
	if !reflect.DeepEqual(*econductorProvider.config, cfg) {
		t.Errorf("Factory: wrong config returned. Want %#v. Got %#v.", cfg, *econductorProvider.config)
	}
}

func TestElementalConductorFactoryValidation(t *testing.T) {
	var tests = []struct {
		host        string
		userLogin   string
		apiKey      string
		authExpires int
	}{
		{"", "", "", 0},
		{"myhost", "", "", 0},
		{"", "myuser", "", 0},
		{"", "", "mykey", 0},
		{"", "", "", 30},
	}
	for _, test := range tests {
		cfg := config.Config{
			ElementalConductor: &config.ElementalConductor{
				Host:        test.host,
				UserLogin:   test.userLogin,
				APIKey:      test.apiKey,
				AuthExpires: test.authExpires,
			},
		}
		provider, err := elementalConductorFactory(&cfg)
		if provider != nil {
			t.Errorf("Unexpected non-nil provider: %#v", provider)
		}
		if err != errElementalConductorInvalidConfig {
			t.Errorf("Wrong error returned. Want errElementalConductorInvalidConfig. Got %#v", err)
		}
	}
}
