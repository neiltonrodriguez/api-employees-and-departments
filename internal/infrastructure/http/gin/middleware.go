package ginapi

import (
	"fmt"
	"net/http"
	"time"

	"api-employees-and-departments/internal/infrastructure/logging"
	"api-employees-and-departments/internal/interfaces/api/dto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CORS middleware to handle Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger middleware for structured HTTP logging with Zap
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Get request ID from context
		requestID, _ := c.Get("RequestID")

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Build full path with query string
		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		// Build structured log fields
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", fullPath),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("body_size", c.Writer.Size()),
		}

		// Add request ID if available
		if requestID != nil {
			fields = append(fields, zap.String("request_id", requestID.(string)))
		}

		// Add error message if any
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// Log based on status code
		msg := "HTTP request completed"
		if statusCode >= 500 {
			logging.Error(msg, fields...)
		} else if statusCode >= 400 {
			logging.Warn(msg, fields...)
		} else {
			logging.Info(msg, fields...)
		}
	}
}

// Recovery middleware to handle panics with structured logging
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get request ID if available
				requestID, _ := c.Get("RequestID")

				// Build log fields
				fields := []zap.Field{
					zap.Any("panic", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
				}

				if requestID != nil {
					fields = append(fields, zap.String("request_id", requestID.(string)))
				}

				// Log the panic with stack trace
				logging.Error("Panic recovered", fields...)

				// Return error response
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "internal_server_error",
					Message: fmt.Sprintf("An unexpected error occurred: %v", err),
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Set("RequestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
