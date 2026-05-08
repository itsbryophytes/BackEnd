package entities

import (
	"time"

	"github.com/google/uuid"
)

type OCRResult struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DocumentID uuid.UUID
	RawText    string
	Confidence float64
	Language   string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
