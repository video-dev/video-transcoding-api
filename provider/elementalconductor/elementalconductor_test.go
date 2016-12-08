package elementalconductor

import (
	"encoding/xml"
	"reflect"
	"testing"
	"time"

	"github.com/NYTimes/encoding-wrapper/elementalconductor"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
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
	expected := &elementalconductor.Client{
		Host:        "elemental-server",
		UserLogin:   "myuser",
		APIKey:      "secret-key",
		AuthExpires: 30,
	}
	if !reflect.DeepEqual(econductorProvider.client, expected) {
		t.Errorf("Factory: wrong client returned. Want %#v. Got %#v.", expected, econductorProvider.client)
	}
	if !reflect.DeepEqual(*econductorProvider.config, *cfg.ElementalConductor) {
		t.Errorf("Factory: wrong config returned. Want %#v. Got %#v.", *cfg.ElementalConductor, *econductorProvider.config)
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
	prov, err := fakeElementalConductorFactory(&elementalConductorConfig)
	if err != nil {
		t.Fatal(err)
	}
	presetProvider, ok := prov.(*elementalConductorProvider)
	if !ok {
		t.Fatal("Could not type assert test provider to elementalConductorProvider")
	}
	source := "http://some.nice/video.mov"
	outputs := []db.TranscodeOutput{
		{
			FileName: "output_720p.webm",
			Preset: db.PresetMap{
				Name:            "webm_720p",
				ProviderMapping: map[string]string{Name: "webm_720p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "webm"},
			},
		},
		{
			FileName: "output_720p.mp4",
			Preset: db.PresetMap{
				Name:            "mp4_720p",
				ProviderMapping: map[string]string{Name: "mp4_720p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "mp4"},
			},
		},
		{
			FileName: "output_1080p.mp4",
			Preset: db.PresetMap{
				Name:            "mp4_1080p",
				ProviderMapping: map[string]string{Name: "mp4_1080p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "mp4"},
			},
		},
	}

	newJob, err := presetProvider.newJob(&db.Job{
		ID:              "job-1",
		SourceMedia:     source,
		Outputs:         outputs,
		StreamingParams: db.StreamingParams{},
	})
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
		OutputGroup: []elementalconductor.OutputGroup{
			{
				Order: 1,
				FileGroupSettings: &elementalconductor.FileGroupSettings{
					Destination: &elementalconductor.Location{
						URI:      "s3://destination/job-1/output_720p",
						Username: "aws-access-key",
						Password: "aws-secret-key",
					},
				},
				Type: elementalconductor.FileOutputGroupType,
				Output: []elementalconductor.Output{
					{
						StreamAssemblyName: "stream_0",
						Order:              1,
						Container:          elementalconductor.Container("webm"),
					},
				},
			},
			{
				Order: 2,
				FileGroupSettings: &elementalconductor.FileGroupSettings{
					Destination: &elementalconductor.Location{
						URI:      "s3://destination/job-1/output_720p",
						Username: "aws-access-key",
						Password: "aws-secret-key",
					},
				},
				Type: elementalconductor.FileOutputGroupType,
				Output: []elementalconductor.Output{
					{
						StreamAssemblyName: "stream_1",
						Order:              1,
						Container:          elementalconductor.MPEG4,
					},
				},
			},
			{
				Order: 3,
				FileGroupSettings: &elementalconductor.FileGroupSettings{
					Destination: &elementalconductor.Location{
						URI:      "s3://destination/job-1/output_1080p",
						Username: "aws-access-key",
						Password: "aws-secret-key",
					},
				},
				Type: elementalconductor.FileOutputGroupType,
				Output: []elementalconductor.Output{
					{
						StreamAssemblyName: "stream_2",
						Order:              1,
						Container:          elementalconductor.MPEG4,
					},
				},
			},
		},
		StreamAssembly: []elementalconductor.StreamAssembly{
			{
				Name:   "stream_0",
				Preset: "webm_720p",
			},
			{
				Name:   "stream_1",
				Preset: "mp4_720p",
			},
			{
				Name:   "stream_2",
				Preset: "mp4_1080p",
			},
		},
	}
	if !reflect.DeepEqual(expectedJob, *newJob) {
		t.Errorf("New job not according to spec.\nWanted %#v.\nGot    %#v.", expectedJob, *newJob)
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
	prov, err := fakeElementalConductorFactory(&elementalConductorConfig)
	if err != nil {
		t.Fatal(err)
	}
	presetProvider, ok := prov.(*elementalConductorProvider)
	if !ok {
		t.Fatal("Could not type assert test provider to elementalConductorProvider")
	}
	source := "http://some.nice/video.mov"
	outputs := []db.TranscodeOutput{
		{
			FileName: "output_hls_360p/video.m3u8",
			Preset: db.PresetMap{
				Name:            "hls_360p",
				ProviderMapping: map[string]string{Name: "hls_360p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "hls"},
			},
		},
		{
			FileName: "output_hls_480p/video.m3u8",
			Preset: db.PresetMap{
				Name:            "hls_480p",
				ProviderMapping: map[string]string{Name: "hls_480p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "ts"},
			},
		},
		{
			FileName: "output_hls_720p/video.m3u8",
			Preset: db.PresetMap{
				Name:            "hls_720p",
				ProviderMapping: map[string]string{Name: "hls_720p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "m3u8"},
			},
		},
		{
			FileName: "output_hls_1080p/video.m3u8",
			Preset: db.PresetMap{
				Name:            "hls_1080p",
				ProviderMapping: map[string]string{Name: "hls_1080p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: ".ts"},
			},
		},
	}
	newJob, err := presetProvider.newJob(&db.Job{
		ID:          "job-2",
		SourceMedia: source,
		Outputs:     outputs,
		StreamingParams: db.StreamingParams{
			Protocol:         "hls",
			SegmentDuration:  3,
			PlaylistFileName: "hls/master.m3u8",
		},
	})
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
		OutputGroup: []elementalconductor.OutputGroup{
			{
				Order: 1,
				AppleLiveGroupSettings: &elementalconductor.AppleLiveGroupSettings{
					Destination: &elementalconductor.Location{
						URI:      "s3://destination/job-2/hls/master",
						Username: "aws-access-key",
						Password: "aws-secret-key",
					},
					SegmentDuration: 3,
					EmitSingleFile:  true,
				},
				Type: elementalconductor.AppleLiveOutputGroupType,
				Output: []elementalconductor.Output{
					{
						StreamAssemblyName: "stream_0",
						NameModifier:       "_0000000001",
						Order:              1,
						Container:          elementalconductor.AppleHTTPLiveStreaming,
					},
					{
						StreamAssemblyName: "stream_1",
						NameModifier:       "_0000000002",
						Order:              2,
						Container:          elementalconductor.AppleHTTPLiveStreaming,
					},
					{
						StreamAssemblyName: "stream_2",
						NameModifier:       "_0000000003",
						Order:              3,
						Container:          elementalconductor.AppleHTTPLiveStreaming,
					},
					{
						StreamAssemblyName: "stream_3",
						NameModifier:       "_0000000004",
						Order:              4,
						Container:          elementalconductor.AppleHTTPLiveStreaming,
					},
				},
			},
		},
		StreamAssembly: []elementalconductor.StreamAssembly{
			{
				Name:   "stream_0",
				Preset: "hls_360p",
			},
			{
				Name:   "stream_1",
				Preset: "hls_480p",
			},
			{
				Name:   "stream_2",
				Preset: "hls_720p",
			},
			{
				Name:   "stream_3",
				Preset: "hls_1080p",
			},
		},
	}
	if !reflect.DeepEqual(&expectedJob, newJob) {
		t.Errorf("New adaptive bitrate job not according to spec.\nWanted %#v.\nGot    %#v.", &expectedJob, newJob)
	}
}

func TestElementalNewJobAdaptiveAndNonAdaptiveStreaming(t *testing.T) {
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
	prov, err := fakeElementalConductorFactory(&elementalConductorConfig)
	if err != nil {
		t.Fatal(err)
	}
	presetProvider, ok := prov.(*elementalConductorProvider)
	if !ok {
		t.Fatal("Could not type assert test provider to elementalConductorProvider")
	}
	source := "http://some.nice/video.mov"
	outputs := []db.TranscodeOutput{
		{
			FileName: "output_720p.webm",
			Preset: db.PresetMap{
				Name:            "webm_720p",
				ProviderMapping: map[string]string{Name: "webm_720p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "webm"},
			},
		},
		{
			FileName: "output_720p.mp4",
			Preset: db.PresetMap{
				Name:            "mp4_720p",
				ProviderMapping: map[string]string{Name: "mp4_720p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "mp4"},
			},
		},
		{
			FileName: "output_1080p.mp4",
			Preset: db.PresetMap{
				Name:            "mp4_1080p",
				ProviderMapping: map[string]string{Name: "mp4_1080p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "mp4"},
			},
		},
		{
			FileName: "output_hls_360p/video.m3u8",
			Preset: db.PresetMap{
				Name:            "hls_360p",
				ProviderMapping: map[string]string{Name: "hls_360p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "m3u8"},
			},
		},
		{
			FileName: "output_hls_480p/video.m3u8",
			Preset: db.PresetMap{
				Name:            "hls_480p",
				ProviderMapping: map[string]string{Name: "hls_480p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "m3u8"},
			},
		},
		{
			FileName: "output_hls_720p/video.m3u8",
			Preset: db.PresetMap{
				Name:            "hls_720p",
				ProviderMapping: map[string]string{Name: "hls_720p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "m3u8"},
			},
		},
		{
			FileName: "output_hls_1080p/video.m3u8",
			Preset: db.PresetMap{
				Name:            "hls_1080p",
				ProviderMapping: map[string]string{Name: "hls_1080p", "other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "m3u8"},
			},
		},
	}
	newJob, err := presetProvider.newJob(&db.Job{
		ID:          "job-3",
		SourceMedia: source,
		Outputs:     outputs,
		StreamingParams: db.StreamingParams{
			Protocol:         "hls",
			SegmentDuration:  3,
			PlaylistFileName: "output_hls/index.m3u8",
		},
	})
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
		OutputGroup: []elementalconductor.OutputGroup{
			{
				Order: 1,
				FileGroupSettings: &elementalconductor.FileGroupSettings{
					Destination: &elementalconductor.Location{
						URI:      "s3://destination/job-3/output_720p",
						Username: "aws-access-key",
						Password: "aws-secret-key",
					},
				},
				Type: elementalconductor.FileOutputGroupType,
				Output: []elementalconductor.Output{
					{
						StreamAssemblyName: "stream_0",
						Order:              1,
						Container:          elementalconductor.Container("webm"),
					},
				},
			},
			{
				Order: 2,
				FileGroupSettings: &elementalconductor.FileGroupSettings{
					Destination: &elementalconductor.Location{
						URI:      "s3://destination/job-3/output_720p",
						Username: "aws-access-key",
						Password: "aws-secret-key",
					},
				},
				Type: elementalconductor.FileOutputGroupType,
				Output: []elementalconductor.Output{
					{
						StreamAssemblyName: "stream_1",
						Order:              1,
						Container:          elementalconductor.MPEG4,
					},
				},
			},
			{
				Order: 3,
				FileGroupSettings: &elementalconductor.FileGroupSettings{
					Destination: &elementalconductor.Location{
						URI:      "s3://destination/job-3/output_1080p",
						Username: "aws-access-key",
						Password: "aws-secret-key",
					},
				},
				Type: elementalconductor.FileOutputGroupType,
				Output: []elementalconductor.Output{
					{
						StreamAssemblyName: "stream_2",
						Order:              1,
						Container:          elementalconductor.MPEG4,
					},
				},
			},
			{
				Order: 4,
				AppleLiveGroupSettings: &elementalconductor.AppleLiveGroupSettings{
					Destination: &elementalconductor.Location{
						URI:      "s3://destination/job-3/output_hls/index",
						Username: "aws-access-key",
						Password: "aws-secret-key",
					},
					SegmentDuration: 3,
					EmitSingleFile:  true,
				},
				Type: elementalconductor.AppleLiveOutputGroupType,
				Output: []elementalconductor.Output{
					{
						StreamAssemblyName: "stream_3",
						Order:              1,
						NameModifier:       "_0000000001",
						Container:          elementalconductor.AppleHTTPLiveStreaming,
					},
					{
						StreamAssemblyName: "stream_4",
						Order:              2,
						NameModifier:       "_0000000002",
						Container:          elementalconductor.AppleHTTPLiveStreaming,
					},
					{
						StreamAssemblyName: "stream_5",
						Order:              3,
						NameModifier:       "_0000000003",
						Container:          elementalconductor.AppleHTTPLiveStreaming,
					},
					{
						StreamAssemblyName: "stream_6",
						Order:              4,
						NameModifier:       "_0000000004",
						Container:          elementalconductor.AppleHTTPLiveStreaming,
					},
				},
			},
		},
		StreamAssembly: []elementalconductor.StreamAssembly{
			{
				Name:   "stream_0",
				Preset: "webm_720p",
			},
			{
				Name:   "stream_1",
				Preset: "mp4_720p",
			},
			{
				Name:   "stream_2",
				Preset: "mp4_1080p",
			},
			{
				Name:   "stream_3",
				Preset: "hls_360p",
			},
			{
				Name:   "stream_4",
				Preset: "hls_480p",
			},
			{
				Name:   "stream_5",
				Preset: "hls_720p",
			},
			{
				Name:   "stream_6",
				Preset: "hls_1080p",
			},
		},
	}
	if !reflect.DeepEqual(&expectedJob, newJob) {
		t.Errorf("New adaptive and non-adaptive bitrate job not according to spec.\nWanted %#v.\nGot    %#v.", &expectedJob, newJob)
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
	outputs := []db.TranscodeOutput{
		{
			Preset: db.PresetMap{
				Name:            "webm_720p",
				ProviderMapping: map[string]string{"other": "not relevant"},
				OutputOpts:      db.OutputOptions{Extension: "webm"},
			},
		},
	}
	newJob, err := presetProvider.newJob(&db.Job{ID: "job-2", SourceMedia: source, Outputs: outputs})
	if err != provider.ErrPresetMapNotFound {
		t.Errorf("Wrong error returned. Want %#v. Got %#v", provider.ErrPresetMapNotFound, err)
	}
	if newJob != nil {
		t.Errorf("Got unexpected non-nil job: %#v.", newJob)
	}
}

func TestJobStatusOutputDestination(t *testing.T) {
	var tests = []struct {
		job            db.Job
		destinationCfg string
		expected       string
	}{
		{
			db.Job{
				ID: "job-123",
			},
			"s3://destination///",
			"s3://destination/job-123",
		},
		{
			db.Job{
				ID: "job-123",
			},
			"s3://destination",
			"s3://destination/job-123",
		},
	}
	elementalConductorConfig := config.Config{
		ElementalConductor: &config.ElementalConductor{
			Host:            "https://mybucket.s3.amazonaws.com/destination-dir/",
			UserLogin:       "myuser",
			APIKey:          "elemental-api-key",
			AuthExpires:     30,
			AccessKeyID:     "aws-access-key",
			SecretAccessKey: "aws-secret-key",
		},
	}
	prov, err := fakeElementalConductorFactory(&elementalConductorConfig)
	if err != nil {
		t.Fatal(err)
	}
	provider, ok := prov.(*elementalConductorProvider)
	if !ok {
		t.Fatal("Could not type assert test provider to elementalConductorProvider")
	}
	for _, test := range tests {
		provider.config.Destination = test.destinationCfg
		got := provider.getOutputDestination(&test.job)
		if got != test.expected {
			t.Errorf("Wrong output destination. Want %q. Got %q", test.expected, got)
		}
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

func TestJobStatus(t *testing.T) {
	elementalConductorConfig := config.ElementalConductor{
		Host:            "https://mybucket.s3.amazonaws.com/destination-dir/",
		UserLogin:       "myuser",
		APIKey:          "elemental-api-key",
		AuthExpires:     30,
		AccessKeyID:     "aws-access-key",
		SecretAccessKey: "aws-secret-key",
		Destination:     "s3://destination",
	}
	submitted := elementalconductor.DateTime{Time: time.Now().UTC()}
	client := newFakeElementalConductorClient(&elementalConductorConfig)
	client.jobs["job-1"] = elementalconductor.Job{
		Href: "whatever",
		Input: elementalconductor.Input{
			InputInfo: &elementalconductor.InputInfo{
				Video: elementalconductor.VideoInputInfo{
					Format: "AVC",
					Width:  "1 920 pixels",
					Height: "1 080 pixels",
				},
			},
		},
		ContentDuration: &elementalconductor.ContentDuration{
			InputDuration: 123,
		},
		OutputGroup: []elementalconductor.OutputGroup{
			{
				Output: []elementalconductor.Output{
					{
						FullURI:            "s3://somebucket/dir/video1.mp4",
						Extension:          "mp4",
						StreamAssemblyName: "stream_0",
						Container:          "mp4",
					},
				},
			},
			{
				Output: []elementalconductor.Output{
					{
						FullURI:            "s3://somebucket/dir/video1.webm",
						Extension:          "webm",
						StreamAssemblyName: "stream_1",
						Container:          "webm",
					},
				},
			},
			{
				Output: []elementalconductor.Output{
					{
						FullURI:            "s3://somebucket/dir/video1.mov",
						Extension:          "mov",
						StreamAssemblyName: "unknown",
						Container:          "mov",
					},
				},
			},
			{
				Type: elementalconductor.AppleLiveOutputGroupType,
				AppleLiveGroupSettings: &elementalconductor.AppleLiveGroupSettings{
					Destination: &elementalconductor.Location{
						URI: "s3://somebucket/dir/video1",
					},
					SegmentDuration: 3,
					EmitSingleFile:  true,
				},
				Output: []elementalconductor.Output{
					{
						FullURI:            "s3://somebucket/dir/video1_720p.m3u8",
						Extension:          "m3u8",
						StreamAssemblyName: "stream_2",
						Container:          "m3u8",
					},
					{
						FullURI:            "s3://somebucket/dir/video1_1080p.m3u8",
						Extension:          "m3u8",
						StreamAssemblyName: "stream_3",
						Container:          "m3u8",
					},
				},
			},
		},
		StreamAssembly: []elementalconductor.StreamAssembly{
			{
				ID:   "2323",
				Name: "stream_0",
				VideoDescription: &elementalconductor.StreamVideoDescription{
					Codec:       "h.264",
					EncoderType: "gpu",
					Height:      "1080",
					Width:       "1920",
				},
			},
			{
				ID:   "2324",
				Name: "stream_1",
				VideoDescription: &elementalconductor.StreamVideoDescription{
					Codec:       "vp8",
					EncoderType: "gpu",
					Height:      "1080",
					Width:       "1920",
				},
			},
			{
				ID:   "2325",
				Name: "stream_2",
				VideoDescription: &elementalconductor.StreamVideoDescription{
					Codec:       "h.264",
					EncoderType: "gpu",
					Height:      "720",
					Width:       "1280",
				},
			},
			{
				ID:   "2326",
				Name: "stream_3",
				VideoDescription: &elementalconductor.StreamVideoDescription{
					Codec:       "h.264",
					EncoderType: "gpu",
					Height:      "1080",
					Width:       "1920",
				},
			},
		},
		PercentComplete: 89,
		Status:          "running",
		Submitted:       submitted,
	}
	prov := elementalConductorProvider{client: client, config: &elementalConductorConfig}
	jobStatus, err := prov.JobStatus(&db.Job{ID: "super-job-1", ProviderJobID: "job-1"})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: "job-1",
		Progress:      89.,
		Status:        provider.StatusStarted,
		Output: provider.JobOutput{
			Destination: "s3://destination/super-job-1",
			Files: []provider.OutputFile{
				{
					Path:      "s3://somebucket/dir/video1.m3u8",
					Container: "m3u8",
				},
				{
					Path:       "s3://somebucket/dir/video1.mp4",
					Container:  "mp4",
					VideoCodec: "h.264",
					Width:      1920,
					Height:     1080,
				},
				{
					Path:       "s3://somebucket/dir/video1.webm",
					Container:  "webm",
					VideoCodec: "vp8",
					Width:      1920,
					Height:     1080,
				},
			},
		},
		SourceInfo: provider.SourceInfo{
			Duration:   123 * time.Second,
			Width:      1920,
			Height:     1080,
			VideoCodec: "AVC",
		},
		ProviderStatus: map[string]interface{}{
			"status":    "running",
			"submitted": submitted,
		},
	}
	if !reflect.DeepEqual(*jobStatus, expectedJobStatus) {
		t.Errorf("wrong job stats\nwant %#v\ngot  %#v", expectedJobStatus, *jobStatus)
	}
}

func TestJobStatusNoDuration(t *testing.T) {
	elementalConductorConfig := config.ElementalConductor{
		Host:            "https://mybucket.s3.amazonaws.com/destination-dir/",
		UserLogin:       "myuser",
		APIKey:          "elemental-api-key",
		AuthExpires:     30,
		AccessKeyID:     "aws-access-key",
		SecretAccessKey: "aws-secret-key",
		Destination:     "s3://destination",
	}
	submitted := elementalconductor.DateTime{Time: time.Now().UTC()}
	client := newFakeElementalConductorClient(&elementalConductorConfig)
	client.jobs["job-1"] = elementalconductor.Job{
		Href: "whatever",
		Input: elementalconductor.Input{
			InputInfo: &elementalconductor.InputInfo{
				Video: elementalconductor.VideoInputInfo{
					Format: "AVC",
					Width:  "1 920 pixels",
					Height: "1 080 pixels",
				},
			},
		},
		OutputGroup: []elementalconductor.OutputGroup{
			{
				Output: []elementalconductor.Output{
					{
						FullURI:            "s3://somebucket/dir/video1.mp4",
						Extension:          "mp4",
						StreamAssemblyName: "stream_0",
						Container:          "mp4",
					},
				},
			},
			{
				Output: []elementalconductor.Output{
					{
						FullURI:            "s3://somebucket/dir/video1.webm",
						Extension:          "webm",
						StreamAssemblyName: "stream_1",
						Container:          "webm",
					},
				},
			},
			{
				Output: []elementalconductor.Output{
					{
						FullURI:            "s3://somebucket/dir/video1.mov",
						Extension:          "mov",
						StreamAssemblyName: "unknown",
						Container:          "mov",
					},
				},
			},
			{
				Output: []elementalconductor.Output{
					{
						FullURI:            "s3://somebucket/dir/video1_720p.m3u8",
						Extension:          "m3u8",
						StreamAssemblyName: "stream_2",
						Container:          "m3u8",
					},
					{
						FullURI:            "s3://somebucket/dir/video1_1080p.m3u8",
						Extension:          "m3u8",
						StreamAssemblyName: "stream_3",
						Container:          "m3u8",
					},
				},
			},
		},
		StreamAssembly: []elementalconductor.StreamAssembly{
			{
				ID:   "2323",
				Name: "stream_0",
				VideoDescription: &elementalconductor.StreamVideoDescription{
					Codec:       "h.264",
					EncoderType: "gpu",
					Height:      "1080",
					Width:       "1920",
				},
			},
			{
				ID:   "2324",
				Name: "stream_1",
				VideoDescription: &elementalconductor.StreamVideoDescription{
					Codec:       "vp8",
					EncoderType: "gpu",
					Height:      "1080",
					Width:       "1920",
				},
			},
			{
				ID:   "2325",
				Name: "stream_2",
				VideoDescription: &elementalconductor.StreamVideoDescription{
					Codec:       "h.264",
					EncoderType: "gpu",
					Height:      "720",
					Width:       "1280",
				},
			},
			{
				ID:   "2326",
				Name: "stream_3",
				VideoDescription: &elementalconductor.StreamVideoDescription{
					Codec:       "h.264",
					EncoderType: "gpu",
					Height:      "1080",
					Width:       "1920",
				},
			},
		},
		PercentComplete: 89,
		Status:          "running",
		Submitted:       submitted,
	}
	prov := elementalConductorProvider{client: client, config: &elementalConductorConfig}
	jobStatus, err := prov.JobStatus(&db.Job{ID: "super-job-1", ProviderJobID: "job-1"})
	if err != nil {
		t.Fatal(err)
	}
	expectedJobStatus := provider.JobStatus{
		ProviderName:  Name,
		ProviderJobID: "job-1",
		Progress:      89.,
		Status:        provider.StatusStarted,
		Output: provider.JobOutput{
			Destination: "s3://destination/super-job-1",
			Files: []provider.OutputFile{
				{
					Path:       "s3://somebucket/dir/video1.mp4",
					Container:  "mp4",
					VideoCodec: "h.264",
					Width:      1920,
					Height:     1080,
				},
				{
					Path:       "s3://somebucket/dir/video1.webm",
					Container:  "webm",
					VideoCodec: "vp8",
					Width:      1920,
					Height:     1080,
				},
				{
					Path:       "s3://somebucket/dir/video1_720p.m3u8",
					Container:  "m3u8",
					VideoCodec: "h.264",
					Width:      1280,
					Height:     720,
				},
				{
					Path:       "s3://somebucket/dir/video1_1080p.m3u8",
					Container:  "m3u8",
					VideoCodec: "h.264",
					Width:      1920,
					Height:     1080,
				},
			},
		},
		SourceInfo: provider.SourceInfo{
			Width:      1920,
			Height:     1080,
			VideoCodec: "AVC",
		},
		ProviderStatus: map[string]interface{}{
			"status":    "running",
			"submitted": submitted,
		},
	}
	if !reflect.DeepEqual(*jobStatus, expectedJobStatus) {
		t.Errorf("wrong job stats\nwant %#v\ngot  %#v", expectedJobStatus, *jobStatus)
	}
}

func TestCancelJob(t *testing.T) {
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
	prov, err := fakeElementalConductorFactory(&elementalConductorConfig)
	if err != nil {
		t.Fatal(err)
	}
	err = prov.CancelJob("idk")
	if err != nil {
		t.Fatal(err)
	}
	client := prov.(*elementalConductorProvider).client.(*fakeElementalConductorClient)
	if client.canceledJobs[0] != "idk" {
		t.Errorf("did not cancel the correct job. Want %q. Got %q", "idk", client.canceledJobs[0])
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
