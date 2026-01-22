package router

import (
	_ "github.com/gadz82/go-api-boilerplate/docs"
	"github.com/gadz82/go-api-boilerplate/internal/delivery/handlers/items"
	v1 "github.com/gadz82/go-api-boilerplate/internal/delivery/http/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(itemHandler *items.ItemHandler, itemPropertyHandler *items.ItemPropertyHandler) *gin.Engine {
	r := gin.Default()
	// Set to specific IPs like []string{"192.168.1.0/24"} if behind a known proxy
	err := r.SetTrustedProxies(nil)
	if err != nil {
		return nil
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		v1.RegisterRoutes(api, itemHandler, itemPropertyHandler)
	}

	return r
}
