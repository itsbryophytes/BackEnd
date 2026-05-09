package controller

import (
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/modules/health_log/dto"
	"github.com/Caknoooo/go-gin-clean-starter/modules/health_log/service"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HealthLogController interface {
	CreateBloodSugar(c *gin.Context)
	GetBloodSugarLogs(c *gin.Context)
	DeleteBloodSugar(c *gin.Context)

	CreateBloodPressure(c *gin.Context)
	GetBloodPressureLogs(c *gin.Context)
	DeleteBloodPressure(c *gin.Context)

	CreateWeight(c *gin.Context)
	GetWeightLogs(c *gin.Context)
	DeleteWeight(c *gin.Context)
}

type healthLogController struct {
	service service.HealthLogService
}

func NewHealthLogController(service service.HealthLogService) HealthLogController {
	return &healthLogController{service: service}
}

func (ctrl *healthLogController) CreateBloodSugar(c *gin.Context) {
	var req dto.CreateBloodSugarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid request", err.Error(), nil))
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))
	res, err := ctrl.service.CreateBloodSugar(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create blood sugar log", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Blood sugar log created", res))
}

func (ctrl *healthLogController) GetBloodSugarLogs(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	res, err := ctrl.service.GetBloodSugarLogs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to get logs", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Blood sugar logs retrieved", res))
}

func (ctrl *healthLogController) DeleteBloodSugar(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", err.Error(), nil))
		return
	}

	if err := ctrl.service.DeleteBloodSugar(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete log", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Blood sugar log deleted", nil))
}

func (ctrl *healthLogController) CreateBloodPressure(c *gin.Context) {
	var req dto.CreateBloodPressureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid request", err.Error(), nil))
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))
	res, err := ctrl.service.CreateBloodPressure(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create BP log", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Blood pressure log created", res))
}

func (ctrl *healthLogController) GetBloodPressureLogs(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	res, err := ctrl.service.GetBloodPressureLogs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to get logs", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Blood pressure logs retrieved", res))
}

func (ctrl *healthLogController) DeleteBloodPressure(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", err.Error(), nil))
		return
	}

	if err := ctrl.service.DeleteBloodPressure(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete log", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Blood pressure log deleted", nil))
}

func (ctrl *healthLogController) CreateWeight(c *gin.Context) {
	var req dto.CreateWeightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid request", err.Error(), nil))
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))
	res, err := ctrl.service.CreateWeight(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to create weight log", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Weight log created", res))
}

func (ctrl *healthLogController) GetWeightLogs(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	res, err := ctrl.service.GetWeightLogs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to get logs", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Weight logs retrieved", res))
}

func (ctrl *healthLogController) DeleteWeight(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", err.Error(), nil))
		return
	}

	if err := ctrl.service.DeleteWeight(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete log", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Weight log deleted", nil))
}
