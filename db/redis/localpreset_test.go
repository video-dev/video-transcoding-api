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
		Preset: map[string]string{
			"videoCodec": "h264",
			"audioCodec": "aac",
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
		"preset_audioCodec": "aac",
		"preset_videoCodec": "h264",
	}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Wrong preset hash returned from Redis. Want %#v. Got %#v", expectedItems, items)
	}
}

func TestNothing(t *testing.T) {
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
	preset := db.LocalPreset{}
	repo.UpdateLocalPreset(&preset)
	repo.DeleteLocalPreset(&preset)
	repo.GetLocalPreset("test")
	repo.ListLocalPresets()

}
