package routes

import (
	"github.com/easyhire/backend/internal/handlers"
	"github.com/easyhire/backend/internal/middleware"
	"github.com/easyhire/internal/pkg/auth"
	"github.com/gin-gonic/gin"
)

func SetupAssessmentRoutes(router *gin.RouterGroup, jwtService *auth.JWTService, assessmentHandler *handlers.AssessmentHandler) {
	// All assessment routes require JWT auth
	assessments := router.Group("/assessments")
	assessments.Use(middleware.AuthMiddleware(jwtService))
	{
		// CRUD
		assessments.POST("", middleware.HRorAdmin(), assessmentHandler.CreateAssessment)
		assessments.GET("", middleware.HRorAdmin(), assessmentHandler.ListAssessments)
		assessments.GET("/:id", assessmentHandler.GetAssessment)
		assessments.PUT("/:id", middleware.HRorAdmin(), assessmentHandler.UpdateAssessment)
		assessments.DELETE("/:id", middleware.HRorAdmin(), assessmentHandler.DeleteAssessment)

		// Invitations
		assessments.POST("/:id/invite", middleware.HRorAdmin(), assessmentHandler.InviteCandidate)
		assessments.POST("/:id/bulk-invite", middleware.HRorAdmin(), assessmentHandler.BulkInvite)

		// Candidate starts session for assessment
		assessments.POST("/:id/start", assessmentHandler.StartSession)
	}

	// Session routes require JWT auth
	sessions := router.Group("/sessions")
	sessions.Use(middleware.AuthMiddleware(jwtService))
	{
		sessions.GET("/:session_id", assessmentHandler.GetSession)
		sessions.POST("/:session_id/answers", assessmentHandler.SubmitAnswer)
		sessions.POST("/:session_id/complete", assessmentHandler.CompleteSession)
	}

	// Invitation token lookup (public)
	invitations := router.Group("/invitations")
	{
		invitations.GET("/:token", assessmentHandler.GetInvitation)
	}
}
