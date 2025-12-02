package middleware

import (
	"net/http"
	"strings"

	"github.com/easyhire/backend/internal/pkg/logger"
	"github.com/easyhire/internal/models"
	"github.com/easyhire/internal/pkg/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT токен и устанавливает user context
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Global().With().Str("middleware", "auth").Logger()

		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Error:   true,
				Code:    "UNAUTHORIZED",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Error:   true,
				Code:    "INVALID_TOKEN_FORMAT",
				Message: "Authorization header format must be 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid JWT token")
			c.JSON(http.StatusUnauthorized, models.APIError{
				Error:   true,
				Code:    "INVALID_TOKEN",
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_name", claims.Name)

		log.Debug().
			Str("user_id", claims.UserID.String()).
			Str("email", claims.Email).
			Str("role", string(claims.Role)).
			Str("path", c.Request.URL.Path).
			Msg("User authenticated")

		c.Next()
	}
}

// RoleMiddleware проверяет, что у пользователя есть одна из разрешенных ролей
func RoleMiddleware(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Error:   true,
				Code:    "UNAUTHORIZED",
				Message: "User not authenticated",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(models.UserRole)
		if !ok {
			c.JSON(http.StatusInternalServerError, models.APIError{
				Error:   true,
				Code:    "INTERNAL_ERROR",
				Message: "Invalid user role type",
			})
			c.Abort()
			return
		}

		// Check if user's role is in allowed roles
		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			logger.Global().Warn().
				Str("user_role", string(role)).
				Str("path", c.Request.URL.Path).
				Msg("Access denied: insufficient permissions")

			c.JSON(http.StatusForbidden, models.APIError{
				Error:   true,
				Code:    "FORBIDDEN",
				Message: "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnly middleware - только для администраторов
func AdminOnly() gin.HandlerFunc {
	return RoleMiddleware(models.RoleAdmin)
}

// HROnly middleware - только для HR
func HROnly() gin.HandlerFunc {
	return RoleMiddleware(models.RoleHR)
}

// ExpertOnly middleware - только для технических экспертов
func ExpertOnly() gin.HandlerFunc {
	return RoleMiddleware(models.RoleTechnicalExpert)
}

// HRorAdmin middleware - для HR и администраторов
func HRorAdmin() gin.HandlerFunc {
	return RoleMiddleware(models.RoleHR, models.RoleAdmin)
}

// OptionalAuth middleware - устанавливает контекст если токен есть, но не блокирует запрос
func OptionalAuth(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_name", claims.Name)

		c.Next()
	}
}
