package items

import (
	"github.com/gin-gonic/gin"
	"github.com/gadz82/go-api-boilerplate/internal/delivery/handlers/items"
	"github.com/gadz82/go-api-boilerplate/internal/delivery/http/middleware"
	"github.com/gadz82/go-api-boilerplate/internal/delivery/http/v1/items/items_properties"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *items.ItemHandler, propertyHandler *items.ItemPropertyHandler) {
	itemGroup := rg.Group("/items")
	{
		// Public routes
		itemGroup.GET("", handler.GetAll)
		itemGroup.GET("/:id", handler.GetByID)
		itemGroup.POST("", handler.Create)

		// Nested property routes
		items_properties.RegisterRoutes(itemGroup, propertyHandler)

		// Authenticated routes
		authorized := itemGroup.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			authorized.PUT("/:id", handler.Update)
			authorized.PATCH("/:id", handler.Patch)
			authorized.DELETE("/:id", handler.Delete)
		}
	}
}
