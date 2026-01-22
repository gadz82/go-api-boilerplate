package items_properties

import (
	"github.com/gin-gonic/gin"
	"github.com/gadz82/go-api-boilerplate/internal/delivery/handlers/items"
	"github.com/gadz82/go-api-boilerplate/internal/delivery/http/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, propertyHandler *items.ItemPropertyHandler) {
	properties := rg.Group("/:id/item_properties")
	{
		properties.GET("", propertyHandler.GetAll)
		properties.GET("/:property_id", propertyHandler.GetByID)
		authorized := properties.Group("/")
		authorized.Use(middleware.AuthMiddleware())
		{
			properties.POST("", propertyHandler.Create)
			properties.PUT("/:property_id", propertyHandler.Update)
			properties.PATCH("/:property_id", propertyHandler.Patch)
			properties.DELETE("/:property_id", propertyHandler.Delete)
		}
	}
}
