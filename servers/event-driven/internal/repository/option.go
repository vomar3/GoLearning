package repository

import (
	"context"
	"event-driven/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type optionRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewOptionRepository(db *pgxpool.Pool) *optionRepositoryImpl {
	return &optionRepositoryImpl{
		db: db,
	}
}

func (o *optionRepositoryImpl) Create(ctx context.Context, poll_id, text string) (string, error) {
	query := `
		INSERT INTO options (poll_id, text)
		VALUES ($1, $2)
		RETURNING id
	`

	var id string

	err := o.db.QueryRow(ctx, query, poll_id, text).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("Create: error to create new option: %w", err)
	}

	return id, nil
}

func (o *optionRepositoryImpl) ListOptions(ctx context.Context, poll_id string) ([]*models.Option, error) {
	query := `
		SELECT id, poll_id, text, votes_count
		FROM options
		WHERE poll_id = $1
		ORDER BY votes_count DESC
	`

	var options []*models.Option

	rows, err := o.db.Query(ctx, query, poll_id)
	if err != nil {
		return nil, fmt.Errorf("ListOptions: error to get options: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		newOption := &models.Option{}

		if err := rows.Scan(&newOption.ID, &newOption.PollID, &newOption.Text, &newOption.VotesCount); err != nil {
			return nil, fmt.Errorf("ListOptions: error to write option: %w", err)
		}

		options = append(options, newOption)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListOptions: error after scanning data: %w", err)
	}

	return options, nil
}
