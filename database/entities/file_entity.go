package entities

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID `gorm:"type:uuid;index"`
	FileURL      string
	FileName     string
	FileType     string
	FileSize     int
	IsPersistent bool      `gorm:"default:true"`
	UploadedAt   time.Time `gorm:"autoCreateTime"`
}
