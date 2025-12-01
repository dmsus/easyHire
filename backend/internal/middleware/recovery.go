package middleware

import (
	"github.com/easyhire/backend/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Global().Error().
					Interface("error", err).
					Str("path", c.Request.URL.Path).
					Msg("Recovered from panic")
				
				c.AbortWithStatus(500)
			}
		}()

		c.Next()
	}
}
