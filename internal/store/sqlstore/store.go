package sqlstore

import (
	"context"
	"database/sql"
)

type DB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Store struct {
	db DB
}

func New(db DB) *Store {
	return &Store{db: db}
}

func (s *Store) Log() {
}
