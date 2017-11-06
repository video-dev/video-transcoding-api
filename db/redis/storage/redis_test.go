package storage

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
)

func TestRedisClientRedisDefaultOptions(t *testing.T) {
	storage, err := NewStorage(&Config{})
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
	defer client.Close()
	_, err = client.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisClientRedisAddr(t *testing.T) {
	port := "49159"
	proc, err := startRedis(port, "not-secret")
	if err != nil {
		t.Fatal(err)
	}
	defer proc.Signal(os.Interrupt)
	cfg := Config{
		RedisAddr: "127.0.0.1:" + port,
		Password:  "not-secret",
	}
	storage, err := NewStorage(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
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
	cfg := Config{
		SentinelAddrs:      "127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381",
		SentinelMasterName: "mymaster",
	}
	storage, err := NewStorage(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
	defer client.Close()
	_, err = client.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSave(t *testing.T) {
	person := Person{
		ID:        "some-id",
		Name:      "gopher",
		Age:       29,
		Weight:    150.03,
		BirthTime: time.Now().Add(-29 * 365 * 24 * time.Hour),
		Address: Address{
			Data:   map[string]string{"first_line": "secret"},
			Number: -2,
			City:   &City{Name: "nyc"},
		},
		NonTagged:        "not relevant",
		unexported:       "not relevant",
		unexportedTagged: "not relevant, believe!",
	}
	storage, err := NewStorage(&Config{})
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Save("person:test", person)
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
	defer client.Close()
	defer client.Del("person:test")
	data, err := client.HGetAll("person:test").Result()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{
		"name":                    "gopher",
		"age":                     "29",
		"weight":                  "150.03",
		"birth":                   person.BirthTime.Format(time.RFC3339Nano),
		"colors":                  "",
		"address_city_name":       "nyc",
		"address_data_first_line": "secret",
		"address_number":          "-2",
		"address_main":            "false",
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Did not save properly.\nWant %#v\nGot  %#v", expected, data)
	}
}

func TestSavePointer(t *testing.T) {
	person := Person{
		ID:              "some-id",
		Name:            "gopher",
		Age:             29,
		Weight:          153.2993,
		BirthTime:       time.Now().Add(-29 * 365 * 24 * time.Hour),
		PreferredColors: []string{"red", "blue", "yellow"},
		Address: Address{
			Data:   map[string]string{"first_line": "secret"},
			Number: -2,
			City:   &City{Name: "nyc"},
		},
		NonTagged:        "not relevant",
		unexported:       "not relevant",
		unexportedTagged: "not relevant, believe!",
	}
	storage, err := NewStorage(&Config{})
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Save("person:test", &person)
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
	defer client.Close()
	defer client.Del("person:test")
	data, err := client.HGetAll("person:test").Result()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{
		"name":                    "gopher",
		"age":                     "29",
		"weight":                  "153.2993",
		"birth":                   person.BirthTime.Format(time.RFC3339Nano),
		"colors":                  "red%%%blue%%%yellow",
		"address_city_name":       "nyc",
		"address_data_first_line": "secret",
		"address_number":          "-2",
		"address_main":            "false",
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Did not save properly.\nWant %#v\nGot  %#v", expected, data)
	}
}

func TestSaveMap(t *testing.T) {
	storage, err := NewStorage(&Config{})
	if err != nil {
		t.Fatal(err)
	}
	input := map[string]string{
		"name": "John",
		"test": "tested",
	}
	err = storage.Save("map:test", input)
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
	defer client.Close()
	defer client.Del("map:test")
	data, err := client.HGetAll("map:test").Result()
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

	storage, err := NewStorage(&Config{})
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
	defer client.Close()
	for _, test := range tests {
		err := storage.Save("some-key", test.input)
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

func TestFieldMap(t *testing.T) {
	storage, err := NewStorage(&Config{})
	if err != nil {
		t.Fatal(err)
	}
	var tests = []struct {
		description string
		hash        interface{}
		expected    map[string]interface{}
	}{
		{
			"Job",
			Job{
				ID:            "job1",
				ProviderJobID: "123abc",
				SourceMedia:   "http://nyt.net/source_here.mp4",
				ProviderName:  "encoding.com",
				StreamingParams: StreamingParams{
					SegmentDuration:  10,
					Protocol:         "hls",
					PlaylistFileName: "hls/playlist.m3u8",
				},
				Outputs: []TranscodeOutput{
					{Preset: PresetMap{Name: "preset-1"}, FileName: "output1.m3u8"},
					{Preset: PresetMap{Name: "preset-2"}, FileName: "output2.m3u8"},
				},
			},
			map[string]interface{}{
				"source":                           "http://nyt.net/source_here.mp4",
				"jobID":                            "job1",
				"providerName":                     "encoding.com",
				"providerJobID":                    "123abc",
				"streamingparams_segmentDuration":  "10",
				"streamingparams_protocol":         "hls",
				"streamingparams_playlistFileName": "hls/playlist.m3u8",
				"creationTime":                     "0001-01-01T00:00:00Z",
			},
		},
		{
			"LocalPreset",
			LocalPreset{
				Name: "this-is-a-localpreset",
				Preset: Preset{
					Name:        "test",
					Description: "test preset",
					Container:   "mp4",
					RateControl: "VBR",
					Video: VideoPreset{
						Profile:       "main",
						ProfileLevel:  "3.1",
						Width:         "640",
						Height:        "360",
						Codec:         "h264",
						Bitrate:       "1000",
						GopSize:       "90",
						GopMode:       "fixed",
						InterlaceMode: "progressive",
					},
					Audio: AudioPreset{
						Codec:   "aac",
						Bitrate: "64000",
					},
				},
			},
			map[string]interface{}{
				"preset_name":                "test",
				"preset_description":         "test preset",
				"preset_container":           "mp4",
				"preset_ratecontrol":         "VBR",
				"preset_video_profilelevel":  "3.1",
				"preset_video_profile":       "main",
				"preset_video_gopmode":       "fixed",
				"preset_video_bitrate":       "1000",
				"preset_video_interlacemode": "progressive",
				"preset_video_codec":         "h264",
				"preset_video_gopsize":       "90",
				"preset_video_height":        "360",
				"preset_video_width":         "640",
				"preset_audio_bitrate":       "64000",
				"preset_audio_codec":         "aac",
			},
		},
	}

	for _, test := range tests {
		result, err := storage.FieldMap(test.hash)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Wrong FieldMap: %s: Want %#v. Got %#v.", test.description, result, test.expected)
		}
	}
}

func TestLoadStruct(t *testing.T) {
	storage, err := NewStorage(&Config{})
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
	defer client.Close()
	date := time.Now().UTC().Add(-29 * 365 * 24 * time.Hour).Truncate(time.Millisecond)
	err = storage.Save("test-key", map[string]string{
		"name":              "Gopher",
		"age":               "29",
		"weight":            "159.332",
		"birth":             date.Format(time.RFC3339Nano),
		"colors":            "red%%%green%%%blue%%%black",
		"address_number":    "-2",
		"address_main":      "true",
		"address_city_name": "New York",
	})
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
	expectedPerson.Address.Main = true
	expectedPerson.Address.Number = -2
	expectedPerson.Name = "Gopher"
	expectedPerson.Age = 29
	expectedPerson.Weight = 159.332
	expectedPerson.BirthTime = date
	expectedPerson.PreferredColors = []string{"red", "green", "blue", "black"}
	err = storage.Load("test-key", &person)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(person, expectedPerson) {
		t.Errorf("Didn't load data to struct\nwant %#v\ngot  %#v", expectedPerson, person)
	}
}

func TestLoadMap(t *testing.T) {
	storage, err := NewStorage(&Config{})
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
	defer client.Close()
	input := map[string]string{"name": "Gopher", "city_name": "New York"}
	err = storage.Save("test-key", input)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Del("test-key")
	person := make(map[string]string)
	err = storage.Load("test-key", &person)
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
		{"dont-know", &Person{}, "not found"},
		{"test-key", Person{}, "please provide a pointer for getting result from the database"},
		{"test-key", &n, "please provider a pointer to a struct or a map for getting result from the database"},
		{"test-key", &InvalidStruct{}, "can only expand values to structs or maps"},
		{"test-key", &invalidMap, "please provide a map[string]string"},
		{"test-key", &InvalidInnerStruct{}, "please provide a map[string]string"},
	}

	storage, err := NewStorage(&Config{})
	if err != nil {
		t.Fatal(err)
	}
	client := storage.RedisClient()
	defer client.Close()
	err = storage.Save("test-key", map[string]string{"name": "Gopher"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Del("test-key")
	for _, test := range tests {
		err := storage.Load(test.key, test.output)
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
				}
				time.Sleep(10e6)
			}
		}(addr)
	}
	wg.Wait()
}
