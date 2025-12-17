package handlers

import (
	"net/http"
	"strconv"

	"github.com/easyhire/backend/internal/models"
	"github.com/easyhire/backend/internal/repository"
	"github.com/easyhire/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AssessmentHandler struct {
	assessmentService services.AssessmentService
}

func NewAssessmentHandler(assessmentService services.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{assessmentService: assessmentService}
}

// CreateAssessment создает новую оценку
func (h *AssessmentHandler) CreateAssessment(c *gin.Context) {
	var req models.CreateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	assessment, err := h.assessmentService.CreateAssessment(c.Request.Context(), req, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, assessment)
}

func (h *AssessmentHandler) GetAssessment(c *gin.Context) {
	id := c.Param("id")

	assessment, err := h.assessmentService.GetAssessment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}

	c.JSON(http.StatusOK, assessment)
}

func (h *AssessmentHandler) UpdateAssessment(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assessment, err := h.assessmentService.UpdateAssessment(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assessment)
}

func (h *AssessmentHandler) DeleteAssessment(c *gin.Context) {
	id := c.Param("id")

	if err := h.assessmentService.DeleteAssessment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AssessmentHandler) ListAssessments(c *gin.Context) {
	filter := repository.AssessmentFilter{Limit: 20, Offset: 0}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Offset = (p - 1) * filter.Limit
		}
	}

	assessments, err := h.assessmentService.ListAssessments(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assessments": assessments,
		"total":       len(assessments),
		"limit":       filter.Limit,
		"offset":      filter.Offset,
	})
}

func (h *AssessmentHandler) InviteCandidate(c *gin.Context) {
	assessmentID := c.Param("id")

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	invitation, err := h.assessmentService.InviteCandidate(
		c.Request.Context(),
		assessmentID,
		req.Email,
		userID.(string), // ✅ invited_by
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, invitation)
}

func (h *AssessmentHandler) BulkInvite(c *gin.Context) {
	assessmentID := c.Param("id")

	var req struct {
		Emails []string `json:"emails" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var results []gin.H
	var failedEmails []string

	for _, email := range req.Emails {
		invitation, err := h.assessmentService.InviteCandidate(
			c.Request.Context(),
			assessmentID,
			email,
			userID.(string),
		)
		if err != nil {
			failedEmails = append(failedEmails, email)
			results = append(results, gin.H{"email": email, "status": "failed", "error": err.Error()})
			continue
		}
		results = append(results, gin.H{"email": email, "status": "invited", "token": invitation.Token})
	}

	resp := gin.H{
		"total":      len(req.Emails),
		"successful": len(req.Emails) - len(failedEmails),
		"failed":     len(failedEmails),
		"results":    results,
	}
	if len(failedEmails) > 0 {
		resp["failed_emails"] = failedEmails
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AssessmentHandler) GetInvitation(c *gin.Context) {
	token := c.Param("token")

	invitation, err := h.assessmentService.GetInvitation(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invitation not found"})
		return
	}
	c.JSON(http.StatusOK, invitation)
}

func (h *AssessmentHandler) StartSession(c *gin.Context) {
	assessmentID := c.Param("id")

	candidateID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	session, err := h.assessmentService.StartSession(c.Request.Context(), assessmentID, candidateID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *AssessmentHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	session, err := h.assessmentService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *AssessmentHandler) SubmitAnswer(c *gin.Context) {
	sessionID := c.Param("session_id")

	var req struct {
		QuestionID string                        `json:"question_id" binding:"required"`
		Payload    models.CandidateAnswerRequest `json:",inline"`
		Answer     string                        `json:"answer"`
		Code       string                        `json:"code"`
		TimeSpent  int                           `json:"time_spent"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := models.CandidateAnswerRequest{
		Answer:    req.Answer,
		Code:      req.Code,
		TimeSpent: req.TimeSpent,
	}

	if err := h.assessmentService.SubmitAnswer(c.Request.Context(), sessionID, req.QuestionID, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "answer submitted"})
}

func (h *AssessmentHandler) CompleteSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	result, err := h.assessmentService.CompleteSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
