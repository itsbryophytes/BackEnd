package repository

import (
	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(file *entities.File) error
	FindByUser(userID string) ([]entities.File, error)
	FindByIDAndUser(fileID, userID string) (*entities.File, error)
	Delete(id string) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db}
}

func (r *fileRepository) Create(file *entities.File) error {
	return r.db.Create(file).Error
}

func (r *fileRepository) FindByUser(userID string) ([]entities.File, error) {
	var files []entities.File
	err := r.db.Where("user_id = ?", userID).Find(&files).Error
	return files, err
}

func (r *fileRepository) FindByIDAndUser(fileID, userID string) (*entities.File, error) {
	var file entities.File
	err := r.db.Where("id = ? AND user_id = ?", fileID, userID).
		First(&file).Error

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *fileRepository) Delete(id string) error {
	return r.db.Delete(&entities.File{}, "id = ?", id).Error
}
