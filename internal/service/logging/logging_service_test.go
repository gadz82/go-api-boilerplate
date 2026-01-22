package logging

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gadz82/go-api-boilerplate/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewLoggingService(t *testing.T) {
	cfg := &config.Config{LoggingLevel: 3}
	logger := NewLoggingService(cfg)
	assert.NotNil(t, logger)
}

func TestLoggingService_Error(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected bool
	}{
		{"Level 1 - Error enabled", LevelError, true},
		{"Level 2 - Error enabled", LevelWarn, true},
		{"Level 3 - Error enabled", LevelInfo, true},
		{"Level 4 - Error enabled", LevelDebug, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			cfg := &config.Config{LoggingLevel: tt.level}
			logger := NewLoggingService(cfg)
			logger.Error("test error %s", "message")

			if tt.expected {
				assert.Contains(t, buf.String(), "[ERROR]")
				assert.Contains(t, buf.String(), "test error message")
			} else {
				assert.Empty(t, buf.String())
			}
		})
	}
}

func TestLoggingService_Warn(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected bool
	}{
		{"Level 1 - Warn disabled", LevelError, false},
		{"Level 2 - Warn enabled", LevelWarn, true},
		{"Level 3 - Warn enabled", LevelInfo, true},
		{"Level 4 - Warn enabled", LevelDebug, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			cfg := &config.Config{LoggingLevel: tt.level}
			logger := NewLoggingService(cfg)
			logger.Warn("test warn %s", "message")

			if tt.expected {
				assert.Contains(t, buf.String(), "[WARN]")
				assert.Contains(t, buf.String(), "test warn message")
			} else {
				assert.Empty(t, buf.String())
			}
		})
	}
}

func TestLoggingService_Info(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected bool
	}{
		{"Level 1 - Info disabled", LevelError, false},
		{"Level 2 - Info disabled", LevelWarn, false},
		{"Level 3 - Info enabled", LevelInfo, true},
		{"Level 4 - Info enabled", LevelDebug, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			cfg := &config.Config{LoggingLevel: tt.level}
			logger := NewLoggingService(cfg)
			logger.Info("test info %s", "message")

			if tt.expected {
				assert.Contains(t, buf.String(), "[INFO]")
				assert.Contains(t, buf.String(), "test info message")
			} else {
				assert.Empty(t, buf.String())
			}
		})
	}
}

func TestLoggingService_Debug(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected bool
	}{
		{"Level 1 - Debug disabled", LevelError, false},
		{"Level 2 - Debug disabled", LevelWarn, false},
		{"Level 3 - Debug disabled", LevelInfo, false},
		{"Level 4 - Debug enabled", LevelDebug, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			cfg := &config.Config{LoggingLevel: tt.level}
			logger := NewLoggingService(cfg)
			logger.Debug("test debug %s", "message")

			if tt.expected {
				assert.Contains(t, buf.String(), "[DEBUG]")
				assert.Contains(t, buf.String(), "test debug message")
			} else {
				assert.Empty(t, buf.String())
			}
		})
	}
}

func TestLoggingService_LogRequest_DebugEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	cfg := &config.Config{LoggingLevel: LevelDebug}
	logger := NewLoggingService(cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestBody := `{"data":{"type":"items","attributes":{"title":"Test"}}}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/items", strings.NewReader(requestBody))

	logger.LogRequest(c)

	// Check that the body was logged
	assert.Contains(t, buf.String(), "[DEBUG]")
	assert.Contains(t, buf.String(), "Request Body")
	assert.Contains(t, buf.String(), "Test")

	// Check that the body can still be read
	body, err := io.ReadAll(c.Request.Body)
	assert.NoError(t, err)
	assert.Equal(t, requestBody, string(body))
}

func TestLoggingService_LogRequest_DebugDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	cfg := &config.Config{LoggingLevel: LevelInfo} // Debug disabled
	logger := NewLoggingService(cfg)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestBody := `{"data":{"type":"items","attributes":{"title":"Test"}}}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/items", strings.NewReader(requestBody))

	logger.LogRequest(c)

	// Check that nothing was logged
	assert.Empty(t, buf.String())
}

func TestLogLevelConstants(t *testing.T) {
	assert.Equal(t, 1, LevelError)
	assert.Equal(t, 2, LevelWarn)
	assert.Equal(t, 3, LevelInfo)
	assert.Equal(t, 4, LevelDebug)
}
