package entities

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title         string
	Content       string
	CoverImageURL string
	AuthorID      uuid.UUID
	Status        string `gorm:"default:draft"`
	PublishedAt   *time.Time
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
