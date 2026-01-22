package storage

import "context"

type Storage interface {
	SaveUrl(ctx context.Context, shortCode string, longUrl string) error
	GetUrl(ctx context.Context, shortCode string) (string, error)
	Close()
}
