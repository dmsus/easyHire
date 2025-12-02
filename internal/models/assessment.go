package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AssessmentStatus string

const (
	AssessmentStatusDraft     AssessmentStatus = "draft"
	AssessmentStatusPublished AssessmentStatus = "published"
	AssessmentStatusActive    AssessmentStatus = "active"
	AssessmentStatusCompleted AssessmentStatus = "completed"
	AssessmentStatusArchived  AssessmentStatus = "archived"
)

type AssessmentRole string

const (
	RoleBackendDeveloper AssessmentRole = "backend_developer"
	RoleFullstackDeveloper AssessmentRole = "fullstack_developer"
	RoleDevOpsEngineer   AssessmentRole = "devops_engineer"
	RoleTeamLead         AssessmentRole = "team_lead"
)

type AssessmentLevel string

const (
	LevelJunior AssessmentLevel = "junior"
	LevelMiddle AssessmentLevel = "middle"
	LevelSenior AssessmentLevel = "senior"
	LevelExpert AssessmentLevel = "expert"
)

type Assessment struct {
	BaseModel
	Title              string                `gorm:"not null;size:255" json:"title"`
	Description        *string               `gorm:"type:text" json:"description,omitempty"`
	CreatedBy          uuid.UUID             `gorm:"type:uuid" json:"created_by,omitempty"`
	Status             AssessmentStatus      `gorm:"type:varchar(50);not null;default:'draft'" json:"status"`
	Role               *AssessmentRole       `gorm:"type:varchar(50)" json:"role,omitempty"`
	TargetLevel        *AssessmentLevel      `gorm:"type:varchar(50)" json:"target_level,omitempty"`
	TimeLimitMinutes   *int                  `gorm:"check:time_limit_minutes >= 1" json:"time_limit_minutes,omitempty"`
	QuestionCount      *int                  `gorm:"check:question_count >= 1" json:"question_count,omitempty"`
	PassingScore       *float64              `gorm:"type:decimal(5,2);check:passing_score >= 0 AND passing_score <= 100" json:"passing_score,omitempty"`
	CompetencyWeights  datatypes.JSONMap     `gorm:"type:jsonb;default:'{}'" json:"competency_weights"`
	LevelDistribution  datatypes.JSONMap     `gorm:"type:jsonb;default:'{}'" json:"level_distribution"`
	Metadata           datatypes.JSONMap     `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	PublishedAt        *time.Time            `json:"published_at,omitempty"`
	
	// Relationships
	Competencies []Competency         `gorm:"many2many:assessment_competencies;" json:"competencies,omitempty"`
	Questions    []Question           `gorm:"many2many:assessment_questions;" json:"questions,omitempty"`
	Results      []Result             `gorm:"foreignKey:AssessmentID" json:"results,omitempty"`
	Invitations  []Invitation         `gorm:"foreignKey:AssessmentID" json:"invitations,omitempty"`
}

func (Assessment) TableName() string {
	return "assessments"
}

type AssessmentCompetency struct {
	AssessmentID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"assessment_id"`
	CompetencyName string    `gorm:"type:varchar(100);primaryKey" json:"competency_name"`
	Weight         float64   `gorm:"type:decimal(5,2);default:1.0" json:"weight"`
}

func (AssessmentCompetency) TableName() string {
	return "assessment_competencies"
}

type AssessmentQuestion struct {
	AssessmentID uuid.UUID `gorm:"type:uuid;primaryKey" json:"assessment_id"`
	QuestionID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"question_id"`
	OrderIndex   int       `gorm:"default:0" json:"order_index"`
	Points       int       `gorm:"default:10" json:"points"`
}

func (AssessmentQuestion) TableName() string {
	return "assessment_questions"
}
