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
	"github.com/easyhire/backend/internal/services"
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
	// =============================================
	// ИНИЦИАЛИЗАЦИЯ СЕРВИСОВ AUTH + ASSESSMENT
	// =============================================
	
	// Initialize JWT service
	jwtCfg := &config.JWTConfig{
		JWTPrivateKey:      cfg.Security.JWTPrivateKey,
		JWTPublicKey:       cfg.Security.JWTPublicKey,
		JWTSecret:          cfg.Security.JWTSecret,
		AccessTokenExpiry:  cfg.Security.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.Security.RefreshTokenExpiry,
	}

	jwtService := services.NewJWTService(jwtCfg)
	if jwtService == nil {
		log.Fatal().Msg("Failed to initialize JWT service")
	}

	// Initialize password service
	passwordService := services.NewPasswordService()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	assessmentRepo := repository.NewAssessmentRepository(db.DB)
	questionRepo := repository.NewQuestionRepository(db.DB)

	// Initialize services
	scoringService := services.NewScoringService()
	emailService := services.NewEmailService()
	
	// Auth service
	authService := services.NewAuthService(userRepo, jwtService, passwordService)
	
	// Assessment service with email integration
	assessmentService := services.NewAssessmentService(
		assessmentRepo,
		questionRepo,
		userRepo,
		scoringService,
		emailService,
	)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db, redisClient)
	authHandler := handlers.NewAuthHandler(authService)
	assessmentHandler := handlers.NewAssessmentHandler(assessmentService)

	log.Info().Msg("✅ All services initialized successfully")
	// =============================================
	// НАСТРОЙКА GIN И РОУТОВ
	// =============================================
	
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.Security.CORSOrigins))

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
		// ========== PUBLIC ROUTES (No Auth) ==========
		// Public auth routes
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
		}

		// Public assessment invitation routes
		public := apiV1.Group("/public")
		{
			public.GET("/assessments/:invite_code", assessmentHandler.GetAssessmentByInviteToken)
			public.POST("/assessments/:invite_code/start", assessmentHandler.StartAssessment)
		}

		// ========== PROTECTED ROUTES (Require Auth) ==========
		protected := apiV1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// Auth routes (protected)
			protectedAuth := protected.Group("/auth")
			{
				protectedAuth.GET("/me", authHandler.GetProfile)
				protectedAuth.POST("/logout", authHandler.Logout)
			}

			// Assessments routes (HR and Admin only)
			assessments := protected.Group("/assessments")
			assessments.Use(middleware.HRorAdmin())
			{
				assessments.GET("", assessmentHandler.GetAssessments)
				assessments.POST("", assessmentHandler.CreateAssessment)
				assessments.GET("/:id", assessmentHandler.GetAssessmentByID)
				assessments.PUT("/:id", assessmentHandler.UpdateAssessment)
				assessments.DELETE("/:id", assessmentHandler.DeleteAssessment)
				assessments.POST("/:id/invite", assessmentHandler.InviteCandidate)
				assessments.POST("/:id/invite/bulk", assessmentHandler.BulkInviteCandidates)
			}

			// Session management routes
			sessions := protected.Group("/sessions")
			{
				sessions.POST("/:id/answers", assessmentHandler.SubmitAnswer)
				sessions.POST("/:id/complete", assessmentHandler.CompleteAssessment)
				sessions.GET("/:id", assessmentHandler.GetSession)
			}

			// Questions routes (HR, Admin, and Technical Experts)
			questions := protected.Group("/questions")
			questions.Use(middleware.RoleMiddleware("hr", "admin", "technical_expert"))
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
	}
	// =============================================
	// ЗАПУСК СЕРВЕРА И GRACEFUL SHUTDOWN
	// =============================================
	
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
