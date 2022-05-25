package sqlstore

import (
	"context"
	"database/sql"

	"github.com/ZombieMInd/go-logger/internal/logger"
)

type DB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Store struct {
	db            DB
	logRepository logger.LogRepository
}

func New(db DB) *Store {
	return &Store{db: db}
}

func (s *Store) Log() logger.LogRepository {
	if s.logRepository == nil {
		s.logRepository = &LogRepository{store: s}
	}

	return s.logRepository
}
