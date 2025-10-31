package cache

import (
	"context"
	"time"
)

// Cache is the interface for caching operations in the domain layer.
// This abstraction allows the domain to be independent of specific cache implementations (Redis, Memcached, etc.)
type Cache interface {
	// Get retrieves a value from cache by key
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value in cache with a TTL
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// DeleteByPattern deletes all keys matching a pattern (e.g., "department:*")
	DeleteByPattern(ctx context.Context, pattern string) error

	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)

	// Ping checks if cache connection is healthy
	Ping(ctx context.Context) error
}

// CacheKeyBuilder helps build consistent cache keys
type CacheKeyBuilder struct {
	prefix string
}

// NewCacheKeyBuilder creates a new cache key builder
func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{prefix: prefix}
}

// Build creates a cache key with the prefix
func (b *CacheKeyBuilder) Build(parts ...string) string {
	key := b.prefix
	for _, part := range parts {
		key += ":" + part
	}
	return key
}
