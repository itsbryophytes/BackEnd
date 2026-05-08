package controller

import (
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/rag"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProcessingController interface {
	RunOCR(c *gin.Context)
	GetStatus(c *gin.Context)
	GetResult(c *gin.Context)
}

type processingController struct {
	db        *gorm.DB
	ragClient rag.Client
}

func NewProcessingController(db *gorm.DB, ragClient rag.Client) ProcessingController {
	return &processingController{db: db, ragClient: ragClient}
}

func (ctrl *processingController) RunOCR(c *gin.Context) {
	resp, err := ctrl.ragClient.GetStaging(c.Request.Context(), c.GetString("user_id"), c.Param("document_id"))
	if err == nil {
		c.JSON(http.StatusOK, resp)
		return
	}

	if dbErr := ctrl.db.Model(&entities.Document{}).
		Where("id = ?", c.Param("document_id")).
		Update("status", "processing").Error; dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": dbErr.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "processing status updated locally", "rag_error": err.Error()})
}

func (ctrl *processingController) GetStatus(c *gin.Context) {
	resp, err := ctrl.ragClient.GetStaging(c.Request.Context(), c.GetString("user_id"), c.Param("document_id"))
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": resp["status"], "document": resp})
		return
	}

	var document entities.Document
	if dbErr := ctrl.db.First(&document, "id = ?", c.Param("document_id")).Error; dbErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": dbErr.Error(), "rag_error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": document.Status})
}

func (ctrl *processingController) GetResult(c *gin.Context) {
	resp, err := ctrl.ragClient.GetStaging(c.Request.Context(), c.GetString("user_id"), c.Param("document_id"))
	if err == nil {
		c.JSON(http.StatusOK, resp)
		return
	}

	var result entities.OCRResult
	if dbErr := ctrl.db.First(&result, "document_id = ?", c.Param("document_id")).Error; dbErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": dbErr.Error(), "rag_error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
