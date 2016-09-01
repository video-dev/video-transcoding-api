package redis

import (
	"strings"
	"time"

	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/db/redis/storage"
)

// NewRepository creates a new Repository that uses Redis for persistence.
func NewRepository(cfg *config.Config) (db.Repository, error) {
	var sentinelAddrs []string
	if cfg.Redis.SentinelAddrs != "" {
		sentinelAddrs = strings.Split(cfg.Redis.SentinelAddrs, ",")
	}
	s, err := storage.NewStorage(storage.NewStorageOptions{
		RedisAddr:          cfg.Redis.RedisAddr,
		SentinelAddrs:      sentinelAddrs,
		SentinelMasterName: cfg.Redis.SentinelMasterName,
		Password:           cfg.Redis.Password,
		PoolSize:           cfg.Redis.PoolSize,
		PoolTimeout:        time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		IdleTimeout:        time.Duration(cfg.Redis.IdleTimeout) * time.Second,
		IdleCheckFrequency: time.Duration(cfg.Redis.IdleCheckFrequency) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &redisRepository{config: cfg, storage: s}, nil
}

type redisRepository struct {
	config  *config.Config
	storage *storage.Storage
}
