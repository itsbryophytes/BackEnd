package file

import (
	"github.com/Caknoooo/go-gin-clean-starter/middlewares"
	authService "github.com/Caknoooo/go-gin-clean-starter/modules/auth/service"
	"github.com/Caknoooo/go-gin-clean-starter/modules/file/controller"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	fileController := do.MustInvoke[controller.FileController](injector)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	routes := server.Group("/api/files")
	{
		routes.POST("/upload", middlewares.Authenticate(jwtService), fileController.UploadFile)
		routes.GET("", middlewares.Authenticate(jwtService), fileController.GetFiles)
		routes.GET("/:id", middlewares.Authenticate(jwtService), fileController.GetFile)
		routes.DELETE("/:id", middlewares.Authenticate(jwtService), fileController.DeleteFile)
	}
}
