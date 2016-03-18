package redis

import (
	"time"

	"github.com/nytm/video-transcoding-api/db"
	"gopkg.in/redis.v3"
)

const jobsSetKey = "jobs"

func (r *redisRepository) CreateJob(job *db.Job) error {
	if job.ID == "" {
		jobID, err := r.generateID()
		if err != nil {
			return err
		}
		job.ID = jobID
	}
	job.CreationTime = time.Now().In(time.UTC)
	fields, err := r.fieldList(job)
	if err != nil {
		return err
	}
	jobKey := r.jobKey(job.ID)
	multi, err := r.redisClient().Watch(jobKey)
	if err != nil {
		return err
	}
	_, err = multi.Exec(func() error {
		multi.HMSet(jobKey, fields[0], fields[1], fields[2:]...)
		multi.ZAddNX(jobsSetKey, redis.Z{Member: jobKey, Score: float64(job.CreationTime.UnixNano())})
		return nil
	})
	return err
}

func (r *redisRepository) DeleteJob(job *db.Job) error {
	return r.delete(r.jobKey(job.ID), db.ErrJobNotFound)
}

func (r *redisRepository) GetJob(id string) (*db.Job, error) {
	job := db.Job{ID: id}
	err := r.load(r.jobKey(id), &job)
	if err == errNotFound {
		return nil, db.ErrJobNotFound
	}
	return &job, err
}

func (r *redisRepository) ListJobs(filter db.JobFilter) ([]db.Job, error) {
	return nil, nil
}

func (r *redisRepository) jobKey(id string) string {
	return "job:" + id
}
