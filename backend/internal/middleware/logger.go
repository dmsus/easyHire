package middleware

import (
	"time"

	"github.com/easyhire/backend/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		log := logger.Global().With().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Int("status", c.Writer.Status()).
			Str("ip", c.ClientIP()).
			Str("user-agent", c.Request.UserAgent()).
			Dur("latency", latency).
			Logger()

		if len(c.Errors) > 0 {
			log.Error().Msg(c.Errors.String())
		} else {
			log.Info().Msg("Request completed")
		}
	}
}
