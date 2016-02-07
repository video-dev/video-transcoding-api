package elementalcloud

import (
	"reflect"
	"testing"

	"github.com/NYTimes/encoding-wrapper/elementalcloud"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

func TestFactoryIsRegistered(t *testing.T) {
	_, err := provider.GetProviderFactory(Name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestElementalCloudFactory(t *testing.T) {
	cfg := config.Config{
		ElementalCloud: &config.ElementalCloud{
			Host:        "elemental-server",
			UserLogin:   "myuser",
			APIKey:      "secret-key",
			AuthExpires: 30,
		},
	}
	provider, err := elementalCloudFactory(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	ecloudProvider, ok := provider.(*elementalCloudProvider)
	if !ok {
		t.Fatalf("Wrong provider returned. Want elementalCloudProvider instance. Got %#v.", provider)
	}
	expected := elementalcloud.Client{
		Host:        "elemental-server",
		UserLogin:   "myuser",
		APIKey:      "secret-key",
		AuthExpires: 30,
	}
	if !reflect.DeepEqual(*ecloudProvider.client, expected) {
		t.Errorf("Factory: wrong client returned. Want %#v. Got %#v.", expected, *ecloudProvider.client)
	}
	if !reflect.DeepEqual(*ecloudProvider.config, cfg) {
		t.Errorf("Factory: wrong config returned. Want %#v. Got %#v.", cfg, *ecloudProvider.config)
	}
}

func TestElementalCloudFactoryValidation(t *testing.T) {
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
			ElementalCloud: &config.ElementalCloud{
				Host:        test.host,
				UserLogin:   test.userLogin,
				APIKey:      test.apiKey,
				AuthExpires: test.authExpires,
			},
		}
		provider, err := elementalCloudFactory(&cfg)
		if provider != nil {
			t.Errorf("Unexpected non-nil provider: %#v", provider)
		}
		if err != errElementalCloudInvalidConfig {
			t.Errorf("Wrong error returned. Want errElementalCloudInvalidConfig. Got %#v", err)
		}
	}
}
