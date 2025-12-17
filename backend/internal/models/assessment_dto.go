package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CreateAssessmentRequest запрос на создание оценки
type CreateAssessmentRequest struct {
	Title            string             `json:"title" binding:"required,min=3,max=255"`
	Description      string             `json:"description"`
	Type             AssessmentType     `json:"type" binding:"required"`
	TargetLevel      DifficultyLevel    `json:"target_level" binding:"required,oneof=junior middle senior expert"`
	TimeLimit        int                `json:"time_limit" binding:"min=300,max=10800"` // seconds: 5 min - 3 hours
	TotalQuestions   int                `json:"total_questions" binding:"min=1,max=100"`
	PassingScore     float64            `json:"passing_score" binding:"min=0,max=100"`
	ShuffleQuestions bool               `json:"shuffle_questions"`
	ShowExplanation  bool               `json:"show_explanation"`
	Competencies     []CompetencyWeight `json:"competencies" binding:"required,min=1"`
	Tags             []string           `json:"tags"`
}

// CompetencyWeight вес компетенции в оценке
//
// ✅ Поддерживает оба формата:
// 1) Новый: {"code":"go_fundamentals","level":"middle","weight":1.2,...}
// 2) Старый: {"competency_id":"go_fundamentals","level":"middle","weight":1.2,...}
type CompetencyWeight struct {
	CompetencyID string          `json:"competency_id" binding:"required,min=2,max=100"`
	Level        DifficultyLevel `json:"level" binding:"required,oneof=junior middle senior expert"`
	Weight       float64         `json:"weight" binding:"min=0.1,max=5.0"`
	MinQuestions int             `json:"min_questions" binding:"min=0,max=100"`
	MaxQuestions int             `json:"max_questions" binding:"min=0,max=100"`
}

// Alias-field: code (для входящих payload)
// В ответе мы всё равно используем competency_id, чтобы не плодить форматы.
type competencyWeightWire struct {
	CompetencyID string          `json:"competency_id"`
	Code         string          `json:"code"`
	Level        DifficultyLevel `json:"level"`
	Weight       float64         `json:"weight"`
	MinQuestions int             `json:"min_questions"`
	MaxQuestions int             `json:"max_questions"`
}

func (c *CompetencyWeight) UnmarshalJSON(data []byte) error {
	var w competencyWeightWire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}

	// Accept both "competency_id" and "code"
	id := strings.TrimSpace(w.CompetencyID)
	if id == "" {
		id = strings.TrimSpace(w.Code)
	}
	w.CompetencyID = id

	// Basic sanity here; validator will run afterwards too.
	if w.CompetencyID == "" {
		return fmt.Errorf("competency_id (or code) is required")
	}
	if w.Level == "" {
		return fmt.Errorf("level is required")
	}

	c.CompetencyID = w.CompetencyID
	c.Level = w.Level
	c.Weight = w.Weight
	c.MinQuestions = w.MinQuestions
	c.MaxQuestions = w.MaxQuestions
	return nil
}

// UpdateAssessmentRequest запрос на обновление оценки
type UpdateAssessmentRequest struct {
	Title            *string           `json:"title"`
	Description      *string           `json:"description"`
	Status           *AssessmentStatus `json:"status"`
	TimeLimit        *int              `json:"time_limit"`
	PassingScore     *float64          `json:"passing_score"`
	ShuffleQuestions *bool             `json:"shuffle_questions"`
	ShowExplanation  *bool             `json:"show_explanation"`
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
	ID             string             `json:"id"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	Status         AssessmentStatus   `json:"status"`
	Type           AssessmentType     `json:"type"`
	TargetLevel    DifficultyLevel    `json:"target_level"`
	TimeLimit      int                `json:"time_limit"`
	TotalQuestions int                `json:"total_questions"`
	PassingScore   float64            `json:"passing_score"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	Competencies   []CompetencyWeight `json:"competencies"`
	Tags           []string           `json:"tags"`
}

// SessionProgress прогресс сессии
type SessionProgress struct {
	SessionID       string        `json:"session_id"`
	AssessmentID    string        `json:"assessment_id"`
	Status          SessionStatus `json:"status"`
	CurrentQuestion int           `json:"current_question"`
	TotalQuestions  int           `json:"total_questions"`
	TimeSpent       int           `json:"time_spent"`
	TimeRemaining   int           `json:"time_remaining"`
	StartedAt       *time.Time    `json:"started_at"`
	CompletedAt     *time.Time    `json:"completed_at"`
}

// AssessmentReport детальный отчет по оценке
type AssessmentReport struct {
	Assessment          AssessmentResponse `json:"assessment"`
	Result              Result             `json:"result"`
	CompetencyBreakdown []CompetencyScore  `json:"competency_breakdown"`
	Answers             []AnswerDetail     `json:"answers,omitempty"`
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
	QuestionID    string  `json:"question_id"`
	QuestionTitle string  `json:"question_title"`
	QuestionType  string  `json:"question_type"`
	Answer        string  `json:"answer"`
	CorrectAnswer string  `json:"correct_answer"`
	IsCorrect     bool    `json:"is_correct"`
	PointsEarned  float64 `json:"points_earned"`
	TimeSpent     int     `json:"time_spent"`
	Explanation   string  `json:"explanation,omitempty"`
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
