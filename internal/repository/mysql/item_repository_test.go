package mysql

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/gadz82/go-api-boilerplate/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&domain.Item{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestItemRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := NewItemRepository(db)
	ctx := context.Background()
	uuidTest := uuid.New().String()
	// Create
	item := &domain.Item{ID: uuidTest, Title: "Test Item", Description: "Test Description"}
	err := repo.Create(ctx, item)
	assert.NoError(t, err)

	// GetByID
	found, err := repo.GetByID(ctx, uuidTest)
	assert.NoError(t, err)
	assert.Equal(t, item.Title, found.Title)

	// GetAll
	items, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, items, 1)

	// Update
	item.Title = "Updated Title"
	err = repo.Update(ctx, item)
	assert.NoError(t, err)

	updated, _ := repo.GetByID(ctx, uuidTest)
	assert.Equal(t, "Updated Title", updated.Title)

	// Delete
	err = repo.Delete(ctx, uuidTest)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, uuidTest)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound) || err != nil)
}
