package storage

import (
	"context"
	"event-driven/internal/config"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(cfg *config.Config) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), cfg.PostgreDSN)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresDB: failed to connect: %w", err)
	}
	if err = db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("NewPostgresDB: failed to ping db: %w", err)
	}
	return db, nil
}
