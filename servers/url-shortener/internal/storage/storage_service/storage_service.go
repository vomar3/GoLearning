package storageservice

import (
	"context"
	"fmt"
	"log/slog"
)

type PostgresSaver interface {
	SaveUrl(ctx context.Context, shortCode string, longUrl string) error
	GetUrl(ctx context.Context, shortCode string) (string, error)
	DeleteUrl(ctx context.Context, shortCode string) error
}

type RedisSaver interface {
	SaveUrl(ctx context.Context, shortCode string, longUrl string) error
	GetUrl(ctx context.Context, shortCode string) (string, error)
	DeleteUrl(ctx context.Context, shortcode string) error
}

type Service struct {
	pg    PostgresSaver
	redis RedisSaver
	log   *slog.Logger
}

func NewService(log *slog.Logger, pg PostgresSaver, redis RedisSaver) *Service {
	return &Service{
		pg:    pg,
		redis: redis,
		log:   log,
	}
}

func (s *Service) SaveUrl(ctx context.Context, url string, alias string) error {
	if err := s.pg.SaveUrl(ctx, alias, url); err != nil {
		return fmt.Errorf("SaveUrl: Failed to save to postgress: %w", err)
	}

	if err := s.redis.SaveUrl(ctx, alias, url); err != nil {
		s.log.Error("failed to save to redis", slog.String("error", err.Error()))
	}

	return nil
}

func (s *Service) GetUrl(ctx context.Context, alias string) (string, error) {
	val, err := s.redis.GetUrl(ctx, alias)
	if err == nil && val != "" {
		return val, nil
	}

	val, err = s.pg.GetUrl(ctx, alias)
	if err != nil {
		return "", err
	}

	if err := s.redis.SaveUrl(ctx, alias, val); err != nil {
		s.log.Error("failed to update cache", slog.String("error", err.Error()))
	}

	return val, nil
}

func (s *Service) DeleteUrl(ctx context.Context, alias string) error {
	if err := s.pg.DeleteUrl(ctx, alias); err != nil {
		return fmt.Errorf("DeleteUrl: failed to delete from DB: %w", err)
	}

	if err := s.redis.DeleteUrl(ctx, alias); err != nil {
		return fmt.Errorf("DeleteUrl: failed to delete from redis: %w", err)
	}

	return nil
}
