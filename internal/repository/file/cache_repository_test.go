package file

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestCacheDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cache_test_*")
	require.NoError(t, err)
	return dir
}

func cleanupTestCacheDir(dir string) {
	os.RemoveAll(dir)
}

func TestFileCacheRepository_SetAndGet(t *testing.T) {
	cacheDir := setupTestCacheDir(t)
	defer cleanupTestCacheDir(cacheDir)

	repo, err := NewCacheRepository(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Set a value
	err = repo.Set(ctx, "test-key", "test-value", 0)
	assert.NoError(t, err)

	// Get the value
	val, err := repo.Get(ctx, "test-key")
	assert.NoError(t, err)
	assert.Equal(t, "test-value", val)
}

func TestFileCacheRepository_GetNotFound(t *testing.T) {
	cacheDir := setupTestCacheDir(t)
	defer cleanupTestCacheDir(cacheDir)

	repo, err := NewCacheRepository(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	_, err = repo.Get(ctx, "non-existent-key")
	assert.ErrorIs(t, err, ErrCacheKeyNotFound)
}

func TestFileCacheRepository_SetWithTTL(t *testing.T) {
	cacheDir := setupTestCacheDir(t)
	defer cleanupTestCacheDir(cacheDir)

	repo, err := NewCacheRepository(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Set a value with short TTL
	err = repo.Set(ctx, "ttl-key", "ttl-value", 100*time.Millisecond)
	assert.NoError(t, err)

	// Get immediately should work
	val, err := repo.Get(ctx, "ttl-key")
	assert.NoError(t, err)
	assert.Equal(t, "ttl-value", val)

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Get after expiration should fail
	_, err = repo.Get(ctx, "ttl-key")
	assert.ErrorIs(t, err, ErrCacheExpired)
}

func TestFileCacheRepository_Delete(t *testing.T) {
	cacheDir := setupTestCacheDir(t)
	defer cleanupTestCacheDir(cacheDir)

	repo, err := NewCacheRepository(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Set a value
	err = repo.Set(ctx, "delete-key", "delete-value", 0)
	assert.NoError(t, err)

	// Delete the value
	err = repo.Delete(ctx, "delete-key")
	assert.NoError(t, err)

	// Get should fail
	_, err = repo.Get(ctx, "delete-key")
	assert.ErrorIs(t, err, ErrCacheKeyNotFound)
}

func TestFileCacheRepository_DeleteNonExistent(t *testing.T) {
	cacheDir := setupTestCacheDir(t)
	defer cleanupTestCacheDir(cacheDir)

	repo, err := NewCacheRepository(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Delete non-existent key should not error
	err = repo.Delete(ctx, "non-existent-key")
	assert.NoError(t, err)
}

func TestFileCacheRepository_Exists(t *testing.T) {
	cacheDir := setupTestCacheDir(t)
	defer cleanupTestCacheDir(cacheDir)

	repo, err := NewCacheRepository(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Check non-existent key
	exists, err := repo.Exists(ctx, "test-key")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Set a value
	err = repo.Set(ctx, "test-key", "test-value", 0)
	assert.NoError(t, err)

	// Check existing key
	exists, err = repo.Exists(ctx, "test-key")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestFileCacheRepository_ExistsExpired(t *testing.T) {
	cacheDir := setupTestCacheDir(t)
	defer cleanupTestCacheDir(cacheDir)

	repo, err := NewCacheRepository(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Set a value with short TTL
	err = repo.Set(ctx, "expire-key", "expire-value", 100*time.Millisecond)
	assert.NoError(t, err)

	// Check immediately should return true
	exists, err := repo.Exists(ctx, "expire-key")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Check after expiration should return false
	exists, err = repo.Exists(ctx, "expire-key")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestFileCacheRepository_Ping(t *testing.T) {
	cacheDir := setupTestCacheDir(t)
	defer cleanupTestCacheDir(cacheDir)

	repo, err := NewCacheRepository(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	err = repo.Ping(ctx)
	assert.NoError(t, err)
}

func TestFileCacheRepository_OverwriteValue(t *testing.T) {
	cacheDir := setupTestCacheDir(t)
	defer cleanupTestCacheDir(cacheDir)

	repo, err := NewCacheRepository(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	// Set initial value
	err = repo.Set(ctx, "overwrite-key", "initial-value", 0)
	assert.NoError(t, err)

	// Overwrite with new value
	err = repo.Set(ctx, "overwrite-key", "new-value", 0)
	assert.NoError(t, err)

	// Get should return new value
	val, err := repo.Get(ctx, "overwrite-key")
	assert.NoError(t, err)
	assert.Equal(t, "new-value", val)
}
