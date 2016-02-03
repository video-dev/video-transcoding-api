package db

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/nytm/video-transcoding-api/config"
	"gopkg.in/redis.v3"
)

type redisRepository struct {
	config *config.Config
	client *redis.Client
	once   sync.Once
}

// NewRedisJobRepository creates a new JobRepository that uses Redis for
// persistence.
func NewRedisJobRepository(cfg *config.Config) (JobRepository, error) {
	repo := &redisRepository{config: cfg}
	repo.client = repo.redisClient()
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
	result := r.redisClient().HMSet(r.jobKey(job), "providerName", job.ProviderName, "providerJobID", job.ProviderJobID)
	return result.Err()
}

func (r *redisRepository) DeleteJob(job *Job) error {
	n, err := r.redisClient().Del(r.jobKey(job)).Result()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrJobNotFound
	}
	return nil
}

func (r *redisRepository) GetJob(id string) (*Job, error) {
	result, err := r.redisClient().HGetAllMap(r.jobKey(&Job{ID: id})).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, ErrJobNotFound
	}
	return &Job{
		ID:            id,
		ProviderJobID: result["providerJobID"],
		ProviderName:  result["providerName"],
	}, nil
}

func (r *redisRepository) jobKey(j *Job) string {
	return "job:" + j.ID
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
	r.once.Do(func() {
		var sentinelAddrs []string
		if r.config.Redis.SentinelAddrs != "" {
			sentinelAddrs = strings.Split(r.config.Redis.SentinelAddrs, ",")
		}
		if len(sentinelAddrs) > 0 {
			r.client = redis.NewFailoverClient(&redis.FailoverOptions{
				SentinelAddrs: sentinelAddrs,
				MasterName:    r.config.Redis.SentinelMasterName,
				Password:      r.config.Redis.Password,
				PoolSize:      r.config.Redis.PoolSize,
				PoolTimeout:   time.Duration(r.config.Redis.PoolTimeout) * time.Second,
			})
		} else {
			redisAddr := r.config.Redis.RedisAddr
			if redisAddr == "" {
				redisAddr = "127.0.0.1:6379"
			}
			r.client = redis.NewClient(&redis.Options{
				Addr:        redisAddr,
				Password:    r.config.Redis.Password,
				PoolSize:    r.config.Redis.PoolSize,
				PoolTimeout: time.Duration(r.config.Redis.PoolTimeout) * time.Second,
			})
		}
	})
	return r.client
}
