package redis

import (
	"reflect"
	"testing"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
)

func TestCreateLocalPreset(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(storage.Config)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	preset := db.LocalPreset{
		Name: "test",
		Preset: db.Preset{
			Name: "test",
		},
	}
	err = repo.CreateLocalPreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).storage.RedisClient()
	defer client.Close()
	items, err := client.HGetAll("localpreset:" + preset.Name).Result()
	if err != nil {
		t.Fatal(err)
	}
	expectedItems := map[string]string{
		"preset_name":    "test",
		"preset_twopass": "false",
	}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Wrong preset hash returned from Redis. Want %#v. Got %#v", expectedItems, items)
	}
}

func TestCreateLocalPresetDuplicate(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	preset := db.LocalPreset{
		Name: "test",
		Preset: db.Preset{
			Name: "test",
		},
	}
	err = repo.CreateLocalPreset(&preset)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.CreateLocalPreset(&preset)
	if err != db.ErrLocalPresetAlreadyExists {
		t.Errorf("Got wrong error. Want %#v. Got %#v", db.ErrLocalPresetAlreadyExists, err)
	}
}

func TestUpdateLocalPreset(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	preset := db.LocalPreset{
		Name: "test",
		Preset: db.Preset{
			Name: "test",
		},
	}
	err = repo.CreateLocalPreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	preset.Preset.Name = "test-different"

	err = repo.UpdateLocalPreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).storage.RedisClient()
	defer client.Close()
	items, err := client.HGetAll("localpreset:" + preset.Name).Result()
	if err != nil {
		t.Fatal(err)
	}
	expectedItems := map[string]string{
		"preset_name":    "test-different",
		"preset_twopass": "false",
	}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Wrong presetmap hash returned from Redis. Want %#v. Got %#v", expectedItems, items)
	}
}

func TestUpdateLocalPresetNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	err = repo.UpdateLocalPreset(&db.LocalPreset{
		Name: "non-existent",
		Preset: db.Preset{
			Name: "test",
		},
	})
	if err != db.ErrLocalPresetNotFound {
		t.Errorf("Wrong error returned by UpdateLocalPreset. Want ErrLocalPresetNotFound. Got %#v.", err)
	}
}

func TestDeleteLocalPreset(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}

	preset := db.LocalPreset{
		Name: "test",
		Preset: db.Preset{
			Name: "test",
		},
	}
	err = repo.CreateLocalPreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeleteLocalPreset(&db.LocalPreset{Name: preset.Name})
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).storage.RedisClient()
	result := client.HGetAll("localpreset:test")
	if len(result.Val()) != 0 {
		t.Errorf("Unexpected value after delete call: %v", result.Val())
	}
}

func TestDeleteLocalPresetNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeleteLocalPreset(&db.LocalPreset{Name: "non-existent"})
	if err != db.ErrLocalPresetNotFound {
		t.Errorf("Wrong error returned by DeleteLocalPreset. Want ErrLocalPresetNotFound. Got %#v.", err)
	}
}
