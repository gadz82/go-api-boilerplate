package redis

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestCacheRepository_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewCacheRepository(db)
	ctx := context.Background()

	mock.ExpectGet("test-key").SetVal("test-value")

	val, err := repo.Get(ctx, "test-key")
	assert.NoError(t, err)
	assert.Equal(t, "test-value", val)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCacheRepository_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewCacheRepository(db)
	ctx := context.Background()

	mock.ExpectSet("test-key", "test-value", 0).SetVal("OK")

	err := repo.Set(ctx, "test-key", "test-value", 0)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCacheRepository_SetWithTTL(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewCacheRepository(db)
	ctx := context.Background()

	ttl := 5 * time.Minute
	mock.ExpectSet("test-key", "test-value", ttl).SetVal("OK")

	err := repo.Set(ctx, "test-key", "test-value", ttl)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCacheRepository_Delete(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewCacheRepository(db)
	ctx := context.Background()

	mock.ExpectDel("test-key").SetVal(1)

	err := repo.Delete(ctx, "test-key")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCacheRepository_Exists(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewCacheRepository(db)
	ctx := context.Background()

	mock.ExpectExists("test-key").SetVal(1)

	exists, err := repo.Exists(ctx, "test-key")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCacheRepository_Exists_NotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewCacheRepository(db)
	ctx := context.Background()

	mock.ExpectExists("test-key").SetVal(0)

	exists, err := repo.Exists(ctx, "test-key")
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCacheRepository_Ping(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewCacheRepository(db)
	ctx := context.Background()

	mock.ExpectPing().SetVal("PONG")

	err := repo.Ping(ctx)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
