package health_log

import (
	"fmt"
	"github.com/Caknoooo/go-gin-clean-starter/middlewares"
	authService "github.com/Caknoooo/go-gin-clean-starter/modules/auth/service"
	"github.com/Caknoooo/go-gin-clean-starter/modules/health_log/controller"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	fmt.Println("Registering HealthLog routes...")
	ctrl := do.MustInvoke[controller.HealthLogController](injector)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	routes := server.Group("/api/health-logs")
	routes.Use(middlewares.Authenticate(jwtService))
	{
		routes.POST("/blood-sugar", ctrl.CreateBloodSugar)
		routes.GET("/blood-sugar", ctrl.GetBloodSugarLogs)
		routes.DELETE("/blood-sugar/:id", ctrl.DeleteBloodSugar)

		routes.POST("/blood-pressure", ctrl.CreateBloodPressure)
		routes.GET("/blood-pressure", ctrl.GetBloodPressureLogs)
		routes.DELETE("/blood-pressure/:id", ctrl.DeleteBloodPressure)

		routes.POST("/weight", ctrl.CreateWeight)
		routes.GET("/weight", ctrl.GetWeightLogs)
		routes.DELETE("/weight/:id", ctrl.DeleteWeight)
	}
}
