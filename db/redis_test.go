package db

import (
	"reflect"
	"testing"

	"github.com/nytm/video-transcoding-api/config"
	"gopkg.in/redis.v3"
)

func TestSaveJob(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	repo, err := NewRedisJobRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	job := Job{Status: "Downloading", ProviderName: "encoding.com"}
	err = repo.SaveJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	if job.ID == "" {
		t.Fatal("Job ID should have been generated on SaveJob")
	}
	client := cfg.RedisClient()
	defer client.Close()
	items, err := client.HGetAll("job:" + job.ID).Result()
	if err != nil {
		t.Fatal(err)
	}
	jobMap := make(map[string]string)
	for i, item := range items {
		switch item {
		case "providerName", "providerJobID", "status":
			jobMap[item] = items[i+1]
		}
	}
	expected := map[string]string{
		"providerName":  "encoding.com",
		"providerJobID": "",
		"status":        "Downloading",
	}
	if !reflect.DeepEqual(jobMap, expected) {
		t.Errorf("Wrong job hash returned from Redis. Want %#v. Got %#v.", expected, jobMap)
	}
}

func TestSaveJobPredefinedID(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	repo, err := NewRedisJobRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	job := Job{ID: "myjob", Status: "Downloaded", ProviderName: "encoding.com"}
	err = repo.SaveJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != "myjob" {
		t.Errorf("Job ID should not be regenerated when it's already defined. Got %q instead of %q.", job.ID, "myjob")
	}
	client := cfg.RedisClient()
	defer client.Close()
	items, err := client.HGetAll("job:myjob").Result()
	if err != nil {
		t.Fatal(err)
	}
	jobMap := make(map[string]string)
	for i, item := range items {
		switch item {
		case "providerName", "providerJobID", "status":
			jobMap[item] = items[i+1]
		}
	}
	expected := map[string]string{
		"providerName":  "encoding.com",
		"providerJobID": "",
		"status":        "Downloaded",
	}
	if !reflect.DeepEqual(jobMap, expected) {
		t.Errorf("Wrong job hash returned from Redis. Want %#v. Got %#v.", expected, jobMap)
	}
}

func cleanRedis() error {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	defer client.Close()
	keys, err := client.Keys("job:*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		_, err = client.Del(keys...).Result()
	}
	return err
}
