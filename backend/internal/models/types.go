package models

// DifficultyLevel уровень сложности
type DifficultyLevel string

const (
    DifficultyJunior  DifficultyLevel = "junior"
    DifficultyMiddle  DifficultyLevel = "middle"
    DifficultySenior  DifficultyLevel = "senior"
    DifficultyExpert  DifficultyLevel = "expert"
)

// AssessmentStatus статус оценки
type AssessmentStatus string

const (
    AssessmentStatusDraft   AssessmentStatus = "draft"
    AssessmentStatusActive  AssessmentStatus = "active"
    AssessmentStatusArchived AssessmentStatus = "archived"
)

// SessionStatus статус сессии
type SessionStatus string

const (
    SessionStatusPending     SessionStatus = "pending"
    SessionStatusInProgress  SessionStatus = "in_progress"
    SessionStatusCompleted   SessionStatus = "completed"
    SessionStatusExpired     SessionStatus = "expired"
)

// AssessmentType тип оценки
type AssessmentType string

const (
    AssessmentTypeTechnical  AssessmentType = "technical"
    AssessmentTypeBehavioral AssessmentType = "behavioral"
    AssessmentTypeMixed      AssessmentType = "mixed"
)

// InvitationStatus статус приглашения
type InvitationStatus string

const (
    InvitationStatusPending  InvitationStatus = "pending"
    InvitationStatusSent     InvitationStatus = "sent"
    InvitationStatusOpened   InvitationStatus = "opened"
    InvitationStatusAccepted InvitationStatus = "accepted"
    InvitationStatusExpired  InvitationStatus = "expired"
)

// QuestionType тип вопроса
type QuestionType string

const (
    QuestionTypeMultipleChoice QuestionType = "multiple_choice"
    QuestionTypeCoding        QuestionType = "coding"
    QuestionTypeArchitecture  QuestionType = "architecture"
    QuestionTypeDebugging     QuestionType = "debugging"
)
