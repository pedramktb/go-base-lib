package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pedramktb/go-base-lib/lifecycle"
)

// NewDB creates a new database connection using the provided connection string.
func NewDB(ctx context.Context, connString string) (*sql.DB, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	go func() {
		if done, err := lifecycle.RegisterCloser(ctx); err == nil {
			defer done(nil)
		}
		<-ctx.Done()
		pool.Close()
	}()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return stdlib.OpenDBFromPool(pool), nil
}
