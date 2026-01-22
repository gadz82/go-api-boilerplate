# Go API Boilerplate

A personal exercise project demonstrating a clean, production-ready Go REST API boilerplate with modern best practices, clean architecture, and comprehensive tooling.

## Table of Contents

- [Features](#features)
- [Software Architecture](#software-architecture)
- [Project Structure](#project-structure)
- [Libraries Used](#libraries-used)
- [Design Patterns](#design-patterns)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Database Migrations](#database-migrations)
- [Swagger Documentation](#swagger-documentation)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)

## Features

- ✅ Clean Architecture with clear separation of concerns
- ✅ Dependency Injection using Uber's fx
- ✅ RESTful API following JSON:API specification
- ✅ Database support (MySQL with SQLite fallback)
- ✅ Database migrations with Goose
- ✅ Caching layer (Redis with file-based fallback)
- ✅ Request validation
- ✅ Swagger/OpenAPI documentation
- ✅ Configurable logging with severity levels
- ✅ Authentication middleware
- ✅ Comprehensive test coverage

## Software Architecture

This project follows **Clean Architecture** principles, organizing code into distinct layers with clear boundaries and dependencies flowing inward.

### Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer (HTTP)                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Handlers  │  │   Router    │  │     Middleware      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                      Service Layer                           │
│  ┌─────────────────────────┐  ┌───────────────────────────┐ │
│  │     Item Service        │  │  ItemProperty Service     │ │
│  └─────────────────────────┘  └───────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                      Domain Layer                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Entities  │  │ Interfaces  │  │    Validators       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                   Infrastructure Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ Repositories│  │   Database  │  │       Cache         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Responsibility |
|-------|----------------|
| **Delivery** | HTTP handlers, routing, middleware, request/response handling |
| **Service** | Business logic, orchestration, use case implementation |
| **Domain** | Entities, interfaces (ports), validation rules |
| **Infrastructure** | Database access, caching, external services |

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── docs/                        # Swagger generated documentation
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   ├── migrator.go          # Goose migration runner
│   │   └── migrations/          # SQL migration files
│   ├── delivery/
│   │   ├── handlers/            # HTTP request handlers
│   │   │   └── items/
│   │   │       ├── item_handler.go
│   │   │       ├── item_property_handler.go
│   │   │       └── schemas.go   # Swagger schema definitions
│   │   └── http/
│   │       ├── middleware/      # HTTP middleware (auth, etc.)
│   │       ├── router/          # Gin router setup
│   │       └── v1/              # API version 1 routes
│   ├── di/
│   │   └── di.go                # Dependency injection container
│   ├── domain/
│   │   ├── item.go              # Item entity and interfaces
│   │   ├── item_property.go     # ItemProperty entity and interfaces
│   │   ├── cache.go             # Cache interface
│   │   └── validator.go         # Validator interface
│   ├── repository/
│   │   ├── mysql/               # MySQL/SQLite implementations
│   │   ├── redis/               # Redis cache implementation
│   │   └── file/                # File-based cache implementation
│   ├── server/
│   │   └── server.go            # HTTP server lifecycle
│   ├── service/
│   │   ├── items/               # Item business logic
│   │   └── logging/             # Logging service
│   └── validation/
│       └── validator.go         # Validation implementation
├── pkg/
│   └── utils/                   # Shared utilities
├── vendor/                      # Vendored dependencies
├── go.mod
├── go.sum
└── README.md
```

## Libraries Used

### Core Framework
| Library | Purpose |
|---------|---------|
| [Gin](https://github.com/gin-gonic/gin) | High-performance HTTP web framework |
| [Uber fx](https://github.com/uber-go/fx) | Dependency injection framework |

### Database & ORM
| Library | Purpose |
|---------|---------|
| [GORM](https://gorm.io/) | ORM library for Go |
| [gorm/mysql](https://gorm.io/driver/mysql) | MySQL driver for GORM |
| [gorm/sqlite](https://gorm.io/driver/sqlite) | SQLite driver for GORM (fallback) |
| [Goose](https://github.com/pressly/goose) | Database migration tool |

### Caching
| Library | Purpose |
|---------|---------|
| [go-redis](https://github.com/redis/go-redis) | Redis client for Go |

### API & Documentation
| Library | Purpose |
|---------|---------|
| [jsonapi](https://github.com/google/jsonapi) | JSON:API specification implementation |
| [swaggo/swag](https://github.com/swaggo/swag) | Swagger documentation generator |
| [gin-swagger](https://github.com/swaggo/gin-swagger) | Swagger UI middleware for Gin |

### Validation & Utilities
| Library | Purpose |
|---------|---------|
| [validator](https://github.com/go-playground/validator) | Struct and field validation |
| [uuid](https://github.com/google/uuid) | UUID generation |
| [godotenv](https://github.com/joho/godotenv) | Environment variable loading |

### Testing
| Library | Purpose |
|---------|---------|
| [testify](https://github.com/stretchr/testify) | Testing toolkit (assertions, mocks) |
| [redismock](https://github.com/go-redis/redismock) | Redis mocking for tests |

## Design Patterns

### 1. Repository Pattern
Abstracts data access logic behind interfaces, allowing easy swapping of implementations (MySQL, SQLite, Redis, File).

```go
type ItemRepository interface {
    GetAll(ctx context.Context) ([]*Item, error)
    GetByID(ctx context.Context, id string) (*Item, error)
    Create(ctx context.Context, item *Item) error
    Update(ctx context.Context, item *Item) error
    Delete(ctx context.Context, id string) error
}
```

### 2. Dependency Injection
Using Uber's fx for constructor-based dependency injection:

```go
fx.Provide(
    config.LoadConfig,
    NewGormDB,
    validation.NewValidator,
    logging.NewLoggingService,
)
```

### 3. Service Layer Pattern
Business logic is encapsulated in service structs that depend on repository interfaces:

```go
type ItemService struct {
    repo  domain.ItemRepository
    cache domain.CacheRepository
}
```

### 4. Middleware Pattern
Cross-cutting concerns (authentication, logging) are handled via Gin middleware:

```go
authorized.Use(middleware.AuthMiddleware())
```

### 5. Interface Segregation
Small, focused interfaces for each concern (Validator, Logger, Repository).

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MySQL (optional, SQLite used as fallback)
- Redis (optional, file-based cache used as fallback)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/gadz82/go-api-boilerplate.git
cd go-api-boilerplate
```

2. Install dependencies:
```bash
go mod download
# or use vendored dependencies
go mod vendor
```

3. Create a `.env` file (optional):
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`.

## Configuration

Configuration is managed via environment variables. Create a `.env` file or set them directly:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_USER` | MySQL username | `root` |
| `DB_PASS` | MySQL password | `root` |
| `DB_HOST` | MySQL host | `127.0.0.1` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_NAME` | Database name | `test` |
| `REDIS_HOST` | Redis host | `127.0.0.1` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | (empty) |
| `CACHE_DIR` | File cache directory | `.cache` |
| `LOGGING_LEVEL` | Log verbosity (1=Error, 2=Warn, 3=Info, 4=Debug) | `3` |

## Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for database migrations.

### Automatic Migrations

Migrations run automatically when the application starts. The migrator is embedded in the application and executes on startup.

### Manual Migration Commands

Install the Goose CLI:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

**For SQLite:**
```bash
# Apply all pending migrations
goose -dir internal/database/migrations sqlite3 ./gorm.db up

# Rollback the last migration
goose -dir internal/database/migrations sqlite3 ./gorm.db down

# Check migration status
goose -dir internal/database/migrations sqlite3 ./gorm.db status
```

**For MySQL:**
```bash
# Apply all pending migrations
goose -dir internal/database/migrations mysql "user:password@tcp(host:port)/dbname?parseTime=true" up

# Rollback the last migration
goose -dir internal/database/migrations mysql "user:password@tcp(host:port)/dbname?parseTime=true" down

# Check migration status
goose -dir internal/database/migrations mysql "user:password@tcp(host:port)/dbname?parseTime=true" status
```

### Creating New Migrations

```bash
goose -dir internal/database/migrations create migration_name sql
```

This creates a new migration file in `internal/database/migrations/` with the proper naming convention.

### Migration File Format

Migrations use Goose's SQL format with statement blocks for MySQL compatibility:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE example (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE example;
-- +goose StatementEnd
```

## Swagger Documentation

API documentation is auto-generated using [swaggo/swag](https://github.com/swaggo/swag).

### Viewing Documentation

Once the server is running, access Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

### Regenerating Documentation

1. Install swag CLI:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate documentation:
```bash
swag init -g cmd/server/main.go -o docs
```

This parses the annotations in your code and generates:
- `docs/docs.go` - Go file with embedded documentation
- `docs/swagger.json` - OpenAPI JSON specification
- `docs/swagger.yaml` - OpenAPI YAML specification

### Adding API Documentation

Add annotations to your handlers:

```go
// GetByID gets an item by ID
// @Summary      Show an item
// @Description  get item by ID
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Item ID"
// @Success      200  {object}  JSONAPIItemResponse
// @Failure      404  {object}  map[string]string
// @Router       /v1/items/{id} [get]
func (h *ItemHandler) GetByID(c *gin.Context) {
    // ...
}
```

## API Endpoints

### Items

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/items` | List all items | No |
| GET | `/api/v1/items/:id` | Get item by ID | No |
| POST | `/api/v1/items` | Create new item | No |
| PUT | `/api/v1/items/:id` | Update item | Yes |
| PATCH | `/api/v1/items/:id` | Partial update | Yes |
| DELETE | `/api/v1/items/:id` | Delete item | Yes |

### Item Properties

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/items/:id/item_properties` | List item properties | No |
| GET | `/api/v1/items/:id/item_properties/:property_id` | Get property by ID | No |
| POST | `/api/v1/items/:id/item_properties` | Create property | Yes |
| PUT | `/api/v1/items/:id/item_properties/:property_id` | Update property | Yes |
| PATCH | `/api/v1/items/:id/item_properties/:property_id` | Partial update | Yes |
| DELETE | `/api/v1/items/:id/item_properties/:property_id` | Delete property | Yes |

### Authentication

Protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer secret-token
```

### Including Related Resources

Use the `include` query parameter to fetch related resources:
```
GET /api/v1/items?include=item_properties
GET /api/v1/items/:id?include=item_properties
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test ./... -v
```

Run tests for a specific package:
```bash
go test ./internal/delivery/handlers/items/... -v
```

Run tests with coverage:
```bash
go test ./... -cover
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Contributing

This is a personal exercise project, but suggestions and improvements are welcome!

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
