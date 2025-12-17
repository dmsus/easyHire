package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/easyhire/backend/internal/models"
	"github.com/easyhire/backend/internal/repository"
	"gorm.io/gorm"
)

type AssessmentService interface {
	// Assessments
	CreateAssessment(ctx context.Context, req models.CreateAssessmentRequest, createdBy string) (*models.Assessment, error)
	GetAssessment(ctx context.Context, id string) (*models.Assessment, error)
	UpdateAssessment(ctx context.Context, id string, req models.UpdateAssessmentRequest) (*models.Assessment, error)
	DeleteAssessment(ctx context.Context, id string) error
	ListAssessments(ctx context.Context, filter repository.AssessmentFilter) ([]models.Assessment, error)

	// Invitations
	InviteCandidate(ctx context.Context, assessmentID, email, invitedBy string) (*models.Invitation, error)
	GetInvitation(ctx context.Context, token string) (*models.Invitation, error)

	// Sessions
	StartSession(ctx context.Context, assessmentID, candidateID string) (*models.AssessmentSession, error)
	GetSession(ctx context.Context, sessionID string) (*models.AssessmentSession, error)
	SubmitAnswer(ctx context.Context, sessionID, questionID string, req models.CandidateAnswerRequest) error
	CompleteSession(ctx context.Context, sessionID string) (*models.Result, error)
}

type assessmentService struct {
	assessmentRepo repository.AssessmentRepository
	questionRepo   repository.QuestionRepository
	emailService   *EmailService
	db             *gorm.DB
}

func NewAssessmentService(
	assessmentRepo repository.AssessmentRepository,
	questionRepo repository.QuestionRepository,
	db *gorm.DB,
) AssessmentService {
	return &assessmentService{
		assessmentRepo: assessmentRepo,
		questionRepo:   questionRepo,
		emailService:   NewEmailService(),
		db:             db,
	}
}

// ==========================
// CREATE ASSESSMENT (CORE #9)
// ==========================

func (s *assessmentService) CreateAssessment(ctx context.Context, req models.CreateAssessmentRequest, createdBy string) (*models.Assessment, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if createdBy == "" {
		return nil, fmt.Errorf("created_by is required")
	}
	if len(req.Competencies) == 0 {
		return nil, fmt.Errorf("at least one competency is required")
	}

	assessment := &models.Assessment{
		Title:            req.Title,
		Description:      req.Description,
		Type:             string(req.Type),
		TargetLevel:      string(req.TargetLevel),
		TimeLimit:        req.TimeLimit,
		TotalQuestions:   req.TotalQuestions,
		PassingScore:     req.PassingScore,
		ShuffleQuestions: req.ShuffleQuestions,
		ShowExplanation:  req.ShowExplanation,
		CreatedBy:        createdBy,
		Status:           models.AssessmentStatus("draft"),
	}

	// One transaction: assessment + competencies
	returnAssessment := (*models.Assessment)(nil)

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Use repo via tx? simplest: direct create for assessment, then repo methods for others are on base DB.
		// To keep it consistent, we create assessment with tx directly.
		if err := tx.Create(assessment).Error; err != nil {
			return fmt.Errorf("create assessment failed: %w", err)
		}

		for _, c := range req.Competencies {
			minQ := c.MinQuestions
			maxQ := c.MaxQuestions
			if minQ <= 0 {
				minQ = 1
			}
			if maxQ <= 0 {
				maxQ = 5
			}
			if maxQ < minQ {
				return fmt.Errorf("competency %s: max_questions must be >= min_questions", c.CompetencyID)
			}

			comp := &models.AssessmentCompetency{
				AssessmentID: assessment.ID,
				CompetencyID: c.CompetencyID,
				Level:        string(c.Level),
				Weight:       c.Weight,
				MinQuestions: minQ,
				MaxQuestions: maxQ,
			}

			if err := tx.Create(comp).Error; err != nil {
				return fmt.Errorf("create competency failed: %w", err)
			}
		}

		// Reload with preloads
		var out models.Assessment
		if err := tx.Preload("Competencies").Preload("Tags").First(&out, "id = ?", assessment.ID).Error; err != nil {
			return fmt.Errorf("reload assessment failed: %w", err)
		}
		returnAssessment = &out
		return nil
	})

	if err != nil {
		return nil, err
	}
	return returnAssessment, nil
}

func (s *assessmentService) GetAssessment(ctx context.Context, id string) (*models.Assessment, error) {
	return s.assessmentRepo.GetAssessmentByID(ctx, id)
}

func (s *assessmentService) UpdateAssessment(ctx context.Context, id string, req models.UpdateAssessmentRequest) (*models.Assessment, error) {
	assessment, err := s.assessmentRepo.GetAssessmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("assessment not found: %w", err)
	}

	if req.Title != nil {
		assessment.Title = *req.Title
	}
	if req.Description != nil {
		assessment.Description = *req.Description
	}
	if req.Status != nil {
		assessment.Status = *req.Status
	}
	if req.TimeLimit != nil {
		assessment.TimeLimit = *req.TimeLimit
	}
	if req.PassingScore != nil {
		assessment.PassingScore = *req.PassingScore
	}
	if req.ShuffleQuestions != nil {
		assessment.ShuffleQuestions = *req.ShuffleQuestions
	}
	if req.ShowExplanation != nil {
		assessment.ShowExplanation = *req.ShowExplanation
	}

	if err := s.assessmentRepo.UpdateAssessment(ctx, assessment); err != nil {
		return nil, fmt.Errorf("update assessment failed: %w", err)
	}
	return s.assessmentRepo.GetAssessmentByID(ctx, id)
}

func (s *assessmentService) DeleteAssessment(ctx context.Context, id string) error {
	return s.assessmentRepo.DeleteAssessment(ctx, id)
}

func (s *assessmentService) ListAssessments(ctx context.Context, filter repository.AssessmentFilter) ([]models.Assessment, error) {
	list, _, err := s.assessmentRepo.ListAssessments(ctx, filter)
	return list, err
}

// ==========================
// INVITATIONS
// ==========================

func (s *assessmentService) InviteCandidate(ctx context.Context, assessmentID, email, invitedBy string) (*models.Invitation, error) {
	if assessmentID == "" {
		return nil, fmt.Errorf("assessment_id is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if invitedBy == "" {
		return nil, fmt.Errorf("invited_by is required")
	}

	// Ensure assessment exists
	assessment, err := s.assessmentRepo.GetAssessmentByID(ctx, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("assessment not found: %w", err)
	}

	token := generateInvitationToken()

	inv := &models.Invitation{
		AssessmentID: assessmentID,
		CandidateID:  nil, // ✅ candidate may not exist yet
		Email:        email,
		Token:        token,
		Status:       models.InvitationStatus("pending"),
		InvitedBy:    invitedBy,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.assessmentRepo.CreateInvitation(ctx, inv); err != nil {
		return nil, fmt.Errorf("create invitation failed: %w", err)
	}

	// Email (dev logs to console)
	s.emailService.SendInvitation(email, token, assessment.Title)
	return inv, nil
}

func (s *assessmentService) GetInvitation(ctx context.Context, token string) (*models.Invitation, error) {
	return s.assessmentRepo.GetInvitationByToken(ctx, token)
}

// ==========================
// SESSIONS
// ==========================

func (s *assessmentService) StartSession(ctx context.Context, assessmentID, candidateID string) (*models.AssessmentSession, error) {
	if assessmentID == "" || candidateID == "" {
		return nil, fmt.Errorf("assessment_id and candidate_id are required")
	}

	// Ensure assessment exists
	_, err := s.assessmentRepo.GetAssessmentByID(ctx, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("assessment not found: %w", err)
	}

	// If there's an active session - return it
	active, err := s.assessmentRepo.GetActiveSession(ctx, assessmentID, candidateID)
	if err == nil && active != nil {
		return active, nil
	}

	now := time.Now()
	session := &models.AssessmentSession{
		AssessmentID: assessmentID,
		CandidateID:  candidateID,
		Status:       models.SessionStatus("in_progress"),
		StartedAt:    &now,
		TimeSpent:    0,
	}

	if err := s.assessmentRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("create session failed: %w", err)
	}
	return session, nil
}

func (s *assessmentService) GetSession(ctx context.Context, sessionID string) (*models.AssessmentSession, error) {
	return s.assessmentRepo.GetSessionByID(ctx, sessionID)
}

func (s *assessmentService) SubmitAnswer(ctx context.Context, sessionID, questionID string, req models.CandidateAnswerRequest) error {
	if sessionID == "" || questionID == "" {
		return fmt.Errorf("session_id and question_id are required")
	}

	session, err := s.assessmentRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}
	if string(session.Status) == "completed" {
		return fmt.Errorf("session already completed")
	}

	now := time.Now()

	// Upsert by (session_id, question_id) if repository supports it
	existing, err := s.assessmentRepo.GetAnswer(ctx, sessionID, questionID)
	if err == nil && existing != nil {
		existing.Answer = req.Answer
		existing.Code = req.Code
		existing.TimeSpent = req.TimeSpent
		existing.SubmittedAt = &now
		return s.assessmentRepo.UpdateAnswer(ctx, existing)
	}

	ans := &models.CandidateAnswer{
		SessionID:   sessionID,
		QuestionID:  questionID,
		Answer:      req.Answer,
		Code:        req.Code,
		TimeSpent:   req.TimeSpent,
		StartedAt:   now, // MVP: real StartedAt can be tracked by client
		SubmittedAt: &now,
	}
	return s.assessmentRepo.SaveAnswer(ctx, ans)
}

func (s *assessmentService) CompleteSession(ctx context.Context, sessionID string) (*models.Result, error) {
	session, err := s.assessmentRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// If already completed, return existing result if any
	if string(session.Status) == "completed" {
		existing, err := s.assessmentRepo.GetResultBySessionID(ctx, sessionID)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	answers, _ := s.assessmentRepo.GetSessionAnswers(ctx, sessionID)

	timeSpent := 0
	for _, a := range answers {
		timeSpent += a.TimeSpent
	}

	now := time.Now()

	// MVP scoring: пока без code-exec (#10) и автоматической проверки,
	// оставляем score/percentage = 0, но фиксируем completion и время.
	result := &models.Result{
		SessionID:   sessionID,
		TotalScore:  0,
		Percentage:  0,
		Level:       "TRAINEE",
		TimeSpent:   timeSpent,
		CompletedAt: now,
	}

	if err := s.assessmentRepo.CreateResult(ctx, result); err != nil {
		return nil, fmt.Errorf("create result failed: %w", err)
	}

	session.Status = models.SessionStatus("completed")
	session.CompletedAt = &now
	session.TimeSpent = timeSpent
	session.Score = result.TotalScore
	session.Percentage = result.Percentage
	session.Level = result.Level

	if err := s.assessmentRepo.UpdateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("update session failed: %w", err)
	}

	return result, nil
}

// 32 hex chars token
func generateInvitationToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
