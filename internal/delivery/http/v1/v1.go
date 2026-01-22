package v1

import (
	"github.com/gin-gonic/gin"
	items2 "github.com/gadz82/go-api-boilerplate/internal/delivery/handlers/items"
	"github.com/gadz82/go-api-boilerplate/internal/delivery/http/v1/items"
)

func RegisterRoutes(rg *gin.RouterGroup, itemHandler *items2.ItemHandler, itemPropertyHandler *items2.ItemPropertyHandler) {
	v1 := rg.Group("/v1")
	{
		items.RegisterRoutes(v1, itemHandler, itemPropertyHandler)
	}
}
