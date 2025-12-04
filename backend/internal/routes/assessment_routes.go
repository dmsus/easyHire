package routes

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAssessmentRoutes(router *gin.RouterGroup, assessmentHandler *handlers.AssessmentHandler) {
	// Assessment management routes
	assessments := router.Group("/assessments")
	assessments.Use(middleware.AuthMiddleware())
	{
		// CRUD operations
		assessments.POST("", assessmentHandler.CreateAssessment)
		assessments.GET("", assessmentHandler.ListAssessments)
		assessments.GET("/:id", assessmentHandler.GetAssessment)
		assessments.PUT("/:id", middleware.RequireRole("hr", "admin"), assessmentHandler.UpdateAssessment)
		assessments.DELETE("/:id", middleware.RequireRole("hr", "admin"), assessmentHandler.DeleteAssessment)
		
		// Assessment actions
		assessments.POST("/:id/publish", middleware.RequireRole("hr", "admin"), assessmentHandler.PublishAssessment)
		assessments.POST("/:id/generate-questions", middleware.RequireRole("hr", "admin"), assessmentHandler.GenerateQuestions)
		
		// Candidate invitations
		assessments.POST("/:id/invite", middleware.RequireRole("hr", "admin"), assessmentHandler.InviteCandidates)
		assessments.POST("/:id/bulk-invite", middleware.RequireRole("hr", "admin"), assessmentHandler.BulkInvite)
		assessments.GET("/:id/invitations", middleware.RequireRole("hr", "admin"), assessmentHandler.ListInvitations)
		
		// Start assessment (for candidates)
		assessments.POST("/:id/start", assessmentHandler.StartAssessment)
	}
	
	// Session management routes
	sessions := router.Group("/sessions")
	sessions.Use(middleware.AuthMiddleware())
	{
		sessions.GET("/:id", assessmentHandler.GetSession)
		sessions.GET("/:id/progress", assessmentHandler.GetSessionProgress)
		sessions.POST("/:id/pause", assessmentHandler.PauseSession)
		sessions.POST("/:id/resume", assessmentHandler.ResumeSession)
		sessions.POST("/:id/complete", assessmentHandler.CompleteAssessment)
		
		// Answer submission
		sessions.POST("/:session_id/questions/:question_id/answer", assessmentHandler.SubmitAnswer)
	}
	
	// Results routes
	results := router.Group("/results")
	results.Use(middleware.AuthMiddleware())
	{
		results.GET("/:session_id", assessmentHandler.GetResult)
	}
	
	// Invitation routes (some public, some protected)
	invitations := router.Group("/invitations")
	{
		// Public route for validating invitations
		invitations.GET("/:token", assessmentHandler.ValidateInvitation)
		
		// Protected route for accepting invitations
		invitations.POST("/:token/accept", middleware.AuthMiddleware(), assessmentHandler.AcceptInvitation)
	}
}
