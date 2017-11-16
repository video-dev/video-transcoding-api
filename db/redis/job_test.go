package redis

import (
	"math"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
	"github.com/go-redis/redis"
	"github.com/kr/pretty"
)

func TestCreateJob(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(storage.Config)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	job := db.Job{
		ID:              "job1",
		SourceMedia:     "http://nyt.net/source_here.mp4",
		ProviderName:    "encoding.com",
		StreamingParams: db.StreamingParams{SegmentDuration: 10, Protocol: "hls", PlaylistFileName: "hls/playlist.m3u8"},
		Outputs: []db.TranscodeOutput{
			{Preset: db.PresetMap{Name: "preset-1"}, FileName: "output1.m3u8"},
			{Preset: db.PresetMap{Name: "preset-2"}, FileName: "output2.m3u8"},
		},
	}
	err = repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	creationTime := job.CreationTime
	if creationTime.IsZero() {
		t.Error("Should set the creation time of the job, but did not")
	}
	if creationTime.Location() != time.UTC {
		t.Errorf("Wrong location for creationTime. Want UTC. Got %#v", creationTime.Location())
	}
	client := repo.(*redisRepository).storage.RedisClient()
	defer client.Close()
	items, err := client.HGetAll("job:" + job.ID).Result()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{
		"source":                           "http://nyt.net/source_here.mp4",
		"jobID":                            "job1",
		"providerName":                     "encoding.com",
		"providerJobID":                    "",
		"streamingparams_segmentDuration":  "10",
		"streamingparams_protocol":         "hls",
		"streamingparams_playlistFileName": "hls/playlist.m3u8",
		"creationTime":                     creationTime.Format(time.RFC3339Nano),
	}
	if !reflect.DeepEqual(items, expected) {
		pretty.Fdiff(os.Stderr, expected, items)
		t.Errorf("Wrong job hash returned from Redis. Want\n %#v.\n Got\n %#v.", expected, items)
	}
	setEntries, err := client.ZRange(jobsSetKey, 0, -1).Result()
	if err != nil {
		t.Fatal(err)
	}
	expectedSetEntries := []string{job.ID}
	if !reflect.DeepEqual(setEntries, expectedSetEntries) {
		pretty.Fdiff(os.Stderr, expectedSetEntries, setEntries)
		t.Errorf("Wrong job set returned from Redis. Want %#v. Got %#v.", expectedSetEntries, setEntries)
	}
}

func TestCreateJobIsSafe(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	jobs := []db.Job{
		{ID: "abcabc", ProviderName: "elastictranscoder"},
		{ID: "abcabc", ProviderJobID: "abf-123", ProviderName: "encoding.com"},
		{ID: "abcabc", ProviderJobID: "abc-213", ProviderName: "encoding.com"},
		{ID: "abcabc", ProviderJobID: "ff12", ProviderName: "encoding.com"},
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := repo.CreateJob(&jobs[i])
			if err != nil && err != redis.TxFailedErr {
				t.Error(err)
			}
		}(i % len(jobs))
	}
	wg.Wait()
}

func TestCreateJobNoID(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	job := db.Job{ProviderName: "elastictranscoder", ProviderJobID: "abc-123"}
	err = repo.CreateJob(&job)
	if err == nil {
		t.Fatal("got unexpected <nil> error")
	}
	if msg := "job id is required"; err.Error() != msg {
		t.Errorf("wrong error message\nWant %q\nGot  %q", msg, err.Error())
	}
}

func TestDeleteJob(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	job := db.Job{ID: "myjob"}
	err = repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeleteJob(&db.Job{ID: job.ID})
	if err != nil {
		t.Fatal(err)
	}
	client := repo.(*redisRepository).storage.RedisClient()
	hGetResult := client.HGetAll("job:myjob")
	if len(hGetResult.Val()) != 0 {
		t.Errorf("Unexpected value after delete call: %v", hGetResult.Val())
	}
	zRangeResult := client.ZRange(jobsSetKey, 0, -1)
	if len(zRangeResult.Val()) != 0 {
		t.Errorf("Unexpected value after delete call: %v", zRangeResult.Val())
	}
}

func TestDeleteJobNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	err = repo.DeleteJob(&db.Job{ID: "myjob"})
	if err != db.ErrJobNotFound {
		t.Errorf("Wrong error returned by DeleteJob. Want ErrJobNotFound. Got %#v.", err)
	}
}

func TestGetJob(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	job := db.Job{ID: "myjob"}
	err = repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	gotJob, err := repo.GetJob(job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*gotJob, job) {
		t.Errorf("Wrong job. Want %#v. Got %#v.", job, *gotJob)
	}
}

func TestGetJobNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(storage.Config)})
	if err != nil {
		t.Fatal(err)
	}
	gotJob, err := repo.GetJob("job:myjob")
	if err != db.ErrJobNotFound {
		t.Errorf("Wrong error returned. Want ErrJobNotFound. Got %#v.", err)
	}
	if gotJob != nil {
		t.Errorf("Unexpected non-nil job: %#v.", gotJob)
	}
}

func TestListJobs(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(storage.Config)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	jobs := []db.Job{
		{
			ID:            "job-1",
			ProviderName:  "encodingcom",
			ProviderJobID: "1",
		},
		{
			ID:            "job-2",
			ProviderName:  "encodingcom",
			ProviderJobID: "2",
		},
		{
			ID:            "job-3",
			ProviderName:  "encodingcom",
			ProviderJobID: "3",
		},
		{
			ID:            "job-4",
			ProviderName:  "encodingcom",
			ProviderJobID: "4",
		},
	}
	expectedJobs := make([]db.Job, len(jobs))
	for i, job := range jobs {
		err = repo.CreateJob(&job)
		if err != nil {
			t.Fatal(err)
		}
		expectedJobs[i] = job
	}
	gotJobs, err := repo.ListJobs(db.JobFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotJobs, expectedJobs) {
		t.Errorf("ListJobs({}): wrong list returned. Want %#v. Got %#v", expectedJobs, gotJobs)
	}
}

func TestListJobsLimit(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(storage.Config)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	jobs := []db.Job{
		{
			ID:            "job-1",
			ProviderName:  "encodingcom",
			ProviderJobID: "1",
		},
		{
			ID:            "job-2",
			ProviderName:  "encodingcom",
			ProviderJobID: "2",
		},
		{
			ID:            "job-3",
			ProviderName:  "encodingcom",
			ProviderJobID: "3",
		},
		{
			ID:            "job-4",
			ProviderName:  "encodingcom",
			ProviderJobID: "4",
		},
	}
	limit := 2
	expectedJobs := make([]db.Job, limit)
	for i, job := range jobs {
		err = repo.CreateJob(&job)
		if err != nil {
			t.Fatal(err)
		}
		if i < limit {
			expectedJobs[i] = job
		}
	}
	gotJobs, err := repo.ListJobs(db.JobFilter{Limit: uint(limit)})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotJobs, expectedJobs) {
		t.Errorf("ListJobs({}): wrong list returned. Want %#v. Got %#v", expectedJobs, gotJobs)
	}
}

func TestListJobsInconsistency(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(storage.Config)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	jobs := []db.Job{
		{
			ID:            "job-1",
			ProviderName:  "encodingcom",
			ProviderJobID: "1",
		},
		{
			ID:            "job-2",
			ProviderName:  "encodingcom",
			ProviderJobID: "2",
		},
		{
			ID:            "job-3",
			ProviderName:  "encodingcom",
			ProviderJobID: "3",
		},
		{
			ID:            "job-4",
			ProviderName:  "encodingcom",
			ProviderJobID: "4",
		},
	}
	redisRepo := repo.(*redisRepository)
	redisRepo.storage.RedisClient().ZAddNX("some-weird-id1", redis.Z{Member: jobs[0], Score: math.Inf(0)})
	redisRepo.storage.RedisClient().ZAddNX("some-weird-id2", redis.Z{Member: jobs[1], Score: math.Inf(0)})
	expectedJobs := make([]db.Job, len(jobs))
	for i, job := range jobs {
		err = repo.CreateJob(&job)
		if err != nil {
			t.Fatal(err)
		}
		expectedJobs[i] = job
	}
	gotJobs, err := repo.ListJobs(db.JobFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotJobs, expectedJobs) {
		t.Errorf("ListJobs({}): wrong list returned. Want %#v. Got %#v", expectedJobs, gotJobs)
	}
}

func TestListJobsFiltering(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(storage.Config)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Truncate(time.Millisecond)
	jobs := []db.Job{
		{
			ID:            "job-1",
			ProviderName:  "encodingcom",
			ProviderJobID: "1",
			CreationTime:  now.Add(-time.Hour),
		},
		{
			ID:            "job-2",
			ProviderName:  "encodingcom",
			ProviderJobID: "2",
			CreationTime:  now.Add(-40 * time.Minute),
		},
		{
			ID:            "job-3",
			ProviderName:  "encodingcom",
			ProviderJobID: "3",
			CreationTime:  now.Add(-10 * time.Minute),
		},
		{
			ID:            "job-4",
			ProviderName:  "encodingcom",
			ProviderJobID: "4",
			CreationTime:  now.Add(-3 * time.Second),
		},
	}
	expectedJobs := make([]db.Job, 0, 3)
	since := now.Add(-59 * time.Minute)
	redisRepo := repo.(*redisRepository)
	for _, job := range jobs {
		err = redisRepo.saveJob(&job)
		if err != nil {
			t.Fatal(err)
		}
		if job.CreationTime.After(since) {
			expectedJobs = append(expectedJobs, job)
		}
	}
	gotJobs, err := repo.ListJobs(db.JobFilter{Since: since})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotJobs, expectedJobs) {
		t.Errorf("ListJobs({}): wrong list returned\nWant %#v\nGot  %#v", expectedJobs, gotJobs)
	}
}

func TestListJobsFilteringAndLimit(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(storage.Config)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Truncate(time.Millisecond)
	jobs := []db.Job{
		{
			ID:            "job-1",
			ProviderName:  "encodingcom",
			ProviderJobID: "1",
			CreationTime:  now.Add(-time.Hour),
		},
		{
			ID:            "job-2",
			ProviderName:  "encodingcom",
			ProviderJobID: "2",
			CreationTime:  now.Add(-40 * time.Minute),
		},
		{
			ID:            "job-3",
			ProviderName:  "encodingcom",
			ProviderJobID: "3",
			CreationTime:  now.Add(-10 * time.Minute),
		},
		{
			ID:            "job-4",
			ProviderName:  "encodingcom",
			ProviderJobID: "4",
			CreationTime:  now.Add(-3 * time.Second),
		},
	}
	limit := 2
	expectedJobs := make([]db.Job, 0, limit)
	since := now.Add(-59 * time.Minute)
	redisRepo := repo.(*redisRepository)
	for _, job := range jobs {
		err = redisRepo.saveJob(&job)
		if err != nil {
			t.Fatal(err)
		}
		if job.CreationTime.After(since) && len(expectedJobs) < limit {
			expectedJobs = append(expectedJobs, job)
		}
	}
	gotJobs, err := repo.ListJobs(db.JobFilter{Since: since, Limit: uint(limit)})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotJobs, expectedJobs) {
		t.Errorf("ListJobs({}): wrong list returned. Want %#v. Got %#v", expectedJobs, gotJobs)
	}
}
