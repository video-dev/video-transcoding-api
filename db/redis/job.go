package redis

import (
	"strconv"
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
	job.CreationTime = time.Now().UTC()
	return r.saveJob(job)
}

func (r *redisRepository) saveJob(job *db.Job) error {
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
		multi.ZAddNX(jobsSetKey, redis.Z{Member: job.ID, Score: float64(job.CreationTime.UnixNano())})
		return nil
	})
	return err
}

func (r *redisRepository) DeleteJob(job *db.Job) error {
	err := r.delete(r.jobKey(job.ID), db.ErrJobNotFound)
	if err != nil {
		return err
	}
	return r.redisClient().ZRem(jobsSetKey, job.ID).Err()
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
	now := time.Now().UTC()
	rangeOpts := redis.ZRangeByScore{
		Min:   strconv.FormatInt(filter.Since.UnixNano(), 10),
		Max:   strconv.FormatInt(now.UnixNano(), 10),
		Count: int64(filter.Limit),
	}
	if rangeOpts.Count == 0 {
		rangeOpts.Count = -1
	}
	jobIDs, err := r.redisClient().ZRangeByScore(jobsSetKey, rangeOpts).Result()
	if err != nil {
		return nil, err
	}
	jobs := make([]db.Job, 0, len(jobIDs))
	for _, id := range jobIDs {
		job, err := r.GetJob(id)
		if err != nil && err != db.ErrJobNotFound {
			return nil, err
		}
		if job != nil {
			jobs = append(jobs, *job)
		}
	}
	return jobs, nil
}

func (r *redisRepository) jobKey(id string) string {
	return "job:" + id
}
