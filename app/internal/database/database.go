package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}
	return db, nil
}
