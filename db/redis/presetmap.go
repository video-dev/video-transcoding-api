package redis

import (
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
	"github.com/go-redis/redis"
)

const presetmapsSetKey = "presetmaps"

func (r *redisRepository) CreatePresetMap(presetMap *db.PresetMap) error {
	if _, err := r.GetPresetMap(presetMap.Name); err == nil {
		return db.ErrPresetMapAlreadyExists
	}
	return r.savePresetMap(presetMap)
}

func (r *redisRepository) UpdatePresetMap(presetMap *db.PresetMap) error {
	if _, err := r.GetPresetMap(presetMap.Name); err == db.ErrPresetMapNotFound {
		return err
	}
	return r.savePresetMap(presetMap)
}

func (r *redisRepository) savePresetMap(presetMap *db.PresetMap) error {
	fields, err := r.storage.FieldMap(presetMap)
	if err != nil {
		return err
	}
	presetMapKey := r.presetMapKey(presetMap.Name)
	return r.storage.RedisClient().Watch(func(tx *redis.Tx) error {
		err := tx.HMSet(presetMapKey, fields).Err()
		if err != nil {
			return err
		}
		return tx.SAdd(presetmapsSetKey, presetMap.Name).Err()
	}, presetMapKey)
}

func (r *redisRepository) DeletePresetMap(presetMap *db.PresetMap) error {
	err := r.storage.Delete(r.presetMapKey(presetMap.Name))
	if err != nil {
		if err == storage.ErrNotFound {
			return db.ErrPresetMapNotFound
		}
		return err
	}
	r.storage.RedisClient().SRem(presetmapsSetKey, presetMap.Name)
	return nil
}

func (r *redisRepository) GetPresetMap(name string) (*db.PresetMap, error) {
	presetMap := db.PresetMap{Name: name, ProviderMapping: make(map[string]string)}
	err := r.storage.Load(r.presetMapKey(name), &presetMap)
	if err == storage.ErrNotFound {
		return nil, db.ErrPresetMapNotFound
	}
	return &presetMap, err
}

func (r *redisRepository) ListPresetMaps() ([]db.PresetMap, error) {
	presetMapNames, err := r.storage.RedisClient().SMembers(presetmapsSetKey).Result()
	if err != nil {
		return nil, err
	}
	presetsMap := make([]db.PresetMap, 0, len(presetMapNames))
	for _, name := range presetMapNames {
		presetMap, err := r.GetPresetMap(name)
		if err != nil && err != db.ErrPresetMapNotFound {
			return nil, err
		}
		if presetMap != nil {
			presetsMap = append(presetsMap, *presetMap)
		}
	}
	return presetsMap, nil
}

func (r *redisRepository) presetMapKey(name string) string {
	return "presetmap:" + name
}
