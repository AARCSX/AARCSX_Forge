package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/AARCSX/AARCSX_Forge/internal/config"
)

// Redis holds the Redis client.
type Redis struct {
	Client *redis.Client
}

// NewRedis creates a new Redis client from config.
func NewRedis(cfg config.Config) (*Redis, error) {
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("verifying redis connection: %w", err)
	}

	return &Redis{Client: client}, nil
}

// Close closes the Redis connection.
func (r *Redis) Close() error {
	return r.Client.Close()
}