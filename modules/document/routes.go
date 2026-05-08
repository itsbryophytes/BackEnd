package document

import (
	"github.com/Caknoooo/go-gin-clean-starter/middlewares"
	authService "github.com/Caknoooo/go-gin-clean-starter/modules/auth/service"
	"github.com/Caknoooo/go-gin-clean-starter/modules/document/controller"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	ctrl := do.MustInvoke[controller.DocumentController](injector)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	routes := server.Group("/api/documents")
	{
		routes.POST("/upload", middlewares.Authenticate(jwtService), ctrl.UploadDocument)
		routes.GET("/pending", middlewares.Authenticate(jwtService), ctrl.GetPendingDocuments)
		routes.POST("", middlewares.Authenticate(jwtService), ctrl.CreateDocument)
		routes.GET("", middlewares.Authenticate(jwtService), ctrl.GetDocuments)
		routes.GET("/:id", middlewares.Authenticate(jwtService), ctrl.GetDocument)
		routes.POST("/:id/confirm", middlewares.Authenticate(jwtService), ctrl.ConfirmDocument)
		routes.POST("/:id/discard", middlewares.Authenticate(jwtService), ctrl.DiscardDocument)
		routes.DELETE("/:id", middlewares.Authenticate(jwtService), ctrl.DeleteDocument)
	}
}
