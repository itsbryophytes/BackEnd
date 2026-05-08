package processing

import (
	"github.com/Caknoooo/go-gin-clean-starter/middlewares"
	authService "github.com/Caknoooo/go-gin-clean-starter/modules/auth/service"
	"github.com/Caknoooo/go-gin-clean-starter/modules/processing/controller"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	ctrl := do.MustInvoke[controller.ProcessingController](injector)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	routes := server.Group("/api/processing")
	{
		routes.POST("/ocr/:document_id", middlewares.Authenticate(jwtService), ctrl.RunOCR)
		routes.GET("/:document_id/status", middlewares.Authenticate(jwtService), ctrl.GetStatus)
		routes.GET("/:document_id/result", middlewares.Authenticate(jwtService), ctrl.GetResult)
	}
}
