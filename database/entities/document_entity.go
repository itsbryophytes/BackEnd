package entities

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID `gorm:"type:uuid;index"`
	FileID       uuid.UUID `gorm:"type:uuid"`
	DocumentType string
	Title        string
	TakenAt      *time.Time
	Status       string    `gorm:"default:processing"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
