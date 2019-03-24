package dbtest // import "github.com/nytimes/video-transcoding-api/db/dbtest"

import (
	"errors"
	"time"

	"github.com/nytimes/video-transcoding-api/db"
)

type fakeRepository struct {
	triggerError bool
	presetmaps   map[string]*db.PresetMap
	localpresets map[string]*db.LocalPreset
	jobs         []*db.Job
}

// NewFakeRepository creates a new instance of the fake repository
// implementation. The underlying fake repository keeps jobs and presets in
// memory.
func NewFakeRepository(triggerError bool) db.Repository {
	return &fakeRepository{
		triggerError: triggerError,
		presetmaps:   make(map[string]*db.PresetMap),
		localpresets: make(map[string]*db.LocalPreset),
	}
}

func (d *fakeRepository) CreateJob(job *db.Job) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if job.CreationTime.IsZero() {
		job.CreationTime = time.Now().UTC()
	}
	d.jobs = append(d.jobs, job)
	return nil
}

func (d *fakeRepository) DeleteJob(job *db.Job) error {
	if d.triggerError {
		return errors.New("database error")
	}
	index, err := d.findJob(job.ID)
	if err != nil {
		return err
	}
	for i := index; i < len(d.jobs)-1; i++ {
		d.jobs[i] = d.jobs[i+1]
	}
	d.jobs = d.jobs[:len(d.jobs)-1]
	return nil
}

func (d *fakeRepository) GetJob(id string) (*db.Job, error) {
	if d.triggerError {
		return nil, errors.New("database error")
	}
	index, err := d.findJob(id)
	if err != nil {
		return nil, err
	}
	return d.jobs[index], nil
}

func (d *fakeRepository) findJob(id string) (int, error) {
	index := -1
	for i, job := range d.jobs {
		if job.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return index, db.ErrJobNotFound
	}
	return index, nil
}

func (d *fakeRepository) ListJobs(filter db.JobFilter) ([]db.Job, error) {
	if d.triggerError {
		return nil, errors.New("database error")
	}
	jobs := make([]db.Job, 0, len(d.jobs))
	var count uint
	for _, job := range d.jobs {
		if job.CreationTime.Before(filter.Since) {
			continue
		}
		if filter.Limit != 0 && count == filter.Limit {
			break
		}
		jobs = append(jobs, *job)
		count++
	}
	return jobs, nil
}

func (d *fakeRepository) CreatePresetMap(presetmap *db.PresetMap) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if presetmap.Name == "" {
		return errors.New("invalid presetmap name")
	}
	if _, ok := d.presetmaps[presetmap.Name]; ok {
		return db.ErrPresetMapAlreadyExists
	}
	d.presetmaps[presetmap.Name] = presetmap
	return nil
}

func (d *fakeRepository) UpdatePresetMap(presetmap *db.PresetMap) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if _, ok := d.presetmaps[presetmap.Name]; !ok {
		return db.ErrPresetMapNotFound
	}
	d.presetmaps[presetmap.Name] = presetmap
	return nil
}

func (d *fakeRepository) GetPresetMap(name string) (*db.PresetMap, error) {
	if d.triggerError {
		return nil, errors.New("database error")
	}
	if presetmap, ok := d.presetmaps[name]; ok {
		return presetmap, nil
	}
	return nil, db.ErrPresetMapNotFound
}

func (d *fakeRepository) DeletePresetMap(presetmap *db.PresetMap) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if _, ok := d.presetmaps[presetmap.Name]; !ok {
		return db.ErrPresetMapNotFound
	}
	delete(d.presetmaps, presetmap.Name)
	return nil
}

func (d *fakeRepository) ListPresetMaps() ([]db.PresetMap, error) {
	if d.triggerError {
		return nil, errors.New("database error")
	}
	presetmaps := make([]db.PresetMap, 0, len(d.presetmaps))
	for _, presetmap := range d.presetmaps {
		presetmaps = append(presetmaps, *presetmap)
	}
	return presetmaps, nil
}

func (d *fakeRepository) CreateLocalPreset(preset *db.LocalPreset) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if preset.Name == "" {
		return errors.New("invalid local preset name")
	}
	if _, ok := d.localpresets[preset.Name]; ok {
		return db.ErrLocalPresetAlreadyExists
	}
	d.localpresets[preset.Name] = preset
	return nil
}

func (d *fakeRepository) UpdateLocalPreset(preset *db.LocalPreset) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if _, ok := d.localpresets[preset.Name]; !ok {
		return db.ErrLocalPresetNotFound
	}
	d.localpresets[preset.Name] = preset

	return nil
}

func (d *fakeRepository) GetLocalPreset(name string) (*db.LocalPreset, error) {
	if d.triggerError {
		return nil, errors.New("database error")
	}
	if localpreset, ok := d.localpresets[name]; ok {
		return localpreset, nil
	}
	return nil, db.ErrLocalPresetNotFound
}

func (d *fakeRepository) DeleteLocalPreset(preset *db.LocalPreset) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if _, ok := d.localpresets[preset.Name]; !ok {
		return db.ErrLocalPresetNotFound
	}
	delete(d.localpresets, preset.Name)
	return nil
}
