package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// RegisterHooks registers the HTTP server lifecycle hooks with the fx application.
// It starts the server on application start and logs when the server stops.
func RegisterHooks(lc fx.Lifecycle, r *gin.Engine) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Server starting on :8080")
			go func() {
				if err := r.Run(":8080"); err != nil {
					log.Printf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Server stopping")
			return nil
		},
	})
}
