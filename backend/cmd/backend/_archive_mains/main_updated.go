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
	"github.com/easyhire/backend/internal/repository"
	"github.com/easyhire/backend/internal/routes"
	"github.com/easyhire/backend/internal/services"
	"github.com/easyhire/backend/internal/pkg/config"
	"github.com/easyhire/backend/pkg/database"
	"github.com/easyhire/backend/internal/pkg/logger"
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
		
		// Run migrations if enabled
		if cfg.Database.AutoMigrate {
			log.Info().Msg("🔄 Running database migrations...")
			if err := database.RunMigrations(db.DB); err != nil {
				log.Error().Err(err).Msg("Failed to run migrations")
			} else {
				log.Info().Msg("✅ Database migrations completed")
			}
		}
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

	// Initialize services
	assessmentRepo := repository.NewAssessmentRepository(db.DB)
	questionRepo := repository.NewQuestionRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	scoringService := services.NewScoringService()
	
	assessmentService := services.NewAssessmentService(
		assessmentRepo,
		questionRepo,
		userRepo,
		scoringService,
	)
	
	assessmentHandler := handlers.NewAssessmentHandler(assessmentService)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Health check handler
	healthHandler := handlers.NewHealthHandler(db, redisClient)

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
		// Setup assessment routes
		routes.SetupAssessmentRoutes(apiV1, assessmentHandler)
		
		// Auth routes will be added separately
		// TODO: Add auth routes from existing auth handler
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
