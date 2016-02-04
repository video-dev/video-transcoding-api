package redis

import (
	"reflect"
	"testing"

	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
)

func TestSavePreset(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRedisRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	preset := db.Preset{
		ProviderMapping: map[string]string{
			"elementalcloud":    "abc123",
			"elastictranscoder": "1281742-93939",
		},
	}
	err = repo.SavePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	if preset.ID == "" {
		t.Fatal("Preset ID should have been generated on SavePreset")
	}
	client := repo.(*redisRepository).redisClient()
	defer client.Close()
	items, err := client.HGetAllMap("preset:" + preset.ID).Result()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(items, preset.ProviderMapping) {
		t.Errorf("Wrong preset hash returned from Redis. Want %#v. Got %#v", preset.ProviderMapping, items)
	}
}

func TestSavePresetPredefinedID(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRedisRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	preset := db.Preset{
		ID:              "mypreset",
		ProviderMapping: map[string]string{"elemental": "123"},
	}
	err = repo.SavePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	if preset.ID != "mypreset" {
		t.Errorf("Preset ID should not be regenerated when it's already defined. Got %q instead of %q.", preset.ID, "preset")
	}
	client := repo.(*redisRepository).redisClient()
	defer client.Close()
	items, err := client.HGetAllMap("preset:mypreset").Result()
	if !reflect.DeepEqual(items, preset.ProviderMapping) {
		t.Errorf("Wrong preset hash returned from Redis. Want %#v. Got %#v.", preset.ProviderMapping, items)
	}
}

func TestDeletePreset(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRedisRepository(&config.Config{Redis: new(config.Redis)})
	if err != nil {
		t.Fatal(err)
	}
	preset := db.Preset{ID: "mypreset", ProviderMapping: map[string]string{"elemental": "abc123"}}
	err = repo.SavePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeletePreset(&db.Preset{ID: preset.ID})
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).redisClient()
	result := client.HGetAllMap("preset:mypreset")
	if len(result.Val()) != 0 {
		t.Errorf("Unexpected value after delete call: %v", result.Val())
	}
}

func TestDeletePresetNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRedisRepository(&config.Config{Redis: new(config.Redis)})
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeletePreset(&db.Preset{ID: "mypreset"})
	if err != db.ErrPresetNotFound {
		t.Errorf("Wrong error returned by DeletePreset. Want ErrPresetNotFound. Got %#v.", err)
	}
}

func TestGetPreset(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRedisRepository(&config.Config{Redis: new(config.Redis)})
	if err != nil {
		t.Fatal(err)
	}
	preset := db.Preset{
		ID: "mypreset",
		ProviderMapping: map[string]string{
			"elementalcloud":    "abc-123",
			"elastictranscoder": "0129291-0001",
			"encoding.com":      "wait what?",
		},
	}
	err = repo.SavePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	gotPreset, err := repo.GetPreset(preset.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*gotPreset, preset) {
		t.Errorf("Wrong preset. Want %#v. Got %#v.", preset, *gotPreset)
	}
}

func TestGetPresetNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRedisRepository(&config.Config{Redis: new(config.Redis)})
	if err != nil {
		t.Fatal(err)
	}
	gotPreset, err := repo.GetPreset("mypreset")
	if err != db.ErrPresetNotFound {
		t.Errorf("Wrong error returned. Want ErrPresetNotFound. Got %#v.", err)
	}
	if gotPreset != nil {
		t.Errorf("Unexpected non-nil preset: %#v.", gotPreset)
	}
}
