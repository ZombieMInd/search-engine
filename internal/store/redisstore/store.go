package redisstore

import (
	"github.com/ZombieMInd/go-logger/internal/logger"
	"github.com/go-redis/redis"
)

type Store struct {
	client        *redis.Client
	logRepository logger.LogRepository
}

func New(c *redis.Client) *Store {
	return &Store{client: c}
}

func (s *Store) Log() logger.LogRepository {
	if s.logRepository == nil {
		s.logRepository = &LogRepository{store: s}
	}

	return s.logRepository
}
