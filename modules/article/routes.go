package article

import (
	"github.com/Caknoooo/go-gin-clean-starter/middlewares"
	"github.com/Caknoooo/go-gin-clean-starter/modules/article/controller"
	authService "github.com/Caknoooo/go-gin-clean-starter/modules/auth/service"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	ctrl := do.MustInvoke[controller.ArticleController](injector)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	routes := server.Group("/api/articles")
	{
		// public
		routes.GET("", ctrl.GetArticles)
		routes.GET("/:id", ctrl.GetArticle)

		// protected (admin nanti bisa ditambah middleware role)
		routes.POST("", middlewares.Authenticate(jwtService), ctrl.CreateArticle)
		routes.PUT("/:id", middlewares.Authenticate(jwtService), ctrl.UpdateArticle)
		routes.DELETE("/:id", middlewares.Authenticate(jwtService), ctrl.DeleteArticle)
		routes.POST("/:id/publish", middlewares.Authenticate(jwtService), ctrl.PublishArticle)
	}
}
