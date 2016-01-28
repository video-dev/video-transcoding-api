package provider

import (
	"reflect"
	"testing"
	"time"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
)

func TestFactory(t *testing.T) {
	cfg := config.Config{
		EncodingCom: config.EncodingCom{
			UserID:  "myuser",
			UserKey: "secret-key",
		},
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
	if !reflect.DeepEqual(*ecomProvider.config, cfg) {
		t.Errorf("Factory: wrong config returned. Want %#v. Got %#v.", cfg, *ecomProvider.config)
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
		cfg := config.Config{
			EncodingCom: config.EncodingCom{UserID: test.userID, UserKey: test.userKey},
		}
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
	provider := encodingComProvider{
		client: client,
		config: &config.Config{
			EncodingCom: config.EncodingCom{
				Destination: "https://mybucket.s3.amazonaws.com/destination-dir/",
			},
		},
	}
	source := "http://some.nice/video.mp4"
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
		TwoPassEncoding:     true,
	}
	jobStatus, err := provider.Transcode(source, profile)
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
		Destination:         []string{provider.config.Destination + "video.mp4"},
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
		TwoPass:             encodingcom.YesNoBoolean(true),
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

func TestJobStatus(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	now := time.Now().In(time.UTC).Truncate(time.Second)
	media := fakeMedia{
		ID:       "mymedia",
		Status:   "Finished",
		Created:  now.Add(-time.Hour),
		Started:  now.Add(-50 * time.Minute),
		Finished: now.Add(-10 * time.Minute),
	}
	server.medias["mymedia"] = &media
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	provider := encodingComProvider{client: client}
	jobStatus, err := provider.JobStatus("mymedia")
	if err != nil {
		t.Fatal(err)
	}
	expected := JobStatus{
		ProviderJobID: "mymedia",
		ProviderName:  "encoding.com",
		Status:        StatusFinished,
		StatusMessage: "",
		ProviderStatus: map[string]interface{}{
			"progress":   100.0,
			"sourcefile": "http://some.source.file",
			"timeleft":   "1",
			"created":    media.Created,
			"started":    media.Started,
			"finished":   media.Finished,
		},
	}
	if !reflect.DeepEqual(*jobStatus, expected) {
		t.Errorf("JobStatus: wrong job returned.\nWant %#v.\nGot  %#v.", expected, *jobStatus)
	}
}

func TestJobStatusMediaNotFound(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	provider := encodingComProvider{client: client}
	jobStatus, err := provider.JobStatus("non-existent-job")
	if err == nil {
		t.Errorf("JobStatus: got unexpected <nil> err.")
	}
	if jobStatus != nil {
		t.Errorf("JobStatus: got unexpected non-nil result: %#v", jobStatus)
	}
}

func TestJobStatusMap(t *testing.T) {
	var tests = []struct {
		encodingComStatus string
		expected          status
	}{
		{"New", StatusQueued},
		{"Downloading", StatusStarted},
		{"Ready to process", StatusStarted},
		{"Waiting for encoder", StatusStarted},
		{"Processing", StatusStarted},
		{"Saving", StatusStarted},
		{"Finished", StatusFinished},
		{"Error", StatusFailed},
		{"Unknown", StatusFailed},
	}
	var p encodingComProvider
	for _, test := range tests {
		got := p.statusMap(test.encodingComStatus)
		if got != test.expected {
			t.Errorf("statusMap(%q): wrong value. Want %q. Got %q", test.encodingComStatus, test.expected, got)
		}
	}
}
