package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SimpleAuthMiddleware - простой middleware для аутентификации (заглушка)
func SimpleAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Authorization header format must be 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		token := parts[1]
		
		// TODO: В реальном приложении здесь должна быть проверка JWT токена
		// Сейчас используем заглушку для тестирования
		
		// Пример: если токен "test-token", считаем это админом
		if token == "test-token" {
			c.Set("user_id", "test-user-id")
			c.Set("user_role", "admin")
			c.Next()
			return
		}
		
		// По умолчанию считаем это кандидатом для тестирования
		c.Set("user_id", "candidate-user-id")
		c.Set("user_role", "candidate")
		
		c.Next()
	}
}

// RequireRole middleware для проверки ролей
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Invalid user role type",
			})
			c.Abort()
			return
		}

		// Check if user's role is in allowed roles
		allowed := false
		for _, allowedRole := range roles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
