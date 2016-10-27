package redis

import (
	"github.com/NYTimes/video-transcoding-api/db"
	"gopkg.in/redis.v4"
)

const localPresetsSetKey = "localpresets"

func (r *redisRepository) CreateLocalPreset(localPreset *db.LocalPreset) error {
	if _, err := r.GetLocalPreset(localPreset.Name); err == nil {
		return db.ErrLocalPresetAlreadyExists
	}
	return r.saveLocalPreset(localPreset)
}

func (r *redisRepository) UpdateLocalPreset(preset *db.LocalPreset) error {
	return nil
}

func (r *redisRepository) saveLocalPreset(localPreset *db.LocalPreset) error {
	fields, err := r.storage.FieldMap(localPreset)
	if err != nil {
		return err
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

func (r *redisRepository) DeleteLocalPreset(preset *db.LocalPreset) error {
	return nil
}

func (r *redisRepository) GetLocalPreset(name string) (*db.LocalPreset, error) {
	return nil, db.ErrPresetMapNotFound
}

func (r *redisRepository) ListLocalPresets() ([]db.LocalPreset, error) {
	r.localPresetKey("nothing")
	return []db.LocalPreset{}, nil
}

func (r *redisRepository) localPresetKey(name string) string {
	return "localpreset:" + name
}
