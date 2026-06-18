package redis

import (
	"context"
	"event-driven/internal/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config, ctx context.Context) (*redis.Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	_, err := rds.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("NewRedisClient: error to ping redis: %w", err)
	}

	return rds, nil
}
