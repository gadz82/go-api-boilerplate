package di

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/gadz82/go-api-boilerplate/internal/config"
	"github.com/gadz82/go-api-boilerplate/internal/database"
	"github.com/gadz82/go-api-boilerplate/internal/delivery/handlers/items"
	"github.com/gadz82/go-api-boilerplate/internal/delivery/http/router"
	"github.com/gadz82/go-api-boilerplate/internal/domain"
	fileRepo "github.com/gadz82/go-api-boilerplate/internal/repository/file"
	repoMysql "github.com/gadz82/go-api-boilerplate/internal/repository/mysql"
	redisRepo "github.com/gadz82/go-api-boilerplate/internal/repository/redis"
	"github.com/gadz82/go-api-boilerplate/internal/server"
	items2 "github.com/gadz82/go-api-boilerplate/internal/service/items"
	"github.com/gadz82/go-api-boilerplate/internal/service/logging"
	"github.com/gadz82/go-api-boilerplate/internal/validation"
	"go.uber.org/fx"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewModule creates the main application module with all dependencies wired together.
func NewModule() fx.Option {
	return fx.Module("app",
		fx.Options(
			provideInfrastructure(),
			provideRepositories(),
			provideServices(),
			provideHandlers(),
			provideHTTP(),
		),
		fx.Invoke(server.RegisterHooks),
	)
}

// provideInfrastructure provides core infrastructure dependencies:
// configuration, database connection, validator, and logging service.
func provideInfrastructure() fx.Option {
	return fx.Provide(
		config.LoadConfig,
		NewGormDB,
		validation.NewValidator,
		logging.NewLoggingService,
	)
}

// provideRepositories provides all repository implementations.
func provideRepositories() fx.Option {
	return fx.Provide(
		repoMysql.NewItemRepository,
		repoMysql.NewItemPropertyRepository,
		NewCacheRepository,
	)
}

// provideServices provides all service layer implementations.
func provideServices() fx.Option {
	return fx.Provide(
		items2.NewItemService,
		items2.NewItemPropertyService,
	)
}

// provideHandlers provides all HTTP handler implementations.
func provideHandlers() fx.Option {
	return fx.Provide(
		items.NewItemHandler,
		items.NewItemPropertyHandler,
	)
}

// provideHTTP provides HTTP-related dependencies: router and middleware.
func provideHTTP() fx.Option {
	return fx.Provide(
		router.NewRouter,
	)
}

// NewGormDB creates a new GORM database connection.
// It attempts to connect to MySQL first, falling back to SQLite for demo purposes.
// Migrations are handled by Goose instead of AutoMigrate.
func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	var dialect string

	dsn := cfg.GetMySQLDSN()
	db, err := gorm.Open(mysqlDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to MySQL: %v. Falling back to SQLite for demo.", err)
		db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		dialect = "sqlite3"
	} else {
		dialect = "mysql"
	}

	// Run Goose migrations
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	migrator := database.NewMigrator(sqlDB, dialect)
	if err := migrator.Up(); err != nil {
		return nil, err
	}
	log.Printf("Database migrations completed successfully (dialect: %s)", dialect)

	return db, nil
}

// NewCacheRepository creates a cache repository.
// It attempts to connect to Redis first, falling back to file-based cache if Redis is unavailable.
func NewCacheRepository(cfg *config.Config) (domain.CacheRepository, error) {
	// Try Redis first
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	// Test Redis connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Redis at %s: %v. Falling back to file-based cache.", cfg.GetRedisAddr(), err)

		// Fall back to file-based cache
		fileCache, err := fileRepo.NewCacheRepository(cfg.CacheDir)
		if err != nil {
			return nil, err
		}
		log.Printf("Using file-based cache in directory: %s", cfg.CacheDir)
		return fileCache, nil
	}

	log.Printf("Connected to Redis at %s", cfg.GetRedisAddr())
	return redisRepo.NewCacheRepository(redisClient), nil
}
