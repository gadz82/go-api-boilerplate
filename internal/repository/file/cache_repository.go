package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gadz82/go-api-boilerplate/internal/domain"
)

// ErrCacheKeyNotFound is returned when a key is not found in the cache.
var ErrCacheKeyNotFound = errors.New("cache key not found")

// ErrCacheExpired is returned when a cached item has expired.
var ErrCacheExpired = errors.New("cache item expired")

// cacheItem represents a cached value with optional expiration.
type cacheItem struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	HasExpiry bool      `json:"has_expiry"`
}

// fileCacheRepository implements CacheRepository using file-based storage.
type fileCacheRepository struct {
	cacheDir string
	mu       sync.RWMutex
}

// NewCacheRepository creates a new file-based cache repository.
// The cacheDir parameter specifies the directory where cache files will be stored.
func NewCacheRepository(cacheDir string) (domain.CacheRepository, error) {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	return &fileCacheRepository{
		cacheDir: cacheDir,
	}, nil
}

// keyToFilename converts a cache key to a safe filename.
func (r *fileCacheRepository) keyToFilename(key string) string {
	// Use a simple hash-like approach to create safe filenames
	safeKey := filepath.Base(key)
	if safeKey == "." || safeKey == "/" {
		safeKey = "default"
	}
	return filepath.Join(r.cacheDir, safeKey+".cache")
}

func (r *fileCacheRepository) Get(ctx context.Context, key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filename := r.keyToFilename(key)
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrCacheKeyNotFound
		}
		return "", err
	}

	var item cacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return "", err
	}

	// Check if item has expired
	if item.HasExpiry && time.Now().After(item.ExpiresAt) {
		// Clean up expired item
		go func() {
			r.mu.Lock()
			defer r.mu.Unlock()
			os.Remove(filename)
		}()
		return "", ErrCacheExpired
	}

	return item.Value, nil
}

func (r *fileCacheRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := cacheItem{
		Value:     value,
		HasExpiry: ttl > 0,
	}

	if ttl > 0 {
		item.ExpiresAt = time.Now().Add(ttl)
	}

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	filename := r.keyToFilename(key)
	return os.WriteFile(filename, data, 0644)
}

func (r *fileCacheRepository) Delete(ctx context.Context, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filename := r.keyToFilename(key)
	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (r *fileCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filename := r.keyToFilename(key)
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	var item cacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return false, nil
	}

	// Check if item has expired
	if item.HasExpiry && time.Now().After(item.ExpiresAt) {
		return false, nil
	}

	return true, nil
}

func (r *fileCacheRepository) Ping(ctx context.Context) error {
	// For file cache, we just verify the cache directory is accessible
	_, err := os.Stat(r.cacheDir)
	return err
}
