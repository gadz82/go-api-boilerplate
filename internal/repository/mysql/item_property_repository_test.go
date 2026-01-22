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

func setupItemPropertyTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&domain.Item{}, &domain.ItemProperty{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestItemPropertyRepository_CRUD(t *testing.T) {
	db := setupItemPropertyTestDB(t)
	itemRepo := NewItemRepository(db)
	propertyRepo := NewItemPropertyRepository(db)
	ctx := context.Background()

	// First create an item to associate properties with
	itemID := uuid.New().String()
	item := &domain.Item{ID: itemID, Title: "Test Item", Description: "Test Description"}
	err := itemRepo.Create(ctx, item)
	assert.NoError(t, err)

	propertyID := uuid.New().String()

	// Create
	property := &domain.ItemProperty{
		ID:     propertyID,
		ItemID: itemID,
		Name:   "Test Property",
		Value:  "Test Value",
	}
	err = propertyRepo.Create(ctx, property)
	assert.NoError(t, err)

	// GetByID
	found, err := propertyRepo.GetByID(ctx, itemID, propertyID)
	assert.NoError(t, err)
	assert.Equal(t, property.Name, found.Name)
	assert.Equal(t, property.Value, found.Value)
	assert.Equal(t, itemID, found.ItemID)

	// GetAllByItemID
	properties, err := propertyRepo.GetAllByItemID(ctx, itemID)
	assert.NoError(t, err)
	assert.Len(t, properties, 1)
	assert.Equal(t, property.Name, properties[0].Name)

	// Update
	property.Name = "Updated Property Name"
	property.Value = "Updated Value"
	err = propertyRepo.Update(ctx, property)
	assert.NoError(t, err)

	updated, err := propertyRepo.GetByID(ctx, itemID, propertyID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Property Name", updated.Name)
	assert.Equal(t, "Updated Value", updated.Value)

	// Delete
	err = propertyRepo.Delete(ctx, itemID, propertyID)
	assert.NoError(t, err)

	_, err = propertyRepo.GetByID(ctx, itemID, propertyID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound) || err != nil)
}

func TestItemPropertyRepository_GetAllByItemID_Empty(t *testing.T) {
	db := setupItemPropertyTestDB(t)
	itemRepo := NewItemRepository(db)
	propertyRepo := NewItemPropertyRepository(db)
	ctx := context.Background()

	// Create an item without properties
	itemID := uuid.New().String()
	item := &domain.Item{ID: itemID, Title: "Test Item", Description: "Test Description"}
	err := itemRepo.Create(ctx, item)
	assert.NoError(t, err)

	// GetAllByItemID should return empty slice
	properties, err := propertyRepo.GetAllByItemID(ctx, itemID)
	assert.NoError(t, err)
	assert.Len(t, properties, 0)
}

func TestItemPropertyRepository_GetByID_NotFound(t *testing.T) {
	db := setupItemPropertyTestDB(t)
	propertyRepo := NewItemPropertyRepository(db)
	ctx := context.Background()

	// Try to get a non-existent property
	_, err := propertyRepo.GetByID(ctx, "non-existent-item", "non-existent-property")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestItemPropertyRepository_MultipleProperties(t *testing.T) {
	db := setupItemPropertyTestDB(t)
	itemRepo := NewItemRepository(db)
	propertyRepo := NewItemPropertyRepository(db)
	ctx := context.Background()

	// Create an item
	itemID := uuid.New().String()
	item := &domain.Item{ID: itemID, Title: "Test Item", Description: "Test Description"}
	err := itemRepo.Create(ctx, item)
	assert.NoError(t, err)

	// Create multiple properties for the same item
	property1 := &domain.ItemProperty{
		ID:     uuid.New().String(),
		ItemID: itemID,
		Name:   "Property 1",
		Value:  "Value 1",
	}
	property2 := &domain.ItemProperty{
		ID:     uuid.New().String(),
		ItemID: itemID,
		Name:   "Property 2",
		Value:  "Value 2",
	}
	property3 := &domain.ItemProperty{
		ID:     uuid.New().String(),
		ItemID: itemID,
		Name:   "Property 3",
		Value:  "Value 3",
	}

	err = propertyRepo.Create(ctx, property1)
	assert.NoError(t, err)
	err = propertyRepo.Create(ctx, property2)
	assert.NoError(t, err)
	err = propertyRepo.Create(ctx, property3)
	assert.NoError(t, err)

	// GetAllByItemID should return all 3 properties
	properties, err := propertyRepo.GetAllByItemID(ctx, itemID)
	assert.NoError(t, err)
	assert.Len(t, properties, 3)
}

func TestItemPropertyRepository_PropertiesIsolatedByItem(t *testing.T) {
	db := setupItemPropertyTestDB(t)
	itemRepo := NewItemRepository(db)
	propertyRepo := NewItemPropertyRepository(db)
	ctx := context.Background()

	// Create two items
	itemID1 := uuid.New().String()
	item1 := &domain.Item{ID: itemID1, Title: "Item 1", Description: "Description 1"}
	err := itemRepo.Create(ctx, item1)
	assert.NoError(t, err)

	itemID2 := uuid.New().String()
	item2 := &domain.Item{ID: itemID2, Title: "Item 2", Description: "Description 2"}
	err = itemRepo.Create(ctx, item2)
	assert.NoError(t, err)

	// Create properties for each item
	property1 := &domain.ItemProperty{
		ID:     uuid.New().String(),
		ItemID: itemID1,
		Name:   "Property for Item 1",
		Value:  "Value 1",
	}
	property2 := &domain.ItemProperty{
		ID:     uuid.New().String(),
		ItemID: itemID2,
		Name:   "Property for Item 2",
		Value:  "Value 2",
	}

	err = propertyRepo.Create(ctx, property1)
	assert.NoError(t, err)
	err = propertyRepo.Create(ctx, property2)
	assert.NoError(t, err)

	// GetAllByItemID should only return properties for the specific item
	propertiesItem1, err := propertyRepo.GetAllByItemID(ctx, itemID1)
	assert.NoError(t, err)
	assert.Len(t, propertiesItem1, 1)
	assert.Equal(t, "Property for Item 1", propertiesItem1[0].Name)

	propertiesItem2, err := propertyRepo.GetAllByItemID(ctx, itemID2)
	assert.NoError(t, err)
	assert.Len(t, propertiesItem2, 1)
	assert.Equal(t, "Property for Item 2", propertiesItem2[0].Name)
}
