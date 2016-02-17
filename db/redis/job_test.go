package redis

import (
	"reflect"
	"sync"
	"testing"

	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"gopkg.in/redis.v3"
)

func TestCreateJob(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	job := db.Job{ProviderName: "encoding.com", AdaptiveStreaming: true}
	err = repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	if job.ID == "" {
		t.Fatal("Job ID should have been generated on CreateJob")
	}
	client := repo.(*redisRepository).redisClient()
	defer client.Close()
	items, err := client.HGetAllMap("job:" + job.ID).Result()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{
		"providerName":      "encoding.com",
		"providerJobID":     "",
		"adaptiveStreaming": "true",
	}
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Wrong job hash returned from Redis. Want %#v. Got %#v.", expected, items)
	}
}

func TestCreateJobPredefinedID(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	cfg.Redis = new(config.Redis)
	repo, err := NewRepository(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	job := db.Job{ID: "myjob", ProviderName: "encoding.com"}
	err = repo.CreateJob(&job)
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != "myjob" {
		t.Errorf("Job ID should not be regenerated when it's already defined. Got %q instead of %q.", job.ID, "myjob")
	}
	client := repo.(*redisRepository).redisClient()
	defer client.Close()
	items, err := client.HGetAllMap("job:myjob").Result()
	expected := map[string]string{
		"providerName":      "encoding.com",
		"providerJobID":     "",
		"adaptiveStreaming": "false",
	}
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Wrong job hash returned from Redis. Want %#v. Got %#v.", expected, items)
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
	repo, err := NewRepository(&config.Config{Redis: new(config.Redis)})
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

func TestDeleteJob(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(config.Redis)})
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
	client := repo.(*redisRepository).redisClient()
	result := client.HGetAllMap("job:myjob")
	if len(result.Val()) != 0 {
		t.Errorf("Unexpected value after delete call: %v", result.Val())
	}
}

func TestDeleteJobNotFound(t *testing.T) {
	err := cleanRedis()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewRepository(&config.Config{Redis: new(config.Redis)})
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
	repo, err := NewRepository(&config.Config{Redis: new(config.Redis)})
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
	repo, err := NewRepository(&config.Config{Redis: new(config.Redis)})
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
