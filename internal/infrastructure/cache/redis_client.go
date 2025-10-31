package cache

import (
	"context"
	"fmt"
	"strconv"

	"api-employees-and-departments/config"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and connects to a Redis client
func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	db, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		return nil, fmt.Errorf("invalid redis db: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       db,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
