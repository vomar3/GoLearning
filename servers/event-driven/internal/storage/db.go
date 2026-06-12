package storage

import (
	"database/sql"
	"event-driven/internal/config"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.PostgreDSN)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresDB: failed to connect to db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("NewPostgresDB: failed to ping db: %w", err)
	}

	return db, nil
}
