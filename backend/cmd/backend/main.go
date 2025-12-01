package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/easyhire/backend/internal/handlers"
	"github.com/easyhire/backend/internal/middleware"
	"github.com/easyhire/backend/internal/pkg/config"
	"github.com/easyhire/backend/internal/pkg/logger"
	"github.com/easyhire/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

var (
	version = "1.0.0"
	build   = "development"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/.env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := logger.InitGlobalLogger(cfg.Server.Mode, true); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	
	log := logger.Global()
	log.Info().
		Str("version", version).
		Str("build", build).
		Msg("🚀 Starting EasyHire Backend API")

	// Initialize database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	// Test Redis connection
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Warn().Err(err).Msg("Redis connection failed, running without cache")
		redisClient = nil
	} else {
		log.Info().Msg("✅ Redis connection established")
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(&cfg.Security))

	// Health check handler
	healthHandler := handlers.NewHealthHandler(db, redisClient)
	
	// Health routes
	router.GET("/", healthHandler.APIRoot)
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/ready", healthHandler.ReadinessCheck)
	router.GET("/live", healthHandler.LivenessCheck)
	
	// Metrics endpoint (Prometheus)
	if cfg.Monitoring.MetricsEnabled {
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))
		log.Info().Int("port", cfg.Monitoring.MetricsPort).Msg("✅ Metrics enabled")
	}

	// API v1 routes group
	apiV1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - TODO"})
			})
			auth.POST("/refresh", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Refresh endpoint - TODO"})
			})
			auth.POST("/logout", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Logout endpoint - TODO"})
			})
		}

		// Assessments routes
		assessments := apiV1.Group("/assessments")
		{
			assessments.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "List assessments - TODO"})
			})
			assessments.POST("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Create assessment - TODO"})
			})
			assessments.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Get assessment - TODO"})
			})
		}

		// Add more API routes here...
	}

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Info().
			Str("host", cfg.Server.Host).
			Int("port", cfg.Server.Port).
			Str("mode", cfg.Server.Mode).
			Msg("🌐 Starting HTTP server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	// Close Redis connection
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Redis connection")
		}
	}

	log.Info().Msg("👋 Server shutdown completed")
}
