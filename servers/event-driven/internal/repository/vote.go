package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAlreadyExists = errors.New("vote already exists")

type voteRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) *voteRepositoryImpl {
	return &voteRepositoryImpl{
		db: db,
	}
}

func (r *voteRepositoryImpl) Cast(ctx context.Context, pollID, optionID, userID string) error {
	// starting transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Cast: failed to check db: %w", err)
	}

	defer tx.Rollback(ctx)

	// checking already voted
	query := `
		SELECT COUNT(*)
		FROM votes
		WHERE poll_id = $1 AND user_id = $2
	`

	var alreadyVoted int

	err = tx.QueryRow(ctx, query, pollID, userID).Scan(&alreadyVoted)
	if err != nil {
		return fmt.Errorf("Cast: error with check already voted: %w", err)
	}

	if alreadyVoted > 0 {
		return ErrAlreadyExists
	}

	// voice applied
	query = `
		INSERT INTO votes (id, poll_id, option_id, user_id, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW())
	`

	_, err = tx.Exec(ctx, pollID, optionID, userID)
	if err != nil {
		return fmt.Errorf("Cast: error to apply the vote: %w", err)
	}

	// update count
	query = `
		UPDATE options
		SET votes_count = votes_count + 1
		WHERE id = $1
	`

	_, err = tx.Exec(ctx, query, optionID)
	if err != nil {
		return fmt.Errorf("Cast: error to increment the counter: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Cast: error with commit the transaction: %w", err)
	}

	return nil
}
