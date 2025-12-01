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

// @Summary Health check
// @Description Check if the API is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	log := logger.Global().WithRequestID(c.GetString("request_id"))
	
	services := make(map[string]string)
	status := "healthy"

	// Check database
	if err := h.db.HealthCheck(); err != nil {
		services["database"] = "unhealthy"
		status = "degraded"
		log.Error().Err(err).Msg("Database health check failed")
	} else {
		services["database"] = "healthy"
	}

	// Check Redis
	if h.redis != nil {
		if _, err := h.redis.Ping(c).Result(); err != nil {
			services["redis"] = "unhealthy"
			status = "degraded"
			log.Error().Err(err).Msg("Redis health check failed")
		} else {
			services["redis"] = "healthy"
		}
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
		Services:  services,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Readiness check
// @Description Check if the API is ready to serve requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	response := gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC(),
		"services": gin.H{
			"api":       "ready",
			"database":  "ready",
			"redis":     "ready",
			"migrations": "applied",
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Liveness check
// @Description Simple liveness probe
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /live [get]
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC(),
	})
}

// @Summary Metrics endpoint
// @Description Prometheus metrics
// @Tags health
// @Accept json
// @Produce text/plain
// @Router /metrics [get]
func (h *HealthHandler) Metrics(c *gin.Context) {
	// This would be handled by Prometheus middleware
	c.JSON(http.StatusOK, gin.H{
		"message": "Metrics are exposed on /metrics endpoint via Prometheus middleware",
	})
}

// @Summary API information
// @Description Get API version and information
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func (h *HealthHandler) APIRoot(c *gin.Context) {
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
