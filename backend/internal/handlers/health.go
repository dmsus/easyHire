package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/easyhire/backend/internal/pkg/logger"
	"github.com/easyhire/backend/pkg/database"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db    *database.Database
	redis *redis.Client
}

func NewHealthHandler(db *database.Database, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	log := logger.Global().With().Str("handler", "health").Logger()
	log.Info().Msg("Health check request")
	
	services := make(map[string]string)
	status := "healthy"

	// Check database
	if h.db != nil {
		if err := h.db.HealthCheck(); err != nil {
			services["database"] = "unhealthy"
			status = "degraded"
			log.Error().Err(err).Msg("Database health check failed")
		} else {
			services["database"] = "healthy"
			log.Debug().Msg("Database health check passed")
		}
	}

	// Check Redis
	if h.redis != nil {
		if _, err := h.redis.Ping(c).Result(); err != nil {
			services["redis"] = "unhealthy"
			status = "degraded"
			log.Error().Err(err).Msg("Redis health check failed")
		} else {
			services["redis"] = "healthy"
			log.Debug().Msg("Redis health check passed")
		}
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
		Services:  services,
	}

	log.Info().Str("status", status).Msg("Health check completed")
	c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	logger.Global().Info().Str("handler", "health").Msg("Readiness check")
	
	response := gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	logger.Global().Info().Str("handler", "health").Msg("Liveness check")
	
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC(),
	})
}

func (h *HealthHandler) APIRoot(c *gin.Context) {
	logger.Global().Info().Str("handler", "root").Msg("API root request")
	
	response := gin.H{
		"name":        "EasyHire API",
		"version":     "1.0.0",
		"description": "AI-Powered Technical Assessment Platform for Go Developers",
		"docs":        "/docs",
		"health":      "/health",
		"ready":       "/ready",
		"live":        "/live",
		"endpoints": gin.H{
			"auth":        "/api/v1/auth",
			"assessments": "/api/v1/assessments",
			"candidates":  "/api/v1/candidates",
			"questions":   "/api/v1/questions",
			"execute":     "/api/v1/execute",
			"results":     "/api/v1/results",
		},
	}

	c.JSON(http.StatusOK, response)
}
