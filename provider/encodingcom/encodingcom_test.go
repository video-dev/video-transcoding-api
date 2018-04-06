package encodingcom

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/NYTimes/encoding-wrapper/encodingcom"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/kr/pretty"
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
	if !reflect.DeepEqual(ecomProvider.client, expected) {
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
			Name: "hls_360p",
			ProviderMapping: map[string]string{
				Name:           "321322",
				"not-relevant": "allthings",
			},
			OutputOpts: db.OutputOptions{Extension: "m3u8"},
		},
		{
			Name: "hls_480p",
			ProviderMapping: map[string]string{
				Name:           "321322",
				"not-relevant": "allthings",
			},
			OutputOpts: db.OutputOptions{Extension: "m3u8"},
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
	outputs := make([]db.TranscodeOutput, len(presets))
	for i, preset := range presets {
		_, err := prov.CreatePreset(db.Preset{
			Name:      preset.ProviderMapping[Name],
			Container: preset.OutputOpts.Extension,
		})
		if err != nil {
			t.Fatal(err)
		}
		fileName := "output-" + preset.Name + "." + preset.OutputOpts.Extension
		if preset.OutputOpts.Extension == "m3u8" {
			fileName = "output-" + preset.Name + "/video.m3u8"
		}
		outputs[i] = db.TranscodeOutput{
			Preset:   preset,
			FileName: fileName,
		}
	}

	jobStatus, err := prov.Transcode(&db.Job{
		ID:          "job-123",
		SourceMedia: source,
		Outputs:     outputs,
		StreamingParams: db.StreamingParams{
			PlaylistFileName: "output_hls/video.m3u8",
			Protocol:         "hls",
			SegmentDuration:  3,
		},
	})
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
	falseYesNoBoolean := encodingcom.YesNoBoolean(false)
	expectedFormats := []encodingcom.Format{
		{
			OutputPreset: "123455",
			Destination:  []string{dest + "job-123/output-webm_720p.webm"},
		},
		{
			OutputPreset: "123456",
			Destination:  []string{dest + "job-123/output-webm_480p.webm"},
		},
		{
			OutputPreset: "321321",
			Destination:  []string{dest + "job-123/output-mp4_1080p.mp4"},
		},
		{
			Output:          []string{hlsOutput},
			Destination:     []string{dest + "job-123/output_hls/video.m3u8"},
			SegmentDuration: 3,
			PackFiles:       &falseYesNoBoolean,
			Stream: []encodingcom.Stream{
				{
					SubPath: "hls_360p",
				},
				{
					SubPath: "hls_480p",
				},
				{
					SubPath: "hls_1080p",
				},
			},
		},
	}
	if !reflect.DeepEqual(media.Request.Format, expectedFormats) {
		t.Errorf("Wrong format.\nWant %#v\nGot  %#v", expectedFormats, media.Request.Format)
	}
	if !reflect.DeepEqual([]string{source}, media.Request.Source) {
		t.Errorf("Wrong source. Want %v. Got %v.", []string{source}, media.Request.Source)
	}

	jobStatus, err = prov.Transcode(&db.Job{
		ID:          "job-123",
		SourceMedia: source,
		Outputs:     outputs,
		StreamingParams: db.StreamingParams{
			PlaylistFileName: "output_hls/video.m3u8",
			Protocol:         "hls",
			SegmentDuration:  3,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if expected := "it worked"; jobStatus.StatusMessage != expected {
		t.Errorf("wrong StatusMessage. Want %q. Got %q", expected, jobStatus.StatusMessage)
	}
	if jobStatus.ProviderName != Name {
		t.Errorf("wrong ProviderName. Want %q. Got %q", Name, jobStatus.ProviderName)
	}
	media, err = server.getMedia(jobStatus.ProviderJobID)
	if err != nil {
		t.Fatal(err)
	}
	expectedFormats = []encodingcom.Format{
		{
			OutputPreset: "123455",
			Destination:  []string{dest + "job-123/output-webm_720p.webm"},
		},
		{
			OutputPreset: "123456",
			Destination:  []string{dest + "job-123/output-webm_480p.webm"},
		},
		{
			OutputPreset: "321321",
			Destination:  []string{dest + "job-123/output-mp4_1080p.mp4"},
		},
		{
			Output:          []string{hlsOutput},
			Destination:     []string{dest + "job-123/output_hls/video.m3u8"},
			SegmentDuration: 3,
			PackFiles:       &falseYesNoBoolean,
			Stream: []encodingcom.Stream{
				{
					SubPath: "hls_360p",
				},
				{
					SubPath: "hls_480p",
				},
				{
					SubPath: "hls_1080p",
				},
			},
		},
	}
	if !reflect.DeepEqual(media.Request.Format, expectedFormats) {
		t.Errorf("Wrong format.\nWant %#v\nGot  %#v.", expectedFormats, media.Request.Format)
	}
	if !reflect.DeepEqual([]string{source}, media.Request.Source) {
		t.Errorf("Wrong source. Want %v. Got %v.", []string{source}, media.Request.Source)
	}

}

func TestEncodingComS3Input(t *testing.T) {
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
	source := "s3://mybucket/directory/video.mp4"
	expectedSource := "https://mybucket.s3.amazonaws.com/directory/video.mp4?nocopy"
	presets := []db.PresetMap{
		{
			Name: "webm_720p",
			ProviderMapping: map[string]string{
				Name:           "123455",
				"not-relevant": "something",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
	}
	outputs := make([]db.TranscodeOutput, len(presets))
	for i, preset := range presets {
		_, err := prov.CreatePreset(db.Preset{
			Name:      preset.ProviderMapping[Name],
			Container: preset.OutputOpts.Extension,
		})
		if err != nil {
			t.Fatal(err)
		}
		outputs[i] = db.TranscodeOutput{
			Preset:   preset,
			FileName: "best-video-ever." + preset.OutputOpts.Extension,
		}
	}

	jobStatus, err := prov.Transcode(&db.Job{
		ID:          "job-123",
		SourceMedia: source,
		Outputs:     outputs,
	})
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
			OutputPreset: "123455",
			Destination:  []string{dest + "job-123/best-video-ever.webm"},
		},
	}
	if !reflect.DeepEqual(media.Request.Format, expectedFormats) {
		t.Errorf("Wrong format.\nWant %#v\nGot  %#v.", expectedFormats, media.Request.Format)
	}
	if !reflect.DeepEqual([]string{expectedSource}, media.Request.Source) {
		t.Errorf("Wrong source. Want %v. Got %v.", []string{expectedSource}, media.Request.Source)
	}
}

func TestEncodingComS3InputWithNoCopy(t *testing.T) {
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
	source := "s3://mybucket/directory/video.mp4?nocopy"
	expectedSource := "https://mybucket.s3.amazonaws.com/directory/video.mp4?nocopy"
	presets := []db.PresetMap{
		{
			Name: "webm_720p",
			ProviderMapping: map[string]string{
				Name:           "123455",
				"not-relevant": "something",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
	}
	outputs := make([]db.TranscodeOutput, len(presets))
	for i, preset := range presets {
		_, err := prov.CreatePreset(db.Preset{
			Name:      preset.ProviderMapping[Name],
			Container: preset.OutputOpts.Extension,
		})
		if err != nil {
			t.Fatal(err)
		}
		outputs[i] = db.TranscodeOutput{
			Preset:   preset,
			FileName: "best-video-ever." + preset.OutputOpts.Extension,
		}
	}

	jobStatus, err := prov.Transcode(&db.Job{
		ID:          "job-123",
		SourceMedia: source,
		Outputs:     outputs,
	})
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
			OutputPreset: "123455",
			Destination:  []string{dest + "job-123/best-video-ever.webm"},
		},
	}
	if !reflect.DeepEqual(media.Request.Format, expectedFormats) {
		t.Errorf("Wrong format.\nWant %#v\nGot  %#v.", expectedFormats, media.Request.Format)
	}
	if !reflect.DeepEqual([]string{expectedSource}, media.Request.Source) {
		t.Errorf("Wrong source. Want %v. Got %v.", []string{expectedSource}, media.Request.Source)
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
	outputs := []db.TranscodeOutput{
		{
			Preset: db.PresetMap{
				Name: "webm_720p",
				ProviderMapping: map[string]string{
					Name:           "123455",
					"not-relevant": "something",
				},
				OutputOpts: db.OutputOptions{Extension: "webm"},
			},
		},
		{
			Preset: db.PresetMap{
				Name: "webm_480p",
				ProviderMapping: map[string]string{
					"not-relevant": "otherthing",
				},
				OutputOpts: db.OutputOptions{Extension: "webm"},
			},
		},
	}
	jobStatus, err := prov.Transcode(&db.Job{
		ID:              "job-2",
		SourceMedia:     source,
		Outputs:         outputs,
		StreamingParams: db.StreamingParams{SegmentDuration: 3},
	})
	expectedErrorString := "Error converting presets to formats on Transcode operation: Error getting preset info: Error returned by the Encoding.com API: {\"Errors\":[\"123455 preset not found\"]}"
	if err.Error() != expectedErrorString {
		t.Errorf("Wrong error\nWant %#v\nGot  %#v", expectedErrorString, err.Error())
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
		ID:   "mymedia",
		Size: "1920x1080",
		Request: request{
			Format: []encodingcom.Format{
				{
					Bitrate:    "2500k",
					Size:       "1920x1080",
					VideoCodec: "VP9",
					Output:     []string{hlsOutput},
				},
			},
		},
		Status:   "Saving",
		Created:  now.Add(-time.Hour),
		Started:  now.Add(-50 * time.Minute),
		Finished: now.Add(-10 * time.Minute),
	}
	server.medias["mymedia"] = &media
	client, err := encodingcom.NewClient(server.URL, "myuser", "secret")
	if err != nil {
		t.Fatal(err)
	}
	prov := encodingComProvider{client: client}
	prov.config = &config.Config{
		EncodingCom: &config.EncodingCom{
			Destination: "https://mybucket.s3.amazonaws.com/dir/",
		},
	}
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: "mymedia"})
	if err != nil {
		t.Fatal(err)
	}
	expected := provider.JobStatus{
		ProviderJobID: "mymedia",
		ProviderName:  "encoding.com",
		Status:        provider.StatusFinished,
		StatusMessage: "",
		Progress:      100,
		ProviderStatus: map[string]interface{}{
			"sourcefile":   "http://some.source.file",
			"timeleft":     "1",
			"created":      media.Created,
			"started":      media.Started,
			"finished":     media.Finished,
			"formatStatus": []string{""},
		},
		SourceInfo: provider.SourceInfo{
			Duration:   183e9,
			Width:      1920,
			Height:     1080,
			VideoCodec: "VP9",
		},
		Output: provider.JobOutput{
			Destination: "s3://mybucket/dir/job-123/",
			Files: []provider.OutputFile{
				{
					Path:       "s3://mybucket/dir/job-123/some_hls_preset/video-0.m3u8",
					VideoCodec: "VP9",
					Width:      1920,
					Height:     1080,
					Container:  "m3u8",
					FileSize:   45674,
				},
				{
					Path:       "s3://mybucket/dir/job-123/video.m3u8",
					VideoCodec: "VP9",
					Width:      1920,
					Height:     1080,
					Container:  "m3u8",
					FileSize:   45674,
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobStatus, expected) {
		t.Errorf("JobStatus: wrong job returned.\nWant %#v\nGot  %#v", expected, *jobStatus)
	}
}

func TestJobStatusMissingDimension(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	now := time.Now().UTC().Truncate(time.Second)
	media := fakeMedia{
		ID:   "mymedia",
		Size: "1920x1080",
		Request: request{
			Format: []encodingcom.Format{
				{
					Bitrate:    "2500k",
					Size:       "0x1080",
					VideoCodec: "VP9",
					Output:     []string{hlsOutput},
				},
			},
		},
		Status:   "Saving",
		Created:  now.Add(-time.Hour),
		Started:  now.Add(-50 * time.Minute),
		Finished: now.Add(-10 * time.Minute),
	}
	server.medias["mymedia"] = &media
	client, err := encodingcom.NewClient(server.URL, "myuser", "secret")
	if err != nil {
		t.Fatal(err)
	}
	prov := encodingComProvider{client: client}
	prov.config = &config.Config{
		EncodingCom: &config.EncodingCom{
			Destination: "https://mybucket.s3.amazonaws.com/dir/",
		},
	}
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: "mymedia"})
	if err != nil {
		t.Fatal(err)
	}
	expected := provider.JobStatus{
		ProviderJobID: "mymedia",
		ProviderName:  "encoding.com",
		Status:        provider.StatusFinished,
		StatusMessage: "",
		Progress:      100,
		ProviderStatus: map[string]interface{}{
			"sourcefile":   "http://some.source.file",
			"timeleft":     "1",
			"created":      media.Created,
			"started":      media.Started,
			"finished":     media.Finished,
			"formatStatus": []string{""},
		},
		SourceInfo: provider.SourceInfo{
			Duration:   183e9,
			Width:      1920,
			Height:     1080,
			VideoCodec: "VP9",
		},
		Output: provider.JobOutput{
			Destination: "s3://mybucket/dir/job-123/",
			Files: []provider.OutputFile{
				{
					Path:       "s3://mybucket/dir/job-123/some_hls_preset/video-0.m3u8",
					VideoCodec: "VP9",
					Width:      1920,
					Height:     1080,
					Container:  "m3u8",
					FileSize:   45674,
				},
				{
					Path:       "s3://mybucket/dir/job-123/video.m3u8",
					VideoCodec: "VP9",
					Width:      1920,
					Height:     1080,
					Container:  "m3u8",
					FileSize:   45674,
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobStatus, expected) {
		t.Errorf("JobStatus: wrong job returned.\nWant %#v\nGot  %#v", expected, *jobStatus)
	}
}

func TestJobStatusRotatedVideo(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	now := time.Now().UTC().Truncate(time.Second)
	media := fakeMedia{
		ID:       "mymedia",
		Size:     "1920x1080",
		Rotation: 90,
		Request: request{
			Format: []encodingcom.Format{
				{
					Bitrate:    "2500k",
					Size:       "0x1080",
					VideoCodec: "VP9",
					Output:     []string{hlsOutput},
				},
			},
		},
		Status:   "Saving",
		Created:  now.Add(-time.Hour),
		Started:  now.Add(-50 * time.Minute),
		Finished: now.Add(-10 * time.Minute),
	}
	server.medias["mymedia"] = &media
	client, err := encodingcom.NewClient(server.URL, "myuser", "secret")
	if err != nil {
		t.Fatal(err)
	}
	prov := encodingComProvider{client: client}
	prov.config = &config.Config{
		EncodingCom: &config.EncodingCom{
			Destination: "https://mybucket.s3.amazonaws.com/dir/",
		},
	}
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: "mymedia"})
	if err != nil {
		t.Fatal(err)
	}
	expected := provider.JobStatus{
		ProviderJobID: "mymedia",
		ProviderName:  "encoding.com",
		Status:        provider.StatusFinished,
		StatusMessage: "",
		Progress:      100,
		ProviderStatus: map[string]interface{}{
			"sourcefile":   "http://some.source.file",
			"timeleft":     "1",
			"created":      media.Created,
			"started":      media.Started,
			"finished":     media.Finished,
			"formatStatus": []string{""},
		},
		SourceInfo: provider.SourceInfo{
			Duration:   183e9,
			Width:      1920,
			Height:     1080,
			VideoCodec: "VP9",
		},
		Output: provider.JobOutput{
			Destination: "s3://mybucket/dir/job-123/",
			Files: []provider.OutputFile{
				{
					Path:       "s3://mybucket/dir/job-123/some_hls_preset/video-0.m3u8",
					VideoCodec: "VP9",
					Width:      608,
					Height:     1080,
					Container:  "m3u8",
					FileSize:   45674,
				},
				{
					Path:       "s3://mybucket/dir/job-123/video.m3u8",
					VideoCodec: "VP9",
					Width:      608,
					Height:     1080,
					Container:  "m3u8",
					FileSize:   45674,
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobStatus, expected) {
		t.Errorf("JobStatus: wrong job returned.\nWant %#v\nGot  %#v", expected, *jobStatus)
	}
}

func TestJobStatusNotFinished(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	now := time.Now().UTC().Truncate(time.Second)
	media := fakeMedia{
		ID: "mymedia",
		Request: request{
			Format: []encodingcom.Format{
				{
					Size:       "1920x1080",
					VideoCodec: "VP9",
					Output:     []string{hlsOutput},
				},
			},
		},
		Status:   "Saving",
		Created:  now,
		Started:  now.Add(10 * time.Minute),
		Finished: now.Add(time.Hour),
	}
	server.medias["mymedia"] = &media
	client, err := encodingcom.NewClient(server.URL, "myuser", "secret")
	if err != nil {
		t.Fatal(err)
	}
	prov := encodingComProvider{client: client}
	prov.config = &config.Config{
		EncodingCom: &config.EncodingCom{
			Destination: "https://mybucket.s3.amazonaws.com/dir/",
		},
	}
	jobStatus, err := prov.JobStatus(&db.Job{ID: "job-123", ProviderJobID: "mymedia"})
	if err != nil {
		t.Fatal(err)
	}
	expected := provider.JobStatus{
		ProviderJobID: "mymedia",
		ProviderName:  "encoding.com",
		Status:        provider.StatusStarted,
		StatusMessage: "",
		Progress:      100,
		ProviderStatus: map[string]interface{}{
			"sourcefile":   "http://some.source.file",
			"timeleft":     "1",
			"created":      media.Created,
			"started":      media.Started,
			"finished":     media.Finished,
			"formatStatus": []string{""},
		},
		Output: provider.JobOutput{
			Destination: "s3://mybucket/dir/job-123/",
			Files: []provider.OutputFile{
				{
					Path:       "s3://mybucket/dir/job-123/some_hls_preset/video-0.m3u8",
					Width:      1920,
					Height:     1080,
					VideoCodec: "VP9",
					Container:  "m3u8",
					FileSize:   45674,
				},
				{
					Path:       "s3://mybucket/dir/job-123/video.m3u8",
					Width:      1920,
					Height:     1080,
					VideoCodec: "VP9",
					Container:  "m3u8",
					FileSize:   45674,
				},
			},
		},
	}
	if !reflect.DeepEqual(*jobStatus, expected) {
		t.Errorf("JobStatus: wrong job returned.\nWant %#v.\nGot  %#v.", expected, *jobStatus)
	}
}

func TestJobStatusInvalidSourceInfo(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	now := time.Now().UTC().Truncate(time.Second)
	media1 := fakeMedia{
		ID:   "media1",
		Size: "1920x1080x900",
		Request: request{
			Format: []encodingcom.Format{
				{
					Size:       "1920x1080",
					VideoCodec: "VP9",
					Output:     []string{"webm"},
				},
			},
		},
		Status:   "Finished",
		Created:  now,
		Started:  now.Add(time.Minute),
		Finished: now.Add(time.Hour),
	}
	server.medias["media1"] = &media1
	media2 := fakeMedia{
		ID:   "media2",
		Size: "πx1080",
		Request: request{
			Format: []encodingcom.Format{
				{
					Size:       "1920x1080",
					VideoCodec: "VP9",
					Output:     []string{"webm"},
				},
			},
		},
		Status:   "Finished",
		Created:  now,
		Started:  now.Add(time.Minute),
		Finished: now.Add(time.Hour),
	}
	server.medias["media2"] = &media2
	media3 := fakeMedia{
		ID:   "media3",
		Size: "π",
		Request: request{
			Format: []encodingcom.Format{
				{
					Size:       "1920x1080",
					VideoCodec: "VP9",
					Output:     []string{"webm"},
				},
			},
		},
		Status:   "Finished",
		Created:  now,
		Started:  now.Add(time.Minute),
		Finished: now.Add(time.Hour),
	}
	server.medias["media3"] = &media3
	var tests = []struct {
		testCase string
		mediaID  string
		errMsg   string
	}{
		{
			"invalid media ID",
			"something",
			`Error returned by the Encoding.com API: {"Errors":["media not found"]}`,
		},
		{
			"invalid height",
			"media1",
			`invalid size returned by the Encoding.com API ("1920x1080x900"): strconv.ParseInt: parsing "1080x900": invalid syntax`,
		},
		{
			"invalid width",
			"media2",
			`invalid size returned by the Encoding.com API ("πx1080"): strconv.ParseInt: parsing "π": invalid syntax`,
		},
		{
			"invalid size",
			"media3",
			`invalid size returned by the Encoding.com API: "π"`,
		},
	}
	client, err := encodingcom.NewClient(server.URL, "myuser", "secret")
	if err != nil {
		t.Fatal(err)
	}
	prov := encodingComProvider{client: client}
	prov.config = &config.Config{
		EncodingCom: &config.EncodingCom{
			Destination: "mybucket",
		},
	}
	for _, test := range tests {
		jobStatus, err := prov.JobStatus(&db.Job{ProviderJobID: test.mediaID})
		if jobStatus != nil {
			t.Errorf("%s: got unexpected non-nil status: %#v", test.testCase, jobStatus)
		}
		if err == nil {
			t.Errorf("%s: got unexpected <nil> error", test.testCase)
			continue
		}
		if err.Error() != test.errMsg {
			t.Errorf("%s: wrong error message\nwant %q\ngot  %q", test.testCase, test.errMsg, err.Error())
		}
	}
}

func TestJobStatusMediaNotFound(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	provider := encodingComProvider{client: client}
	jobStatus, err := provider.JobStatus(&db.Job{ProviderJobID: "non-existent-job"})
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
		{"Waiting for encoder", provider.StatusQueued},
		{"Processing", provider.StatusStarted},
		{"Saving", provider.StatusStarted},
		{"Finished", provider.StatusFinished},
		{"Error", provider.StatusFailed},
		{"Unknown", provider.StatusUnknown},
		{"new", provider.StatusQueued},
		{"downloading", provider.StatusStarted},
		{"ready to process", provider.StatusStarted},
		{"waiting for encoder", provider.StatusQueued},
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

func TestCreatePreset(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	prov := encodingComProvider{client: client}
	presetName, err := prov.CreatePreset(db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
		Container:   "mp4",
		Description: "my nice preset",
		Name:        "mp4_1080p",
		RateControl: "VBR",
		Video: db.VideoPreset{
			Profile:      "main",
			ProfileLevel: "3.1",
			Bitrate:      "3500000",
			Codec:        "h264",
			GopMode:      "fixed",
			GopSize:      "90",
			Height:       "1080",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	fakePreset := server.presets[presetName]
	expectedFormat := encodingcom.Format{
		AudioCodec:   "dolby_aac",
		AudioBitrate: "128k",
		AudioVolume:  100,
		Output:       []string{"mp4"},
		Profile:      "main",
		TwoPass:      false,
		VideoCodec:   "libx264",
		Bitrate:      "3500k",
		Gop:          "cgop",
		Keyframe:     []string{"90"},
		Size:         "0x1080",
		Destination:  []string{"ftp://username:password@yourftphost.com/video/encoded/test.flv"},
	}
	if !reflect.DeepEqual(fakePreset.Request.Format[0], expectedFormat) {
		pretty.Fdiff(os.Stderr, fakePreset.Request.Format[0], expectedFormat)
		t.Errorf("wrong format provided\nWant %#v\nGot  %#v", expectedFormat, fakePreset.Request.Format[0])

	}
}

func TestCreatePresetHLS(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	prov := encodingComProvider{client: client}
	presetName, err := prov.CreatePreset(db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
		Container:   "m3u8",
		Description: "my nice preset",
		Name:        "mp4_1080p",
		RateControl: "VBR",
		Video: db.VideoPreset{
			Profile:      "main",
			ProfileLevel: "3.1",
			Bitrate:      "3500000",
			Codec:        "h264",
			GopMode:      "fixed",
			GopSize:      "90",
			Height:       "1080",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	fakePreset := server.presets[presetName]
	falseYesNoBoolean := encodingcom.YesNoBoolean(false)
	expectedFormat := encodingcom.Format{
		Output:      []string{hlsOutput},
		Destination: []string{"ftp://username:password@yourftphost.com/video/encoded/test.flv"},
		PackFiles:   &falseYesNoBoolean,
		Stream: []encodingcom.Stream{
			{
				AudioBitrate: "128k",
				AudioCodec:   "dolby_aac",
				AudioVolume:  100,
				Bitrate:      "3500k",
				Keyframe:     "90",
				Profile:      "main",
				Size:         "0x1080",
				TwoPass:      false,
				VideoCodec:   "libx264",
			},
		},
	}
	if !reflect.DeepEqual(fakePreset.Request.Format[0], expectedFormat) {
		t.Errorf("wrong format provided\nWant %#v\nGot  %#v", expectedFormat, fakePreset.Request.Format[0])
	}
}

func TestCreatePresetTwoPass(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	prov := encodingComProvider{client: client}
	presetName, err := prov.CreatePreset(db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
		Container:   "mp4",
		Description: "my nice preset",
		Name:        "mp4_1080p",
		RateControl: "VBR",
		TwoPass:     true,
		Video: db.VideoPreset{
			Profile:      "main",
			ProfileLevel: "3.1",
			Bitrate:      "3500000",
			Codec:        "h264",
			GopMode:      "fixed",
			GopSize:      "90",
			Height:       "1080",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	fakePreset := server.presets[presetName]
	expectedFormat := encodingcom.Format{
		AudioCodec:   "dolby_aac",
		AudioBitrate: "128k",
		AudioVolume:  100,
		Output:       []string{"mp4"},
		Profile:      "main",
		TwoPass:      true,
		VideoCodec:   "libx264",
		Bitrate:      "3500k",
		Gop:          "cgop",
		Keyframe:     []string{"90"},
		Size:         "0x1080",
		Destination:  []string{"ftp://username:password@yourftphost.com/video/encoded/test.flv"},
	}
	if !reflect.DeepEqual(fakePreset.Request.Format[0], expectedFormat) {
		pretty.Fdiff(os.Stderr, fakePreset.Request.Format[0], expectedFormat)
		t.Errorf("wrong format provided\nWant %#v\nGot  %#v", expectedFormat, fakePreset.Request.Format[0])

	}
}

func TestPresetToFormat(t *testing.T) {
	falseYesNoBoolean := encodingcom.YesNoBoolean(false)
	var tests = []struct {
		givenTestCase  string
		givenPreset    db.Preset
		expectedFormat encodingcom.Format
	}{
		{
			"HLS preset",
			db.Preset{
				Container: "m3u8",
				Audio: db.AudioPreset{
					Codec: "aac",
				},
				Video: db.VideoPreset{
					Profile:      "Main",
					ProfileLevel: "3.1",
					Codec:        "h264",
				},
			},
			encodingcom.Format{
				Output:               []string{"advanced_hls"},
				Destination:          []string{"ftp://username:password@yourftphost.com/video/encoded/test.flv"},
				VideoCodecParameters: encodingcom.VideoCodecParameters{},
				Stream: []encodingcom.Stream{
					{
						AudioCodec:  "dolby_aac",
						AudioVolume: 100,
						Size:        "0x0",
						VideoCodec:  "libx264",
						Profile:     "Main",
					},
				},
				PackFiles: &falseYesNoBoolean,
			},
		},
		{
			"WEBM vp8 preset",
			db.Preset{
				Container: "webm",
				Audio: db.AudioPreset{
					Codec: "vorbis",
				},
				Video: db.VideoPreset{
					Codec:   "vp8",
					GopSize: "90",
				},
			},
			encodingcom.Format{
				Output:      []string{"webm"},
				Destination: []string{"ftp://username:password@yourftphost.com/video/encoded/test.flv"},
				AudioCodec:  "libvorbis",
				AudioVolume: 100,
				Gop:         "cgop",
				Keyframe:    []string{"90"},
				VideoCodec:  "libvpx",
				Size:        "0x0",
			},
		},
		{
			"WEBM vp9 preset",
			db.Preset{
				Container: "webm",
				Audio: db.AudioPreset{
					Codec: "vorbis",
				},
				Video: db.VideoPreset{
					Codec:   "vp9",
					GopSize: "90",
				},
			},
			encodingcom.Format{
				Output:      []string{"webm"},
				Destination: []string{"ftp://username:password@yourftphost.com/video/encoded/test.flv"},
				AudioCodec:  "libvorbis",
				AudioVolume: 100,
				Gop:         "cgop",
				Keyframe:    []string{"90"},
				VideoCodec:  "libvpx-vp9",
				Size:        "0x0",
			},
		},
		{
			"MP4 preset",
			db.Preset{
				Container: "mp4",
				Audio: db.AudioPreset{
					Codec: "aac",
				},
				Video: db.VideoPreset{
					Profile:      "Main",
					ProfileLevel: "3.1",
					Codec:        "h264",
					GopSize:      "90",
				},
			},
			encodingcom.Format{
				Output:      []string{"mp4"},
				Destination: []string{"ftp://username:password@yourftphost.com/video/encoded/test.flv"},
				AudioCodec:  "dolby_aac",
				AudioVolume: 100,
				Gop:         "cgop",
				Keyframe:    []string{"90"},
				Size:        "0x0",
				VideoCodec:  "libx264",
				Profile:     "Main",
			},
		},
	}
	var p encodingComProvider
	for _, test := range tests {
		resultingFormat := p.presetToFormat(test.givenPreset)
		if !reflect.DeepEqual(resultingFormat, test.expectedFormat) {
			t.Errorf("%s: presetToFormat: wrong value. Want %#v. Got %#v", test.givenTestCase, test.expectedFormat, resultingFormat)
			pretty.Fdiff(os.Stderr, resultingFormat, test.expectedFormat)
		}
	}
}

func TestGetPreset(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	prov := encodingComProvider{client: client}
	presetName, err := prov.CreatePreset(db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
		Container:   "mp4",
		Description: "my nice preset",
		Name:        "mp4_1080p",
		RateControl: "VBR",
		Video: db.VideoPreset{
			Profile:      "main",
			ProfileLevel: "3.1",
			Bitrate:      "3500000",
			Codec:        "h264",
			GopMode:      "fixed",
			GopSize:      "90",
			Width:        "1920",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	preset, err := prov.GetPreset(presetName)
	if err != nil {
		t.Fatal(err)
	}
	expectedPreset := &encodingcom.Preset{
		Name: presetName,
		Format: convertFormat(encodingcom.Format{
			AudioCodec:   "dolby_aac",
			AudioBitrate: "128k",
			AudioVolume:  100,
			Output:       []string{"mp4"},
			Profile:      "main",
			TwoPass:      false,
			VideoCodec:   "libx264",
			Bitrate:      "3500k",
			Gop:          "cgop",
			Keyframe:     []string{"90"},
			Size:         "1920x0",
		}),
		Output: "mp4",
		Type:   encodingcom.UserPresets,
	}
	if !reflect.DeepEqual(preset, expectedPreset) {
		t.Errorf("GetPreset(%q): wrong preset returned.\nWant %#v\nGot  %#v", presetName, expectedPreset, preset)
	}
}

func TestGetPresetNotFound(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	prov := encodingComProvider{client: client}
	preset, err := prov.GetPreset("some-id")
	if preset != nil {
		t.Errorf("unexpected non-nil preset: %#v", preset)
	}
	if err == nil {
		t.Fatal("unexpected <nil> error")
	}
}

func TestDeletePreset(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	prov := encodingComProvider{client: client}
	presetName, err := prov.CreatePreset(db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
		Container:   "mp4",
		Description: "my nice preset",
		Name:        "mp4_1080p",
		RateControl: "VBR",
		Video: db.VideoPreset{
			Profile:      "Main",
			ProfileLevel: "3.1",
			Bitrate:      "3500000",
			Codec:        "h264",
			GopMode:      "fixed",
			GopSize:      "90",
			Width:        "1920",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = prov.DeletePreset(presetName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = prov.GetPreset(presetName)
	if err == nil {
		t.Error("did not delete the preset")
	}
}

func TestDeletePresetNotFound(t *testing.T) {
	server := newEncodingComFakeServer()
	defer server.Close()
	client, _ := encodingcom.NewClient(server.URL, "myuser", "secret")
	prov := encodingComProvider{client: client}
	err := prov.DeletePreset("some-preset")
	if err == nil {
		t.Error("unexpected <nil> error")
	}
}

func TestCancelJob(t *testing.T) {
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
	client, err := encodingcom.NewClient(server.URL, "user", "pass")
	if err != nil {
		t.Fatal(err)
	}
	prov := encodingComProvider{client: client}
	err = prov.CancelJob("mymedia")
	if err != nil {
		t.Fatal(err)
	}
	if media.Status != "Canceled" {
		t.Errorf("wrong status. Want %q. Got %q", "Canceled", media.Status)
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
