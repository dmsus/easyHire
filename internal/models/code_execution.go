package models

import (
	"github.com/google/uuid"
)

type ExecutionStatus string

const (
	ExecutionStatusPending  ExecutionStatus = "pending"
	ExecutionStatusRunning  ExecutionStatus = "running"
	ExecutionStatusSuccess  ExecutionStatus = "success"
	ExecutionStatusError    ExecutionStatus = "error"
	ExecutionStatusTimeout  ExecutionStatus = "timeout"
)

type CodeExecution struct {
	BaseModelWithoutSoftDelete
	ResultID           uuid.UUID       `gorm:"type:uuid;not null" json:"result_id"`
	QuestionID         uuid.UUID       `gorm:"type:uuid;not null" json:"question_id"`
	CandidateCode      string          `gorm:"type:text;not null" json:"candidate_code"`
	Language           string          `gorm:"type:varchar(50);default:'go'" json:"language"`
	Status             ExecutionStatus `gorm:"type:varchar(50);not null" json:"status"`
	ExitCode           *int            `json:"exit_code,omitempty"`
	Stdout             *string         `gorm:"type:text" json:"stdout,omitempty"`
	Stderr             *string         `gorm:"type:text" json:"stderr,omitempty"`
	ExecutionTimeMS    *int            `json:"execution_time_ms,omitempty"`
	MemoryUsedKB       *int            `json:"memory_used_kb,omitempty"`
	DockerContainerID  *string         `gorm:"type:varchar(100)" json:"docker_container_id,omitempty"`
	ErrorMessage       *string         `gorm:"type:text" json:"error_message,omitempty"`
	Logs               *string         `gorm:"type:text" json:"logs,omitempty"`
	
	// Relationships
	Result   *Result   `gorm:"foreignKey:ResultID" json:"result,omitempty"`
	Question *Question `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (CodeExecution) TableName() string {
	return "code_executions"
}
