package db

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"

	"github.com/nytm/video-transcoding-api/config"
	"gopkg.in/redis.v3"
)

type redisRepository struct {
	config *config.Config
	client *redis.Client
}

// NewRedisJobRepository creates a new JobRepository that uses Redis for
// persistence.
func NewRedisJobRepository(cfg *config.Config) (JobRepository, error) {
	return &redisRepository{config: cfg}, nil
}

func (r *redisRepository) SaveJob(job *Job) error {
	if job.ID == "" {
		jobID, err := r.generateID()
		if err != nil {
			return err
		}
		job.ID = jobID
	}
	jobKey := "job:" + job.ID
	multi, err := r.redisClient().Watch(jobKey)
	if err != nil {
		return err
	}
	_, err = multi.Exec(func() error {
		multi.HSet(jobKey, "providerName", job.ProviderName)
		multi.HSet(jobKey, "providerJobID", job.ProviderJobID)
		multi.HSet(jobKey, "status", job.Status)
		return nil
	})
	return err
}

func (r *redisRepository) DeleteJob(job *Job) error {
	return nil
}

func (r *redisRepository) GetJob(id string) (*Job, error) {
	return nil, nil
}

func (r *redisRepository) generateID() (string, error) {
	var raw [8]byte
	n, err := rand.Read(raw[:])
	if err != nil {
		return "", err
	}
	if n != 8 {
		return "", io.ErrShortWrite
	}
	return fmt.Sprintf("%x", raw), nil
}

func (r *redisRepository) redisClient() *redis.Client {
	if r.client == nil {
		var sentinelAddrs []string
		if r.config.SentinelAddrs != "" {
			sentinelAddrs = strings.Split(r.config.SentinelAddrs, ",")
		}
		if len(sentinelAddrs) > 0 {
			r.client = redis.NewFailoverClient(&redis.FailoverOptions{
				SentinelAddrs: sentinelAddrs,
				MasterName:    r.config.SentinelMasterName,
				Password:      r.config.Password,
			})
		} else {
			redisAddr := r.config.RedisAddr
			if redisAddr == "" {
				redisAddr = "127.0.0.1:6379"
			}
			r.client = redis.NewClient(&redis.Options{
				Addr:     redisAddr,
				Password: r.config.Password,
			})
		}
	}
	return r.client
}
