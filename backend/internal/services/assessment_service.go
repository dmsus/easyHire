package services

import (
    "context"
    "fmt"
    "time"
    
    "github.com/easyhire/backend/internal/models"
    "github.com/easyhire/backend/internal/repository"
)

type AssessmentService interface {
    // Basic assessment operations
    CreateAssessment(ctx context.Context, req models.CreateAssessmentRequest, createdBy string) (*models.Assessment, error)
    GetAssessment(ctx context.Context, id string) (*models.Assessment, error)
    UpdateAssessment(ctx context.Context, id string, req models.UpdateAssessmentRequest) (*models.Assessment, error)
    DeleteAssessment(ctx context.Context, id string) error
    ListAssessments(ctx context.Context, filter repository.AssessmentFilter) ([]models.Assessment, error)
    
    // Candidate invitations
    InviteCandidate(ctx context.Context, assessmentID, email string) (*models.Invitation, error)
    GetInvitation(ctx context.Context, token string) (*models.Invitation, error)
    
    // Assessment sessions
    StartSession(ctx context.Context, assessmentID, candidateID string) (*models.AssessmentSession, error)
    SubmitAnswer(ctx context.Context, sessionID, questionID string, answer string) error
    CompleteSession(ctx context.Context, sessionID string) (*models.Result, error)
    GetSession(ctx context.Context, sessionID string) (*models.AssessmentSession, error)
}

type assessmentService struct {
    assessmentRepo repository.AssessmentRepository
    questionRepo   repository.QuestionRepository
}

func NewAssessmentService(
    assessmentRepo repository.AssessmentRepository,
    questionRepo repository.QuestionRepository,
) AssessmentService {
    return &assessmentService{
        assessmentRepo: assessmentRepo,
        questionRepo:   questionRepo,
    }
}

func (s *assessmentService) CreateAssessment(ctx context.Context, req models.CreateAssessmentRequest, createdBy string) (*models.Assessment, error) {
    // Validate request
    if req.Title == "" {
        return nil, fmt.Errorf("title is required")
    }
    if createdBy == "" {
        return nil, fmt.Errorf("createdBy is required")
    }
    
    // Create assessment
    assessment := &models.Assessment{
        Title:           req.Title,
        Description:     req.Description,
        Type:            string(req.Type),
        TargetLevel:     string(req.TargetLevel),
        TimeLimit:       req.TimeLimit,
        TotalQuestions:  req.TotalQuestions,
        PassingScore:    req.PassingScore,
        ShuffleQuestions: req.ShuffleQuestions,
        ShowExplanation: req.ShowExplanation,
        CreatedBy:       createdBy,
        Status:          models.AssessmentStatusDraft,
    }
    
    // Save to database
    if err := s.assessmentRepo.CreateAssessment(ctx, assessment); err != nil {
        return nil, fmt.Errorf("failed to create assessment: %w", err)
    }
    
    return assessment, nil
}

func (s *assessmentService) GetAssessment(ctx context.Context, id string) (*models.Assessment, error) {
    return s.assessmentRepo.GetAssessmentByID(ctx, id)
}

func (s *assessmentService) UpdateAssessment(ctx context.Context, id string, req models.UpdateAssessmentRequest) (*models.Assessment, error) {
    assessment, err := s.assessmentRepo.GetAssessmentByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("assessment not found: %w", err)
    }
    
    // Update fields
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
    
    // Save changes
    if err := s.assessmentRepo.UpdateAssessment(ctx, assessment); err != nil {
        return nil, fmt.Errorf("failed to update assessment: %w", err)
    }
    
    return assessment, nil
}

func (s *assessmentService) DeleteAssessment(ctx context.Context, id string) error {
    return s.assessmentRepo.DeleteAssessment(ctx, id)
}

func (s *assessmentService) ListAssessments(ctx context.Context, filter repository.AssessmentFilter) ([]models.Assessment, error) {
    assessments, _, err := s.assessmentRepo.ListAssessments(ctx, filter)
    return assessments, err
}

func (s *assessmentService) InviteCandidate(ctx context.Context, assessmentID, email string) (*models.Invitation, error) {
    // Generate unique token
    token := generateInvitationToken()
    
    invitation := &models.Invitation{
        AssessmentID: assessmentID,
        Email:        email,
        Token:        token,
        Status:       models.InvitationStatusPending,
        ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // 7 days
        InvitedBy:    "system", // TODO: Get from context
    }
    
    if err := s.assessmentRepo.CreateInvitation(ctx, invitation); err != nil {
        return nil, fmt.Errorf("failed to create invitation: %w", err)
    }
    
    return invitation, nil
}

func (s *assessmentService) GetInvitation(ctx context.Context, token string) (*models.Invitation, error) {
    return s.assessmentRepo.GetInvitationByToken(ctx, token)
}

func (s *assessmentService) StartSession(ctx context.Context, assessmentID, candidateID string) (*models.AssessmentSession, error) {
    // Check if assessment exists
    _, err := s.assessmentRepo.GetAssessmentByID(ctx, assessmentID)
    if err != nil {
        return nil, fmt.Errorf("assessment not found: %w", err)
    }
    
    // Create session
    now := time.Now()
    session := &models.AssessmentSession{
        AssessmentID: assessmentID,
        CandidateID:  candidateID,
        Status:       models.SessionStatusInProgress,
        StartedAt:    &now,
    }
    
    if err := s.assessmentRepo.CreateSession(ctx, session); err != nil {
        return nil, fmt.Errorf("failed to create session: %w", err)
    }
    
    return session, nil
}

func (s *assessmentService) SubmitAnswer(ctx context.Context, sessionID, questionID string, answer string) error {
    _, err := s.assessmentRepo.GetSessionByID(ctx, sessionID)
    if err != nil {
        return fmt.Errorf("session not found: %w", err)
    }
    
    // Create or update answer
    candidateAnswer := &models.CandidateAnswer{
        SessionID:   sessionID,
        QuestionID:  questionID,
        Answer:      answer,
        TimeSpent:   60, // TODO: Calculate actual time
        StartedAt:   time.Now().Add(-60 * time.Second),
        SubmittedAt: &time.Time{},
    }
    *candidateAnswer.SubmittedAt = time.Now()
    
    return s.assessmentRepo.SaveAnswer(ctx, candidateAnswer)
}

func (s *assessmentService) CompleteSession(ctx context.Context, sessionID string) (*models.Result, error) {
    session, err := s.assessmentRepo.GetSessionByID(ctx, sessionID)
    if err != nil {
        return nil, fmt.Errorf("session not found: %w", err)
    }
    
    // Mark session as completed
    now := time.Now()
    session.Status = models.SessionStatusCompleted
    session.CompletedAt = &now
    
    // Create result
    result := &models.Result{
        SessionID:    sessionID,
        TotalScore:   0, // TODO: Calculate score
        Percentage:   0,
        Level:        "pending",
        TimeSpent:    session.TimeSpent,
        CompletedAt:  now,
    }
    
    if err := s.assessmentRepo.UpdateSession(ctx, session); err != nil {
        return nil, fmt.Errorf("failed to update session: %w", err)
    }
    
    if err := s.assessmentRepo.CreateResult(ctx, result); err != nil {
        return nil, fmt.Errorf("failed to create result: %w", err)
    }
    
    return result, nil
}

func (s *assessmentService) GetSession(ctx context.Context, sessionID string) (*models.AssessmentSession, error) {
    return s.assessmentRepo.GetSessionByID(ctx, sessionID)
}

// Helper function
func generateInvitationToken() string {
    return fmt.Sprintf("inv_%d", time.Now().UnixNano())
}
