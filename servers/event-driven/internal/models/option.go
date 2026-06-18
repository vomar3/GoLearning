package models

import "context"

type Option struct {
	ID         string
	PollID     string
	Text       string
	VotesCount int64
}

type OptionRepository interface {
	Create(context.Context, string, string) (string, error)
	ListOptions(context.Context, string) ([]*Option, error)
}
