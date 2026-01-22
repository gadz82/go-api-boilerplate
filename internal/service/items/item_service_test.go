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

// MockItemRepository is a mock of ItemRepository
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) GetAll(ctx context.Context) ([]*domain.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Item), args.Error(1)
}

func (m *MockItemRepository) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockItemRepository) Create(ctx context.Context, item *domain.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) Update(ctx context.Context, item *domain.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCacheRepository is a mock of CacheRepository
type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestItemService_GetAllItems_CacheMiss(t *testing.T) {
	repo := new(MockItemRepository)
	cache := new(MockCacheRepository)
	svc := NewItemService(repo, cache)

	expectedItems := []*domain.Item{{ID: "1", Title: "Test"}}

	// Cache miss scenario
	cache.On("Get", mock.Anything, "items:list").Return("", errors.New("cache miss"))
	repo.On("GetAll", mock.Anything).Return(expectedItems, nil)
	cache.On("Set", mock.Anything, "items:list", mock.Anything, 5*time.Minute).Return(nil)

	items, err := svc.GetAllItems(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemService_GetAllItems_CacheHit(t *testing.T) {
	repo := new(MockItemRepository)
	cache := new(MockCacheRepository)
	svc := NewItemService(repo, cache)

	cachedJSON := `[{"ID":"1","Title":"Test","Description":"","ItemProperties":null}]`

	// Cache hit scenario - repo should NOT be called
	cache.On("Get", mock.Anything, "items:list").Return(cachedJSON, nil)

	items, err := svc.GetAllItems(context.Background())

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "1", items[0].ID)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetAll", mock.Anything)
}

func TestItemService_GetItemByID_CacheMiss(t *testing.T) {
	repo := new(MockItemRepository)
	cache := new(MockCacheRepository)
	svc := NewItemService(repo, cache)

	expectedItem := &domain.Item{ID: "1", Title: "Test"}

	// Cache miss scenario
	cache.On("Get", mock.Anything, "item:1").Return("", errors.New("cache miss"))
	repo.On("GetByID", mock.Anything, "1").Return(expectedItem, nil)
	cache.On("Set", mock.Anything, "item:1", mock.Anything, 5*time.Minute).Return(nil)

	item, err := svc.GetItemByID(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, expectedItem, item)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemService_GetItemByID_CacheHit(t *testing.T) {
	repo := new(MockItemRepository)
	cache := new(MockCacheRepository)
	svc := NewItemService(repo, cache)

	cachedJSON := `{"ID":"1","Title":"Test","Description":"","ItemProperties":null}`

	// Cache hit scenario - repo should NOT be called
	cache.On("Get", mock.Anything, "item:1").Return(cachedJSON, nil)

	item, err := svc.GetItemByID(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", item.ID)
	assert.Equal(t, "Test", item.Title)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestItemService_CreateItem(t *testing.T) {
	repo := new(MockItemRepository)
	cache := new(MockCacheRepository)
	svc := NewItemService(repo, cache)

	item := &domain.Item{Title: "New Item"}
	repo.On("Create", mock.Anything, item).Return(nil)
	// Cache invalidation for items list
	cache.On("Delete", mock.Anything, "items:list").Return(nil)

	err := svc.CreateItem(context.Background(), item)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemService_UpdateItem(t *testing.T) {
	repo := new(MockItemRepository)
	cache := new(MockCacheRepository)
	svc := NewItemService(repo, cache)

	item := &domain.Item{ID: "1", Title: "Updated"}
	repo.On("Update", mock.Anything, item).Return(nil)
	// Cache invalidation for single item and items list
	cache.On("Delete", mock.Anything, "item:1").Return(nil)
	cache.On("Delete", mock.Anything, "items:list").Return(nil)

	err := svc.UpdateItem(context.Background(), item)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemService_DeleteItem(t *testing.T) {
	repo := new(MockItemRepository)
	cache := new(MockCacheRepository)
	svc := NewItemService(repo, cache)

	repo.On("Delete", mock.Anything, "1").Return(nil)
	// Cache invalidation for single item and items list
	cache.On("Delete", mock.Anything, "item:1").Return(nil)
	cache.On("Delete", mock.Anything, "items:list").Return(nil)

	err := svc.DeleteItem(context.Background(), "1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestItemService_GetItemByID_Error(t *testing.T) {
	repo := new(MockItemRepository)
	cache := new(MockCacheRepository)
	svc := NewItemService(repo, cache)

	// Cache miss, then repo returns error
	cache.On("Get", mock.Anything, "item:1").Return("", errors.New("cache miss"))
	repo.On("GetByID", mock.Anything, "1").Return(nil, errors.New("not found"))

	item, err := svc.GetItemByID(context.Background(), "1")

	assert.Error(t, err)
	assert.Nil(t, item)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}
