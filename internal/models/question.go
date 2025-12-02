package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeCoding         QuestionType = "coding"
	QuestionTypeArchitecture   QuestionType = "architecture"
	QuestionTypeDebugging      QuestionType = "debugging"
)

type QuestionLevel string

const (
	QuestionLevelJunior QuestionLevel = "junior"
	QuestionLevelMiddle QuestionLevel = "middle"
	QuestionLevelSenior QuestionLevel = "senior"
	QuestionLevelExpert QuestionLevel = "expert"
)

type ValidationStatus string

const (
	ValidationStatusPending   ValidationStatus = "pending"
	ValidationStatusApproved  ValidationStatus = "approved"
	ValidationStatusRejected  ValidationStatus = "rejected"
	ValidationStatusNeedsReview ValidationStatus = "needs_review"
)

type QuestionContent struct {
	Text        string   `json:"text"`
	CodeSnippet *string  `json:"code_snippet,omitempty"`
	Options     []string `json:"options,omitempty"`
}

type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Hidden         bool   `json:"hidden" default:"false"`
}

type QuestionSolution struct {
	Code        string `json:"code"`
	Explanation string `json:"explanation"`
}

type Question struct {
	BaseModelWithoutSoftDelete
	Type              QuestionType     `gorm:"type:varchar(50);not null" json:"type"`
	Competency        *string          `gorm:"type:varchar(100)" json:"competency,omitempty"`
	Level             QuestionLevel    `gorm:"type:varchar(50);not null" json:"level"`
	Content           datatypes.JSON   `gorm:"type:jsonb;not null" json:"content"`
	TestCases         datatypes.JSON   `gorm:"type:jsonb;default:'[]'" json:"test_cases"`
	Solution          datatypes.JSON   `gorm:"type:jsonb" json:"solution,omitempty"`
	Explanation       *string          `gorm:"type:text" json:"explanation,omitempty"`
	AIGenerated       bool             `gorm:"default:false" json:"ai_generated"`
	ValidationStatus  ValidationStatus `gorm:"type:varchar(50);default:'pending'" json:"validation_status"`
	ValidatedBy       *uuid.UUID       `gorm:"type:uuid" json:"validated_by,omitempty"`
	ValidatedAt       *time.Time       `json:"validated_at,omitempty"`
	Metadata          datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	
	// Relationships
	Assessments []Assessment `gorm:"many2many:assessment_questions;" json:"assessments,omitempty"`
	Results     []Result     `gorm:"foreignKey:QuestionID" json:"results,omitempty"`
}

func (Question) TableName() string {
	return "questions"
}
