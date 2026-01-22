package items

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gadz82/go-api-boilerplate/internal/domain"
	"github.com/gadz82/go-api-boilerplate/internal/service/logging"
	"github.com/gadz82/go-api-boilerplate/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

// newTestValidator returns the real validator for integration-style tests
func newTestValidator() domain.Validator {
	return validation.NewValidator()
}

// newTestLogger returns a mock logger for testing
func newTestLogger() logging.Logger {
	return &MockLogger{}
}

func TestItemHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemService)
	validator := newTestValidator()
	logger := newTestLogger()
	handler := NewItemHandler(svc, validator, logger)

	expectedItems := []*domain.Item{{ID: "1", Title: "Test"}}
	svc.On("GetAllItems", mock.Anything).Return(expectedItems, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/items", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestItemHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemService)
	validator := newTestValidator()
	logger := newTestLogger()
	handler := NewItemHandler(svc, validator, logger)

	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	item := &domain.Item{ID: testUUID, Title: "Test"}
	svc.On("GetItemByID", mock.Anything, testUUID).Return(item, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: testUUID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/"+testUUID, nil)

	handler.GetByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestItemHandler_GetByID_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemService)
	validator := newTestValidator()
	logger := newTestLogger()
	handler := NewItemHandler(svc, validator, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/invalid-uuid", nil)

	handler.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemService)
	validator := newTestValidator()
	logger := newTestLogger()
	handler := NewItemHandler(svc, validator, logger)

	testUUID := "550e8400-e29b-41d4-a716-446655440001"
	svc.On("GetItemByID", mock.Anything, testUUID).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: testUUID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/"+testUUID, nil)

	handler.GetByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestItemHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemService)
	validator := newTestValidator()
	logger := newTestLogger()
	handler := NewItemHandler(svc, validator, logger)

	item := &domain.Item{Title: "New Item"}
	svc.On("CreateItem", mock.Anything, mock.MatchedBy(func(i *domain.Item) bool {
		// Verify that UUID is auto-generated and title is preserved
		return i.Title == item.Title && i.ID != "" && isValidUUID(i.ID)
	})).Return(nil)

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, item)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/items", &buf)

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	svc.AssertExpectations(t)
}

func TestItemHandler_Create_IgnoresProvidedID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemService)
	validator := newTestValidator()
	logger := newTestLogger()
	handler := NewItemHandler(svc, validator, logger)

	// Item with a provided ID that should be ignored
	providedID := "550e8400-e29b-41d4-a716-446655440000"
	item := &domain.Item{ID: providedID, Title: "New Item"}
	svc.On("CreateItem", mock.Anything, mock.MatchedBy(func(i *domain.Item) bool {
		// Verify that the provided ID is ignored and a new UUID is generated
		return i.Title == item.Title && i.ID != providedID && isValidUUID(i.ID)
	})).Return(nil)

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, item)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/items", &buf)

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	svc.AssertExpectations(t)
}

func TestItemHandler_Create_MissingTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemService)
	validator := newTestValidator()
	logger := newTestLogger()
	handler := NewItemHandler(svc, validator, logger)

	// Item without required title
	item := &domain.Item{ID: "550e8400-e29b-41d4-a716-446655440000"}

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, item)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/items", &buf)

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemService)
	validator := newTestValidator()
	logger := newTestLogger()
	handler := NewItemHandler(svc, validator, logger)

	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	svc.On("DeleteItem", mock.Anything, testUUID).Return(nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: testUUID}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/items/"+testUUID, nil)

	r.DELETE("/items/:id", handler.Delete)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusNoContent, w.Code)
	svc.AssertExpectations(t)
}

func TestItemHandler_Delete_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemService)
	validator := newTestValidator()
	logger := newTestLogger()
	handler := NewItemHandler(svc, validator, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/items/invalid-uuid", nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
