package entities

import (
	"time"

	"github.com/google/uuid"
)

type HealthMetric struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID           uuid.UUID
	MetricType       string
	Value            float64
	Unit             string
	MeasuredAt       time.Time
	SourceDocumentID uuid.UUID
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}
