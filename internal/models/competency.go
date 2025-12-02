package models

import (
	"time"
)

type Competency struct {
	Name        string    `gorm:"type:varchar(100);primaryKey" json:"name"`
	Category    string    `gorm:"type:varchar(50)" json:"category,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	BaseWeight  float64   `gorm:"type:decimal(3,2);default:1.0" json:"base_weight"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
}

func (Competency) TableName() string {
	return "competencies"
}
