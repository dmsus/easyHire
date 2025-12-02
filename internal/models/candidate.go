package models

import (
	"time"
)

type Candidate struct {
	BaseModel
	Email       string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	FirstName   string     `gorm:"size:100" json:"first_name"`
	LastName    string     `gorm:"size:100" json:"last_name"`
	Phone       *string    `gorm:"size:20" json:"phone,omitempty"`
	LinkedInURL *string    `gorm:"type:text" json:"linkedin_url,omitempty"`
	GitHubURL   *string    `gorm:"type:text" json:"github_url,omitempty"`
	ResumeURL   *string    `gorm:"type:text" json:"resume_url,omitempty"`
	Status      string     `gorm:"size:50;default:'new'" json:"status"`
	InvitedAt   *time.Time `json:"invited_at,omitempty"`
	LastContact *time.Time `json:"last_contact,omitempty"`
	Notes       *string    `gorm:"type:text" json:"notes,omitempty"`
	Metadata    map[string]interface{} `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	
	// Relationships
	Invitations []Invitation `gorm:"foreignKey:CandidateID" json:"invitations,omitempty"`
	Assessments []Assessment `gorm:"many2many:candidate_assessments;" json:"assessments,omitempty"`
	Results     []Result     `gorm:"foreignKey:CandidateID" json:"results,omitempty"`
}

func (Candidate) TableName() string {
	return "candidates"
}
