package handlers

import (
	"net/http"
	"time"

	"github.com/easyhire/backend/internal/pkg/logger"
	"github.com/easyhire/backend/pkg/database"
	"github.com/easyhire/internal/models"
	"github.com/easyhire/internal/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db              *database.Database
	jwtService      *auth.JWTService
	passwordService *auth.PasswordService
}

func NewAuthHandler(db *database.Database, jwtService *auth.JWTService, passwordService *auth.PasswordService) *AuthHandler {
	return &AuthHandler{
		db:              db,
		jwtService:      jwtService,
		passwordService: passwordService,
	}
}

// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.LoginResponse
// @Failure 400 {object} models.APIError
// @Failure 409 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	log := logger.Global().With().Str("handler", "auth").Str("method", "register").Logger()

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid registration request")
		c.JSON(http.StatusBadRequest, models.APIError{
			Error:   true,
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request data",
			Details: err.Error(),
		})
		return
	}

	// Validate password strength
	if strengthErrors := h.passwordService.ValidateStrength(req.Password); len(strengthErrors) > 0 {
		c.JSON(http.StatusBadRequest, models.APIError{
			Error:   true,
			Code:    "PASSWORD_WEAK",
			Message: "Password does not meet requirements",
			Details: strengthErrors,
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	result := h.db.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		log.Warn().Str("email", req.Email).Msg("User already exists")
		c.JSON(http.StatusConflict, models.APIError{
			Error:   true,
			Code:    "USER_EXISTS",
			Message: "User with this email already exists",
		})
		return
	}

	// Hash password
	passwordHash, err := h.passwordService.Hash(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		c.JSON(http.StatusInternalServerError, models.APIError{
			Error:   true,
			Code:    "INTERNAL_ERROR",
			Message: "Failed to process registration",
		})
		return
	}

	// Set role - default to candidate unless specified
	role := models.RoleCandidate
	if req.Role != "" {
		switch req.Role {
		case string(models.RoleAdmin):
			role = models.RoleAdmin
		case string(models.RoleHR):
			role = models.RoleHR
		case string(models.RoleTechnicalExpert):
			role = models.RoleTechnicalExpert
		default:
			role = models.RoleCandidate
		}
	}

	// Create user
	user := models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Role:         role,
		Company:      &req.Company,
		IsActive:     true,
	}

	if err := h.db.DB.Create(&user).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		c.JSON(http.StatusInternalServerError, models.APIError{
			Error:   true,
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create user account",
		})
		return
	}

	// Generate tokens
	tokenPair, err := h.jwtService.GenerateTokenPair(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate tokens")
		c.JSON(http.StatusInternalServerError, models.APIError{
			Error:   true,
			Code:    "INTERNAL_ERROR",
			Message: "Failed to generate authentication tokens",
		})
		return
	}

	// Update last login
	h.db.DB.Model(&user).Update("last_login_at", time.Now())

	log.Info().Str("email", user.Email).Str("role", string(user.Role)).Msg("User registered successfully")

	response := models.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    tokenPair.ExpiresAt,
		User: models.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			Company:   user.Company,
			AvatarURL: user.AvatarURL,
			IsActive:  user.IsActive,
		},
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.APIError
// @Failure 401 {object} models.APIError
// @Failure 404 {object} models.APIError
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	log := logger.Global().With().Str("handler", "auth").Str("method", "login").Logger()

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid login request")
		c.JSON(http.StatusBadRequest, models.APIError{
			Error:   true,
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request data",
		})
		return
	}

	// Find user
	var user models.User
	result := h.db.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Warn().Str("email", req.Email).Msg("User not found")
			c.JSON(http.StatusUnauthorized, models.APIError{
				Error:   true,
				Code:    "INVALID_CREDENTIALS",
				Message: "Invalid email or password",
			})
			return
		}

		log.Error().Err(result.Error).Msg("Database error")
		c.JSON(http.StatusInternalServerError, models.APIError{
			Error:   true,
			Code:    "INTERNAL_ERROR",
			Message: "Failed to authenticate",
		})
		return
	}

	// Check if user is active
	if !user.IsActive {
		log.Warn().Str("email", user.Email).Msg("Inactive user tried to login")
		c.JSON(http.StatusUnauthorized, models.APIError{
			Error:   true,
			Code:    "ACCOUNT_DISABLED",
			Message: "Account is disabled",
		})
		return
	}

	// Verify password
	if err := h.passwordService.Compare(user.PasswordHash, req.Password); err != nil {
		log.Warn().Str("email", user.Email).Msg("Invalid password")
		c.JSON(http.StatusUnauthorized, models.APIError{
			Error:   true,
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
		})
		return
	}

	// Generate tokens
	tokenPair, err := h.jwtService.GenerateTokenPair(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate tokens")
		c.JSON(http.StatusInternalServerError, models.APIError{
			Error:   true,
			Code:    "INTERNAL_ERROR",
			Message: "Failed to generate authentication tokens",
		})
		return
	}

	// Update last login
	h.db.DB.Model(&user).Update("last_login_at", time.Now())

	log.Info().Str("email", user.Email).Str("role", string(user.Role)).Msg("User logged in successfully")

	response := models.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    tokenPair.ExpiresAt,
		User: models.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			Company:   user.Company,
			AvatarURL: user.AvatarURL,
			IsActive:  user.IsActive,
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Refresh access token
// @Description Refresh expired access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.APIError
// @Failure 401 {object} models.APIError
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	log := logger.Global().With().Str("handler", "auth").Str("method", "refresh").Logger()

	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid refresh request")
		c.JSON(http.StatusBadRequest, models.APIError{
			Error:   true,
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request data",
		})
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid refresh token")
		c.JSON(http.StatusUnauthorized, models.APIError{
			Error:   true,
			Code:    "INVALID_TOKEN",
			Message: "Invalid or expired refresh token",
		})
		return
	}

	// Get user from database
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Error().Err(err).Str("subject", claims.Subject).Msg("Invalid user ID in token")
		c.JSON(http.StatusUnauthorized, models.APIError{
			Error:   true,
			Code:    "INVALID_TOKEN",
			Message: "Invalid token payload",
		})
		return
	}

	var user models.User
	result := h.db.DB.First(&user, "id = ?", userID)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("user_id", userID.String()).Msg("User not found")
		c.JSON(http.StatusUnauthorized, models.APIError{
			Error:   true,
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	// Check if user is active
	if !user.IsActive {
		log.Warn().Str("email", user.Email).Msg("Inactive user tried to refresh token")
		c.JSON(http.StatusUnauthorized, models.APIError{
			Error:   true,
			Code:    "ACCOUNT_DISABLED",
			Message: "Account is disabled",
		})
		return
	}

	// Generate new token pair
	tokenPair, err := h.jwtService.GenerateTokenPair(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate new tokens")
		c.JSON(http.StatusInternalServerError, models.APIError{
			Error:   true,
			Code:    "INTERNAL_ERROR",
			Message: "Failed to generate authentication tokens",
		})
		return
	}

	log.Info().Str("email", user.Email).Msg("Token refreshed successfully")

	response := models.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    tokenPair.ExpiresAt,
		User: models.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			Company:   user.Company,
			AvatarURL: user.AvatarURL,
			IsActive:  user.IsActive,
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get current user profile
// @Description Get current authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserInfo
// @Failure 401 {object} models.APIError
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Error:   true,
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	var user models.User
	result := h.db.DB.First(&user, "id = ?", userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.APIError{
			Error:   true,
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	response := models.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Company:   user.Company,
		AvatarURL: user.AvatarURL,
		IsActive:  user.IsActive,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Logout user
// @Description Invalidate user's refresh token (client should discard tokens)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled client-side
	// Server can blacklist refresh token if using token rotation
	// For now, just return success

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully. Please discard your tokens.",
	})
}
