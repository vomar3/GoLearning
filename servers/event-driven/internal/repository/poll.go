package repository

import (
	"context"
	"database/sql"
	"errors"
	"event-driven/internal/models"
	"fmt"
)

type pollRepositoryImpl struct {
	db *sql.DB
}

func NewPollRepository(db *sql.DB) models.PollRepository {
	return &pollRepositoryImpl{
		db: db,
	}
}

func (r *pollRepositoryImpl) Create(ctx context.Context, title, description string) (string, error) {
	query := `
	INSERT INTO polls (title, description) 
	VALUES ($1, $2) 
	RETURNING id
	`

	var newID string

	// Getting only one row after insertion
	err := r.db.QueryRowContext(ctx, query, title, description).Scan(&newID)
	if err != nil {
		return "", fmt.Errorf("Create: row insertion error: %w", err)
	}

	return newID, nil
}

func (r *pollRepositoryImpl) GetByID(ctx context.Context, id string) (*models.Poll, error) {
	query := `
		SELECT id, title, description, is_active, created_at
		FROM polls
		WHERE id = $1
		LIMIT 1
	`

	var data models.Poll

	err := r.db.QueryRowContext(ctx, query, id).Scan(&data.ID, &data.Title, &data.Description, &data.IsActive, &data.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("GetByID: User with %s ID not found", id)
		}

		return nil, fmt.Errorf("GetByID: error getting data by id: %w", err)
	}

	return &data, nil
}

func (r *pollRepositoryImpl) Update(ctx context.Context, id, title, description string, isActive bool) (*models.Poll, error) {
	query := `
		UPDATE polls
		SET title = $1, description = $2, is_active = $3
		WHERE id = $4
		RETURNING id, title, description, is_active, created_at
	`

	var data models.Poll

	err := r.db.QueryRowContext(ctx, query, title, description, isActive, id).Scan(&data.ID, &data.Title, &data.Description, &data.IsActive, &data.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Update: User with %s ID not found", id)
		}

		return nil, fmt.Errorf("Update: error updating by data: %w", err)
	}

	return &data, nil
}

func (r *pollRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `
		DELETE 
		FROM polls
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Delete: error deleting by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Delete: failed to get the number of deleted rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Delete: User with %s ID not found", id)
	}

	return nil
}

func (r *pollRepositoryImpl) List(ctx context.Context, limit, offset int, onlyActive bool) ([]*models.Poll, error) {
	query := `
		SELECT id, title, description, is_active, created_at
		FROM polls
		WHERE ($3 = false OR is_active = true)
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`

	var polls []*models.Poll

	rows, err := r.db.QueryContext(ctx, query, limit, offset, onlyActive)
	if err != nil {
		return nil, fmt.Errorf("List: error with getting data: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		p := &models.Poll{}

		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("List: error scanning data: %w", err)
		}

		polls = append(polls, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("List: error after scanning data: %w", err)
	}

	return polls, nil
}
