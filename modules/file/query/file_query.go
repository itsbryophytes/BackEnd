package query

import (
	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"gorm.io/gorm"
)

type FileQuery struct {
	db *gorm.DB
}

func NewFileQuery(db *gorm.DB) *FileQuery {
	return &FileQuery{db}
}

func (q *FileQuery) GetUserFiles(userID string) ([]entities.File, error) {
	var files []entities.File
	err := q.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&files).Error

	return files, err
}
