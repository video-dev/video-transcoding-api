package redis

import (
	"github.com/NYTimes/video-transcoding-api/db"
)

const localPresetsSetKey = "localpresets"

func (r *redisRepository) CreateLocalPreset(preset *db.LocalPreset) error {
	return nil
}

func (r *redisRepository) UpdateLocalPreset(preset *db.LocalPreset) error {
	return nil
}

func (r *redisRepository) saveLocalPreset(preset *db.LocalPreset) error {
	return nil
}

func (r *redisRepository) DeleteLocalPreset(preset *db.LocalPreset) error {
	return nil
}

func (r *redisRepository) GetLocalPreset(name string) (*db.LocalPreset, error) {
	return &db.LocalPreset{}, nil
}

func (r *redisRepository) ListLocalPresets() ([]db.LocalPreset, error) {
	return []db.LocalPreset{}, nil
}

func (r *redisRepository) localPresetKey(name string) string {
	return ""
}
