package redis

import (
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"github.com/nytm/video-transcoding-api/db/redis/storage"
)

// NewRepository creates a new Repository that uses Redis for persistence.
func NewRepository(cfg *config.Config) (db.Repository, error) {
	s, err := storage.NewStorage(cfg.Redis)
	if err != nil {
		return nil, err
	}
	return &redisRepository{config: cfg, storage: s}, nil
}

type redisRepository struct {
	config  *config.Config
	storage *storage.Storage
}
