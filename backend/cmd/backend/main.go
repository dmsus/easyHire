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
	"github.com/easyhire/internal/models"
	"github.com/easyhire/internal/pkg/auth"
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
	var db *database.Database
	if cfg.Database.Host != "" {
		db, err = database.NewDatabase(&cfg.Database)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Error().Err(err).Msg("Failed to close database connection")
			}
		}()
		log.Info().Msg("✅ Database connection established")
	} else {
		log.Warn().Msg("⚠️ Database configuration missing, running without database")
	}

	// Initialize Redis
	var redisClient *redis.Client
	if cfg.Redis.Host != "" {
		redisClient = redis.NewClient(&redis.Options{
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
	} else {
		log.Warn().Msg("⚠️ Redis configuration missing, running without cache")
	}

	// Initialize JWT service
	jwtCfg := &auth.JWTConfig{
		JWTPrivateKey:      cfg.Security.JWTPrivateKey,
		JWTPublicKey:       cfg.Security.JWTPublicKey,
		JWTSecret:          cfg.Security.JWTSecret,
		AccessTokenExpiry:  cfg.Security.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.Security.RefreshTokenExpiry,
	}

	jwtService, err := auth.NewJWTService(jwtCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize JWT service")
	}

	// Initialize password service
	passwordService := auth.NewPasswordService()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.Security.CORSOrigins))

	// Health check handler
	healthHandler := handlers.NewHealthHandler(db, redisClient)

	// Auth handler
	authHandler := handlers.NewAuthHandler(db, jwtService, passwordService)

	// Health routes (public)
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
		// Public auth routes
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)

			// Protected auth routes
			protectedAuth := authGroup.Group("")
			protectedAuth.Use(middleware.AuthMiddleware(jwtService))
			{
				protectedAuth.GET("/me", authHandler.GetProfile)
				protectedAuth.POST("/logout", authHandler.Logout)
			}
		}

		// Protected routes (require authentication)
		protected := apiV1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// Assessments routes (HR and Admin only)
			assessments := protected.Group("/assessments")
			assessments.Use(middleware.HRorAdmin())
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
				assessments.PUT("/:id", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Update assessment - TODO"})
				})
				assessments.DELETE("/:id", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Delete assessment - TODO"})
				})
			}

			// Candidates routes (HR and Admin only)
			candidates := protected.Group("/candidates")
			candidates.Use(middleware.HRorAdmin())
			{
				candidates.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "List candidates - TODO"})
				})
				candidates.POST("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Create candidate - TODO"})
				})
				candidates.GET("/:id", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Get candidate - TODO"})
				})
			}

			// Questions routes (HR, Admin, and Technical Experts)
			questions := protected.Group("/questions")
			questions.Use(middleware.RoleMiddleware(
				models.RoleHR,
				models.RoleAdmin,
				models.RoleTechnicalExpert,
			))
			{
				questions.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "List questions - TODO"})
				})
				questions.POST("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Create question - TODO"})
				})
				questions.GET("/:id", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Get question - TODO"})
				})
			}

			// Results routes (accessible by all authenticated users)
			results := protected.Group("/results")
			{
				results.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "List results - TODO"})
				})
				results.GET("/:id", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Get result - TODO"})
				})
			}

			// Admin-only routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/users", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "List all users - TODO"})
				})
				admin.GET("/stats", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Admin statistics - TODO"})
				})
			}
		}

		// Public assessment invitation routes (no auth required)
		public := apiV1.Group("/public")
		{
			public.GET("/assessments/:invite_code", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Get assessment by invite code - TODO"})
			})
			public.POST("/assessments/:invite_code/start", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Start assessment - TODO"})
			})
		}
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
