package items

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/gadz82/go-api-boilerplate/internal/domain"
	"github.com/gadz82/go-api-boilerplate/internal/service/logging"
)

// isValidUUID checks if a string is a valid UUID v4
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

type ItemHandler struct {
	Service   domain.ItemService
	Validator domain.Validator
	Logger    logging.Logger
}

func NewItemHandler(service domain.ItemService, validator domain.Validator, logger logging.Logger) *ItemHandler {
	return &ItemHandler{Service: service, Validator: validator, Logger: logger}
}

// GetAll gets all items
// @Summary      List items
// @Description  get items
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        include  query     string  false  "Include related resources (e.g. item_properties)"
// @Success      200  {object}  JSONAPIItemListResponse "Items"
// @Failure      500  {object}  map[string]string
// @Router       /v1/items [get]
func (h *ItemHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	if c.Query("include") == "item_properties" {
		ctx = context.WithValue(ctx, "include_properties", true)
	}

	items, err := h.Service.GetAllItems(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// GetByID gets an item by ID
// @Summary      Show an item
// @Description  get item by ID
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Item ID"
// @Param        include  query     string  false  "Include related resources (e.g. item_properties)"
// @Success      200  {object}  JSONAPIItemResponse "Item"
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/items/{id} [get]
func (h *ItemHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID format
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	ctx := c.Request.Context()
	if c.Query("include") == "item_properties" {
		ctx = context.WithValue(ctx, "include_properties", true)
	}

	item, err := h.Service.GetItemByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Create creates a new item
// @Summary      Create an item
// @Description  Create a new item (ID is auto-generated)
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        item  body      JSONAPIItem  true  "Item data"
// @Success      201   {object}  JSONAPIItemResponse "Created Item"
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /v1/items [post]
func (h *ItemHandler) Create(c *gin.Context) {
	h.Logger.LogRequest(c)

	item := new(domain.Item)
	if err := jsonapi.UnmarshalPayload(c.Request.Body, item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-generate UUID, ignoring any ID provided in the request
	item.ID = uuid.New().String()

	// Set CreatedAt to current timestamp
	now := time.Now()
	item.CreatedAt = &now

	// Validate the item using the injected validator
	if validationErrors := h.Validator.Validate(item); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	if err := h.Service.CreateItem(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Update updates an item
// @Summary      Update an item
// @Description  Update an item by ID (ID in request body is ignored, path parameter is used)
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "Item ID (UUID format)"
// @Param        item  body      JSONAPIItem true  "Item data"
// @Success      200   {object}  JSONAPIItemResponse "Updated Item"
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /v1/items/{id} [put]
func (h *ItemHandler) Update(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID format
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	h.Logger.LogRequest(c)

	item := new(domain.Item)
	if err := jsonapi.UnmarshalPayload(c.Request.Body, item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Use ID from path parameter, ignoring any ID in request body
	item.ID = id

	// Validate the item using the injected validator
	if validationErrors := h.Validator.Validate(item); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	if err := h.Service.UpdateItem(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *ItemHandler) Patch(c *gin.Context) {
	// For simplicity in this boilerplate, PATCH is handled similarly to Update
	h.Update(c)
}

// Delete deletes an item
// @Summary      Delete an item
// @Description  Delete an item by ID
// @Tags         items
// @Param        id   path      string  true  "Item ID (UUID format)"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/items/{id} [delete]
func (h *ItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID format
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := h.Service.DeleteItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
