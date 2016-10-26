package redis

import (
	"reflect"
	"testing"

	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/db/redis/storage"
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
	preset := db.LocalPreset{}
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
	expectedItems := map[string]string{"name": "test"}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Wrong preset hash returned from Redis. Want %#v. Got %#v", expectedItems, items)
	}
}
