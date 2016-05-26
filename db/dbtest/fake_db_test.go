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
	now := time.Now().UTC()
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

func TestCreatePresetMap(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.PresetMap{Name: "mypreset"}
	err := repo.CreatePresetMap(&preset)
	if err != nil {
		t.Fatal(err)
	}
	expectedPresetMaps := map[string]*db.PresetMap{"mypreset": &preset}
	presets := repo.(*fakeRepository).presets
	if !reflect.DeepEqual(presets, expectedPresetMaps) {
		t.Errorf("Wrong internal preset registry. Want %#v. Got %#v", expectedPresetMaps, presets)
	}
}

func TestCreatePresetMapEmptyName(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.PresetMap{}
	err := repo.CreatePresetMap(&preset)
	if err == nil {
		t.Fatal("got unexpected <nil> error")
	}
	expectedMsg := "invalid preset name"
	if err.Error() != expectedMsg {
		t.Errorf("CreatePresetMap: wrong error message. Want %q. Got %q", expectedMsg, err.Error())
	}
}

func TestCreatePresetMapDuplicate(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.PresetMap{Name: "mypreset"}
	err := repo.CreatePresetMap(&preset)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.CreatePresetMap(&preset)
	if err != db.ErrPresetMapAlreadyExists {
		t.Errorf("CreatePresetMap: wrong error returned. Want %#v. Got %#v", db.ErrPresetMapAlreadyExists, err)
	}
}

func TestCreatePresetMapDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	preset := db.PresetMap{}
	err := repo.CreatePresetMap(&preset)
	if err == nil {
		t.Fatal("got unexpected <nil> error")
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("CreatePresetMap: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestUpdatePresetMap(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.PresetMap{Name: "mypreset"}
	err := repo.CreatePresetMap(&preset)
	if err != nil {
		t.Fatal(err)
	}
	newPresetMap := preset
	newPresetMap.ProviderMapping = map[string]string{"some": "provider"}
	err = repo.UpdatePresetMap(&newPresetMap)
	if err != nil {
		t.Fatal(err)
	}
	expectedPresetMaps := map[string]*db.PresetMap{"mypreset": &newPresetMap}
	presets := repo.(*fakeRepository).presets
	if !reflect.DeepEqual(presets, expectedPresetMaps) {
		t.Errorf("Wrong internal preset registry. Want %#v. Got %#v", expectedPresetMaps, presets)
	}
}

func TestUpdatePresetMapNotFound(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.PresetMap{Name: "mypreset"}
	err := repo.UpdatePresetMap(&preset)
	if err != db.ErrPresetMapNotFound {
		t.Errorf("UpdatePresetMap: wrong error. Want %#v. Got %#v", db.ErrPresetMapNotFound, err)
	}
}

func TestUpdatePresetMapDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	preset := db.PresetMap{Name: "mypreset"}
	err := repo.UpdatePresetMap(&preset)
	if err == nil {
		t.Fatal("Unexpected <nil> error")
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("UpdatePresetMap: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestGetPresetMap(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.PresetMap{Name: "mypreset"}
	err := repo.CreatePresetMap(&preset)
	if err != nil {
		t.Fatal(err)
	}
	gotPresetMap, err := repo.GetPresetMap(preset.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*gotPresetMap, preset) {
		t.Errorf("GetPresetMap: wrong preset returned. Want %#v. Got %#v", preset, *gotPresetMap)
	}
}

func TestGetPresetMapNotFound(t *testing.T) {
	repo := NewFakeRepository(false)
	preset, err := repo.GetPresetMap("some-preset")
	if preset != nil {
		t.Errorf("GetPresetMap: unexpected non-nil preset: %#v", *preset)
	}
	if err != db.ErrPresetMapNotFound {
		t.Errorf("GetPresetMap: wrong error. Want ErrPresetMapNotFound. Got %#v", err)
	}
}

func TestGetPresetMapDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	preset, err := repo.GetPresetMap("some-preset")
	if preset != nil {
		t.Errorf("GetPresetMap: unexpected non-nil preset: %#v", *preset)
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("GetPresetMap: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestDeletePresetMap(t *testing.T) {
	repo := NewFakeRepository(false)
	preset1 := db.PresetMap{Name: "mypreset"}
	err := repo.CreatePresetMap(&preset1)
	if err != nil {
		t.Fatal(err)
	}
	preset2 := db.PresetMap{Name: "theirpreset"}
	err = repo.CreatePresetMap(&preset2)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeletePresetMap(&preset1)
	if err != nil {
		t.Fatal(err)
	}
	expectedPresetMaps := map[string]*db.PresetMap{"theirpreset": &preset2}
	presets := repo.(*fakeRepository).presets
	if !reflect.DeepEqual(presets, expectedPresetMaps) {
		t.Errorf("Wrong internal preset registry. Want %#v. Got %#v", expectedPresetMaps, presets)
	}
}

func TestDeletePresetMapNotFound(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.PresetMap{Name: "mypreset"}
	err := repo.DeletePresetMap(&preset)
	if err != db.ErrPresetMapNotFound {
		t.Errorf("DeletePresetMap: wrong error. Want %#v. Got %#v", db.ErrPresetMapNotFound, err)
	}
}

func TestDeletePresetMapDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	preset := db.PresetMap{Name: "mypreset"}
	err := repo.DeletePresetMap(&preset)
	if err == nil {
		t.Fatal("Unexpected <nil> error")
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("DeletePresetMap: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}

func TestListPresetMaps(t *testing.T) {
	repo := NewFakeRepository(false)
	preset := db.PresetMap{Name: "mypreset"}
	err := repo.CreatePresetMap(&preset)
	if err != nil {
		t.Fatal(err)
	}
	expectedPresetMaps := []db.PresetMap{preset}
	presets, err := repo.ListPresetMaps()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(presets, expectedPresetMaps) {
		t.Errorf("ListPresetMaps: wrong list returned. Want %#v. Got %#v", expectedPresetMaps, presets)
	}
}

func TestListPresetMapsDBError(t *testing.T) {
	repo := NewFakeRepository(true)
	presets, err := repo.ListPresetMaps()
	if len(presets) > 0 {
		t.Errorf("ListPresetMaps: got unexpected non-empty list: %#v", presets)
	}
	if err.Error() != dbErrorMsg {
		t.Errorf("ListPresetMaps: wrong error message. Want %q. Got %q", dbErrorMsg, err.Error())
	}
}
