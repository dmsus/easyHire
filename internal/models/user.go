package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserRole string

const (
	RoleAdmin          UserRole = "admin"
	RoleHR             UserRole = "hr"
	RoleCandidate      UserRole = "candidate"
	RoleTechnicalExpert UserRole = "technical_expert"
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
	Metadata     datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	
	// Relationships
	Assessments []Assessment `gorm:"foreignKey:CreatedBy" json:"assessments,omitempty"`
	Invitations []Invitation `gorm:"foreignKey:InvitedBy" json:"invitations,omitempty"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// User DTOs for API requests/responses
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Company  string `json:"company,omitempty"`
	Role     string `json:"role,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      UserRole  `json:"role"`
	Company   *string   `json:"company,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	IsActive  bool      `json:"is_active"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type UpdateProfileRequest struct {
	Name      string `json:"name,omitempty"`
	Company   string `json:"company,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}
