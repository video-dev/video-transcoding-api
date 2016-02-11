package redis

import "github.com/nytm/video-transcoding-api/db"

const presetsSetKey = "presets"

func (r *redisRepository) SavePreset(preset *db.Preset) error {
	if preset.ID == "" {
		id, err := r.generateID()
		if err != nil {
			return err
		}
		preset.ID = id
	} else if _, err := r.GetPreset(preset.ID); err == nil {
		return db.ErrPresetAlreadyExists
	}
	fields, err := r.fieldList(preset)
	if err != nil {
		return err
	}
	presetKey := r.presetKey(preset.ID)
	multi, err := r.redisClient().Watch(presetKey)
	_, err = multi.Exec(func() error {
		multi.HMSet(presetKey, fields[0], fields[1], fields[2:]...)
		multi.SAdd(presetsSetKey, preset.ID)
		return nil
	})
	return err
}

func (r *redisRepository) DeletePreset(preset *db.Preset) error {
	err := r.delete(r.presetKey(preset.ID), db.ErrPresetNotFound)
	if err != nil {
		return err
	}
	r.redisClient().SRem(presetsSetKey, preset.ID)
	return nil
}

func (r *redisRepository) GetPreset(id string) (*db.Preset, error) {
	preset := db.Preset{ID: id, ProviderMapping: make(map[string]string)}
	err := r.load(r.presetKey(id), &preset)
	if err == errNotFound {
		return nil, db.ErrPresetNotFound
	}
	return &preset, err
}

func (r *redisRepository) ListPresets() ([]db.Preset, error) {
	presetIDs, err := r.redisClient().SMembers(presetsSetKey).Result()
	if err != nil {
		return nil, err
	}
	presets := make([]db.Preset, 0, len(presetIDs))
	for _, id := range presetIDs {
		preset, err := r.GetPreset(id)
		if err != nil && err != db.ErrPresetNotFound {
			return nil, err
		}
		if preset != nil {
			presets = append(presets, *preset)
		}
	}
	return presets, nil
}

func (r *redisRepository) presetKey(id string) string {
	return "preset:" + id
}
