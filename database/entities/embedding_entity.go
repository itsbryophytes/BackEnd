package entities

import (
	"time"

	"github.com/google/uuid"
)

type Embedding struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID
	DocumentID uuid.UUID
	Content    string
	Embedding  []float32 `gorm:"type:vector(1536)"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
