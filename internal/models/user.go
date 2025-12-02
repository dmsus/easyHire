package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleRecruiter UserRole = "recruiter"
	RoleCandidate UserRole = "candidate"
)

type User struct {
	BaseModel
	Email        string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string     `gorm:"not null;size:255" json:"-"`
	Name         string     `gorm:"not null;size:255" json:"name"`
	Role         UserRole   `gorm:"type:varchar(50);not null;default:'candidate'" json:"role"`
	Company      *string    `gorm:"size:255" json:"company,omitempty"`
	AvatarURL    *string    `gorm:"type:text" json:"avatar_url,omitempty"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}
