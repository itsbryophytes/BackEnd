package controller

import (
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/modules/profile/dto"
	"github.com/Caknoooo/go-gin-clean-starter/modules/profile/service"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ProfileController interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

type profileController struct {
	service service.ProfileService
}

func NewProfileController(service service.ProfileService) ProfileController {
	return &profileController{service: service}
}

func (ctrl *profileController) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	profile, err := ctrl.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PROFILE, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PROFILE, profile))
}

func (ctrl *profileController) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req dto.HealthProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Bad Request", err.Error(), nil))
		return
	}

	profile, err := ctrl.service.UpsertProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PROFILE, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_PROFILE, profile))
}
