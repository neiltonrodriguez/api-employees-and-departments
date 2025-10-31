package cache

import (
	"context"
	"time"

	domainCache "api-employees-and-departments/internal/domain/cache"

	"github.com/redis/go-redis/v9"
)

// RedisCache is an adapter that implements domain.Cache using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new RedisCache that implements the domain Cache interface
func NewRedisCache(client *redis.Client) domainCache.Cache {
	return &RedisCache{
		client: client,
	}
}

// Get retrieves a value from Redis by key
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key not found
	}
	return val, err
}

// Set stores a value in Redis with a TTL
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a value from Redis
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeleteByPattern deletes all keys matching a pattern (e.g., "department:*")
func (r *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	// Scan for all matching keys
	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	// Delete all found keys
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists checks if a key exists in Redis
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Ping checks if Redis connection is healthy
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
