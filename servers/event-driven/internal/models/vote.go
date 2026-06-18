package models

import (
	"context"
	"time"
)

type Vote struct {
	ID        string
	PollId    string
	OptionId  string
	UserId    string
	CreatedAt time.Time
}

type VoteRepository interface {
	Cast(ctx context.Context, pollID, optionID, userID string) error
}
