package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/NYTimes/video-transcoding-api/db"
	"github.com/NYTimes/video-transcoding-api/db/redis/storage"
	"github.com/go-redis/redis"
)

const jobsSetKey = "jobs"

func (r *redisRepository) CreateJob(job *db.Job) error {
	if job.ID == "" {
		return errors.New("job id is required")
	}
	job.CreationTime = time.Now().UTC().Truncate(time.Millisecond)
	return r.saveJob(job)
}

func (r *redisRepository) saveJob(job *db.Job) error {
	fields, err := r.storage.FieldMap(job)
	if err != nil {
		return err
	}
	jobKey := r.jobKey(job.ID)
	return r.storage.RedisClient().Watch(func(tx *redis.Tx) error {
		err := tx.HMSet(jobKey, fields).Err()
		if err != nil {
			return err
		}
		return tx.ZAddNX(jobsSetKey, redis.Z{Member: job.ID, Score: float64(job.CreationTime.UnixNano())}).Err()
	}, jobKey)
}

func (r *redisRepository) DeleteJob(job *db.Job) error {
	err := r.storage.Delete(r.jobKey(job.ID))
	if err != nil {
		if err == storage.ErrNotFound {
			return db.ErrJobNotFound
		}
		return err
	}
	return r.storage.RedisClient().ZRem(jobsSetKey, job.ID).Err()
}

func (r *redisRepository) GetJob(id string) (*db.Job, error) {
	job := db.Job{ID: id}
	err := r.storage.Load(r.jobKey(id), &job)
	if err == storage.ErrNotFound {
		return nil, db.ErrJobNotFound
	}
	return &job, err
}

func (r *redisRepository) ListJobs(filter db.JobFilter) ([]db.Job, error) {
	now := time.Now().UTC()
	rangeOpts := redis.ZRangeBy{
		Min:   strconv.FormatInt(filter.Since.UnixNano(), 10),
		Max:   strconv.FormatInt(now.UnixNano(), 10),
		Count: int64(filter.Limit),
	}
	if rangeOpts.Count == 0 {
		rangeOpts.Count = -1
	}
	jobIDs, err := r.storage.RedisClient().ZRangeByScore(jobsSetKey, rangeOpts).Result()
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
