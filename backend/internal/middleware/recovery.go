package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/easyhire/backend/internal/pkg/logger"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log := logger.Global().WithRequestID(c.GetString("request_id"))
				
				// Log the panic
				log.Error().
					Interface("error", err).
					Str("stack", string(debug.Stack())).
					Msg("PANIC recovered")

				// Return error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   true,
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "An internal server error occurred",
					"request_id": c.GetString("request_id"),
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
