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

func TestEncodingComTranscode(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	provider := encodingComProvider{client: client}
	source := "http://some.nice/video.mp4"
	destination := "http://some.nice.transcoded/video.mp4"
	profile := Profile{
		Output:              "webm",
		Size:                Size{Height: 360},
		AudioCodec:          "libvorbis",
		AudioBitRate:        "64k",
		AudioChannelsNumber: "2",
		AudioSampleRate:     48000,
		BitRate:             "900k",
		FrameRate:           "30",
		KeepAspectRatio:     true,
		VideoCodec:          "libvpx",
		KeyFrame:            "90",
		AudioVolume:         100,
	}
	jobStatus, err := provider.Transcode(source, destination, profile)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "it worked"; jobStatus.StatusMessage != expected {
		t.Errorf("wrong StatusMessage. Want %q. Got %q", expected, jobStatus.StatusMessage)
	}
	media, err := server.getMedia(jobStatus.ProviderJobID)
	if err != nil {
		t.Fatal(err)
	}
	expectedFormat := encodingcom.Format{
		Output:              []string{"webm"},
		Destination:         []string{destination},
		Size:                "0x360",
		AudioCodec:          "libvorbis",
		AudioBitrate:        "64k",
		AudioChannelsNumber: "2",
		AudioSampleRate:     48000,
		Bitrate:             "900k",
		Framerate:           "30",
		KeepAspectRatio:     encodingcom.YesNoBoolean(true),
		VideoCodec:          "libvpx",
		Keyframe:            []string{"90"},
		AudioVolume:         100,
		Rotate:              "def",
	}
	if !reflect.DeepEqual(*media.Request.Format, expectedFormat) {
		t.Errorf("Wrong format. Want %#v. Got %#v.", expectedFormat, *media.Request.Format)
	}
	if !reflect.DeepEqual([]string{source}, media.Request.Source) {
		t.Errorf("Wrong source. Want %v. Got %v.", []string{source}, media.Request.Source)
	}
}

func TestProfileToFormatRotation(t *testing.T) {
	var tests = []struct {
		r        rotation
		expected string
	}{
		{Rotate0Degrees, "0"},
		{Rotate90Degrees, "90"},
		{Rotate180Degrees, "180"},
		{Rotate270Degrees, "270"},
		{rotation{}, "def"},
	}
	var p encodingComProvider
	for _, test := range tests {
		profile := Profile{Rotate: test.r}
		format := p.profileToFormat(profile)
		if format.Rotate != test.expected {
			t.Errorf("profileToFormat: expected rotate to be %q. Got %q.", test.expected, format.Rotate)
		}
	}
}
