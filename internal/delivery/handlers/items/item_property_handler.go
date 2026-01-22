package items

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/gadz82/go-api-boilerplate/internal/domain"
)

type ItemPropertyHandler struct {
	Service   domain.ItemPropertyService
	Validator domain.Validator
}

func NewItemPropertyHandler(service domain.ItemPropertyService, validator domain.Validator) *ItemPropertyHandler {
	return &ItemPropertyHandler{Service: service, Validator: validator}
}

// GetAll gets all item properties
// @Summary      List item properties
// @Description  get item properties for a specific item
// @Tags         item_properties
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Item ID (UUID format)"
// @Success      200  {object}  JSONAPIItemPropertyListResponse "Item Properties"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/items/{id}/properties [get]
func (h *ItemPropertyHandler) GetAll(c *gin.Context) {
	itemID := c.Param("id")

	// Validate UUID format for item ID
	if !isValidUUID(itemID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for item ID"})
		return
	}

	properties, err := h.Service.GetItemPropertiesByItemID(c.Request.Context(), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, properties); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// GetByID gets an item property by ID
// @Summary      Show an item property
// @Description  get item property by ID for a specific item
// @Tags         item_properties
// @Accept       json
// @Produce      json
// @Param        id           path      string  true  "Item ID (UUID format)"
// @Param        property_id  path      string  true  "Property ID (UUID format)"
// @Success      200          {object}  JSONAPIItemPropertyResponse "Item Property"
// @Failure      400          {object}  map[string]string
// @Failure      404          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /v1/items/{id}/properties/{property_id} [get]
func (h *ItemPropertyHandler) GetByID(c *gin.Context) {
	itemID := c.Param("id")
	id := c.Param("property_id")

	// Validate UUID format for item ID
	if !isValidUUID(itemID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for item ID"})
		return
	}

	// Validate UUID format for property ID
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for property ID"})
		return
	}

	property, err := h.Service.GetItemPropertyByID(c.Request.Context(), itemID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item property not found"})
		return
	}

	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Create creates a new item property
// @Summary      Create an item property
// @Description  Create a new item property for a specific item (ID is auto-generated)
// @Tags         item_properties
// @Accept       json
// @Produce      json
// @Param        id        path      string               true  "Item ID (UUID format)"
// @Param        property  body      JSONAPIItemProperty true  "Property data"
// @Success      201       {object}  JSONAPIItemPropertyResponse "Created Item Property"
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /v1/items/{id}/properties [post]
func (h *ItemPropertyHandler) Create(c *gin.Context) {
	itemID := c.Param("id")

	// Validate UUID format for item ID
	if !isValidUUID(itemID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for item ID"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	log.Printf("Request Body: %s", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	property := new(domain.ItemProperty)
	if err := jsonapi.UnmarshalPayload(c.Request.Body, property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-generate UUID for property ID, ignoring any ID provided in the request
	property.ID = uuid.New().String()
	property.ItemID = itemID

	// Validate the property using the injected validator
	if validationErrors := h.Validator.Validate(property); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	if err := h.Service.CreateItemProperty(c.Request.Context(), property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Update updates an item property
// @Summary      Update an item property
// @Description  Update an item property by ID for a specific item (ID in request body is ignored, path parameter is used)
// @Tags         item_properties
// @Accept       json
// @Produce      json
// @Param        id           path      string               true  "Item ID (UUID format)"
// @Param        property_id  path      string               true  "Property ID (UUID format)"
// @Param        property     body      JSONAPIItemProperty true  "Property data"
// @Success      200          {object}  JSONAPIItemPropertyResponse "Updated Item Property"
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /v1/items/{id}/properties/{property_id} [put]
func (h *ItemPropertyHandler) Update(c *gin.Context) {
	itemID := c.Param("id")
	id := c.Param("property_id")

	// Validate UUID format for item ID
	if !isValidUUID(itemID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for item ID"})
		return
	}

	// Validate UUID format for property ID
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for property ID"})
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	log.Printf("Request Body: %s", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	property := new(domain.ItemProperty)
	if err := jsonapi.UnmarshalPayload(c.Request.Body, property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Use IDs from path parameters, ignoring any IDs in request body
	property.ID = id
	property.ItemID = itemID

	// Validate the property using the injected validator
	if validationErrors := h.Validator.Validate(property); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	if err := h.Service.UpdateItemProperty(c.Request.Context(), property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *ItemPropertyHandler) Patch(c *gin.Context) {
	h.Update(c)
}

// Delete deletes an item property
// @Summary      Delete an item property
// @Description  Delete an item property by ID for a specific item
// @Tags         item_properties
// @Param        id           path      string  true  "Item ID (UUID format)"
// @Param        property_id  path      string  true  "Property ID (UUID format)"
// @Success      204          {object}  nil
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /v1/items/{id}/properties/{property_id} [delete]
func (h *ItemPropertyHandler) Delete(c *gin.Context) {
	itemID := c.Param("id")
	id := c.Param("property_id")

	// Validate UUID format for item ID
	if !isValidUUID(itemID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for item ID"})
		return
	}

	// Validate UUID format for property ID
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for property ID"})
		return
	}

	if err := h.Service.DeleteItemProperty(c.Request.Context(), itemID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
