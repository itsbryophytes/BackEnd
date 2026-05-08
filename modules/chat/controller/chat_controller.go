package controller

import (
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/modules/chat/dto"
	"github.com/Caknoooo/go-gin-clean-starter/modules/chat/service"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ChatController interface {
	SendMessage(ctx *gin.Context)
	GetSessions(ctx *gin.Context)
	GetHistory(ctx *gin.Context)
	DeleteSession(ctx *gin.Context)
	ClearHistory(ctx *gin.Context)
}

type chatController struct {
	chatService service.ChatService
}

func NewChatController(chatService service.ChatService) ChatController {
	return &chatController{chatService: chatService}
}

func (c *chatController) SendMessage(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	var req dto.SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Bad Request", err.Error(), nil))
		return
	}
	reply, err := c.chatService.SendMessage(ctx.Request.Context(), userID, req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.BuildResponseFailed("RAG request failed", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess(dto.MESSAGE_SEND_CHAT_SUCCESS, reply))
}

func (c *chatController) GetSessions(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	sessions, err := c.chatService.GetSessions(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to get sessions", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess(dto.MESSAGE_GET_SESSIONS_SUCCESS, sessions))
}

func (c *chatController) GetHistory(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	sessionID := ctx.Param("session_id")
	history, err := c.chatService.GetHistory(ctx.Request.Context(), userID, sessionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.BuildResponseFailed("Session not found", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess(dto.MESSAGE_GET_HISTORY_SUCCESS, history))
}

func (c *chatController) DeleteSession(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	sessionID := ctx.Param("session_id")
	if err := c.chatService.DeleteSession(ctx.Request.Context(), userID, sessionID); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to delete session", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess(dto.MESSAGE_DELETE_SESSION_SUCCESS, nil))
}

func (c *chatController) ClearHistory(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if err := c.chatService.ClearAllHistory(ctx.Request.Context(), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Failed to clear history", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess("Chat history cleared", nil))
}
