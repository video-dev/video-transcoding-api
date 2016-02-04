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
	job.ID = "12345"
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
	if preset.ID == "" {
		preset.ID = "12345"
	}
	d.presets[preset.ID] = preset
	return nil
}

func (d *fakeDB) GetPreset(id string) (*db.Preset, error) {
	if preset, ok := d.presets[id]; ok {
		return preset, nil
	}
	return nil, db.ErrPresetNotFound
}

func (d *fakeDB) DeletePreset(preset *db.Preset) error {
	if _, ok := d.presets[preset.ID]; !ok {
		return db.ErrPresetNotFound
	}
	delete(d.presets, preset.ID)
	return nil
}
