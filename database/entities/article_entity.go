package entities

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	CoverImageURL string     `json:"cover_image_url"`
	AuthorID      string     `json:"author_id"`
	Status        string     `gorm:"default:draft" json:"status"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
