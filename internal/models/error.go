package models

import (
	"time"

	"gorm.io/datatypes"
)

type APIErrorCode string

const (
	ErrorCodeValidationError APIErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound        APIErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized    APIErrorCode = "UNAUTHORIZED"
	ErrorCodeInternalError   APIErrorCode = "INTERNAL_ERROR"
	ErrorCodeBadRequest      APIErrorCode = "BAD_REQUEST"
)

type APIError struct {
	Error     bool         `json:"error" example:"true"`
	Code      APIErrorCode `json:"code" example:"VALIDATION_ERROR"`
	Message   string       `json:"message" example:"Invalid request parameters"`
	Details   interface{}  `json:"details,omitempty" gorm:"-"`
	RequestID *string      `json:"request_id,omitempty" example:"req_123456789"`
	Timestamp time.Time    `json:"timestamp" example:"2024-01-15T10:30:00Z"`
	Path      string       `json:"path,omitempty" gorm:"-"`
	Status    int          `json:"status,omitempty" gorm:"-"`
}

// For database storage (if needed)
type StoredError struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	Error     bool           `json:"error"`
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Details   datatypes.JSON `json:"details,omitempty"`
	RequestID *string        `json:"request_id,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
	Path      string         `json:"path,omitempty"`
	Status    int            `json:"status,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

func (StoredError) TableName() string {
	return "error_logs"
}

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Errors []ValidationErrorDetail `json:"errors"`
}
