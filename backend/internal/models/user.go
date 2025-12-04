package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type UserRole string

const (
    UserRoleAdmin      UserRole = "admin"
    UserRoleRecruiter  UserRole = "recruiter"
    UserRoleCandidate  UserRole = "candidate"
    UserRoleExpert     UserRole = "expert"
)

type User struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
    
    Email        string   `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string   `gorm:"not null" json:"-"`
    FirstName    string   `gorm:"type:varchar(100)" json:"first_name"`
    LastName     string   `gorm:"type:varchar(100)" json:"last_name"`
    Role         UserRole `gorm:"type:varchar(20);not null;default:'candidate'" json:"role"`
    
    // Для кандидатов
    CurrentCompany string `gorm:"type:varchar(255)" json:"current_company"`
    JobTitle       string `gorm:"type:varchar(255)" json:"job_title"`
    Experience     int    `json:"experience"`
    
    // Для рекрутеров
    Company   string `gorm:"type:varchar(255)" json:"company"`
    Position  string `gorm:"type:varchar(255)" json:"position"`
    IsActive  bool   `gorm:"default:true" json:"is_active"`
    LastLogin *time.Time `json:"last_login"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

// User DTOs
type UserCreateRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Role      string `json:"role"`
    Company   string `json:"company,omitempty"`
}

type UserUpdateRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Company   string `json:"company"`
    Position  string `json:"position"`
    IsActive  *bool  `json:"is_active"`
}
