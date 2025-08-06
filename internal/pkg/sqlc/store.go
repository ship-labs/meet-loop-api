// Package sqlc
package sqlc

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	pool *pgxpool.Pool
}

func (s *Store) ExecuteTransaction(ctx context.Context, f func() error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	if err := f(); err != nil {
		if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
			return fmt.Errorf("executing provided function: %w, rolling back transaction: %w", err, rollBackErr)
		}
		return fmt.Errorf("executing provided function: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil

}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(pool),
		pool:    pool,
	}
}
