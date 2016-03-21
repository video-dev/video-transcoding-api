package dbtest

import (
	"reflect"
	"testing"
	"time"

	"github.com/nytm/video-transcoding-api/db"
)

const dbErrorMsg = "database error"

func TestCreateJob(t *testing.T) {
	repo := NewFakeRepository(false)
	job := db.Job{ID: "j-123", ProviderName: "myprovider"}
	err := repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	if job.CreationTime.IsZero() {
		t.Error("Did not set the CreationTime")
	}
	if job.CreationTime.Location() != time.UTC {
		t.Errorf("Did not set CreationTime to UTC: %#v", job.CreationTime.Location())
	}
}

func TestCreateJobPredefinedDate(t *testing.T) {
	repo := NewFakeRepository(false)
	creationTime := time.Date(1983, 2, 19, 20, 15, 53, 0, time.UTC)
	job := db.Job{ID: "j-123", ProviderName: "myprovider", CreationTime: creationTime}
	err := repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	if job.CreationTime != creationTime {
		t.Errorf("Wrong CreationTime. Want %s. Got %s", creationTime, job.CreationTime)
	}
	if job.CreationTime.Location() != time.UTC {
		t.Errorf("Did not set CreationTime to UTC: %#v", job.CreationTime.Location())
	}
}

func TestCreateJobNoID(t *testing.T) {
	repo := NewFakeRepository(false)
	job := db.Job{ProviderName: "myprovider"}
	err := repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	if job.CreationTime.IsZero() {
		t.Error("Did not set the CreationTime")
	}
	if job.CreationTime.Location() != time.UTC {
		t.Errorf("Did not set CreationTime to UTC: %#v", job.CreationTime.Location())
	}
	if job.ID != "12345" {
		t.Errorf("Did not generate an ID. Want 12345. Got %q", job.ID)
	}
}

func TestCreateJobDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	job := db.Job{ProviderName: "myprovider"}
	err := repo.CreateJob(&job)
	if err == nil {
		t.Fatal("Got unexpected <nil> error")
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("Got wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestGetJob(t *testing.T) {
	repo := NewFakeRepository(false)
	job := db.Job{ID: "j-123", ProviderName: "myprovider"}
	err := repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	retrievedJob, err := repo.GetJob(job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*retrievedJob, job) {
		t.Errorf("Wrong job returned. Want %#v. Got %#v", job, *retrievedJob)
	}
}

func TestGetJobNotFound(t *testing.T) {
	repo := NewFakeRepository(false)
	job, err := repo.GetJob("some-job")
	if job != nil {
		t.Errorf("Got unexpected non-nil job: %#v", job)
	}
	if err != db.ErrJobNotFound {
		t.Errorf("Wrong error returned. Want %#v. Got %#v", db.ErrJobNotFound, err)
	}
}

func TestGetJobDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	job, err := repo.GetJob("some-job")
	if job != nil {
		t.Errorf("Got unexpected non-nil job: %#v", job)
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("Wrong error message returned. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestDeleteJob(t *testing.T) {
	repo := NewFakeRepository(false)
	job := db.Job{ID: "j-123", ProviderName: "myprovider"}
	err := repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.CreateJob(&db.Job{ID: "j-124", ProviderName: "myprovider"})
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeleteJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.GetJob(job.ID)
	if err != db.ErrJobNotFound {
		t.Errorf("Got wrong error. Want db.ErrJobNotFound. Got %#v", err)
	}
	_, err = repo.GetJob("j-124")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteJobNotFound(t *testing.T) {
	repo := NewFakeRepository(false)
	err := repo.DeleteJob(&db.Job{ID: "some-job"})
	if err != db.ErrJobNotFound {
		t.Errorf("Wrong error returned. Want %#v. Got %#v", db.ErrJobNotFound, err)
	}
}

func TestDeleteJobDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	err := repo.DeleteJob(&db.Job{ID: "some-job"})
	if err.Error() != dbErrorMsg {
		t.Errorf("Wrong error message returned. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestListJobs(t *testing.T) {
	jobs := []db.Job{
		{ID: "job-1", ProviderName: "encodingcom"},
		{ID: "job-2", ProviderName: "encodingcom"},
		{ID: "job-3", ProviderName: "encodingcom"},
	}
	repo := NewFakeRepository(false)
	for i, job := range jobs {
		job := job
		err := repo.CreateJob(&job)
		if err != nil {
			t.Fatal(err)
		}
		jobs[i] = job
	}
	gotJobs, err := repo.ListJobs(db.JobFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotJobs, jobs) {
		t.Errorf("ListJobs: wrong list returned. Want %#v. Got %#v.", jobs, gotJobs)
	}
}

func TestListJobsFilter(t *testing.T) {
	now := time.Now().In(time.UTC)
	jobs := []db.Job{
		{ID: "job-1", ProviderName: "encodingcom", CreationTime: now.Add(-2 * time.Hour)},
		{ID: "job-2", ProviderName: "encodingcom", CreationTime: now.Add(-1 * time.Hour)},
		{ID: "job-3", ProviderName: "encodingcom", CreationTime: now.Add(-30 * time.Minute)},
	}
	repo := NewFakeRepository(false)
	for i, job := range jobs {
		job := job
		err := repo.CreateJob(&job)
		if err != nil {
			t.Fatal(err)
		}
		jobs[i] = job
	}
	gotJobs, err := repo.ListJobs(db.JobFilter{Since: now.Add(-90 * time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotJobs, jobs[1:]) {
		t.Errorf("ListJobs: wrong list returned. Want %#v. Got %#v", jobs[1:], gotJobs)
	}
}

func TestListJobsLimit(t *testing.T) {
	jobs := []db.Job{
		{ID: "job-1", ProviderName: "encodingcom"},
		{ID: "job-2", ProviderName: "encodingcom"},
		{ID: "job-3", ProviderName: "encodingcom"},
	}
	repo := NewFakeRepository(false)
	for i, job := range jobs {
		job := job
		err := repo.CreateJob(&job)
		if err != nil {
			t.Fatal(err)
		}
		jobs[i] = job
	}
	gotJobs, err := repo.ListJobs(db.JobFilter{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotJobs, jobs[:2]) {
		t.Errorf("ListJobs: wrong list returned. Want %#v. Got %#v", jobs, gotJobs)
	}
}

func TestListJobsDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	jobs, err := repo.ListJobs(db.JobFilter{})
	if len(jobs) > 0 {
		t.Errorf("Got unexpected non-empty job list: %v", jobs)
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("Wrong error message returned. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestCreatePreset(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.Preset{Name: "mypreset"}
	err := repo.CreatePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	expectedPresets := map[string]*db.Preset{"mypreset": &preset}
	presets := repo.(*fakeRepository).presets
	if !reflect.DeepEqual(presets, expectedPresets) {
		t.Errorf("Wrong internal preset registry. Want %#v. Got %#v", expectedPresets, presets)
	}
}

func TestCreatePresetEmptyName(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.Preset{}
	err := repo.CreatePreset(&preset)
	if err == nil {
		t.Fatal("got unexpected <nil> error")
	}
	expectedMsg := "invalid preset name"
	if err.Error() != expectedMsg {
		t.Errorf("CreatePreset: wrong error message. Want %q. Got %q", expectedMsg, err.Error())
	}
}

func TestCreatePresetDuplicate(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.Preset{Name: "mypreset"}
	err := repo.CreatePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.CreatePreset(&preset)
	if err != db.ErrPresetAlreadyExists {
		t.Errorf("CreatePreset: wrong error returned. Want %#v. Got %#v", db.ErrPresetAlreadyExists, err)
	}
}

func TestCreatePresetDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	preset := db.Preset{}
	err := repo.CreatePreset(&preset)
	if err == nil {
		t.Fatal("got unexpected <nil> error")
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("CreatePreset: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestUpdatePreset(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.Preset{Name: "mypreset"}
	err := repo.CreatePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	newPreset := preset
	newPreset.ProviderMapping = map[string]string{"some": "provider"}
	err = repo.UpdatePreset(&newPreset)
	if err != nil {
		t.Fatal(err)
	}
	expectedPresets := map[string]*db.Preset{"mypreset": &newPreset}
	presets := repo.(*fakeRepository).presets
	if !reflect.DeepEqual(presets, expectedPresets) {
		t.Errorf("Wrong internal preset registry. Want %#v. Got %#v", expectedPresets, presets)
	}
}

func TestUpdatePresetNotFound(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.Preset{Name: "mypreset"}
	err := repo.UpdatePreset(&preset)
	if err != db.ErrPresetNotFound {
		t.Errorf("UpdatePreset: wrong error. Want %#v. Got %#v", db.ErrPresetNotFound, err)
	}
}

func TestUpdatePresetDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	preset := db.Preset{Name: "mypreset"}
	err := repo.UpdatePreset(&preset)
	if err == nil {
		t.Fatal("Unexpected <nil> error")
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("UpdatePreset: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestGetPreset(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.Preset{Name: "mypreset"}
	err := repo.CreatePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	gotPreset, err := repo.GetPreset(preset.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*gotPreset, preset) {
		t.Errorf("GetPreset: wrong preset returned. Want %#v. Got %#v", preset, *gotPreset)
	}
}

func TestGetPresetNotFound(t *testing.T) {
	repo := NewFakeRepository(false)
	preset, err := repo.GetPreset("some-preset")
	if preset != nil {
		t.Errorf("GetPreset: unexpected non-nil preset: %#v", *preset)
	}
	if err != db.ErrPresetNotFound {
		t.Errorf("GetPreset: wrong error. Want ErrPresetNotFound. Got %#v", err)
	}
}

func TestGetPresetDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	preset, err := repo.GetPreset("some-preset")
	if preset != nil {
		t.Errorf("GetPreset: unexpected non-nil preset: %#v", *preset)
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("GetPreset: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestDeletePreset(t *testing.T) {
	repo := NewFakeRepository(false)
	preset1 := db.Preset{Name: "mypreset"}
	err := repo.CreatePreset(&preset1)
	if err != nil {
		t.Fatal(err)
	}
	preset2 := db.Preset{Name: "theirpreset"}
	err = repo.CreatePreset(&preset2)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeletePreset(&preset1)
	if err != nil {
		t.Fatal(err)
	}
	expectedPresets := map[string]*db.Preset{"theirpreset": &preset2}
	presets := repo.(*fakeRepository).presets
	if !reflect.DeepEqual(presets, expectedPresets) {
		t.Errorf("Wrong internal preset registry. Want %#v. Got %#v", expectedPresets, presets)
	}
}

func TestDeletePresetNotFound(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.Preset{Name: "mypreset"}
	err := repo.DeletePreset(&preset)
	if err != db.ErrPresetNotFound {
		t.Errorf("DeletePreset: wrong error. Want %#v. Got %#v", db.ErrPresetNotFound, err)
	}
}

func TestDeletePresetDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	preset := db.Preset{Name: "mypreset"}
	err := repo.DeletePreset(&preset)
	if err == nil {
		t.Fatal("Unexpected <nil> error")
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("DeletePreset: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestListPresets(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.Preset{Name: "mypreset"}
	err := repo.CreatePreset(&preset)
	if err != nil {
		t.Fatal(err)
	}
	expectedPresets := []db.Preset{preset}
	presets, err := repo.ListPresets()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(presets, expectedPresets) {
		t.Errorf("ListPresets: wrong list returned. Want %#v. Got %#v", expectedPresets, presets)
	}
}

func TestListPresetsDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	presets, err := repo.ListPresets()
	if len(presets) > 0 {
		t.Errorf("ListPresets: got unexpected non-empty list: %#v", presets)
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("ListPresets: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}
