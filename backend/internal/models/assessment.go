package models

import (
    "time"
)

// Assessment представляет тестовое задание
type Assessment struct {
    BaseModel
    Title           string                 `gorm:"type:varchar(255);not null" json:"title"`
    Description     string                 `gorm:"type:text" json:"description"`
    Type            string                 `gorm:"type:varchar(50)" json:"type"`
    TargetLevel     string                 `gorm:"type:varchar(50)" json:"target_level"`
    TimeLimit       int                    `gorm:"not null;default:3600" json:"time_limit"` // в секундах
    TotalQuestions  int                    `gorm:"not null;default:20" json:"total_questions"`
    PassingScore    float64                `gorm:"default:70" json:"passing_score"`
    ShuffleQuestions bool                  `gorm:"default:true" json:"shuffle_questions"`
    ShowExplanation bool                   `gorm:"default:true" json:"show_explanation"`
    CreatedBy       string                 `gorm:"type:uuid;not null" json:"created_by"`
    Status          AssessmentStatus       `gorm:"type:varchar(20);default:'draft'" json:"status"`
    
    // Relationships
    Competencies   []AssessmentCompetency `gorm:"foreignKey:AssessmentID" json:"competencies"`
    Tags           []AssessmentTag       `gorm:"foreignKey:AssessmentID" json:"tags"`
    Questions      []AssessmentQuestion  `gorm:"foreignKey:AssessmentID" json:"questions"`
    Sessions       []AssessmentSession   `gorm:"foreignKey:AssessmentID" json:"sessions"`
    Invitations    []Invitation          `gorm:"foreignKey:AssessmentID" json:"invitations"`
}

// AssessmentQuestion связь между оценкой и вопросами
type AssessmentQuestion struct {
    BaseModel
    AssessmentID string `gorm:"type:uuid;not null;index" json:"assessment_id"`
    QuestionID   string `gorm:"type:uuid;not null;index" json:"question_id"`
    Order        int    `gorm:"not null" json:"order"`
    
    // Relationship
    Assessment Assessment `gorm:"foreignKey:AssessmentID"`
    Question   Question   `gorm:"foreignKey:QuestionID"`
}

// AssessmentCompetency связь компетенций с оценкой
type AssessmentCompetency struct {
    BaseModel
    AssessmentID   string  `gorm:"type:uuid;not null;index" json:"assessment_id"`
    CompetencyID   string  `gorm:"type:varchar(100);not null" json:"competency_id"`
    Weight         float64 `gorm:"default:1.0" json:"weight"`
    MinQuestions   int     `gorm:"default:1" json:"min_questions"`
    MaxQuestions   int     `gorm:"default:5" json:"max_questions"`
    
    // Relationship
    Assessment Assessment `gorm:"foreignKey:AssessmentID"`
}

// AssessmentTag теги оценки
type AssessmentTag struct {
    BaseModel
    AssessmentID string `gorm:"type:uuid;not null;index" json:"assessment_id"`
    Name         string `gorm:"type:varchar(100);not null" json:"name"`
    
    // Relationship
    Assessment Assessment `gorm:"foreignKey:AssessmentID"`
}

// AssessmentSession сессия кандидата
type AssessmentSession struct {
    BaseModel
    AssessmentID string         `gorm:"type:uuid;not null;index" json:"assessment_id"`
    CandidateID  string         `gorm:"type:uuid;not null;index" json:"candidate_id"`
    Status       SessionStatus  `gorm:"type:varchar(20);default:'pending'" json:"status"`
    StartedAt    *time.Time     `gorm:"type:timestamp" json:"started_at"`
    CompletedAt  *time.Time     `gorm:"type:timestamp" json:"completed_at"`
    TimeSpent    int            `gorm:"default:0" json:"time_spent"`
    Score        float64        `json:"score"`
    Percentage   float64        `json:"percentage"`
    Level        string         `gorm:"type:varchar(50)" json:"level"`
    
    // Relationships
    Assessment Assessment        `gorm:"foreignKey:AssessmentID"`
    Answers    []CandidateAnswer `gorm:"foreignKey:SessionID" json:"answers"`
    Result     *Result           `gorm:"foreignKey:SessionID" json:"result"`
}

// CandidateAnswer ответ кандидата
type CandidateAnswer struct {
    BaseModel
    SessionID    string     `gorm:"type:uuid;not null;index" json:"session_id"`
    QuestionID   string     `gorm:"type:uuid;not null;index" json:"question_id"`
    Answer       string     `gorm:"type:text" json:"answer"`
    Code         string     `gorm:"type:text" json:"code"`
    TimeSpent    int        `gorm:"not null" json:"time_spent"`
    StartedAt    time.Time  `gorm:"type:timestamp;not null" json:"started_at"`
    SubmittedAt  *time.Time `gorm:"type:timestamp" json:"submitted_at"`
    IsCorrect    bool       `gorm:"default:false" json:"is_correct"`
    Score        float64    `json:"score"`
    
    // Relationships
    Session  AssessmentSession `gorm:"foreignKey:SessionID"`
    Question Question          `gorm:"foreignKey:QuestionID"`
}

// Result результат оценки
type Result struct {
    BaseModel
    SessionID    string     `gorm:"type:uuid;not null;uniqueIndex" json:"session_id"`
    TotalScore   float64    `gorm:"not null" json:"total_score"`
    Percentage   float64    `gorm:"not null" json:"percentage"`
    Level        string     `gorm:"type:varchar(50);not null" json:"level"`
    TimeSpent    int        `gorm:"not null" json:"time_spent"`
    CompletedAt  time.Time  `gorm:"type:timestamp;not null" json:"completed_at"`
    
    // Relationships
    Session AssessmentSession `gorm:"foreignKey:SessionID"`
}

// Invitation приглашение кандидату
type Invitation struct {
    BaseModel
    AssessmentID string           `gorm:"type:uuid;not null;index" json:"assessment_id"`
    CandidateID  string           `gorm:"type:uuid;not null;index" json:"candidate_id"`
    Email        string           `gorm:"type:varchar(255);not null" json:"email"`
    Token        string           `gorm:"type:varchar(64);uniqueIndex;not null" json:"token"`
    Status       InvitationStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
    InvitedBy    string           `gorm:"type:uuid;not null" json:"invited_by"`
    OpenedAt     *time.Time       `gorm:"type:timestamp" json:"opened_at"`
    ExpiresAt    time.Time        `gorm:"type:timestamp;not null" json:"expires_at"`
    
    // Relationships
    Assessment Assessment `gorm:"foreignKey:AssessmentID"`
}

// ScoreResult результат оценки
type ScoreResult struct {
    TotalScore float64 `json:"total_score"`
    Percentage float64 `json:"percentage"`
    Level      string  `json:"level"`
}
