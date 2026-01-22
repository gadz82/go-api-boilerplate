package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/gadz82/go-api-boilerplate/internal/domain"
)

type cacheRepository struct {
	client redis.Cmdable
}

// NewCacheRepository creates a new Redis-based cache repository.
func NewCacheRepository(client redis.Cmdable) domain.CacheRepository {
	return &cacheRepository{client: client}
}

func (r *cacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *cacheRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *cacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *cacheRepository) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
