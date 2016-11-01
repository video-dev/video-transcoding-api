package zencoder

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/brandscreen/zencoder"
	"github.com/kr/pretty"
	redisDriver "gopkg.in/redis.v4"
)

func TestFactoryIsRegistered(t *testing.T) {
	_, err := provider.GetProviderFactory(Name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestZencoderFactory(t *testing.T) {
	cfg := config.Config{
		Zencoder: &config.Zencoder{
			APIKey: "api-key-here",
		},
	}
	prov, err := zencoderFactory(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	zencoderProvider, ok := prov.(*zencoderProvider)
	if !ok {
		t.Fatalf("Wrong provider returned. Want zencoderProvider instance. Got %#v.", prov)
	}
	expected := zencoder.NewZencoder("api-key-here")
	if !reflect.DeepEqual(zencoderProvider.client, expected) {
		t.Errorf("Factory: wrong client returned. Want %#v. Got %#v.", expected, zencoderProvider.client)
	}
	if !reflect.DeepEqual(zencoderProvider.config, &cfg) {
		t.Errorf("Factory: wrong config returned. Want %#v. Got %#v.", &cfg, zencoderProvider.config)
	}
}

func TestZencoderFactoryValidation(t *testing.T) {
	cfg := config.Config{Zencoder: &config.Zencoder{APIKey: "api-key"}}
	prov, err := zencoderFactory(&cfg)
	if prov == nil {
		t.Errorf("Unexpected nil provider: %#v", prov)
	}
	if err != nil {
		t.Errorf("Unexpected Error returned. Got %#v", err)
	}

	cfg = config.Config{Zencoder: &config.Zencoder{APIKey: ""}}
	prov, err = zencoderFactory(&cfg)
	if prov != nil {
		t.Errorf("Unexpected non-nil provider: %#v", prov)
	}
	if err != errZencoderInvalidConfig {
		t.Errorf("Wrong error returned. Want errZencoderInvalidConfig. Got %#v", err)
	}
}

func TestZencoderCapabilities(t *testing.T) {
	var prov zencoderProvider
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

func TestZencoderCreatePreset(t *testing.T) {
	cleanLocalPresets()
	cfg := config.Config{
		Zencoder: &config.Zencoder{APIKey: "api-key-here"},
		Redis:    new(storage.Config),
	}
	preset := db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
		Container:    "mp4",
		Description:  "my nice preset",
		Name:         "mp4_1080p",
		Profile:      "main",
		ProfileLevel: "3.1",
		RateControl:  "VBR",
		Video: db.VideoPreset{
			Bitrate: "3500000",
			Codec:   "h264",
			GopMode: "fixed",
			GopSize: "90",
			Height:  "1080",
		},
	}
	provider, err := zencoderFactory(&cfg)
	repo, err := redis.NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	presetName, err := provider.CreatePreset(preset)
	if err != nil {
		t.Fatal(err)
	}
	expected := &db.LocalPreset{
		Name:   "mp4_1080p",
		Preset: preset,
	}
	res, err := repo.GetLocalPreset(presetName)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Got wrong preset. Want %#v. Got %#v", expected, res)
	}
}

func TestCreatePresetError(t *testing.T) {
	cleanLocalPresets()
	cfg := config.Config{
		Zencoder: &config.Zencoder{APIKey: "api-key-here"},
		Redis:    new(storage.Config),
	}
	preset := db.Preset{}
	provider, err := zencoderFactory(&cfg)

	_, err = provider.CreatePreset(preset)
	if !reflect.DeepEqual(err, errors.New("preset name missing")) {
		t.Errorf("Got wrong error. Want %#v. Got %#v", errors.New("preset name missing"), err)
	}
}

func TestGetPreset(t *testing.T) {
	cleanLocalPresets()
	cfg := config.Config{
		Zencoder: &config.Zencoder{APIKey: "api-key-here"},
		Redis:    new(storage.Config),
	}
	preset := db.Preset{
		Name: "get_preset",
		Video: db.VideoPreset{
			Bitrate: "3500000",
			Codec:   "h264",
			GopMode: "fixed",
			GopSize: "90",
			Height:  "1080",
		},
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
	}
	provider, err := zencoderFactory(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	presetName, err := provider.CreatePreset(preset)
	if err != nil {
		t.Fatal(err)
	}
	expected := &db.LocalPreset{
		Name:   "get_preset",
		Preset: preset,
	}
	res, err := provider.GetPreset(presetName)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Got wrong preset. Want %#v. Got %#v", expected, res)
	}
}

func TestZencoderDeletePreset(t *testing.T) {
	cleanLocalPresets()
	cfg := config.Config{
		Zencoder: &config.Zencoder{APIKey: "api-key-here"},
		Redis:    new(storage.Config),
	}
	preset := db.Preset{
		Name: "get_preset",
		Video: db.VideoPreset{
			Bitrate: "3500000",
			Codec:   "h264",
			GopMode: "fixed",
			GopSize: "90",
			Height:  "1080",
		},
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
	}
	prov, err := zencoderFactory(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	presetName, err := prov.CreatePreset(preset)
	if err != nil {
		t.Fatal(err)
	}
	err = prov.DeletePreset(presetName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = prov.GetPreset(presetName)
	if err != db.ErrLocalPresetNotFound {
		t.Errorf("Got wrong error. Want errLocalPresetNotFound. Got %#v", err)
	}
}

func TestZencoderTranscode(t *testing.T) {
	cleanLocalPresets()
	cfg := config.Config{
		Zencoder: &config.Zencoder{APIKey: "api-key-here"},
		Redis:    new(storage.Config),
	}
	fakeZencoder := &FakeZencoder{}
	dbRepo, err := redis.NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	prov := &zencoderProvider{
		config: &cfg,
		client: fakeZencoder,
		db:     dbRepo,
	}
	preset := db.Preset{
		Audio: db.AudioPreset{
			Bitrate: "128000",
			Codec:   "aac",
		},
		Container:    "mp4",
		Description:  "my nice preset",
		Name:         "mp4_1080p",
		Profile:      "main",
		ProfileLevel: "3.1",
		RateControl:  "VBR",
		Video: db.VideoPreset{
			Bitrate: "3500000",
			Codec:   "h264",
			GopMode: "fixed",
			GopSize: "90",
			Height:  "1080",
			Width:   "720",
		},
	}
	_, err = prov.CreatePreset(preset)
	if err != nil {
		t.Fatal(err)
	}
	outputs := []provider.TranscodeOutput{
		{
			FileName: "output-720p.mp4",
			Preset: db.PresetMap{
				Name: "mp4_1080p",
				ProviderMapping: map[string]string{
					Name:    "93239832-0001",
					"other": "irrelevant",
				},
				OutputOpts: db.OutputOptions{Extension: "mp4"},
			},
		},
	}
	transcodeProfile := provider.TranscodeProfile{
		SourceMedia:     "dir/file.mov",
		Outputs:         outputs,
		StreamingParams: provider.StreamingParams{},
	}
	jobStatus, err := prov.Transcode(&db.Job{ID: "job-123"}, transcodeProfile)
	if err != nil {
		t.Fatal(err)
	}
	if jobStatus.ProviderJobID != "123" {
		t.Errorf("Got wrong jobStatus ID. Expected 123, got %#v", jobStatus.ProviderJobID)
	}
}

func TestZencoderBuildOutput(t *testing.T) {
	prov := &zencoderProvider{}
	var tests = []struct {
		Description string
		Preset      db.Preset
		Expected    map[string]interface{}
	}{
		{
			"Test with mp4 preset",
			db.Preset{
				Name:         "mp4_1080p",
				Description:  "my nice preset",
				Container:    "mp4",
				Profile:      "main",
				ProfileLevel: "3.1",
				RateControl:  "CBR",
				Video: db.VideoPreset{
					Bitrate: "3500000",
					Codec:   "h264",
					GopMode: "fixed",
					GopSize: "90",
					Height:  "1080",
					Width:   "1920",
				},
				Audio: db.AudioPreset{
					Bitrate: "128000",
					Codec:   "aac",
				},
			},
			map[string]interface{}{
				"label":                   "mp4_1080p:my nice preset",
				"format":                  "mp4",
				"video_codec":             "h264",
				"h264_profile":            "main",
				"h264_level":              "3.1",
				"audio_codec":             "aac",
				"width":                   float64(1920),
				"height":                  float64(1080),
				"video_bitrate":           float64(3500),
				"audio_bitrate":           float64(128),
				"keyframe_interval":       float64(90),
				"fixed_keyframe_interval": true,
				"constant_bitrate":        true,
				"deinterlace":             "on",
			},
		},
		{
			"Test with webm preset",
			db.Preset{
				Name:        "webm_1080p",
				Description: "my vp8 preset",
				Container:   "webm",
				Video: db.VideoPreset{
					Bitrate: "3500000",
					Codec:   "vp8",
					GopSize: "90",
					Height:  "1080",
					Width:   "1920",
				},
				Audio: db.AudioPreset{
					Bitrate: "128000",
					Codec:   "aac",
				},
			},
			map[string]interface{}{
				"label":             "webm_1080p:my vp8 preset",
				"format":            "webm",
				"video_codec":       "vp8",
				"audio_codec":       "aac",
				"width":             float64(1920),
				"height":            float64(1080),
				"video_bitrate":     float64(3500),
				"audio_bitrate":     float64(128),
				"keyframe_interval": float64(90),
				"deinterlace":       "on",
			},
		},
	}

	for _, test := range tests {
		res, err := prov.buildOutput(test.Preset)
		if err != nil {
			t.Fatal(err)
		}
		resultJSON, err := json.Marshal(res)
		if err != nil {
			t.Fatal(err)
		}
		result := make(map[string]interface{})
		err = json.Unmarshal(resultJSON, &result)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(result, test.Expected) {
			pretty.Fdiff(os.Stderr, test.Expected, result)
			t.Errorf("Failed to build output. Want\n %+v. Got\n %+v.", test.Expected, result)
		}
	}
}

func TestZencoderHealthcheck(t *testing.T) {
	cfg := config.Config{
		Zencoder: &config.Zencoder{APIKey: "api-key-here"},
		Redis:    new(storage.Config),
	}
	fakeZencoder := &FakeZencoder{}
	dbRepo, err := redis.NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	prov := &zencoderProvider{
		config: &cfg,
		client: fakeZencoder,
		db:     dbRepo,
	}

	err = prov.Healthcheck()
	if err != nil {
		t.Fatal(err)
	}
}

func TestZencoderCancelJob(t *testing.T) {
	cfg := config.Config{
		Zencoder: &config.Zencoder{APIKey: "api-key-here"},
		Redis:    new(storage.Config),
	}
	fakeZencoder := &FakeZencoder{}
	dbRepo, err := redis.NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	prov := &zencoderProvider{
		config: &cfg,
		client: fakeZencoder,
		db:     dbRepo,
	}

	err = prov.CancelJob("123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestZencoderJobStatus(t *testing.T) {
	cfg := config.Config{
		Zencoder: &config.Zencoder{APIKey: "api-key-here"},
		Redis:    new(storage.Config),
	}
	fakeZencoder := &FakeZencoder{}
	dbRepo, err := redis.NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	prov := &zencoderProvider{
		config: &cfg,
		client: fakeZencoder,
		db:     dbRepo,
	}
	jobStatus, err := prov.JobStatus(&db.Job{
		ProviderJobID: "1234567890",
	})
	if err != nil {
		t.Fatal(err)
	}
	resultJSON, err := json.Marshal(jobStatus)
	if err != nil {
		t.Fatal(err)
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(resultJSON, &result)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"providerName":  "zencoder",
		"providerJobId": "1234567890",
		"status":        "Transcoding",
		"progress":      float64(10),
		"sourceInfo": map[string]interface{}{
			"duration":   float64(10000000),
			"height":     float64(1080),
			"width":      float64(1920),
			"videoCodec": "ProRes422",
		},
		"output": map[string]interface{}{
			"files": []interface{}{
				map[string]interface{}{
					"path":       "http://nyt.net/output1.mp4",
					"container":  "mp4",
					"videoCodec": "h264",
					"height":     float64(1080),
					"width":      float64(1920),
				},
				map[string]interface{}{
					"height":     float64(720),
					"width":      float64(1080),
					"path":       "http://nyt.net/output2.webm",
					"container":  "webm",
					"videoCodec": "vp8",
				},
			},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		pretty.Fdiff(os.Stderr, expected, result)
		t.Errorf("Wrong JobStatus returned. Want %#v. Got %#v.", expected, result)
	}
}

func cleanLocalPresets() error {
	client := redisDriver.NewClient(&redisDriver.Options{Addr: "127.0.0.1:6379"})
	defer client.Close()
	err := deleteKeys("localpreset:*", client)
	err = deleteKeys("localpresets", client)
	return err
}

func deleteKeys(pattern string, client *redisDriver.Client) error {
	keys, err := client.Keys(pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		_, err = client.Del(keys...).Result()
	}
	return err
}
