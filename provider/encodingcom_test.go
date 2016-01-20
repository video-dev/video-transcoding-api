package provider

import (
	"reflect"
	"testing"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
)

func TestFactory(t *testing.T) {
	cfg := config.Config{
		EncodingComUserID:  "myuser",
		EncodingComUserKey: "secret-key",
	}
	provider, err := EncodingComProvider(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	ecomProvider, ok := provider.(*encodingComProvider)
	if !ok {
		t.Fatalf("Wrong provider returned. Want encodingComProvider instance. Got %#v.", provider)
	}
	expected := encodingcom.Client{
		Endpoint: "https://manage.encoding.com",
		UserID:   "myuser",
		UserKey:  "secret-key",
	}
	if !reflect.DeepEqual(*ecomProvider.client, expected) {
		t.Errorf("Factory: wrong client returned. Want %#v. Got %#v.", expected, *ecomProvider.client)
	}
}

func TestFactoryValidation(t *testing.T) {
	var tests = []struct {
		userID  string
		userKey string
	}{
		{"", ""},
		{"", "mykey"},
		{"myuser", ""},
	}
	for _, test := range tests {
		cfg := config.Config{EncodingComUserID: test.userID, EncodingComUserKey: test.userKey}
		provider, err := EncodingComProvider(&cfg)
		if provider != nil {
			t.Errorf("Unexpected non-nil provider: %#v", provider)
		}
		if err != ErrMissingData {
			t.Errorf("Wrong error returned. Want ErrMissingData. Got %#v", err)
		}
	}
}
