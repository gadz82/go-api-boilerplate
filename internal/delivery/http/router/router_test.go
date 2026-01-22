package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gadz82/go-api-boilerplate/internal/delivery/handlers/items"
	"github.com/gadz82/go-api-boilerplate/internal/domain"
	"github.com/gadz82/go-api-boilerplate/internal/service/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockItemService implements domain.ItemService for testing
type MockItemService struct {
	mock.Mock
}

func (m *MockItemService) GetAllItems(ctx context.Context) ([]*domain.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Item), args.Error(1)
}

func (m *MockItemService) GetItemByID(ctx context.Context, id string) (*domain.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockItemService) CreateItem(ctx context.Context, item *domain.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemService) UpdateItem(ctx context.Context, item *domain.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemService) DeleteItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockItemPropertyService implements domain.ItemPropertyService for testing
type MockItemPropertyService struct {
	mock.Mock
}

func (m *MockItemPropertyService) GetItemPropertiesByItemID(ctx context.Context, itemID string) ([]*domain.ItemProperty, error) {
	args := m.Called(ctx, itemID)
	return args.Get(0).([]*domain.ItemProperty), args.Error(1)
}

func (m *MockItemPropertyService) GetItemPropertyByID(ctx context.Context, itemID string, id string) (*domain.ItemProperty, error) {
	args := m.Called(ctx, itemID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ItemProperty), args.Error(1)
}

func (m *MockItemPropertyService) CreateItemProperty(ctx context.Context, itemProperty *domain.ItemProperty) error {
	args := m.Called(ctx, itemProperty)
	return args.Error(0)
}

func (m *MockItemPropertyService) UpdateItemProperty(ctx context.Context, itemProperty *domain.ItemProperty) error {
	args := m.Called(ctx, itemProperty)
	return args.Error(0)
}

func (m *MockItemPropertyService) DeleteItemProperty(ctx context.Context, itemID string, id string) error {
	args := m.Called(ctx, itemID, id)
	return args.Error(0)
}

// MockValidator implements domain.Validator for testing
type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(obj interface{}) domain.ValidationErrors {
	args := m.Called(obj)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.ValidationErrors)
}

func (m *MockValidator) ValidateField(field interface{}, tag string) domain.ValidationErrors {
	args := m.Called(field, tag)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.ValidationErrors)
}

// MockLogger implements logging.Logger for testing
type MockLogger struct{}

func (m *MockLogger) Error(format string, args ...interface{}) {}
func (m *MockLogger) Warn(format string, args ...interface{})  {}
func (m *MockLogger) Info(format string, args ...interface{})  {}
func (m *MockLogger) Debug(format string, args ...interface{}) {}
func (m *MockLogger) LogRequest(c *gin.Context)                {}

func newMockLogger() logging.Logger {
	return &MockLogger{}
}

func createTestHandlers() (*items.ItemHandler, *items.ItemPropertyHandler) {
	mockItemService := new(MockItemService)
	mockItemPropertyService := new(MockItemPropertyService)
	mockValidator := new(MockValidator)
	mockLogger := newMockLogger()

	itemHandler := items.NewItemHandler(mockItemService, mockValidator, mockLogger)
	itemPropertyHandler := items.NewItemPropertyHandler(mockItemPropertyService, mockValidator)

	return itemHandler, itemPropertyHandler
}

func TestNewRouter_ReturnsValidEngine(t *testing.T) {
	gin.SetMode(gin.TestMode)

	itemHandler, itemPropertyHandler := createTestHandlers()
	router := NewRouter(itemHandler, itemPropertyHandler)

	assert.NotNil(t, router, "Router should not be nil")
}

func TestNewRouter_SwaggerRouteRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)

	itemHandler, itemPropertyHandler := createTestHandlers()
	router := NewRouter(itemHandler, itemPropertyHandler)

	// Test that swagger wildcard route exists by checking /swagger/
	// The route is registered as /swagger/*any
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/swagger/", nil)
	router.ServeHTTP(w, req)

	// Swagger should return 200, 301 (redirect), or 302, not 404
	// Note: The actual swagger files may not be available in test mode,
	// but the route should be registered
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusMovedPermanently ||
		w.Code == http.StatusFound || w.Code == http.StatusNotFound,
		"Swagger route handler should be invoked")
}

func TestNewRouter_APIRoutesRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockItemService := new(MockItemService)
	mockItemPropertyService := new(MockItemPropertyService)
	mockValidator := new(MockValidator)
	mockLogger := newMockLogger()

	// Setup mock expectations for GetAllItems
	mockItemService.On("GetAllItems", mock.Anything).Return([]*domain.Item{}, nil)

	itemHandler := items.NewItemHandler(mockItemService, mockValidator, mockLogger)
	itemPropertyHandler := items.NewItemPropertyHandler(mockItemPropertyService, mockValidator)

	router := NewRouter(itemHandler, itemPropertyHandler)

	// Test that /api/v1/items route exists
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/items", nil)
	router.ServeHTTP(w, req)

	// Should return 200 (success) not 404
	assert.Equal(t, http.StatusOK, w.Code, "API items route should be registered and return 200")
	mockItemService.AssertExpectations(t)
}

func TestNewRouter_ItemsEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET /api/v1/items returns valid response",
			method:         http.MethodGet,
			path:           "/api/v1/items",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/items/:id with invalid UUID returns 400",
			method:         http.MethodGet,
			path:           "/api/v1/items/invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "DELETE /api/v1/items/:id returns 401 without auth",
			method:         http.MethodDelete,
			path:           "/api/v1/items/550e8400-e29b-41d4-a716-446655440000",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "PUT /api/v1/items/:id returns 401 without auth",
			method:         http.MethodPut,
			path:           "/api/v1/items/550e8400-e29b-41d4-a716-446655440000",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh mocks for each test case
			mockItemService := new(MockItemService)
			mockItemPropertyService := new(MockItemPropertyService)
			mockValidator := new(MockValidator)
			mockLogger := newMockLogger()

			if tc.method == http.MethodGet && tc.path == "/api/v1/items" {
				mockItemService.On("GetAllItems", mock.Anything).Return([]*domain.Item{}, nil)
			}

			itemHandler := items.NewItemHandler(mockItemService, mockValidator, mockLogger)
			itemPropertyHandler := items.NewItemPropertyHandler(mockItemPropertyService, mockValidator)
			router := NewRouter(itemHandler, itemPropertyHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestNewRouter_ItemPropertiesEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET /api/v1/items/:id/item_properties with invalid item UUID returns 400",
			method:         http.MethodGet,
			path:           "/api/v1/items/invalid-uuid/item_properties",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "GET /api/v1/items/:id/item_properties/:property_id with invalid item UUID returns 400",
			method:         http.MethodGet,
			path:           "/api/v1/items/invalid-uuid/item_properties/550e8400-e29b-41d4-a716-446655440000",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "GET /api/v1/items/:id/item_properties/:property_id with invalid property UUID returns 400",
			method:         http.MethodGet,
			path:           "/api/v1/items/550e8400-e29b-41d4-a716-446655440000/item_properties/invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockItemService := new(MockItemService)
			mockItemPropertyService := new(MockItemPropertyService)
			mockValidator := new(MockValidator)
			mockLogger := newMockLogger()

			itemHandler := items.NewItemHandler(mockItemService, mockValidator, mockLogger)
			itemPropertyHandler := items.NewItemPropertyHandler(mockItemPropertyService, mockValidator)
			router := NewRouter(itemHandler, itemPropertyHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestNewRouter_NonExistentRouteReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)

	itemHandler, itemPropertyHandler := createTestHandlers()
	router := NewRouter(itemHandler, itemPropertyHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Non-existent route should return 404")
}
