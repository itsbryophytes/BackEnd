package chat

import (
	"github.com/Caknoooo/go-gin-clean-starter/middlewares"
	authService "github.com/Caknoooo/go-gin-clean-starter/modules/auth/service"
	"github.com/Caknoooo/go-gin-clean-starter/modules/chat/controller"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	ctrl := do.MustInvoke[controller.ChatController](injector)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	routes := server.Group("/api/chat")
	{
		routes.POST("", middlewares.CORSMiddleware(), middlewares.Authenticate(jwtService), ctrl.SendMessage)
		routes.GET("/sessions", middlewares.CORSMiddleware(), middlewares.Authenticate(jwtService), ctrl.GetSessions)
		routes.GET("/sessions/:session_id", middlewares.CORSMiddleware(), middlewares.Authenticate(jwtService), ctrl.GetHistory)
		routes.DELETE("/sessions/:session_id", middlewares.CORSMiddleware(), middlewares.Authenticate(jwtService), ctrl.DeleteSession)
		routes.DELETE("/history", middlewares.CORSMiddleware(), middlewares.Authenticate(jwtService), ctrl.ClearHistory)
	}
}
