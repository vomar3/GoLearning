package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("url not found")

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(connStr string) (*DB, error) {
	val, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create DB: %w", err)
	}

	if err := val.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_code TEXT UNIQUE NOT NULL,
			long_url TEXT NOT NULL
		);
	`

	if _, err := val.Exec(context.Background(), query); err != nil {
		return nil, fmt.Errorf("NewDB: failed to create tabe: %w", err)
	}

	return &DB{pool: val}, nil
}

func (d *DB) SaveUrl(ctx context.Context, shortCode string, longUrl string) error {
	query := `
		INSERT INTO urls (long_url, short_code)
		VALUES ($1, $2)
	`

	_, err := d.pool.Exec(ctx, query, longUrl, shortCode)
	if err != nil {
		return fmt.Errorf("SaveUrl: failed to save url: %w", err)
	}

	return nil
}

func (d *DB) GetUrl(ctx context.Context, shortCode string) (string, error) {
	query := `
		SELECT long_url FROM urls WHERE short_code = $1
	`

	var longUrl string
	err := d.pool.QueryRow(ctx, query, shortCode).Scan(&longUrl)

	if err != nil {
		return "", fmt.Errorf("GetUrl: failed to get url: %w", err)
	}

	return longUrl, nil
}

func (d *DB) DeleteUrl(ctx context.Context, shortCode string) error {
	query := `
		DELETE FROM urls WHERE short_code = $1
	`

	tag, err := d.pool.Exec(ctx, query, shortCode)
	if err != nil {
		return fmt.Errorf("DeleteUrl: failed to delete url: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (d *DB) Close() {
	d.pool.Close()
}
