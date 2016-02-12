package redis

import "github.com/nytm/video-transcoding-api/db"

const presetsSetKey = "presets"

func (r *redisRepository) CreatePreset(preset *db.Preset) error {
	if _, err := r.GetPreset(preset.Name); err == nil {
		return db.ErrPresetAlreadyExists
	}
	return r.savePreset(preset)
}

func (r *redisRepository) UpdatePreset(preset *db.Preset) error {
	if _, err := r.GetPreset(preset.Name); err == db.ErrPresetNotFound {
		return err
	}
	return r.savePreset(preset)
}

func (r *redisRepository) savePreset(preset *db.Preset) error {
	fields, err := r.fieldList(preset)
	if err != nil {
		return err
	}
	presetKey := r.presetKey(preset.Name)
	multi, err := r.redisClient().Watch(presetKey)
	_, err = multi.Exec(func() error {
		multi.HMSet(presetKey, fields[0], fields[1], fields[2:]...)
		multi.SAdd(presetsSetKey, preset.Name)
		return nil
	})
	return err
}

func (r *redisRepository) DeletePreset(preset *db.Preset) error {
	err := r.delete(r.presetKey(preset.Name), db.ErrPresetNotFound)
	if err != nil {
		return err
	}
	r.redisClient().SRem(presetsSetKey, preset.Name)
	return nil
}

func (r *redisRepository) GetPreset(name string) (*db.Preset, error) {
	preset := db.Preset{Name: name, ProviderMapping: make(map[string]string)}
	err := r.load(r.presetKey(name), &preset)
	if err == errNotFound {
		return nil, db.ErrPresetNotFound
	}
	return &preset, err
}

func (r *redisRepository) ListPresets() ([]db.Preset, error) {
	presetNames, err := r.redisClient().SMembers(presetsSetKey).Result()
	if err != nil {
		return nil, err
	}
	presets := make([]db.Preset, 0, len(presetNames))
	for _, name := range presetNames {
		preset, err := r.GetPreset(name)
		if err != nil && err != db.ErrPresetNotFound {
			return nil, err
		}
		if preset != nil {
			presets = append(presets, *preset)
		}
	}
	return presets, nil
}

func (r *redisRepository) presetKey(name string) string {
	return "preset:" + name
}
