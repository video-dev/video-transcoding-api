package service

import (
	"errors"

	"github.com/nytm/video-transcoding-api/db"
)

type fakeDB struct {
	triggerDBError bool
	presets        map[string]*db.Preset
	jobs           map[string]*db.Job
}

func newFakeDB(triggerDBError bool) db.Repository {
	return &fakeDB{
		triggerDBError: triggerDBError,
		presets:        make(map[string]*db.Preset),
		jobs:           make(map[string]*db.Job),
	}
}

func (d *fakeDB) SaveJob(job *db.Job) error {
	if d.triggerDBError {
		return errors.New("database error")
	}
	if job.ID == "" {
		job.ID = "12345"
	}
	d.jobs[job.ID] = job
	return nil
}

func (d *fakeDB) DeleteJob(job *db.Job) error {
	if _, ok := d.jobs[job.ID]; !ok {
		return db.ErrJobNotFound
	}
	delete(d.jobs, job.ID)
	return nil
}

func (d *fakeDB) GetJob(id string) (*db.Job, error) {
	if job, ok := d.jobs[id]; ok {
		return job, nil
	}
	return nil, db.ErrJobNotFound
}

func (d *fakeDB) SavePreset(preset *db.Preset) error {
	if d.triggerDBError {
		return errors.New("database error")
	}
	if _, ok := d.presets[preset.Name]; ok {
		return db.ErrPresetAlreadyExists
	}
	d.presets[preset.Name] = preset
	return nil
}

func (d *fakeDB) GetPreset(name string) (*db.Preset, error) {
	if preset, ok := d.presets[name]; ok {
		return preset, nil
	}
	return nil, db.ErrPresetNotFound
}

func (d *fakeDB) DeletePreset(preset *db.Preset) error {
	if _, ok := d.presets[preset.Name]; !ok {
		return db.ErrPresetNotFound
	}
	delete(d.presets, preset.Name)
	return nil
}

func (d *fakeDB) ListPresets() ([]db.Preset, error) {
	presets := make([]db.Preset, 0, len(d.presets))
	for _, preset := range d.presets {
		presets = append(presets, *preset)
	}
	return presets, nil
}
