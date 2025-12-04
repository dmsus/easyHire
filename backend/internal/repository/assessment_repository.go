package repository

import (
    "context"
    
    "github.com/easyhire/backend/internal/models"
    "gorm.io/gorm"
)

type AssessmentFilter struct {
    ID        string
    CreatedBy string
    Status    string
    Search    string
    Limit     int
    Offset    int
}

type AssessmentRepository interface {
    // Assessment CRUD
    CreateAssessment(ctx context.Context, assessment *models.Assessment) error
    GetAssessmentByID(ctx context.Context, id string) (*models.Assessment, error)
    GetAssessmentWithQuestions(ctx context.Context, id string) (*models.Assessment, error)
    UpdateAssessment(ctx context.Context, assessment *models.Assessment) error
    DeleteAssessment(ctx context.Context, id string) error
    ListAssessments(ctx context.Context, filter AssessmentFilter) ([]models.Assessment, int64, error)
    
    // Questions
    AddQuestionToAssessment(ctx context.Context, assessmentID, questionID string, order int) error
    RemoveQuestionFromAssessment(ctx context.Context, assessmentID, questionID string) error
    GetAssessmentQuestions(ctx context.Context, assessmentID string) ([]models.AssessmentQuestion, error)
    
    // Sessions
    CreateSession(ctx context.Context, session *models.AssessmentSession) error
    GetSessionByID(ctx context.Context, sessionID string) (*models.AssessmentSession, error)
    GetActiveSession(ctx context.Context, assessmentID, candidateID string) (*models.AssessmentSession, error)
    UpdateSession(ctx context.Context, session *models.AssessmentSession) error
    GetSessionAnswers(ctx context.Context, sessionID string) ([]models.CandidateAnswer, error)
    
    // Answers
    SaveAnswer(ctx context.Context, answer *models.CandidateAnswer) error
    GetAnswer(ctx context.Context, sessionID, questionID string) (*models.CandidateAnswer, error)
    UpdateAnswer(ctx context.Context, answer *models.CandidateAnswer) error
    
    // Results
    CreateResult(ctx context.Context, result *models.Result) error
    GetResultBySessionID(ctx context.Context, sessionID string) (*models.Result, error)
    
    // Invitations
    CreateInvitation(ctx context.Context, invitation *models.Invitation) error
    BulkCreateInvitations(ctx context.Context, invitations []models.Invitation) error
    GetInvitationByToken(ctx context.Context, token string) (*models.Invitation, error)
    UpdateInvitation(ctx context.Context, invitation *models.Invitation) error
    GetInvitationsByAssessment(ctx context.Context, assessmentID string) ([]models.Invitation, error)
}

type assessmentRepository struct {
    db *gorm.DB
}

func NewAssessmentRepository(db *gorm.DB) AssessmentRepository {
    return &assessmentRepository{db: db}
}

func (r *assessmentRepository) CreateAssessment(ctx context.Context, assessment *models.Assessment) error {
    return r.db.WithContext(ctx).Create(assessment).Error
}

func (r *assessmentRepository) GetAssessmentByID(ctx context.Context, id string) (*models.Assessment, error) {
    var assessment models.Assessment
    result := r.db.WithContext(ctx).
        Preload("Competencies").
        Preload("Tags").
        First(&assessment, "id = ?", id)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &assessment, nil
}

func (r *assessmentRepository) GetAssessmentWithQuestions(ctx context.Context, id string) (*models.Assessment, error) {
    var assessment models.Assessment
    result := r.db.WithContext(ctx).
        Preload("Competencies").
        Preload("Tags").
        Preload("Questions").
        First(&assessment, "id = ?", id)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &assessment, nil
}

func (r *assessmentRepository) UpdateAssessment(ctx context.Context, assessment *models.Assessment) error {
    return r.db.WithContext(ctx).Save(assessment).Error
}

func (r *assessmentRepository) DeleteAssessment(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Delete(&models.Assessment{}, "id = ?", id).Error
}

func (r *assessmentRepository) ListAssessments(ctx context.Context, filter AssessmentFilter) ([]models.Assessment, int64, error) {
    var assessments []models.Assessment
    var total int64
    
    query := r.db.WithContext(ctx).Model(&models.Assessment{})
    
    if filter.ID != "" {
        query = query.Where("id = ?", filter.ID)
    }
    
    if filter.CreatedBy != "" {
        query = query.Where("created_by = ?", filter.CreatedBy)
    }
    
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    
    if filter.Search != "" {
        query = query.Where("title ILIKE ? OR description ILIKE ?", 
            "%"+filter.Search+"%", "%"+filter.Search+"%")
    }
    
    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Apply pagination
    if filter.Limit > 0 {
        query = query.Limit(filter.Limit)
    }
    if filter.Offset > 0 {
        query = query.Offset(filter.Offset)
    }
    
    // Execute query
    result := query.Find(&assessments)
    return assessments, total, result.Error
}

func (r *assessmentRepository) AddQuestionToAssessment(ctx context.Context, assessmentID, questionID string, order int) error {
    assessmentQuestion := &models.AssessmentQuestion{
        AssessmentID: assessmentID,
        QuestionID:   questionID,
        Order:        order,
    }
    
    return r.db.WithContext(ctx).Create(assessmentQuestion).Error
}

func (r *assessmentRepository) RemoveQuestionFromAssessment(ctx context.Context, assessmentID, questionID string) error {
    return r.db.WithContext(ctx).
        Delete(&models.AssessmentQuestion{}, "assessment_id = ? AND question_id = ?", assessmentID, questionID).
        Error
}

func (r *assessmentRepository) GetAssessmentQuestions(ctx context.Context, assessmentID string) ([]models.AssessmentQuestion, error) {
    var questions []models.AssessmentQuestion
    
    result := r.db.WithContext(ctx).
        Where("assessment_id = ?", assessmentID).
        Order("\"order\"").
        Find(&questions)
    
    return questions, result.Error
}

func (r *assessmentRepository) CreateSession(ctx context.Context, session *models.AssessmentSession) error {
    return r.db.WithContext(ctx).Create(session).Error
}

func (r *assessmentRepository) GetSessionByID(ctx context.Context, sessionID string) (*models.AssessmentSession, error) {
    var session models.AssessmentSession
    result := r.db.WithContext(ctx).
        Preload("Answers").
        First(&session, "id = ?", sessionID)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &session, nil
}

func (r *assessmentRepository) GetActiveSession(ctx context.Context, assessmentID, candidateID string) (*models.AssessmentSession, error) {
    var session models.AssessmentSession
    result := r.db.WithContext(ctx).
        Where("assessment_id = ? AND candidate_id = ? AND status IN (?)", 
            assessmentID, candidateID, []string{"pending", "in_progress"}).
        First(&session)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &session, nil
}

func (r *assessmentRepository) UpdateSession(ctx context.Context, session *models.AssessmentSession) error {
    return r.db.WithContext(ctx).Save(session).Error
}

func (r *assessmentRepository) GetSessionAnswers(ctx context.Context, sessionID string) ([]models.CandidateAnswer, error) {
    var answers []models.CandidateAnswer
    result := r.db.WithContext(ctx).
        Where("session_id = ?", sessionID).
        Order("created_at").
        Find(&answers)
    
    return answers, result.Error
}

func (r *assessmentRepository) SaveAnswer(ctx context.Context, answer *models.CandidateAnswer) error {
    return r.db.WithContext(ctx).Create(answer).Error
}

func (r *assessmentRepository) GetAnswer(ctx context.Context, sessionID, questionID string) (*models.CandidateAnswer, error) {
    var answer models.CandidateAnswer
    result := r.db.WithContext(ctx).
        Where("session_id = ? AND question_id = ?", sessionID, questionID).
        First(&answer)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &answer, nil
}

func (r *assessmentRepository) UpdateAnswer(ctx context.Context, answer *models.CandidateAnswer) error {
    return r.db.WithContext(ctx).Save(answer).Error
}

func (r *assessmentRepository) CreateResult(ctx context.Context, result *models.Result) error {
    return r.db.WithContext(ctx).Create(result).Error
}

func (r *assessmentRepository) GetResultBySessionID(ctx context.Context, sessionID string) (*models.Result, error) {
    var result models.Result
    err := r.db.WithContext(ctx).
        Where("session_id = ?", sessionID).
        First(&result).Error
    
    return &result, err
}

func (r *assessmentRepository) CreateInvitation(ctx context.Context, invitation *models.Invitation) error {
    return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *assessmentRepository) BulkCreateInvitations(ctx context.Context, invitations []models.Invitation) error {
    return r.db.WithContext(ctx).Create(&invitations).Error
}

func (r *assessmentRepository) GetInvitationByToken(ctx context.Context, token string) (*models.Invitation, error) {
    var invitation models.Invitation
    result := r.db.WithContext(ctx).
        Where("token = ?", token).
        First(&invitation)
    
    return &invitation, result.Error
}

func (r *assessmentRepository) UpdateInvitation(ctx context.Context, invitation *models.Invitation) error {
    return r.db.WithContext(ctx).Save(invitation).Error
}

func (r *assessmentRepository) GetInvitationsByAssessment(ctx context.Context, assessmentID string) ([]models.Invitation, error) {
    var invitations []models.Invitation
    result := r.db.WithContext(ctx).
        Where("assessment_id = ?", assessmentID).
        Find(&invitations)
    
    return invitations, result.Error
}
