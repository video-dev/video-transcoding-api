package redis

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

func TestSave(t *testing.T) {
	person := Person{
		ID:   "some-id",
		Name: "gopher",
		Address: Address{
			Data: map[string]string{"first_line": "secret"},
			City: &City{Name: "nyc"},
		},
		NonTagged:        "not relevant",
		unexported:       "not relevant",
		unexportedTagged: "not relevant, believe!",
	}
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	redisRepo := repo.(*redisRepository)
	err = redisRepo.save("person:test", person)
	if err != nil {
		t.Fatal(err)
	}
	client := redisRepo.redisClient()
	defer client.Close()
	defer client.Del("person:test")
	data, err := client.HGetAllMap("person:test").Result()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{
		"name":       "gopher",
		"first_line": "secret",
		"city_name":  "nyc",
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Did not save properly. Want %#v. Got %#v", expected, data)
	}
}

func TestSavePointer(t *testing.T) {
	person := Person{
		ID:   "some-id",
		Name: "gopher",
		Address: Address{
			Data: map[string]string{"first_line": "secret"},
			City: &City{Name: "nyc"},
		},
		NonTagged:        "not relevant",
		unexported:       "not relevant",
		unexportedTagged: "not relevant, believe!",
	}
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	redisRepo := repo.(*redisRepository)
	err = redisRepo.save("person:test", &person)
	if err != nil {
		t.Fatal(err)
	}
	client := redisRepo.redisClient()
	defer client.Close()
	defer client.Del("person:test")
	data, err := client.HGetAllMap("person:test").Result()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{
		"name":       "gopher",
		"first_line": "secret",
		"city_name":  "nyc",
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Did not save properly. Want %#v. Got %#v", expected, data)
	}
}

func TestSaveMap(t *testing.T) {
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	redisRepo := repo.(*redisRepository)
	input := map[string]string{
		"name": "John",
		"test": "tested",
	}
	err = redisRepo.save("map:test", input)
	if err != nil {
		t.Fatal(err)
	}
	client := redisRepo.redisClient()
	defer client.Close()
	defer client.Del("map:test")
	data, err := client.HGetAllMap("map:test").Result()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(data, input) {
		t.Errorf("Did not save properly. Want %#v. Got %#v", input, data)
	}
}

func TestSaveErrors(t *testing.T) {
	var tests = []struct {
		input  interface{}
		errMsg string
	}{
		{nil, "no fields provided"},
		{10, "please provide a map or a struct"},
		{map[int]int{10: 12}, "please provide a map[string]string"},
		{map[string]int{"10": 12}, "please provide a map[string]string"},
		{map[string]string{}, "please provide a map[string]string with at least one item"},
		{struct {
			Name string `redis-hash:",expand"`
		}{}, "can only expand structs and maps"},
		{struct {
			Data map[int]int `redis-hash:",expand"`
		}{}, "please provide a map[string]string"},
		{struct {
			Inner struct {
				Data map[int]int `redis-hash:",expand"`
			} `redis-hash:",expand"`
		}{}, "please provide a map[string]string"},
	}

	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	redisRepo := repo.(*redisRepository)
	client := redisRepo.redisClient()
	defer client.Close()
	for _, test := range tests {
		err := redisRepo.save("some-key", test.input)
		if err == nil {
			client.Del("some-key")
			t.Error("Got unexpected nil error")
			continue
		}
		if err.Error() != test.errMsg {
			t.Errorf("Got wrong error message. Want %q. Got %q", test.errMsg, err.Error())
		}
	}
}

func TestLoadStruct(t *testing.T) {
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	redisRepo := repo.(*redisRepository)
	client := redisRepo.redisClient()
	defer client.Close()
	err = redisRepo.save("test-key", map[string]string{"name": "Gopher", "city_name": "New York"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Del("test-key")
	person := Person{
		Address:          Address{City: new(City), Data: make(map[string]string)},
		unexported:       "don't change",
		unexportedTagged: "don't change",
	}
	expectedPerson := person
	expectedPerson.Address.City = &City{Name: "New York"}
	expectedPerson.Name = "Gopher"
	err = redisRepo.load("test-key", &person)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(person, expectedPerson) {
		t.Errorf("Didn't load data to struct. Want %#v. Got %#v.", expectedPerson, person)
	}
}

func TestLoadMap(t *testing.T) {
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	redisRepo := repo.(*redisRepository)
	client := redisRepo.redisClient()
	defer client.Close()
	input := map[string]string{"name": "Gopher", "city_name": "New York"}
	err = redisRepo.save("test-key", input)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Del("test-key")
	person := make(map[string]string)
	err = redisRepo.load("test-key", &person)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(person, input) {
		t.Errorf("Didn't load data to map. Want %#v. Got %#v.", input, person)
	}
}

func TestLoadErrors(t *testing.T) {
	var n int
	var invalidMap map[string]int
	var tests = []struct {
		key    string
		output interface{}
		errMsg string
	}{
		{"dont-know", nil, "not found"},
		{"test-key", Person{}, "please provide a pointer for getting result from the database"},
		{"test-key", &n, "please provider a pointer to a struct or a map for getting result from the database"},
		{"test-key", &InvalidStruct{}, "can only expand values to structs or maps"},
		{"test-key", &invalidMap, "please provide a map[string]string"},
		{"test-key", &InvalidInnerStruct{}, "please provide a map[string]string"},
	}

	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	redisRepo := repo.(*redisRepository)
	client := redisRepo.redisClient()
	defer client.Close()
	err = redisRepo.save("test-key", map[string]string{"name": "Gopher"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Del("test-key")
	for _, test := range tests {
		err := redisRepo.load(test.key, test.output)
		if err == nil {
			t.Errorf("Got unexpected nil error, want %q", test.errMsg)
			continue
		}
		if err.Error() != test.errMsg {
			t.Errorf("Got wrong error message. Want %q. Got %q", test.errMsg, err.Error())
		}
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
	err = deleteKeys("preset:*", client)
	if err != nil {
		return err
	}
	return deleteKeys(presetsSetKey, client)
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