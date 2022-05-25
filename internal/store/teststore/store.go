package teststore

import (
	"github.com/ZombieMInd/go-logger/internal/logger"
	"github.com/google/uuid"
)

type Store struct {
	logRepository logger.LogRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Log() logger.LogRepository {
	if s.logRepository == nil {
		s.logRepository = &LogRepository{
			store: s,
			logs:  make(map[uuid.UUID]*logger.LogRequest),
		}
	}

	return s.logRepository
}
