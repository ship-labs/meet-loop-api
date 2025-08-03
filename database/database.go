// Package database provides functions for connecting to the database.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Dial(ctx context.Context, address string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("opening connection to %s: %w", address, err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging %s: %w", address, err)
	}

	return pool, nil
}
