package redis

import (
	"github.com/nytm/video-transcoding-api/db"
	"gopkg.in/redis.v4"
)

const presetsSetKey = "presets"

func (r *redisRepository) CreatePresetMap(preset *db.PresetMap) error {
	if _, err := r.GetPresetMap(preset.Name); err == nil {
		return db.ErrPresetMapAlreadyExists
	}
	return r.savePresetMap(preset)
}

func (r *redisRepository) UpdatePresetMap(preset *db.PresetMap) error {
	if _, err := r.GetPresetMap(preset.Name); err == db.ErrPresetMapNotFound {
		return err
	}
	return r.savePresetMap(preset)
}

func (r *redisRepository) savePresetMap(preset *db.PresetMap) error {
	fields, err := r.fieldList(preset)
	if err != nil {
		return err
	}
	presetKey := r.presetKey(preset.Name)
	return r.redisClient().Watch(func(tx *redis.Tx) error {
		err := tx.HMSet(presetKey, fields).Err()
		if err != nil {
			return err
		}
		return tx.SAdd(presetsSetKey, preset.Name).Err()
	}, presetKey)
}

func (r *redisRepository) DeletePresetMap(preset *db.PresetMap) error {
	err := r.delete(r.presetKey(preset.Name), db.ErrPresetMapNotFound)
	if err != nil {
		return err
	}
	r.redisClient().SRem(presetsSetKey, preset.Name)
	return nil
}

func (r *redisRepository) GetPresetMap(name string) (*db.PresetMap, error) {
	preset := db.PresetMap{Name: name, ProviderMapping: make(map[string]string)}
	err := r.load(r.presetKey(name), &preset)
	if err == errNotFound {
		return nil, db.ErrPresetMapNotFound
	}
	return &preset, err
}

func (r *redisRepository) ListPresetMaps() ([]db.PresetMap, error) {
	presetNames, err := r.redisClient().SMembers(presetsSetKey).Result()
	if err != nil {
		return nil, err
	}
	presetsMap := make([]db.PresetMap, 0, len(presetNames))
	for _, name := range presetNames {
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

func (r *redisRepository) presetKey(name string) string {
	return "preset:" + name
}
