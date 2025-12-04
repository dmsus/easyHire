package models

import (
	"time"
)

// CreateAssessmentRequest запрос на создание оценки
type CreateAssessmentRequest struct {
	Title           string                `json:"title" binding:"required,min=3,max=255"`
	Description     string                `json:"description"`
	Type            AssessmentType        `json:"type" binding:"required"`
	TargetLevel     DifficultyLevel       `json:"target_level" binding:"required"`
	TimeLimit       int                   `json:"time_limit" binding:"min=300,max=10800"` // 5 мин - 3 часа
	TotalQuestions  int                   `json:"total_questions" binding:"min=1,max=100"`
	PassingScore    float64               `json:"passing_score" binding:"min=0,max=100"`
	ShuffleQuestions bool                 `json:"shuffle_questions"`
	ShowExplanation bool                  `json:"show_explanation"`
	Competencies    []CompetencyWeight    `json:"competencies" binding:"required,min=1"`
	Tags            []string              `json:"tags"`
}

// CompetencyWeight вес компетенции в оценке
type CompetencyWeight struct {
	CompetencyID string  `json:"competency_id" binding:"required,uuid"`
	Weight       float64 `json:"weight" binding:"min=0.1,max=5.0"`
	MinQuestions int     `json:"min_questions" binding:"min=0,max=20"`
	MaxQuestions int     `json:"max_questions" binding:"min=0,max=20"`
}

// UpdateAssessmentRequest запрос на обновление оценки
type UpdateAssessmentRequest struct {
	Title           *string               `json:"title"`
	Description     *string               `json:"description"`
	Status          *AssessmentStatus     `json:"status"`
	TimeLimit       *int                  `json:"time_limit"`
	PassingScore    *float64              `json:"passing_score"`
	ShuffleQuestions *bool                `json:"shuffle_questions"`
	ShowExplanation *bool                 `json:"show_explanation"`
}

// InviteCandidatesRequest запрос на приглашение кандидатов
type InviteCandidatesRequest struct {
	Emails  []string `json:"emails" binding:"required,min=1"`
	Message string   `json:"message"`
}

// CandidateAnswerRequest запрос с ответом кандидата
type CandidateAnswerRequest struct {
	Answer    string `json:"answer"`
	Code      string `json:"code"`
	TimeSpent int    `json:"time_spent" binding:"min=0"`
}

// AssessmentResponse ответ с данными оценки
type AssessmentResponse struct {
	ID              string                `json:"id"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	Status          AssessmentStatus      `json:"status"`
	Type            AssessmentType        `json:"type"`
	TargetLevel     DifficultyLevel       `json:"target_level"`
	TimeLimit       int                   `json:"time_limit"`
	TotalQuestions  int                   `json:"total_questions"`
	PassingScore    float64               `json:"passing_score"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	Competencies    []CompetencyWeight    `json:"competencies"`
	Tags            []string              `json:"tags"`
}

// SessionProgress прогресс сессии
type SessionProgress struct {
	SessionID       string                `json:"session_id"`
	AssessmentID    string                `json:"assessment_id"`
	Status          SessionStatus         `json:"status"`
	CurrentQuestion int                   `json:"current_question"`
	TotalQuestions  int                   `json:"total_questions"`
	TimeSpent       int                   `json:"time_spent"`
	TimeRemaining   int                   `json:"time_remaining"`
	StartedAt       *time.Time            `json:"started_at"`
	CompletedAt     *time.Time            `json:"completed_at"`
}

// AssessmentReport детальный отчет по оценке
type AssessmentReport struct {
	Assessment      AssessmentResponse    `json:"assessment"`
	Result          Result                `json:"result"`
	CompetencyBreakdown []CompetencyScore `json:"competency_breakdown"`
	Answers         []AnswerDetail        `json:"answers,omitempty"`
}

// CompetencyScore счет по компетенции
type CompetencyScore struct {
	CompetencyID   string  `json:"competency_id"`
	CompetencyName string  `json:"competency_name"`
	Achieved       float64 `json:"achieved"`
	Possible       float64 `json:"possible"`
	Percentage     float64 `json:"percentage"`
}

// AnswerDetail детали ответа
type AnswerDetail struct {
	QuestionID     string    `json:"question_id"`
	QuestionTitle  string    `json:"question_title"`
	QuestionType   string    `json:"question_type"`
	Answer         string    `json:"answer"`
	CorrectAnswer  string    `json:"correct_answer"`
	IsCorrect      bool      `json:"is_correct"`
	PointsEarned   float64   `json:"points_earned"`
	TimeSpent      int       `json:"time_spent"`
	Explanation    string    `json:"explanation,omitempty"`
}

// BulkInviteRequest запрос на массовое приглашение
type BulkInviteRequest struct {
	Candidates []CandidateInfo `json:"candidates" binding:"required,min=1"`
}

// CandidateInfo информация о кандидате для массового приглашения
type CandidateInfo struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// BulkInviteResult результат массового приглашения
type BulkInviteResult struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	FailedEmails []string `json:"failed_emails"`
}
