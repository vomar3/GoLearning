package storage

import (
	"context"
	"fmt"
	"order-system/internal/models"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB(ctx context.Context, connString string) (*DB, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("db.NewDB: unable to connect to database: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Init(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			item TEXT NOT NULL,
			price INT NOT NULL,
			status TEXT NOT NULL DEFAULT 'Processed',
			created_at TIMESTAMP DEFAULT NOW()
		);
	`

	if _, err := db.conn.Exec(ctx, query); err != nil {
		return fmt.Errorf("db.Init: failed to create table: %w", err)
	}

	return nil
}

func (db *DB) SaveOrder(ctx context.Context, order models.OrderRequest) error {
	query := `
		INSERT INTO orders (id, item, price, status)
		VALUES (@id, @item, @price, @status)
	`

	args := pgx.NamedArgs{
		"id":     order.ID,
		"item":   order.Item,
		"price":  order.Price,
		"status": "Processed",
	}

	if _, err := db.conn.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("db.SaveOrder: failed to save order: %w", err)
	}

	return nil
}

func (db *DB) GetOrder(ctx context.Context, id string) (models.OrderRequest, error) {
	query := `
		SELECT id, item, price FROM orders
		WHERE id = @id
	`

	var model models.OrderRequest

	args := pgx.NamedArgs{
		"id": id,
	}

	err := db.conn.QueryRow(ctx, query, args).Scan(
		&model.ID,
		&model.Item,
		&model.Price,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.OrderRequest{}, pgx.ErrNoRows
		}

		return models.OrderRequest{}, fmt.Errorf("Error with get order with id: %w", err)
	}

	return model, nil
}

func (db *DB) DeleteOrder(ctx context.Context, id string) error {
	query := `
		DELETE FROM orders
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	tag, err := db.conn.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.DeleteOrder: failed to delete order: %w", err)
	}

	if count := tag.RowsAffected(); count == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *DB) Close(ctx context.Context) error {
	return db.conn.Close(ctx)
}
