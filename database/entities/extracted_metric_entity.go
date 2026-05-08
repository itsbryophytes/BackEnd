package entities

import (
	"time"

	"github.com/google/uuid"
)

type ExtractedMetric struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DocumentID     uuid.UUID
	Name           string
	Value          float64
	Unit           string
	ReferenceRange string
	Flag           string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}
