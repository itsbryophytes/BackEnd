package metric

import (
	"github.com/Caknoooo/go-gin-clean-starter/middlewares"
	authService "github.com/Caknoooo/go-gin-clean-starter/modules/auth/service"
	"github.com/Caknoooo/go-gin-clean-starter/modules/metric/controller"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	ctrl := do.MustInvoke[controller.MetricController](injector)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	routes := server.Group("/api/metrics")
	{
		routes.GET("", middlewares.Authenticate(jwtService), ctrl.GetMetrics)
		routes.POST("/confirm", middlewares.Authenticate(jwtService), ctrl.ConfirmMetrics)
		routes.PUT("/:id", middlewares.Authenticate(jwtService), ctrl.UpdateMetric)
		routes.DELETE("/:id", middlewares.Authenticate(jwtService), ctrl.DeleteMetric)
	}
}
