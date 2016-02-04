package redis

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nytm/video-transcoding-api/config"
	"gopkg.in/redis.v3"
)

func TestRedisClientRedisDefaultConfig(t *testing.T) {
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
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
	repo, err := NewRepository(&cfg)
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
	repo, err := NewRepository(&cfg)
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
