package logging

import (
	"bytes"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gadz82/go-api-boilerplate/internal/config"
)

// Log levels from least to most verbose
const (
	LevelError = 1
	LevelWarn  = 2
	LevelInfo  = 3
	LevelDebug = 4
)

// Logger defines the interface for the logging service
type Logger interface {
	Error(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	LogRequest(c *gin.Context)
}

// LoggingService is the concrete implementation of Logger
type LoggingService struct {
	level int
}

// NewLoggingService creates a new logging service with the configured log level
func NewLoggingService(cfg *config.Config) Logger {
	return &LoggingService{
		level: cfg.LoggingLevel,
	}
}

// Error logs error messages (level 1)
func (l *LoggingService) Error(format string, args ...interface{}) {
	if l.level >= LevelError {
		log.Printf("[ERROR] "+format, args...)
	}
}

// Warn logs warning messages (level 2)
func (l *LoggingService) Warn(format string, args ...interface{}) {
	if l.level >= LevelWarn {
		log.Printf("[WARN] "+format, args...)
	}
}

// Info logs info messages (level 3)
func (l *LoggingService) Info(format string, args ...interface{}) {
	if l.level >= LevelInfo {
		log.Printf("[INFO] "+format, args...)
	}
}

// Debug logs debug messages (level 4)
func (l *LoggingService) Debug(format string, args ...interface{}) {
	if l.level >= LevelDebug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// LogRequest logs the request body at debug level and restores the body for further reading
func (l *LoggingService) LogRequest(c *gin.Context) {
	if l.level >= LevelDebug {
		body, _ := io.ReadAll(c.Request.Body)
		log.Printf("[DEBUG] Request Body: %s", string(body))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}
}
