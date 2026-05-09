package controller

import (
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/rag"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MetricController interface {
	GetMetrics(c *gin.Context)
	ConfirmMetrics(c *gin.Context)
	UpdateMetric(c *gin.Context)
	DeleteMetric(c *gin.Context)
}

type metricController struct {
	db        *gorm.DB
	ragClient rag.Client
}

func NewMetricController(db *gorm.DB, ragClient rag.Client) MetricController {
	return &metricController{db: db, ragClient: ragClient}
}

func (ctrl *metricController) GetMetrics(c *gin.Context) {
	userID := c.GetString("user_id")
	if results, err := ctrl.ragClient.GetResults(c.Request.Context(), userID); err == nil {
		c.JSON(http.StatusOK, results)
		return
	}

	var metrics []entities.HealthMetric
	query := ctrl.db.Order("measured_at DESC")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": metrics})
}

func (ctrl *metricController) ConfirmMetrics(c *gin.Context) {
	documentID := c.Query("document_id")
	if documentID == "" {
		var req struct {
			DocumentID string `json:"document_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		documentID = req.DocumentID
	}

	resp, err := ctrl.ragClient.ConfirmDocument(c.Request.Context(), c.GetString("user_id"), documentID, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ctrl *metricController) UpdateMetric(c *gin.Context) {
	var metric entities.HealthMetric
	if err := ctrl.db.First(&metric, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var req entities.HealthMetric
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	metric.MetricType = req.MetricType
	metric.Value = req.Value
	metric.Unit = req.Unit
	metric.MeasuredAt = req.MeasuredAt

	if err := ctrl.db.Save(&metric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metric)
}

func (ctrl *metricController) DeleteMetric(c *gin.Context) {
	if err := ctrl.db.Delete(&entities.HealthMetric{}, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
