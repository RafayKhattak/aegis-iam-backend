package repository

import (
	"context"
	"errors"

	"github.com/RafayKhattak/aegis-iam-backend/internal/repository/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store wraps the generated sqlc querier with a pooled PostgreSQL connection.
type Store struct {
	Pool *pgxpool.Pool
	db.Querier
}

// NewStore initializes the database pool and verifies connectivity.
func NewStore(ctx context.Context, connString string) (*Store, error) {
	if connString == "" {
		return nil, errors.New("database connection string is required")
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Store{
		Pool:    pool,
		Querier: db.New(pool),
	}, nil
}

// Close gracefully releases database resources.
func (s *Store) Close() {
	if s == nil || s.Pool == nil {
		return
	}

	s.Pool.Close()
}
