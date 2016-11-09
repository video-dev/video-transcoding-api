package redis

import (
	"reflect"
	"testing"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
)

func TestCreatePresetMap(t *testing.T) {
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
	presetmap := db.PresetMap{
		Name: "mypreset",
		ProviderMapping: map[string]string{
			"elementalconductor": "abc123",
			"elastictranscoder":  "1281742-93939",
		},
		OutputOpts: db.OutputOptions{Extension: "ts"},
	}
	err = repo.CreatePresetMap(&presetmap)
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).storage.RedisClient()
	defer client.Close()
	items, err := client.HGetAll("presetmap:" + presetmap.Name).Result()
	if err != nil {
		t.Fatal(err)
	}
	expectedItems := map[string]string{
		"pmapping_elementalconductor": "abc123",
		"pmapping_elastictranscoder":  "1281742-93939",
		"output_extension":            "ts",
		"presetmap_name":              "mypreset",
	}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Wrong presetmap hash returned from Redis. Want %#v. Got %#v", expectedItems, items)
	}
}

func TestCreatePresetMapDuplicate(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	presetmap := db.PresetMap{
		Name:            "mypreset",
		ProviderMapping: map[string]string{"elemental": "123"},
	}
	err = repo.CreatePresetMap(&presetmap)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.CreatePresetMap(&presetmap)
	if err != db.ErrPresetMapAlreadyExists {
		t.Errorf("Got wrong error. Want %#v. Got %#v", db.ErrPresetMapAlreadyExists, err)
	}
}

func TestUpdatePresetMap(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	presetmap := db.PresetMap{Name: "mypresetmap", ProviderMapping: map[string]string{"elemental": "abc123"}}
	err = repo.CreatePresetMap(&presetmap)
	if err != nil {
		t.Fatal(err)
	}
	presetmap.ProviderMapping = map[string]string{
		"elemental":         "abc1234",
		"elastictranscoder": "def123",
	}
	presetmap.OutputOpts = db.OutputOptions{Extension: "mp4"}
	err = repo.UpdatePresetMap(&presetmap)
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).storage.RedisClient()
	defer client.Close()
	items, err := client.HGetAll("presetmap:" + presetmap.Name).Result()
	if err != nil {
		t.Fatal(err)
	}
	expectedItems := map[string]string{
		"pmapping_elemental":         "abc1234",
		"pmapping_elastictranscoder": "def123",
		"output_extension":           "mp4",
		"presetmap_name":             "mypresetmap",
	}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Wrong presetmap hash returned from Redis. Want %#v. Got %#v", expectedItems, items)
	}
}

func TestUpdatePresetMapNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	err = repo.UpdatePresetMap(&db.PresetMap{Name: "mypresetmap"})
	if err != db.ErrPresetMapNotFound {
		t.Errorf("Wrong error returned by UpdatePresetMap. Want ErrPresetMapNotFound. Got %#v.", err)
	}
}

func TestDeletePresetMap(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	presetmap := db.PresetMap{Name: "mypresetmap", ProviderMapping: map[string]string{"elemental": "abc123"}}
	err = repo.CreatePresetMap(&presetmap)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeletePresetMap(&db.PresetMap{Name: presetmap.Name})
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).storage.RedisClient()
	result := client.HGetAll("presetmap:mypresetmap")
	if len(result.Val()) != 0 {
		t.Errorf("Unexpected value after delete call: %v", result.Val())
	}
}

func TestDeletePresetMapNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeletePresetMap(&db.PresetMap{Name: "mypresetmap"})
	if err != db.ErrPresetMapNotFound {
		t.Errorf("Wrong error returned by DeletePresetMap. Want ErrPresetMapNotFound. Got %#v.", err)
	}
}

func TestGetPresetMap(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	presetmap := db.PresetMap{
		Name: "mypresetmap",
		ProviderMapping: map[string]string{
			"elementalconductor": "abc-123",
			"elastictranscoder":  "0129291-0001",
			"encoding.com":       "wait what?",
		},
		OutputOpts: db.OutputOptions{Extension: "ts"},
	}
	err = repo.CreatePresetMap(&presetmap)
	if err != nil {
		t.Fatal(err)
	}
	gotPresetMap, err := repo.GetPresetMap(presetmap.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*gotPresetMap, presetmap) {
		t.Errorf("Wrong preset. Want %#v. Got %#v.", presetmap, *gotPresetMap)
	}
}

func TestGetPresetMapNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	gotPresetMap, err := repo.GetPresetMap("mypresetmap")
	if err != db.ErrPresetMapNotFound {
		t.Errorf("Wrong error returned. Want ErrPresetMapNotFound. Got %#v.", err)
	}
	if gotPresetMap != nil {
		t.Errorf("Unexpected non-nil presetmap: %#v.", gotPresetMap)
	}
}

func TestListPresetMaps(t *testing.T) {
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
	presetmaps := []db.PresetMap{
		{
			Name: "presetmap-1",
			ProviderMapping: map[string]string{
				"elementalconductor": "abc123",
				"elastictranscoder":  "1281742-93939",
			},
			OutputOpts: db.OutputOptions{Extension: "mp4"},
		},
		{
			Name: "presetmap-2",
			ProviderMapping: map[string]string{
				"elementalconductor": "abc124",
				"elastictranscoder":  "1281743-93939",
			},
			OutputOpts: db.OutputOptions{Extension: "webm"},
		},
		{
			Name: "presetmap-3",
			ProviderMapping: map[string]string{
				"elementalconductor": "abc125",
				"elastictranscoder":  "1281744-93939",
			},
			OutputOpts: db.OutputOptions{Extension: "ts"},
		},
	}
	for i := range presetmaps {
		err = repo.CreatePresetMap(&presetmaps[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	gotPresetMaps, err := repo.ListPresetMaps()
	if err != nil {
		t.Fatal(err)
	}

	// Why? The "list" of IDs is a set on Redis, so we need to make sure
	// that order is not important before invoking reflect.DeepEqual.
	expected := presetListToMap(presetmaps)
	got := presetListToMap(gotPresetMaps)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ListPresetMaps(): wrong list. Want %#v. Got %#v.", presetmaps, gotPresetMaps)
	}
}

func presetListToMap(presetmaps []db.PresetMap) map[string]db.PresetMap {
	result := make(map[string]db.PresetMap, len(presetmaps))
	for _, presetmap := range presetmaps {
		result[presetmap.Name] = presetmap
	}
	return result
}
