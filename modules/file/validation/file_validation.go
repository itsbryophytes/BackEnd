package validation

import (
	"errors"

	"github.com/Caknoooo/go-gin-clean-starter/modules/file/dto"
)

func ValidateUpload(req dto.UploadFileRequest) error {
	if req.FileName == "" {
		return errors.New("file_name is required")
	}
	if req.FileURL == "" {
		return errors.New("file_url is required")
	}
	return nil
}
