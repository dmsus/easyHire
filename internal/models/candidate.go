package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CandidateStatus string

const (
	CandidateStatusActive    CandidateStatus = "active"
	CandidateStatusInvited   CandidateStatus = "invited"
	CandidateStatusCompleted CandidateStatus = "completed"
	CandidateStatusArchived  CandidateStatus = "archived"
)

type Candidate struct {
	BaseModelWithoutSoftDelete
	Email                       string           `gorm:"not null;size:255;uniqueIndex" json:"email"`
	FirstName                   *string          `gorm:"size:100" json:"first_name,omitempty"`
	LastName                    *string          `gorm:"size:100" json:"last_name,omitempty"`
	Phone                       *string          `gorm:"size:50" json:"phone,omitempty"`
	Status                      CandidateStatus  `gorm:"type:varchar(50);default:'active'" json:"status"`
	InvitedAt                   *time.Time       `json:"invited_at,omitempty"`
	CompletedAssessmentsCount   int              `gorm:"default:0" json:"completed_assessments_count"`
	AverageScore                *float64         `gorm:"type:decimal(5,2)" json:"average_score,omitempty"`
	LastAssessmentAt            *time.Time       `json:"last_assessment_at,omitempty"`
	Metadata                    datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	
	// Relationships
	Results      []Result     `gorm:"foreignKey:CandidateID" json:"results,omitempty"`
	Invitations  []Invitation `gorm:"foreignKey:CandidateID" json:"invitations,omitempty"`
}

func (Candidate) TableName() string {
	return "candidates"
}
