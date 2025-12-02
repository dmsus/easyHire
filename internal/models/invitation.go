package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type InvitationStatus string

const (
	InvitationStatusSent     InvitationStatus = "sent"
	InvitationStatusOpened   InvitationStatus = "opened"
	InvitationStatusStarted  InvitationStatus = "started"
	InvitationStatusCompleted InvitationStatus = "completed"
	InvitationStatusExpired  InvitationStatus = "expired"
)

type Invitation struct {
	BaseModelWithoutSoftDelete
	AssessmentID uuid.UUID        `gorm:"type:uuid;not null" json:"assessment_id"`
	CandidateID  uuid.UUID        `gorm:"type:uuid;not null" json:"candidate_id"`
	InvitedBy    uuid.UUID        `gorm:"type:uuid;not null" json:"invited_by"`
	Token        string           `gorm:"type:varchar(100);not null;uniqueIndex" json:"token"`
	Status       InvitationStatus `gorm:"type:varchar(50);default:'sent'" json:"status"`
	EmailSent    bool             `gorm:"default:false" json:"email_sent"`
	SentAt       *time.Time       `json:"sent_at,omitempty"`
	ExpiresAt    *time.Time       `json:"expires_at,omitempty"`
	Metadata     datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	
	// Relationships
	Assessment *Assessment `gorm:"foreignKey:AssessmentID" json:"assessment,omitempty"`
	Candidate  *Candidate  `gorm:"foreignKey:CandidateID" json:"candidate,omitempty"`
	Inviter    *User       `gorm:"foreignKey:InvitedBy" json:"inviter,omitempty"`
}

func (Invitation) TableName() string {
	return "invitations"
}
