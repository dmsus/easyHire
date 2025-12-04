package repository

import (
    "context"
    
    "github.com/easyhire/backend/internal/models"
    "gorm.io/gorm"
)

type QuestionFilter struct {
    ID           string
    CompetencyID string
    Level        string
    Type         string
    IsActive     *bool
    Search       string
    Limit        int
    Offset       int
}

type QuestionRepository interface {
    // Basic CRUD
    CreateQuestion(ctx context.Context, question *models.Question) error
    GetQuestionByID(ctx context.Context, id string) (*models.Question, error)
    UpdateQuestion(ctx context.Context, question *models.Question) error
    DeleteQuestion(ctx context.Context, id string) error
    ListQuestions(ctx context.Context, filter QuestionFilter) ([]models.Question, int64, error)
    
    // Specialized
    GetRandomQuestions(ctx context.Context, filter QuestionFilter, count int) ([]models.Question, error)
    GetQuestionsByCompetency(ctx context.Context, competencyID string, level string, limit int) ([]models.Question, error)
    BulkCreateQuestions(ctx context.Context, questions []models.Question) error
}

type questionRepository struct {
    db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
    return &questionRepository{db: db}
}

func (r *questionRepository) CreateQuestion(ctx context.Context, question *models.Question) error {
    result := r.db.WithContext(ctx).Create(question)
    return result.Error
}

func (r *questionRepository) GetQuestionByID(ctx context.Context, id string) (*models.Question, error) {
    var question models.Question
    result := r.db.WithContext(ctx).
        Preload("Tags").
        Preload("Options").
        Preload("TestCases").
        First(&question, "id = ?", id)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &question, nil
}

func (r *questionRepository) UpdateQuestion(ctx context.Context, question *models.Question) error {
    result := r.db.WithContext(ctx).Save(question)
    return result.Error
}

func (r *questionRepository) DeleteQuestion(ctx context.Context, id string) error {
    result := r.db.WithContext(ctx).Delete(&models.Question{}, "id = ?", id)
    return result.Error
}

func (r *questionRepository) ListQuestions(ctx context.Context, filter QuestionFilter) ([]models.Question, int64, error) {
    var questions []models.Question
    var total int64
    
    query := r.db.WithContext(ctx).Model(&models.Question{})
    
    if filter.ID != "" {
        query = query.Where("id = ?", filter.ID)
    }
    
    if filter.CompetencyID != "" {
        query = query.Where("competency = ?", filter.CompetencyID)
    }
    
    if filter.Level != "" {
        query = query.Where("difficulty = ?", filter.Level)
    }
    
    if filter.Type != "" {
        query = query.Where("type = ?", filter.Type)
    }
    
    if filter.Search != "" {
        query = query.Where("title ILIKE ? OR description ILIKE ?", 
            "%"+filter.Search+"%", "%"+filter.Search+"%")
    }
    
    if filter.IsActive != nil {
        query = query.Where("is_active = ?", *filter.IsActive)
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
    result := query.Find(&questions)
    if result.Error != nil {
        return nil, 0, result.Error
    }
    
    return questions, total, nil
}

func (r *questionRepository) GetRandomQuestions(ctx context.Context, filter QuestionFilter, count int) ([]models.Question, error) {
    var questions []models.Question
    
    query := r.db.WithContext(ctx).Model(&models.Question{})
    
    if filter.CompetencyID != "" {
        query = query.Where("competency = ?", filter.CompetencyID)
    }
    
    if filter.Level != "" {
        query = query.Where("difficulty = ?", filter.Level)
    }
    
    if filter.Type != "" {
        query = query.Where("type = ?", filter.Type)
    }
    
    // Get random questions
    result := query.Where("is_active = ?", true).
        Order("RANDOM()").
        Limit(count).
        Find(&questions)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return questions, nil
}

func (r *questionRepository) GetQuestionsByCompetency(ctx context.Context, competencyID string, level string, limit int) ([]models.Question, error) {
    var questions []models.Question
    
    query := r.db.WithContext(ctx).
        Where("competency = ? AND difficulty = ? AND is_active = ?", 
            competencyID, level, true)
    
    if limit > 0 {
        query = query.Limit(limit)
    }
    
    result := query.Find(&questions)
    return questions, result.Error
}

func (r *questionRepository) BulkCreateQuestions(ctx context.Context, questions []models.Question) error {
    if len(questions) == 0 {
        return nil
    }
    
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        for i := range questions {
            if err := tx.Create(&questions[i]).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
