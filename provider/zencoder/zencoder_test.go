package zencoder

import (
	"errors"
	"reflect"
	"testing"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
	"github.com/NYTimes/video-transcoding-api/provider"
	"github.com/brandscreen/zencoder"
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
	provider, err := zencoderFactory(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	zencoderProvider, ok := provider.(*zencoderProvider)
	if !ok {
		t.Fatalf("Wrong provider returned. Want zencoderProvider instance. Got %#v.", provider)
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
	provider, err := zencoderFactory(&cfg)
	if provider == nil {
		t.Errorf("Unexpected nil provider: %#v", provider)
	}
	if err != nil {
		t.Errorf("Unexpected Error returned. Got %#v", err)
	}

	cfg = config.Config{Zencoder: &config.Zencoder{APIKey: ""}}
	provider, err = zencoderFactory(&cfg)
	if provider != nil {
		t.Errorf("Unexpected non-nil provider: %#v", provider)
	}
	if err != errZencoderInvalidConfig {
		t.Errorf("Wrong error returned. Want errZencoderInvalidConfig. Got %#v", err)
	}
}

func TestCapabilities(t *testing.T) {
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

func TestCreatePreset(t *testing.T) {
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
