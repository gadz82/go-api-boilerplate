package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database configuration
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string

	// Redis configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Cache configuration
	CacheDir string

	// Logging configuration
	// LoggingLevel defines the verbosity of logs:
	// 1 = Error, 2 = Warn, 3 = Info, 4 = Debug
	LoggingLevel int
}

func LoadConfig() *Config {
	// Try to load .env file but don't fail if it's missing
	_ = godotenv.Load()

	return &Config{
		// Database
		DBUser: getEnv("DB_USER", "root"),
		DBPass: getEnv("DB_PASS", "root"),
		DBHost: getEnv("DB_HOST", "127.0.0.1"),
		DBPort: getEnv("DB_PORT", "3306"),
		DBName: getEnv("DB_NAME", "test"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// Cache
		CacheDir: getEnv("CACHE_DIR", ".cache"),

		// Logging (default to 3=Info)
		LoggingLevel: getEnvInt("LOGGING_LEVEL", 3),
	}
}

func (c *Config) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

// GetRedisAddr returns the Redis address in host:port format.
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
