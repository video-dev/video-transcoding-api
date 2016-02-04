package redis

import "github.com/nytm/video-transcoding-api/db"

func (r *redisRepository) SaveJob(job *db.Job) error {
	if job.ID == "" {
		jobID, err := r.generateID()
		if err != nil {
			return err
		}
		job.ID = jobID
	}
	return r.save(r.jobKey(job.ID), job)
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

func (r *redisRepository) jobKey(id string) string {
	return "job:" + id
}
