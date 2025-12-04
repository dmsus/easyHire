package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoggerMiddleware логирует запросы
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Str("ip", c.ClientIP()).
			Str("user-agent", c.Request.UserAgent()).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Msg("HTTP request")
	}
}

// RecoveryMiddleware обрабатывает паники
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Interface("error", err).Msg("Recovered from panic")
				c.AbortWithStatusJSON(500, gin.H{
					"error":   true,
					"message": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}

// CORSMiddleware настройки CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
