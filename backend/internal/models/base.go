package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel содержит общие поля для всех моделей
type BaseModel struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp;index" json:"deleted_at,omitempty"`
}

// BeforeCreate хука для GORM
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = GenerateUUID()
	}
	return nil
}

// GenerateUUID генерирует UUID (заглушка, в реальности используйте github.com/google/uuid)
func GenerateUUID() string {
	// В реальном приложении используйте:
	// return uuid.New().String()
	return "generated-uuid-placeholder"
}
