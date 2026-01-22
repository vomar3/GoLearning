package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(addr string, password string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("NewClient: failed to connect to redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) GetUrl(ctx context.Context, shortcode string) (string, error) {
	val, err := c.rdb.Get(ctx, shortcode).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("GetUrl: url not found in cache")
	} else if err != nil {
		return "", err
	}

	return val, nil
}

func (c *Client) SaveUrl(ctx context.Context, shortcode string, longUrl string) error {
	err := c.rdb.Set(ctx, shortcode, longUrl, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("SaveUrl: failed to set: %w", err)
	}

	return nil
}

func (c *Client) DeleteUrl(ctx context.Context, shortcode string) error {
	err := c.rdb.Del(ctx, shortcode).Err()
	if err != nil {
		return fmt.Errorf("DeleteUrl: failed to delete: %w", err)
	}

	return nil
}
