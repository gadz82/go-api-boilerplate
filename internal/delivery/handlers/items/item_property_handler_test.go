package items

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/gadz82/go-api-boilerplate/internal/domain"
)

// MockItemPropertyService implements domain.ItemPropertyService for testing
type MockItemPropertyService struct {
	mock.Mock
}

func (m *MockItemPropertyService) GetItemPropertiesByItemID(ctx context.Context, itemID string) ([]*domain.ItemProperty, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

// Test GetAll
func TestItemPropertyHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	expectedProperties := []*domain.ItemProperty{
		{ID: "550e8400-e29b-41d4-a716-446655440001", ItemID: itemID, Name: "color", Value: "red"},
	}
	svc.On("GetItemPropertiesByItemID", mock.Anything, itemID).Return(expectedProperties, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: itemID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/"+itemID+"/properties", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestItemPropertyHandler_GetAll_InvalidItemUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/invalid-uuid/properties", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_GetAll_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	svc.On("GetItemPropertiesByItemID", mock.Anything, itemID).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: itemID}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/"+itemID+"/properties", nil)

	handler.GetAll(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// Test GetByID
func TestItemPropertyHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	propertyID := "550e8400-e29b-41d4-a716-446655440001"
	property := &domain.ItemProperty{ID: propertyID, ItemID: itemID, Name: "color", Value: "red"}
	svc.On("GetItemPropertyByID", mock.Anything, itemID, propertyID).Return(property, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: itemID},
		{Key: "property_id", Value: propertyID},
	}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/"+itemID+"/properties/"+propertyID, nil)

	handler.GetByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestItemPropertyHandler_GetByID_InvalidItemUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: "invalid-uuid"},
		{Key: "property_id", Value: "550e8400-e29b-41d4-a716-446655440001"},
	}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/invalid-uuid/properties/550e8400-e29b-41d4-a716-446655440001", nil)

	handler.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_GetByID_InvalidPropertyUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"},
		{Key: "property_id", Value: "invalid-uuid"},
	}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/550e8400-e29b-41d4-a716-446655440000/properties/invalid-uuid", nil)

	handler.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	propertyID := "550e8400-e29b-41d4-a716-446655440001"
	svc.On("GetItemPropertyByID", mock.Anything, itemID, propertyID).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: itemID},
		{Key: "property_id", Value: propertyID},
	}
	c.Request, _ = http.NewRequest(http.MethodGet, "/items/"+itemID+"/properties/"+propertyID, nil)

	handler.GetByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	svc.AssertExpectations(t)
}

// Test Create
func TestItemPropertyHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	property := &domain.ItemProperty{Name: "color", Value: "red"}
	svc.On("CreateItemProperty", mock.Anything, mock.MatchedBy(func(p *domain.ItemProperty) bool {
		return p.Name == property.Name && p.Value == property.Value && p.ItemID == itemID && isValidUUID(p.ID)
	})).Return(nil)

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: itemID}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/items/"+itemID+"/properties", &buf)

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	svc.AssertExpectations(t)
}

func TestItemPropertyHandler_Create_InvalidItemUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	property := &domain.ItemProperty{Name: "color", Value: "red"}

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/items/invalid-uuid/properties", &buf)

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_Create_MissingName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	property := &domain.ItemProperty{Value: "red"} // Missing Name

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: itemID}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/items/"+itemID+"/properties", &buf)

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_Create_MissingValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	property := &domain.ItemProperty{Name: "color"} // Missing Value

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: itemID}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/items/"+itemID+"/properties", &buf)

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_Create_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	property := &domain.ItemProperty{Name: "color", Value: "red"}
	svc.On("CreateItemProperty", mock.Anything, mock.Anything).Return(errors.New("database error"))

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: itemID}}
	c.Request, _ = http.NewRequest(http.MethodPost, "/items/"+itemID+"/properties", &buf)

	handler.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// Test Update
func TestItemPropertyHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	propertyID := "550e8400-e29b-41d4-a716-446655440001"
	property := &domain.ItemProperty{Name: "color", Value: "blue"}
	svc.On("UpdateItemProperty", mock.Anything, mock.MatchedBy(func(p *domain.ItemProperty) bool {
		return p.Name == property.Name && p.Value == property.Value && p.ItemID == itemID && p.ID == propertyID
	})).Return(nil)

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: itemID},
		{Key: "property_id", Value: propertyID},
	}
	c.Request, _ = http.NewRequest(http.MethodPut, "/items/"+itemID+"/properties/"+propertyID, &buf)

	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestItemPropertyHandler_Update_InvalidItemUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	property := &domain.ItemProperty{Name: "color", Value: "blue"}

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: "invalid-uuid"},
		{Key: "property_id", Value: "550e8400-e29b-41d4-a716-446655440001"},
	}
	c.Request, _ = http.NewRequest(http.MethodPut, "/items/invalid-uuid/properties/550e8400-e29b-41d4-a716-446655440001", &buf)

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_Update_InvalidPropertyUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	property := &domain.ItemProperty{Name: "color", Value: "blue"}

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"},
		{Key: "property_id", Value: "invalid-uuid"},
	}
	c.Request, _ = http.NewRequest(http.MethodPut, "/items/550e8400-e29b-41d4-a716-446655440000/properties/invalid-uuid", &buf)

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_Update_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	propertyID := "550e8400-e29b-41d4-a716-446655440001"
	property := &domain.ItemProperty{Name: "", Value: ""} // Missing required fields

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: itemID},
		{Key: "property_id", Value: propertyID},
	}
	c.Request, _ = http.NewRequest(http.MethodPut, "/items/"+itemID+"/properties/"+propertyID, &buf)

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_Update_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	propertyID := "550e8400-e29b-41d4-a716-446655440001"
	property := &domain.ItemProperty{Name: "color", Value: "blue"}
	svc.On("UpdateItemProperty", mock.Anything, mock.Anything).Return(errors.New("database error"))

	var buf bytes.Buffer
	err := jsonapi.MarshalPayload(&buf, property)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: itemID},
		{Key: "property_id", Value: propertyID},
	}
	c.Request, _ = http.NewRequest(http.MethodPut, "/items/"+itemID+"/properties/"+propertyID, &buf)

	handler.Update(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

// Test Delete
func TestItemPropertyHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	propertyID := "550e8400-e29b-41d4-a716-446655440001"
	svc.On("DeleteItemProperty", mock.Anything, itemID, propertyID).Return(nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: itemID},
		{Key: "property_id", Value: propertyID},
	}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/items/"+itemID+"/properties/"+propertyID, nil)

	r.DELETE("/items/:id/properties/:property_id", handler.Delete)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusNoContent, w.Code)
	svc.AssertExpectations(t)
}

func TestItemPropertyHandler_Delete_InvalidItemUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: "invalid-uuid"},
		{Key: "property_id", Value: "550e8400-e29b-41d4-a716-446655440001"},
	}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/items/invalid-uuid/properties/550e8400-e29b-41d4-a716-446655440001", nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_Delete_InvalidPropertyUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"},
		{Key: "property_id", Value: "invalid-uuid"},
	}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/items/550e8400-e29b-41d4-a716-446655440000/properties/invalid-uuid", nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemPropertyHandler_Delete_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(MockItemPropertyService)
	validator := newTestValidator()
	handler := NewItemPropertyHandler(svc, validator)

	itemID := "550e8400-e29b-41d4-a716-446655440000"
	propertyID := "550e8400-e29b-41d4-a716-446655440001"
	svc.On("DeleteItemProperty", mock.Anything, itemID, propertyID).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "id", Value: itemID},
		{Key: "property_id", Value: propertyID},
	}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/items/"+itemID+"/properties/"+propertyID, nil)

	handler.Delete(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
