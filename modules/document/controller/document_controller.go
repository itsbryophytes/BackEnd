package controller

import (
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/rag"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentController interface {
	UploadDocument(c *gin.Context)
	CreateDocument(c *gin.Context)
	GetDocuments(c *gin.Context)
	GetPendingDocuments(c *gin.Context)
	GetDocument(c *gin.Context)
	ConfirmDocument(c *gin.Context)
	DiscardDocument(c *gin.Context)
	DeleteDocument(c *gin.Context)
}

type documentController struct {
	db        *gorm.DB
	ragClient rag.Client
}

func NewDocumentController(db *gorm.DB, ragClient rag.Client) DocumentController {
	return &documentController{db: db, ragClient: ragClient}
}

func (ctrl *documentController) UploadDocument(c *gin.Context) {
	userID := c.GetString("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	opened, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer opened.Close()

	contentType := file.Header.Get("Content-Type")
	resp, err := ctrl.ragClient.UploadDocument(c.Request.Context(), rag.UploadRequest{
		UserID:       userID,
		DocumentType: c.DefaultPostForm("document_type", "lab_result"),
		FileName:     file.Filename,
		ContentType:  contentType,
		File:         opened,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if documentID, ok := resp["document_id"].(string); ok {
		ctrl.upsertDocumentFromRAG(userID, documentID, file.Filename, c.DefaultPostForm("document_type", "lab_result"))
	}

	c.JSON(http.StatusCreated, resp)
}

func (ctrl *documentController) CreateDocument(c *gin.Context) {
	var document entities.Document
	if err := c.ShouldBindJSON(&document); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if document.ID == uuid.Nil {
		document.ID = uuid.New()
	}
	if document.UserID == uuid.Nil {
		if parsed, err := uuid.Parse(c.GetString("user_id")); err == nil {
			document.UserID = parsed
		}
	}

	if err := ctrl.db.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, document)
}

func (ctrl *documentController) GetDocuments(c *gin.Context) {
	userID := c.GetString("user_id")
	ragDocs, err := ctrl.ragClient.ListDocuments(c.Request.Context(), userID)
	if err == nil {
		c.JSON(http.StatusOK, ragDocs)
		return
	}

	var documents []entities.Document
	query := ctrl.db.Order("created_at DESC")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if dbErr := query.Find(&documents).Error; dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": dbErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents, "rag_error": err.Error()})
}

func (ctrl *documentController) GetPendingDocuments(c *gin.Context) {
	resp, err := ctrl.ragClient.ListStaging(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *documentController) GetDocument(c *gin.Context) {
	resp, err := ctrl.ragClient.GetStaging(c.Request.Context(), c.GetString("user_id"), c.Param("id"))
	if err == nil {
		c.JSON(http.StatusOK, resp)
		return
	}

	var document entities.Document
	query := ctrl.db.Where("id = ?", c.Param("id"))
	if userID := c.GetString("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if dbErr := query.First(&document).Error; dbErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": dbErr.Error(), "rag_error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, document)
}

func (ctrl *documentController) ConfirmDocument(c *gin.Context) {
	resp, err := ctrl.ragClient.ConfirmDocument(c.Request.Context(), c.GetString("user_id"), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *documentController) DiscardDocument(c *gin.Context) {
	removeFromRAG := c.DefaultQuery("remove_from_rag", "true") != "false"
	resp, err := ctrl.ragClient.DiscardDocument(c.Request.Context(), rag.DiscardRequest{
		UserID:        c.GetString("user_id"),
		DocumentID:    c.Param("id"),
		RemoveFromRAG: removeFromRAG,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *documentController) DeleteDocument(c *gin.Context) {
	userID := c.GetString("user_id")
	ragResp, ragErr := ctrl.ragClient.DeleteDocument(c.Request.Context(), userID, c.Param("id"))

	query := ctrl.db.Where("id = ?", c.Param("id"))
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Delete(&entities.Document{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if ragErr != nil {
		c.JSON(http.StatusOK, gin.H{"message": "deleted locally", "rag_error": ragErr.Error()})
		return
	}
	c.JSON(http.StatusOK, ragResp)
}

func (ctrl *documentController) upsertDocumentFromRAG(userID string, documentID string, filename string, documentType string) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return
	}
	parsedDocumentID, err := uuid.Parse(documentID)
	if err != nil {
		return
	}

	document := entities.Document{
		ID:           parsedDocumentID,
		UserID:       parsedUserID,
		Title:        filename,
		DocumentType: documentType,
		Status:       "pending_confirmation",
	}
	ctrl.db.Where("id = ?", parsedDocumentID).Assign(document).FirstOrCreate(&document)
}
