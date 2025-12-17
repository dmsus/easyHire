package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/routes"
	"backend/internal/services"
	"backend/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database (in-memory для тестирования)
	db, err := database.NewDatabase(&database.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "easyhire",
		SSLMode:  "disable",
	})
	
	if err != nil {
		log.Printf("⚠️ Database connection failed, using in-memory mode: %v", err)
		// Продолжаем без БД для тестирования
	} else {
		defer db.Close()
		log.Println("✅ Database connection established")
	}

	// Initialize repositories
	var assessmentRepo repository.AssessmentRepository
	var questionRepo repository.QuestionRepository
	var userRepo repository.UserRepository
	
	if db != nil {
		assessmentRepo = repository.NewAssessmentRepository(db.DB)
		questionRepo = repository.NewQuestionRepository(db.DB)
		userRepo = repository.NewUserRepository(db.DB)
	} else {
		// Используем заглушки для тестирования
		log.Println("⚠️ Using mock repositories")
		// В реальном приложении нужно создать mock репозитории
	}

	// Initialize services
	scoringService := services.NewScoringService()
	assessmentService := services.NewAssessmentService(
		assessmentRepo,
		questionRepo,
		userRepo,
		scoringService,
	)
	
	assessmentHandler := handlers.NewAssessmentHandler(assessmentService)

	// Set Gin mode
	gin.SetMode(gin.DebugMode)

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "EasyHire Assessment Engine",
			"version": "1.0.0",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// API v1 routes group
	apiV1 := router.Group("/api/v1")
	{
		// Setup assessment routes
		routes.SetupAssessmentRoutes(apiV1, assessmentHandler)
		
		// Test endpoint
		apiV1.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Assessment engine API is working",
				"endpoints": []string{
					"POST   /api/v1/assessments",
					"GET    /api/v1/assessments",
					"GET    /api/v1/assessments/:id",
					"POST   /api/v1/assessments/:id/start",
					"POST   /api/v1/sessions/:session_id/complete",
				},
			})
		})
	}

	// Start server
	port := ":8080"
	log.Printf("🚀 Starting Assessment Engine on http://localhost%s", port)
	
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
