package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear any existing env vars that might interfere
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("CACHE_DIR")
	os.Unsetenv("LOGGING_LEVEL")

	cfg := LoadConfig()

	assert.Equal(t, "root", cfg.DBUser)
	assert.Equal(t, "root", cfg.DBPass)
	assert.Equal(t, "127.0.0.1", cfg.DBHost)
	assert.Equal(t, "3306", cfg.DBPort)
	assert.Equal(t, "test", cfg.DBName)
	assert.Equal(t, "127.0.0.1", cfg.RedisHost)
	assert.Equal(t, "6379", cfg.RedisPort)
	assert.Equal(t, "", cfg.RedisPassword)
	assert.Equal(t, ".cache", cfg.CacheDir)
	assert.Equal(t, 3, cfg.LoggingLevel)
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Set custom env vars
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASS", "testpass")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3307")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("REDIS_HOST", "redis.local")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_PASSWORD", "redispass")
	os.Setenv("CACHE_DIR", "/tmp/cache")
	os.Setenv("LOGGING_LEVEL", "4")

	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("CACHE_DIR")
		os.Unsetenv("LOGGING_LEVEL")
	}()

	cfg := LoadConfig()

	assert.Equal(t, "testuser", cfg.DBUser)
	assert.Equal(t, "testpass", cfg.DBPass)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "3307", cfg.DBPort)
	assert.Equal(t, "testdb", cfg.DBName)
	assert.Equal(t, "redis.local", cfg.RedisHost)
	assert.Equal(t, "6380", cfg.RedisPort)
	assert.Equal(t, "redispass", cfg.RedisPassword)
	assert.Equal(t, "/tmp/cache", cfg.CacheDir)
	assert.Equal(t, 4, cfg.LoggingLevel)
}

func TestLoadConfig_InvalidLoggingLevel(t *testing.T) {
	os.Setenv("LOGGING_LEVEL", "invalid")
	defer os.Unsetenv("LOGGING_LEVEL")

	cfg := LoadConfig()

	// Should fall back to default when invalid
	assert.Equal(t, 3, cfg.LoggingLevel)
}

func TestConfig_GetMySQLDSN(t *testing.T) {
	cfg := &Config{
		DBUser: "myuser",
		DBPass: "mypass",
		DBHost: "myhost",
		DBPort: "3306",
		DBName: "mydb",
	}

	dsn := cfg.GetMySQLDSN()

	expected := "myuser:mypass@tcp(myhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	assert.Equal(t, expected, dsn)
}

func TestConfig_GetRedisAddr(t *testing.T) {
	cfg := &Config{
		RedisHost: "redis.example.com",
		RedisPort: "6379",
	}

	addr := cfg.GetRedisAddr()

	assert.Equal(t, "redis.example.com:6379", addr)
}

func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default")
	assert.Equal(t, "test_value", result)

	// Test with non-existing env var
	result = getEnv("NON_EXISTING_VAR", "default")
	assert.Equal(t, "default", result)
}

func TestGetEnvInt(t *testing.T) {
	// Test with valid int
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	result := getEnvInt("TEST_INT", 10)
	assert.Equal(t, 42, result)

	// Test with invalid int
	os.Setenv("TEST_INVALID_INT", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_INT")

	result = getEnvInt("TEST_INVALID_INT", 10)
	assert.Equal(t, 10, result)

	// Test with non-existing env var
	result = getEnvInt("NON_EXISTING_INT", 10)
	assert.Equal(t, 10, result)
}
