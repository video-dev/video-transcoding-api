package redis

import "github.com/nytm/video-transcoding-api/db"

func (r *redisRepository) SavePreset(preset *db.Preset) error {
	if preset.ID == "" {
		id, err := r.generateID()
		if err != nil {
			return err
		}
		preset.ID = id
	}
	return r.save(r.presetKey(preset.ID), preset)
}

func (r *redisRepository) DeletePreset(preset *db.Preset) error {
	return r.delete(r.presetKey(preset.ID), db.ErrPresetNotFound)
}

func (r *redisRepository) GetPreset(id string) (*db.Preset, error) {
	preset := db.Preset{ID: id, ProviderMapping: make(map[string]string)}
	err := r.load(r.presetKey(id), &preset)
	if err == errNotFound {
		return nil, db.ErrPresetNotFound
	}
	return &preset, err
}

func (r *redisRepository) presetKey(id string) string {
	return "preset:" + id
}
