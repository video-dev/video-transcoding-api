package dbtest

import (
	"errors"
	"time"

	"github.com/nytm/video-transcoding-api/db"
)

type fakeRepository struct {
	triggerError bool
	presets      map[string]*db.PresetMap
	jobs         []*db.Job
}

// NewFakeRepository creates a new instance of the fake repository
// implementation. The underlying fake repository keeps jobs and presets in
// memory.
func NewFakeRepository(triggerError bool) db.Repository {
	return &fakeRepository{
		triggerError: triggerError,
		presets:      make(map[string]*db.PresetMap),
	}
}

func (d *fakeRepository) CreateJob(job *db.Job) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if job.ID == "" {
		job.ID = "12345"
	}
	if job.CreationTime.IsZero() {
		job.CreationTime = time.Now().In(time.UTC)
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

func (d *fakeRepository) CreatePresetMap(preset *db.PresetMap) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if preset.Name == "" {
		return errors.New("invalid preset name")
	}
	if _, ok := d.presets[preset.Name]; ok {
		return db.ErrPresetMapAlreadyExists
	}
	d.presets[preset.Name] = preset
	return nil
}

func (d *fakeRepository) UpdatePresetMap(preset *db.PresetMap) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if _, ok := d.presets[preset.Name]; !ok {
		return db.ErrPresetMapNotFound
	}
	d.presets[preset.Name] = preset
	return nil
}

func (d *fakeRepository) GetPresetMap(name string) (*db.PresetMap, error) {
	if d.triggerError {
		return nil, errors.New("database error")
	}
	if preset, ok := d.presets[name]; ok {
		return preset, nil
	}
	return nil, db.ErrPresetMapNotFound
}

func (d *fakeRepository) DeletePresetMap(preset *db.PresetMap) error {
	if d.triggerError {
		return errors.New("database error")
	}
	if _, ok := d.presets[preset.Name]; !ok {
		return db.ErrPresetMapNotFound
	}
	delete(d.presets, preset.Name)
	return nil
}

func (d *fakeRepository) ListPresetMaps() ([]db.PresetMap, error) {
	if d.triggerError {
		return nil, errors.New("database error")
	}
	presets := make([]db.PresetMap, 0, len(d.presets))
	for _, preset := range d.presets {
		presets = append(presets, *preset)
	}
	return presets, nil
}
