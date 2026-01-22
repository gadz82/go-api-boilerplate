package domain

import (
	"context"
	"time"
)

// CacheRepository defines the interface for cache operations.
// Implementations can use Redis, file-based storage, or any other caching mechanism.
type CacheRepository interface {
	// Get retrieves a value from the cache by key.
	// Returns an error if the key doesn't exist or the cache is unavailable.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value in the cache with the given key and TTL.
	// If ttl is 0, the value will not expire.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete removes a value from the cache by key.
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in the cache.
	Exists(ctx context.Context, key string) (bool, error)

	// Ping checks if the cache backend is available.
	Ping(ctx context.Context) error
}
