package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ResultLevel string

const (
	ResultLevelTrainee ResultLevel = "trainee"
	ResultLevelJunior  ResultLevel = "junior"
	ResultLevelMiddle  ResultLevel = "middle"
	ResultLevelSenior  ResultLevel = "senior"
	ResultLevelExpert  ResultLevel = "expert"
)

type CompetencyBreakdown struct {
	Competency string  `json:"competency"`
	Score      float64 `json:"score"`
	Level      string  `json:"level"`
}

type QuestionResult struct {
	QuestionID       uuid.UUID `json:"question_id"`
	Correct          bool      `json:"correct"`
	Score            float64   `json:"score"`
	TimeSpentSeconds int       `json:"time_spent_seconds"`
}

type Feedback struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Severity string `json:"severity"` // info, warning, critical
}

type Result struct {
	BaseModelWithoutSoftDelete
	AssessmentID       uuid.UUID            `gorm:"type:uuid;not null" json:"assessment_id"`
	CandidateID        uuid.UUID            `gorm:"type:uuid;not null" json:"candidate_id"`
	Score              float64              `gorm:"type:decimal(5,2);not null;check:score >= 0 AND score <= 100" json:"score"`
	FibonacciScore     *float64             `gorm:"type:decimal(5,2)" json:"fibonacci_score,omitempty"`
	Level              ResultLevel          `gorm:"type:varchar(50);not null" json:"level"`
	Passed             bool                 `gorm:"not null" json:"passed"`
	TimeSpentMinutes   *int                 `json:"time_spent_minutes,omitempty"`
	StartedAt          *time.Time           `json:"started_at,omitempty"`
	CompletedAt        *time.Time           `json:"completed_at,omitempty"`
	CompetencyBreakdown datatypes.JSON      `gorm:"type:jsonb;default:'[]'" json:"competency_breakdown"`
	QuestionResults     datatypes.JSON      `gorm:"type:jsonb;default:'[]'" json:"question_results"`
	Feedback           datatypes.JSON       `gorm:"type:jsonb;default:'[]'" json:"feedback"`
	Metadata           datatypes.JSONMap    `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	
	// Relationships
	Assessment        *Assessment          `gorm:"foreignKey:AssessmentID" json:"assessment,omitempty"`
	Candidate         *Candidate           `gorm:"foreignKey:CandidateID" json:"candidate,omitempty"`
	CodeExecutions    []CodeExecution      `gorm:"foreignKey:ResultID" json:"code_executions,omitempty"`
}

func (Result) TableName() string {
	return "results"
}
