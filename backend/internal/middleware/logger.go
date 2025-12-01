package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/easyhire/backend/internal/pkg/logger"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		
		// Set logger with request context
		log := logger.Global().WithRequestID(requestID)
		c.Set("logger", log)

		// Log request
		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Request started")

		// Process request
		c.Next()

		// Log response
		latency := time.Since(start)
		status := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		logEntry := log.Info().
			Int("status", status).
			Dur("latency", latency).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path)

		if errorMessage != "" {
			logEntry.Str("error", errorMessage)
		}

		if status >= 400 && status < 500 {
			logEntry.Msg("Client error")
		} else if status >= 500 {
			logEntry.Msg("Server error")
		} else {
			logEntry.Msg("Request completed")
		}
	}
}
