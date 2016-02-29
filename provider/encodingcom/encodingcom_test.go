package encodingcom

import (
	"reflect"
	"testing"
	"time"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

func TestFactoryIsRegistered(t *testing.T) {
	_, err := provider.GetProviderFactory(Name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncodingComFactory(t *testing.T) {
	cfg := config.Config{
		EncodingCom: &config.EncodingCom{
			UserID:  "myuser",
			UserKey: "secret-key",
		},
	}
	provider, err := encodingComFactory(&cfg)
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

func TestEncodingComFactoryValidation(t *testing.T) {
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
			EncodingCom: &config.EncodingCom{UserID: test.userID, UserKey: test.userKey},
		}
		provider, err := encodingComFactory(&cfg)
		if provider != nil {
			t.Errorf("Unexpected non-nil provider: %#v", provider)
		}
		if err != errEncodingComInvalidConfig {
			t.Errorf("Wrong error returned. Want errEncodingComInvalidConfig. Got %#v", err)
		}
	}
}

func TestEncodingComTranscode(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	prov := encodingComProvider{
		client: client,
		config: &config.Config{
			EncodingCom: &config.EncodingCom{
				Destination: "https://mybucket.s3.amazonaws.com/destination-dir/",
			},
		},
	}
	source := "http://some.nice/video.mp4"
	profile := provider.Profile{
		Output:              []string{"webm", "hls"},
		Size:                provider.Size{Height: 360},
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
	jobStatus, err := prov.TranscodeWithProfiles(source, []provider.Profile{profile})
	if err != nil {
		t.Fatal(err)
	}
	if expected := "it worked"; jobStatus.StatusMessage != expected {
		t.Errorf("wrong StatusMessage. Want %q. Got %q", expected, jobStatus.StatusMessage)
	}
	if jobStatus.ProviderName != Name {
		t.Errorf("wrong ProviderName. Want %q. Got %q", Name, jobStatus.ProviderName)
	}
	media, err := server.getMedia(jobStatus.ProviderJobID)
	if err != nil {
		t.Fatal(err)
	}
	expectedDestination := []string{
		prov.config.EncodingCom.Destination + "video_360p.webm",
		prov.config.EncodingCom.Destination + "video_hls/video.m3u8",
	}
	expectedFormat := encodingcom.Format{
		Output:              []string{"webm", "advanced_hls"},
		Destination:         expectedDestination,
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
	if !reflect.DeepEqual(media.Request.Format[0], expectedFormat) {
		t.Errorf("Wrong format. Want %#v. Got %#v.", expectedFormat, media.Request.Format[0])
	}
	if !reflect.DeepEqual([]string{source}, media.Request.Source) {
		t.Errorf("Wrong source. Want %v. Got %v.", []string{source}, media.Request.Source)
	}
}

func TestProfileToFormatRotation(t *testing.T) {
	var tests = []struct {
		r        provider.Rotation
		expected string
	}{
		{provider.Rotate0Degrees, "0"},
		{provider.Rotate90Degrees, "90"},
		{provider.Rotate180Degrees, "180"},
		{provider.Rotate270Degrees, "270"},
		{provider.Rotation{}, "def"},
	}
	var p encodingComProvider
	for _, test := range tests {
		profile := provider.Profile{Rotate: test.r}
		formats := p.profilesToFormats("sourceFile", []provider.Profile{profile})
		for _, format := range formats {
			if format.Rotate != test.expected {
				t.Errorf("profileToFormat: expected rotate to be %q. Got %q.", test.expected, format.Rotate)
			}
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
	prov := encodingComProvider{client: client}
	jobStatus, err := prov.JobStatus("mymedia")
	if err != nil {
		t.Fatal(err)
	}
	expected := provider.JobStatus{
		ProviderJobID: "mymedia",
		ProviderName:  "encoding.com",
		Status:        provider.StatusFinished,
		StatusMessage: "",
		ProviderStatus: map[string]interface{}{
			"progress":          100.0,
			"sourcefile":        "http://some.source.file",
			"timeleft":          "1",
			"created":           media.Created,
			"started":           media.Started,
			"finished":          media.Finished,
			"destinationStatus": []encodingcom.DestinationStatus(nil),
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
		expected          provider.Status
	}{
		{"New", provider.StatusQueued},
		{"Downloading", provider.StatusStarted},
		{"Ready to process", provider.StatusStarted},
		{"Waiting for encoder", provider.StatusStarted},
		{"Processing", provider.StatusStarted},
		{"Saving", provider.StatusStarted},
		{"Finished", provider.StatusFinished},
		{"Error", provider.StatusFailed},
		{"Unknown", provider.StatusUnknown},
		{"new", provider.StatusQueued},
		{"downloading", provider.StatusStarted},
		{"ready to process", provider.StatusStarted},
		{"waiting for encoder", provider.StatusStarted},
		{"processing", provider.StatusStarted},
		{"saving", provider.StatusStarted},
		{"finished", provider.StatusFinished},
		{"error", provider.StatusFailed},
		{"unknown", provider.StatusUnknown},
	}
	var p encodingComProvider
	for _, test := range tests {
		got := p.statusMap(test.encodingComStatus)
		if got != test.expected {
			t.Errorf("statusMap(%q): wrong value. Want %q. Got %q", test.encodingComStatus, test.expected, got)
		}
	}
}

func TestHealthcheck(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	provider := encodingComProvider{
		client: client,
		config: &config.Config{
			EncodingCom: &config.EncodingCom{StatusEndpoint: server.URL},
		},
	}
	var tests = []struct {
		apiStatus   encodingcom.APIStatusResponse
		expectedMsg string
	}{
		{
			encodingcom.APIStatusResponse{Status: "Ok", StatusCode: "ok"},
			"",
		},
		{
			encodingcom.APIStatusResponse{
				Status:     "Investigation",
				StatusCode: "queue_slow",
				Incident:   "Our encoding queue is processing slower than normal.  Check back for updates.",
			},
			"Status code: queue_slow.\nIncident: Our encoding queue is processing slower than normal.  Check back for updates.\nStatus: Investigation",
		},
		{
			encodingcom.APIStatusResponse{
				Status:     "Maintenance",
				StatusCode: "deploy",
				Incident:   "We are currently working within a scheduled maintenance window.  Check back for updates.",
			},
			"Status code: deploy.\nIncident: We are currently working within a scheduled maintenance window.  Check back for updates.\nStatus: Maintenance",
		},
	}
	for _, test := range tests {
		server.SetAPIStatus(&test.apiStatus)
		err := provider.Healthcheck()
		if test.expectedMsg != "" {
			if got := err.Error(); got != test.expectedMsg {
				t.Errorf("Wrong error returned. Want %q. Got %q", test.expectedMsg, got)
			}
		} else if err != nil {
			t.Errorf("Got unexpected non-nil error: %#v", err)
		}
	}
}
