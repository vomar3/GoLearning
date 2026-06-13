package models

import (
	"context"
	"time"
)

type Poll struct {
	ID          string
	Title       string
	Description *string // may be null
	IsActive    bool
	CreatedAt   time.Time
}

type PollRepository interface {
	Create(ctx context.Context, title, description string) (string, error)
	GetByID(ctx context.Context, id string) (*Poll, error)
	Update(ctx context.Context, id, title, description string, isActive bool) (*Poll, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int, onlyActive bool) ([]*Poll, error)
}
