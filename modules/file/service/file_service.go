package service

import (
	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/Caknoooo/go-gin-clean-starter/modules/file/dto"
	"github.com/Caknoooo/go-gin-clean-starter/modules/file/repository"
	"github.com/google/uuid"
)

type FileService interface {
	Upload(userID string, req dto.UploadFileRequest) (*entities.File, error)
	GetFiles(userID string) ([]entities.File, error)
	GetFileByID(userID, fileID string) (*entities.File, error)
	DeleteFile(userID, fileID string) error
}

type fileService struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return &fileService{repo}
}

func (s *fileService) Upload(userID string, req dto.UploadFileRequest) (*entities.File, error) {

	file := &entities.File{
		ID:       uuid.New(),
		UserID:   uuid.MustParse(userID),
		FileName: req.FileName,
		FileURL:  req.FileURL,
	}

	err := s.repo.Create(file)
	return file, err
}

func (s *fileService) GetFiles(userID string) ([]entities.File, error) {
	return s.repo.FindByUser(userID)
}

func (s *fileService) GetFileByID(userID, fileID string) (*entities.File, error) {
	return s.repo.FindByIDAndUser(fileID, userID)
}

func (s *fileService) DeleteFile(userID, fileID string) error {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return err
	}

	return s.repo.Delete(file.ID.String())
}
