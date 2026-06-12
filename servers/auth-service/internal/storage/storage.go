package storage

import (
	apperrors "auth-service/errors"
	"auth-service/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, connstr string) (*DB, error) {
	database, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, fmt.Errorf("NewDB: failed to pool connect to db: %w", err)
	}

	if err = database.Ping(ctx); err != nil {
		database.Close()
		return nil, fmt.Errorf("NewDB: failed to Ping db: %w", err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`

	if _, err = database.Exec(ctx, query); err != nil {
		database.Close()
		return nil, fmt.Errorf("NewDB: failed to create db: %w", err)
	}

	return &DB{pool: database}, nil
}

func (db *DB) CreateUser(ctx context.Context, email string, password string) (int, error) {
	var id int
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: Error with generate password: %w", err)
	}

	query := `
		INSERT INTO users (email, password_hash)
		VALUES (@email, @password_hash)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"email":         email,
		"password_hash": string(passwordHash),
	}

	if err = db.pool.QueryRow(ctx, query, args).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, apperrors.ErrUserAlreadyExists
			}
		}

		return 0, fmt.Errorf("CreateUser: Error with paste in db: %w", err)
	}

	return id, nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (models.Authorization, error) {
	var user models.Authorization

	query := `
		SELECT id, email, password_hash
		FROM users
		WHERE email = @email
	`

	args := pgx.NamedArgs{
		"email": email,
	}

	if err := db.pool.QueryRow(ctx, query, args).Scan(&user.UserID, &user.Email, &user.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Authorization{}, apperrors.ErrUserNotFound
		}

		return models.Authorization{}, fmt.Errorf("GetUserByEmail: Failed to find user: %w", err)
	}

	return user, nil
}

func (db *DB) Close() {
	db.pool.Close()
}
