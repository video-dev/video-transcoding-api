package encodingcom

import (
	"reflect"
	"testing"
	"time"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
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
	expected := &encodingcom.Client{
		Endpoint: "https://manage.encoding.com",
		UserID:   "myuser",
		UserKey:  "secret-key",
	}
	if !reflect.DeepEqual(ecomProvider.client.(*encodingcom.Client), expected) {
		t.Errorf("Factory: wrong client returned. Want %#v. Got %#v.", expected, ecomProvider.client)
	}
	if !reflect.DeepEqual(ecomProvider.config, &cfg) {
		t.Errorf("Factory: wrong config returned. Want %#v. Got %#v.", &cfg, ecomProvider.config)
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
	presets := []db.PresetMap{
		{
			Name: "webm_720p",
			ProviderMapping: map[string]string{
				Name:           "123455",
				"not-relevant": "something",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
		{
			Name: "webm_480p",
			ProviderMapping: map[string]string{
				Name:           "123456",
				"not-relevant": "otherthing",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
		{
			Name: "mp4_1080p",
			ProviderMapping: map[string]string{
				Name:           "321321",
				"not-relevant": "allthings",
			},
			OutputOpts: db.OutputOptions{Extension: "mp4"},
		},
		{
			Name: "hls_1080p",
			ProviderMapping: map[string]string{
				Name:           "321322",
				"not-relevant": "allthings",
			},
			OutputOpts: db.OutputOptions{Extension: "m3u8"},
		},
	}

	transcodeProfile := provider.TranscodeProfile{
		SourceMedia: source,
		Presets:     presets,
		StreamingParams: provider.StreamingParams{
			Protocol:        "hls",
			SegmentDuration: 3,
		},
	}

	jobStatus, err := prov.Transcode(&db.Job{ID: "job-123"}, transcodeProfile)
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
	dest := prov.config.EncodingCom.Destination
	expectedFormats := []encodingcom.Format{
		{
			Output:          []string{"123455"},
			Destination:     []string{dest + "job-123/webm_720p/video.webm"},
			SegmentDuration: 3,
		},
		{
			Output:          []string{"123456"},
			Destination:     []string{dest + "job-123/webm_480p/video.webm"},
			SegmentDuration: 3,
		},
		{
			Output:          []string{"321321"},
			Destination:     []string{dest + "job-123/mp4_1080p/video.mp4"},
			SegmentDuration: 3,
		},
		{
			Output:          []string{"321322"},
			Destination:     []string{dest + "job-123/hls_1080p/video/master.m3u8"},
			SegmentDuration: 3,
		},
	}
	if !reflect.DeepEqual(media.Request.Format, expectedFormats) {
		t.Errorf("Wrong format.\nWant %#v\nGot  %#v.", expectedFormats, media.Request.Format)
	}
	if !reflect.DeepEqual([]string{source}, media.Request.Source) {
		t.Errorf("Wrong source. Want %v. Got %v.", []string{source}, media.Request.Source)
	}
}

func TestEncodingComTranscodePresetNotFound(t *testing.T) {
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
	presets := []db.PresetMap{
		{
			Name: "webm_720p",
			ProviderMapping: map[string]string{
				Name:           "123455",
				"not-relevant": "something",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
		{
			Name: "webm_480p",
			ProviderMapping: map[string]string{
				"not-relevant": "otherthing",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
	}

	transcodeProfile := provider.TranscodeProfile{
		SourceMedia: source,
		Presets:     presets,
		StreamingParams: provider.StreamingParams{
			SegmentDuration: 3,
		},
	}

	jobStatus, err := prov.Transcode(&db.Job{ID: "job-2"}, transcodeProfile)
	if err != provider.ErrPresetMapNotFound {
		t.Errorf("Wrong error. Want %#v. Got %#v", provider.ErrPresetMapNotFound, err)
	}
	if jobStatus != nil {
		t.Errorf("Got unexpected non-nil JobStatus: %#v", jobStatus)
	}
}

func TestJobStatus(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	now := time.Now().UTC().Truncate(time.Second)
	media := fakeMedia{
		ID:       "mymedia",
		Status:   "Finished",
		Created:  now.Add(-time.Hour),
		Started:  now.Add(-50 * time.Minute),
		Finished: now.Add(-10 * time.Minute),
	}
	server.medias["mymedia"] = &media
	client := newEncodingComFakeClient(media)
	prov := encodingComProvider{client: client}
	prov.config = &config.Config{
		EncodingCom: &config.EncodingCom{
			Destination: "mybucket",
		},
	}
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
			"progress":   100.0,
			"sourcefile": "http://some.source.file",
			"timeleft":   "1",
			"created":    media.Created,
			"started":    media.Started,
			"finished":   media.Finished,
			"destinationStatus": []encodingcom.DestinationStatus{
				{
					Name:   "s3://mybucket/dir/file.mp4",
					Status: "Saved",
				},
			},
		},
		OutputDestination: "s3://mybucket/dir",
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

func TestCapabilities(t *testing.T) {
	var prov encodingComProvider
	expected := provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls", "webm"},
		Destinations:  []string{"akamai", "s3"},
	}
	cap := prov.Capabilities()
	if !reflect.DeepEqual(cap, expected) {
		t.Errorf("Capabilities: want %#v. Got %#v", expected, cap)
	}
}
