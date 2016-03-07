package elementalconductor

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/NYTimes/encoding-wrapper/elementalconductor"
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

func TestElementalNewJob(t *testing.T) {
	elementalConductorConfig := config.Config{
		ElementalConductor: &config.ElementalConductor{
			Host:            "https://mybucket.s3.amazonaws.com/destination-dir/",
			UserLogin:       "myuser",
			APIKey:          "elemental-api-key",
			AuthExpires:     30,
			AccessKeyID:     "aws-access-key",
			SecretAccessKey: "aws-secret-key",
			Destination:     "s3://destination",
		},
	}
	prov, err := elementalConductorFactory(&elementalConductorConfig)
	if err != nil {
		t.Fatal(err)
	}
	presetProvider, ok := prov.(*elementalConductorProvider)
	if !ok {
		t.Fatal("Could not type assert test provider to elementalConductorProvider")
	}
	source := "http://some.nice/video.mov"
	presets := []db.Preset{
		{
			Name:            "webm_720p",
			ProviderMapping: map[string]string{Name: "10", "other": "not relevant"},
			OutputOpts:      db.OutputOptions{Extension: "webm"},
		},
		{
			Name:            "mp4_720p",
			ProviderMapping: map[string]string{Name: "15", "other": "not relevant"},
			OutputOpts:      db.OutputOptions{Extension: "mp4"},
		},
		{
			Name:            "mp4_1080p",
			ProviderMapping: map[string]string{Name: "20", "other": "not relevant"},
			OutputOpts:      db.OutputOptions{Extension: ""},
		},
	}

	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     source,
		Presets:         presets,
		StreamingParams: provider.StreamingParams{},
	}
	newJob, err := presetProvider.newJob(transcodeProfile)
	if err != nil {
		t.Error(err)
	}
	expectedJob := elementalconductor.Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Input: elementalconductor.Input{
			FileInput: elementalconductor.Location{
				URI:      "http://some.nice/video.mov",
				Username: "aws-access-key",
				Password: "aws-secret-key",
			},
		},
		Priority: 50,
		OutputGroup: elementalconductor.OutputGroup{
			Order: 1,
			FileGroupSettings: elementalconductor.FileGroupSettings{
				Destination: elementalconductor.Location{
					URI:      "s3://destination/video",
					Username: "aws-access-key",
					Password: "aws-secret-key",
				},
			},
			Type: elementalconductor.FileOutputGroupType,
			Output: []elementalconductor.Output{
				{
					StreamAssemblyName: "stream_0",
					NameModifier:       "_webm_720p",
					Order:              0,
					Container:          elementalconductor.Container("webm"),
				},
				{
					StreamAssemblyName: "stream_1",
					NameModifier:       "_mp4_720p",
					Order:              1,
					Container:          elementalconductor.MPEG4,
				},
				{
					StreamAssemblyName: "stream_2",
					NameModifier:       "_mp4_1080p",
					Order:              2,
					Container:          defaultContainer,
				},
			},
		},
		StreamAssembly: []elementalconductor.StreamAssembly{
			{
				Name:   "stream_0",
				Preset: "10",
			},
			{
				Name:   "stream_1",
				Preset: "15",
			},
			{
				Name:   "stream_2",
				Preset: "20",
			},
		},
	}
	if !reflect.DeepEqual(&expectedJob, newJob) {
		t.Errorf("New job not according to spec.\nWanted %#v.\nGot    %#v.", &expectedJob, newJob)
	}
}

func TestElementalNewJobAdaptiveStreaming(t *testing.T) {
	elementalConductorConfig := config.Config{
		ElementalConductor: &config.ElementalConductor{
			Host:            "https://mybucket.s3.amazonaws.com/destination-dir/",
			UserLogin:       "myuser",
			APIKey:          "elemental-api-key",
			AuthExpires:     30,
			AccessKeyID:     "aws-access-key",
			SecretAccessKey: "aws-secret-key",
			Destination:     "s3://destination",
		},
	}
	prov, err := elementalConductorFactory(&elementalConductorConfig)
	if err != nil {
		t.Fatal(err)
	}
	presetProvider, ok := prov.(*elementalConductorProvider)
	if !ok {
		t.Fatal("Could not type assert test provider to elementalConductorProvider")
	}
	source := "http://some.nice/video.mov"
	presets := []db.Preset{
		{
			Name:            "hls_360p",
			ProviderMapping: map[string]string{Name: "15", "other": "not relevant"},
			OutputOpts:      db.OutputOptions{Extension: "hls"},
		},
		{
			Name:            "hls_480p",
			ProviderMapping: map[string]string{Name: "20", "other": "not relevant"},
			OutputOpts:      db.OutputOptions{Extension: "ts"},
		},
		{
			Name:            "hls_720p",
			ProviderMapping: map[string]string{Name: "25", "other": "not relevant"},
			OutputOpts:      db.OutputOptions{Extension: "m3u8"},
		},
		{
			Name:            "hls_1080p",
			ProviderMapping: map[string]string{Name: "30", "other": "not relevant"},
			OutputOpts:      db.OutputOptions{Extension: ".ts"},
		},
	}
	transcodeProfile := provider.TranscodeProfile{
		SourceMedia: source,
		Presets:     presets,
		StreamingParams: provider.StreamingParams{
			SegmentDuration: 3,
		},
	}
	newJob, err := presetProvider.newJob(transcodeProfile)
	if err != nil {
		t.Error(err)
	}
	expectedJob := elementalconductor.Job{
		XMLName: xml.Name{
			Local: "job",
		},
		Input: elementalconductor.Input{
			FileInput: elementalconductor.Location{
				URI:      "http://some.nice/video.mov",
				Username: "aws-access-key",
				Password: "aws-secret-key",
			},
		},
		Priority: 50,
		OutputGroup: elementalconductor.OutputGroup{
			Order: 1,
			AppleLiveGroupSettings: elementalconductor.AppleLiveGroupSettings{
				Destination: elementalconductor.Location{
					URI:      "s3://destination/video",
					Username: "aws-access-key",
					Password: "aws-secret-key",
				},
				SegmentDuration: 3,
			},
			Type: elementalconductor.AppleLiveOutputGroupType,
			Output: []elementalconductor.Output{
				{
					StreamAssemblyName: "stream_0",
					NameModifier:       "_hls_360p",
					Order:              0,
					Container:          elementalconductor.AppleHTTPLiveStreaming,
				},
				{
					StreamAssemblyName: "stream_1",
					NameModifier:       "_hls_480p",
					Order:              1,
					Container:          elementalconductor.AppleHTTPLiveStreaming,
				},
				{
					StreamAssemblyName: "stream_2",
					NameModifier:       "_hls_720p",
					Order:              2,
					Container:          elementalconductor.AppleHTTPLiveStreaming,
				},
				{
					StreamAssemblyName: "stream_3",
					NameModifier:       "_hls_1080p",
					Order:              3,
					Container:          elementalconductor.AppleHTTPLiveStreaming,
				},
			},
		},
		StreamAssembly: []elementalconductor.StreamAssembly{
			{
				Name:   "stream_0",
				Preset: "15",
			},
			{
				Name:   "stream_1",
				Preset: "20",
			},
			{
				Name:   "stream_2",
				Preset: "25",
			},
			{
				Name:   "stream_3",
				Preset: "30",
			},
		},
	}
	if !reflect.DeepEqual(&expectedJob, newJob) {
		t.Errorf("New adaptive bitrate job not according to spec.\nWanted %#v.\nGot    %#v.", &expectedJob, newJob)
	}
}

func TestElementalNewJobPresetNotFound(t *testing.T) {
	elementalConductorConfig := config.Config{
		ElementalConductor: &config.ElementalConductor{
			Host:            "https://mybucket.s3.amazonaws.com/destination-dir/",
			UserLogin:       "myuser",
			APIKey:          "elemental-api-key",
			AuthExpires:     30,
			AccessKeyID:     "aws-access-key",
			SecretAccessKey: "aws-secret-key",
			Destination:     "s3://destination",
		},
	}
	prov, err := elementalConductorFactory(&elementalConductorConfig)
	if err != nil {
		t.Fatal(err)
	}
	presetProvider, ok := prov.(*elementalConductorProvider)
	if !ok {
		t.Fatal("Could not type assert test provider to elementalConductorProvider")
	}
	source := "http://some.nice/video.mov"
	presets := []db.Preset{
		{
			Name:            "webm_720p",
			ProviderMapping: map[string]string{"other": "not relevant"},
			OutputOpts:      db.OutputOptions{Extension: "webm"},
		},
	}
	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     source,
		Presets:         presets,
		StreamingParams: provider.StreamingParams{},
	}
	newJob, err := presetProvider.newJob(transcodeProfile)
	if err != provider.ErrPresetNotFound {
		t.Errorf("Wrong error returned. Want %#v. Got %#v", provider.ErrPresetNotFound, err)
	}
	if newJob != nil {
		t.Errorf("Got unexpected non-nil job: %#v.", newJob)
	}
}

func TestJobStatusMap(t *testing.T) {
	var tests = []struct {
		elementalConductorStatus string
		expected                 provider.Status
	}{
		{"pending", provider.StatusQueued},
		{"preprocessing", provider.StatusStarted},
		{"running", provider.StatusStarted},
		{"postprocessing", provider.StatusStarted},
		{"complete", provider.StatusFinished},
		{"cancelled", provider.StatusCanceled},
		{"error", provider.StatusFailed},
		{"unknown", provider.StatusUnknown},
		{"someotherstatus", provider.StatusUnknown},
	}
	var p elementalConductorProvider
	for _, test := range tests {
		got := p.statusMap(test.elementalConductorStatus)
		if got != test.expected {
			t.Errorf("statusMap(%q): wrong value. Want %q. Got %q", test.elementalConductorStatus, test.expected, got)
		}
	}
}

func TestHealthcheck(t *testing.T) {
	server := NewElementalServer(nil, nil)
	defer server.Close()
	prov := elementalConductorProvider{
		client: elementalconductor.NewClient(server.URL, "", "", 0, "", "", ""),
	}
	var tests = []struct {
		minNodes    int
		nodes       []elementalconductor.Node
		expectedMsg string
	}{
		{
			2,
			[]elementalconductor.Node{
				{
					Product: elementalconductor.ProductConductorFile,
					Status:  "active",
				},
				{
					Product: elementalconductor.ProductServer,
					Status:  "starting",
				},
				{
					Product: elementalconductor.ProductServer,
					Status:  "active",
				},
				{
					Product: elementalconductor.ProductServer,
					Status:  "active",
				},
			},
			"",
		},
		{
			3,
			[]elementalconductor.Node{
				{
					Product: elementalconductor.ProductConductorFile,
					Status:  "active",
				},
				{
					Product: elementalconductor.ProductServer,
					Status:  "starting",
				},
				{
					Product: elementalconductor.ProductServer,
					Status:  "active",
				},
				{
					Product: elementalconductor.ProductServer,
					Status:  "error",
				},
			},
			"there are not enough active nodes. 3 nodes required to be active, but found only 1",
		},
		{
			2,
			[]elementalconductor.Node{
				{
					Product: elementalconductor.ProductConductorFile,
					Status:  "active",
				},
				{
					Product: elementalconductor.ProductConductorFile,
					Status:  "active",
				},
				{
					Product: elementalconductor.ProductServer,
					Status:  "active",
				},
			},
			"there are not enough active nodes. 2 nodes required to be active, but found only 1",
		},
	}
	for _, test := range tests {
		server.SetCloudConfig(&elementalconductor.CloudConfig{MinNodes: test.minNodes})
		server.SetNodes(test.nodes)
		err := prov.Healthcheck()
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
	var prov elementalConductorProvider
	expected := provider.Capabilities{
		InputFormats:  []string{"prores", "h264"},
		OutputFormats: []string{"mp4", "hls"},
		Destinations:  []string{"akamai", "s3"},
	}
	cap := prov.Capabilities()
	if !reflect.DeepEqual(cap, expected) {
		t.Errorf("Capabilities: want %#v. Got %#v", expected, cap)
	}
}
