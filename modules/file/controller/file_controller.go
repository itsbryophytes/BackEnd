package controller

import (
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/modules/file/dto"
	"github.com/Caknoooo/go-gin-clean-starter/modules/file/service"
	"github.com/gin-gonic/gin"
)

type FileController interface {
	UploadFile(c *gin.Context)
	GetFiles(c *gin.Context)
	GetFile(c *gin.Context)
	DeleteFile(c *gin.Context)
}

type fileController struct {
	service service.FileService
}

func NewFileController(s service.FileService) FileController {
	return &fileController{s}
}

func (ctrl *fileController) UploadFile(c *gin.Context) {
	var req dto.UploadFileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")

	file, err := ctrl.service.Upload(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, file)
}

func (ctrl *fileController) GetFiles(c *gin.Context) {
	userID := c.GetString("user_id")

	files, err := ctrl.service.GetFiles(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, files)
}

func (ctrl *fileController) GetFile(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	file, err := ctrl.service.GetFileByID(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, file)
}

func (ctrl *fileController) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	err := ctrl.service.DeleteFile(userID, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
