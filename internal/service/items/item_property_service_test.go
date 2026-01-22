package items

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/gadz82/go-api-boilerplate/internal/domain"
)

// MockItemPropertyRepository is a mock of ItemPropertyRepository
type MockItemPropertyRepository struct {
	mock.Mock
}

func (m *MockItemPropertyRepository) GetAllByItemID(ctx context.Context, itemID string) ([]*domain.ItemProperty, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ItemProperty), args.Error(1)
}

func (m *MockItemPropertyRepository) GetByID(ctx context.Context, itemID string, id string) (*domain.ItemProperty, error) {
	args := m.Called(ctx, itemID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ItemProperty), args.Error(1)
}

func (m *MockItemPropertyRepository) Create(ctx context.Context, itemProperty *domain.ItemProperty) error {
	args := m.Called(ctx, itemProperty)
	return args.Error(0)
}

func (m *MockItemPropertyRepository) Update(ctx context.Context, itemProperty *domain.ItemProperty) error {
	args := m.Called(ctx, itemProperty)
	return args.Error(0)
}

func (m *MockItemPropertyRepository) Delete(ctx context.Context, itemID string, id string) error {
	args := m.Called(ctx, itemID, id)
	return args.Error(0)
}

func TestItemPropertyService_GetItemPropertiesByItemID_CacheMiss(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	expectedProperties := []*domain.ItemProperty{
		{ID: "prop-1", ItemID: itemID, Name: "color", Value: "red"},
		{ID: "prop-2", ItemID: itemID, Name: "size", Value: "large"},
	}

	// Cache miss scenario
	cache.On("Get", mock.Anything, "item_properties:list:item-123").Return("", errors.New("cache miss"))
	repo.On("GetAllByItemID", mock.Anything, itemID).Return(expectedProperties, nil)
	cache.On("Set", mock.Anything, "item_properties:list:item-123", mock.Anything, 5*time.Minute).Return(nil)

	properties, err := svc.GetItemPropertiesByItemID(context.Background(), itemID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProperties, properties)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemPropertyService_GetItemPropertiesByItemID_CacheHit(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	cachedJSON := `[{"ID":"prop-1","ItemID":"item-123","Name":"color","Value":"red"}]`

	// Cache hit scenario - repo should NOT be called
	cache.On("Get", mock.Anything, "item_properties:list:item-123").Return(cachedJSON, nil)

	properties, err := svc.GetItemPropertiesByItemID(context.Background(), itemID)

	assert.NoError(t, err)
	assert.Len(t, properties, 1)
	assert.Equal(t, "prop-1", properties[0].ID)
	assert.Equal(t, "color", properties[0].Name)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetAllByItemID", mock.Anything, mock.Anything)
}

func TestItemPropertyService_GetItemPropertiesByItemID_RepoError(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"

	// Cache miss, then repo error
	cache.On("Get", mock.Anything, "item_properties:list:item-123").Return("", errors.New("cache miss"))
	repo.On("GetAllByItemID", mock.Anything, itemID).Return(nil, errors.New("database error"))

	properties, err := svc.GetItemPropertiesByItemID(context.Background(), itemID)

	assert.Error(t, err)
	assert.Nil(t, properties)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemPropertyService_GetItemPropertyByID_CacheMiss(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	propID := "prop-1"
	expectedProperty := &domain.ItemProperty{ID: propID, ItemID: itemID, Name: "color", Value: "red"}

	// Cache miss scenario
	cache.On("Get", mock.Anything, "item_property:item-123:prop-1").Return("", errors.New("cache miss"))
	repo.On("GetByID", mock.Anything, itemID, propID).Return(expectedProperty, nil)
	cache.On("Set", mock.Anything, "item_property:item-123:prop-1", mock.Anything, 5*time.Minute).Return(nil)

	property, err := svc.GetItemPropertyByID(context.Background(), itemID, propID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProperty, property)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemPropertyService_GetItemPropertyByID_CacheHit(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	propID := "prop-1"
	cachedJSON := `{"ID":"prop-1","ItemID":"item-123","Name":"color","Value":"red"}`

	// Cache hit scenario - repo should NOT be called
	cache.On("Get", mock.Anything, "item_property:item-123:prop-1").Return(cachedJSON, nil)

	property, err := svc.GetItemPropertyByID(context.Background(), itemID, propID)

	assert.NoError(t, err)
	assert.Equal(t, propID, property.ID)
	assert.Equal(t, "color", property.Name)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
}

func TestItemPropertyService_GetItemPropertyByID_NotFound(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	propID := "prop-nonexistent"

	// Cache miss, then repo returns not found
	cache.On("Get", mock.Anything, "item_property:item-123:prop-nonexistent").Return("", errors.New("cache miss"))
	repo.On("GetByID", mock.Anything, itemID, propID).Return(nil, errors.New("not found"))

	property, err := svc.GetItemPropertyByID(context.Background(), itemID, propID)

	assert.Error(t, err)
	assert.Nil(t, property)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemPropertyService_CreateItemProperty(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	property := &domain.ItemProperty{ID: "prop-1", ItemID: itemID, Name: "color", Value: "red"}

	repo.On("Create", mock.Anything, property).Return(nil)
	// Cache invalidation for item properties list
	cache.On("Delete", mock.Anything, "item_properties:list:item-123").Return(nil)

	err := svc.CreateItemProperty(context.Background(), property)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemPropertyService_CreateItemProperty_RepoError(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	property := &domain.ItemProperty{ID: "prop-1", ItemID: itemID, Name: "color", Value: "red"}

	repo.On("Create", mock.Anything, property).Return(errors.New("database error"))

	err := svc.CreateItemProperty(context.Background(), property)

	assert.Error(t, err)
	repo.AssertExpectations(t)
	// Cache should NOT be invalidated on error
	cache.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestItemPropertyService_UpdateItemProperty(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	propID := "prop-1"
	property := &domain.ItemProperty{ID: propID, ItemID: itemID, Name: "color", Value: "blue"}

	repo.On("Update", mock.Anything, property).Return(nil)
	// Cache invalidation for single property
	cache.On("Delete", mock.Anything, "item_property:item-123:prop-1").Return(nil)
	// Cache invalidation for item properties list
	cache.On("Delete", mock.Anything, "item_properties:list:item-123").Return(nil)

	err := svc.UpdateItemProperty(context.Background(), property)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemPropertyService_UpdateItemProperty_RepoError(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	propID := "prop-1"
	property := &domain.ItemProperty{ID: propID, ItemID: itemID, Name: "color", Value: "blue"}

	repo.On("Update", mock.Anything, property).Return(errors.New("database error"))

	err := svc.UpdateItemProperty(context.Background(), property)

	assert.Error(t, err)
	repo.AssertExpectations(t)
	// Cache should NOT be invalidated on error
	cache.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestItemPropertyService_DeleteItemProperty(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	propID := "prop-1"

	repo.On("Delete", mock.Anything, itemID, propID).Return(nil)
	// Cache invalidation for single property
	cache.On("Delete", mock.Anything, "item_property:item-123:prop-1").Return(nil)
	// Cache invalidation for item properties list
	cache.On("Delete", mock.Anything, "item_properties:list:item-123").Return(nil)

	err := svc.DeleteItemProperty(context.Background(), itemID, propID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemPropertyService_DeleteItemProperty_RepoError(t *testing.T) {
	repo := new(MockItemPropertyRepository)
	cache := new(MockCacheRepository)
	svc := NewItemPropertyService(repo, cache)

	itemID := "item-123"
	propID := "prop-1"

	repo.On("Delete", mock.Anything, itemID, propID).Return(errors.New("database error"))

	err := svc.DeleteItemProperty(context.Background(), itemID, propID)

	assert.Error(t, err)
	repo.AssertExpectations(t)
	// Cache should NOT be invalidated on error
	cache.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}
