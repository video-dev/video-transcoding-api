package db

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nytm/video-transcoding-api/config"
	"gopkg.in/redis.v3"
)

func TestSaveJob(t *testing.T) {
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
	job := Job{ProviderName: "encoding.com"}
	err = repo.SaveJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	if job.ID == "" {
		t.Fatal("Job ID should have been generated on SaveJob")
	}
	client := repo.(*redisRepository).redisClient()
	defer client.Close()
	items, err := client.HGetAllMap("job:" + job.ID).Result()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{
		"providerName":  "encoding.com",
		"providerJobID": "",
	}
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Wrong job hash returned from Redis. Want %#v. Got %#v.", expected, items)
	}
}

func TestSaveJobPredefinedID(t *testing.T) {
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
	job := Job{ID: "myjob", ProviderName: "encoding.com"}
	err = repo.SaveJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != "myjob" {
		t.Errorf("Job ID should not be regenerated when it's already defined. Got %q instead of %q.", job.ID, "myjob")
	}
	client := repo.(*redisRepository).redisClient()
	defer client.Close()
	items, err := client.HGetAllMap("job:myjob").Result()
	expected := map[string]string{
		"providerName":  "encoding.com",
		"providerJobID": "",
	}
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Wrong job hash returned from Redis. Want %#v. Got %#v.", expected, items)
	}
}

func TestSaveJobIsSafe(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	jobs := []Job{
		{ID: "abcabc", ProviderName: "elastictranscoder"},
		{ID: "abcabc", ProviderJobID: "abf-123", ProviderName: "encoding.com"},
		{ID: "abcabc", ProviderJobID: "abc-213", ProviderName: "encoding.com"},
		{ID: "abcabc", ProviderJobID: "ff12", ProviderName: "encoding.com"},
	}
	repo, err := NewRedisRepository(&config.Config{Redis: new(config.Redis)})
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := repo.SaveJob(&jobs[i])
			if err != nil && err != redis.TxFailedErr {
				t.Error(err)
			}
		}(i % len(jobs))
	}
	wg.Wait()
}

func TestDeleteJob(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRedisRepository(&config.Config{Redis: new(config.Redis)})
	if err != nil {
		t.Fatal(err)
	}
	job := Job{ID: "myjob"}
	err = repo.SaveJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeleteJob(&Job{ID: job.ID})
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).redisClient()
	result := client.HGetAllMap("job:myjob")
	if len(result.Val()) != 0 {
		t.Errorf("Unexpected value after delete call: %v", result.Val())
	}
}

func TestDeleteJobNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRedisRepository(&config.Config{Redis: new(config.Redis)})
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeleteJob(&Job{ID: "myjob"})
	if err != ErrJobNotFound {
		t.Errorf("Wrong error returned by DeleteJob. Want ErrJobNotFound. Got %#v.", err)
	}
}

func TestGetJob(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRedisRepository(&config.Config{Redis: new(config.Redis)})
	if err != nil {
		t.Fatal(err)
	}
	job := Job{ID: "myjob"}
	err = repo.SaveJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	gotJob, err := repo.GetJob(job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*gotJob, job) {
		t.Errorf("Wrong job. Want %#v. Got %#v.", job, *gotJob)
	}
}

func TestGetJobNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRedisRepository(&config.Config{Redis: new(config.Redis)})
	if err != nil {
		t.Fatal(err)
	}
	gotJob, err := repo.GetJob("job:myjob")
	if err != ErrJobNotFound {
		t.Errorf("Wrong error returned. Want ErrJobNotFound. Got %#v.", err)
	}
	if gotJob != nil {
		t.Errorf("Unexpected non-nil job: %#v.", gotJob)
	}
}

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
	preset := Preset{
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
	preset := Preset{
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
	preset := Preset{ID: "mypreset", ProviderMapping: map[string]string{"elemental": "abc123"}}
	err = repo.SavePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeletePreset(&Preset{ID: preset.ID})
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
	err = repo.DeletePreset(&Preset{ID: "mypreset"})
	if err != ErrPresetNotFound {
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
	preset := Preset{
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
	if err != ErrPresetNotFound {
		t.Errorf("Wrong error returned. Want ErrPresetNotFound. Got %#v.", err)
	}
	if gotPreset != nil {
		t.Errorf("Unexpected non-nil preset: %#v.", gotPreset)
	}
}

func cleanRedis() error {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	defer client.Close()
	err := deleteKeys("job:*", client)
	if err != nil {
		return err
	}
	return deleteKeys("preset:*", client)
}

func deleteKeys(pattern string, client *redis.Client) error {
	keys, err := client.Keys(pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		_, err = client.Del(keys...).Result()
	}
	return err
}

func TestRedisClientRedisDefaultConfig(t *testing.T) {
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRedisRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).redisClient()
	defer client.Close()
	_, err = client.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisClientRedisAddr(t *testing.T) {
	proc, err := startRedis("49153", "not-secret")
	if err != nil {
		t.Fatal(err)
	}
	defer proc.Signal(os.Interrupt)
	cfg := config.Config{
		Redis: &config.Redis{
			RedisAddr: "127.0.0.1:49153",
			Password:  "not-secret",
		},
	}
	repo, err := NewRedisRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).redisClient()
	defer client.Close()
	_, err = client.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisClientRedisSentinel(t *testing.T) {
	cleanup, err := startSentinels([]string{"26379", "26380", "26381"})
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	cfg := config.Config{
		Redis: &config.Redis{
			SentinelAddrs:      "127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381",
			SentinelMasterName: "mymaster",
		},
	}
	repo, err := NewRedisRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).redisClient()
	defer client.Close()
	_, err = client.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}
}

func startRedis(port, password string) (*os.Process, error) {
	configLines := []string{"port " + port}
	if password != "" {
		configLines = append(configLines, "requirepass "+password)
	}
	cmd := exec.Command("redis-server", "-")
	cmd.Dir = os.TempDir()
	cmd.Stdin = strings.NewReader(strings.Join(configLines, "\n"))
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	waitListening(30, "127.0.0.1:"+port)
	return cmd.Process, nil
}

func startSentinels(ports []string) (func(), error) {
	processes := make([]*os.Process, len(ports))
	tempFiles := make([]string, len(ports))
	addrs := make([]string, len(ports))
	configTemplate, err := ioutil.ReadFile("testdata/sentinel.conf")
	if err != nil {
		return nil, err
	}
	for i, port := range ports {
		f, err := ioutil.TempFile("", "sentinel")
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(f, string(configTemplate), port)
		f.Close()
		tempFiles[i] = f.Name()
		cmd := exec.Command("redis-server", f.Name(), "--sentinel")
		cmd.Dir = os.TempDir()
		err = cmd.Start()
		if err != nil {
			for j := 0; j < i; j++ {
				processes[j].Signal(os.Interrupt)
			}
			return nil, err
		}
		processes[i] = cmd.Process
		addrs[i] = "127.0.0.1:" + port
	}
	waitListening(30, addrs...)
	return func() {
		for i, process := range processes {
			process.Signal(os.Interrupt)
			os.Remove(tempFiles[i])
		}
	}, nil
}

func waitListening(maxTries int, addrs ...string) {
	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			for i := 0; i < maxTries; i++ {
				if conn, err := net.Dial("tcp", addr); err == nil {
					conn.Close()
					return
				} else {
					time.Sleep(10e6)
				}
			}
		}(addr)
	}
	wg.Wait()
}
