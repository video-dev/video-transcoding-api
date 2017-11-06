package redis

import (
	"errors"

	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
	"github.com/go-redis/redis"
)

const localPresetsSetKey = "localpresets"

func (r *redisRepository) CreateLocalPreset(localPreset *db.LocalPreset) error {
	if _, err := r.GetLocalPreset(localPreset.Name); err == nil {
		return db.ErrLocalPresetAlreadyExists
	}
	return r.saveLocalPreset(localPreset)
}

func (r *redisRepository) UpdateLocalPreset(localPreset *db.LocalPreset) error {
	if _, err := r.GetLocalPreset(localPreset.Name); err == db.ErrLocalPresetNotFound {
		return err
	}
	return r.saveLocalPreset(localPreset)
}

func (r *redisRepository) saveLocalPreset(localPreset *db.LocalPreset) error {
	fields, err := r.storage.FieldMap(localPreset)
	if err != nil {
		return err
	}
	if localPreset.Name == "" {
		return errors.New("preset name missing")
	}
	localPresetKey := r.localPresetKey(localPreset.Name)
	return r.storage.RedisClient().Watch(func(tx *redis.Tx) error {
		err := tx.HMSet(localPresetKey, fields).Err()
		if err != nil {
			return err
		}
		return tx.SAdd(localPresetsSetKey, localPreset.Name).Err()
	}, localPresetKey)
}

func (r *redisRepository) DeleteLocalPreset(localPreset *db.LocalPreset) error {
	err := r.storage.Delete(r.localPresetKey(localPreset.Name))
	if err != nil {
		if err == storage.ErrNotFound {
			return db.ErrLocalPresetNotFound
		}
		return err
	}
	r.storage.RedisClient().SRem(localPresetsSetKey, localPreset.Name)
	return nil
}

func (r *redisRepository) GetLocalPreset(name string) (*db.LocalPreset, error) {
	localPreset := db.LocalPreset{Name: name, Preset: db.Preset{}}
	err := r.storage.Load(r.localPresetKey(name), &localPreset)
	if err == storage.ErrNotFound {
		return nil, db.ErrLocalPresetNotFound
	}
	return &localPreset, err
}

func (r *redisRepository) localPresetKey(name string) string {
	return "localpreset:" + name
}
