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
		routes.GET("", ctrl.GetArticles)
		routes.GET("/:id", ctrl.GetArticle)

		routes.POST("", middlewares.Authenticate(jwtService), middlewares.AuthorizeRole("admin"), ctrl.CreateArticle)
		routes.PUT("/:id", middlewares.Authenticate(jwtService), middlewares.AuthorizeRole("admin"), ctrl.UpdateArticle)
		routes.DELETE("/:id", middlewares.Authenticate(jwtService), middlewares.AuthorizeRole("admin"), ctrl.DeleteArticle)
		routes.POST("/:id/publish", middlewares.Authenticate(jwtService), middlewares.AuthorizeRole("admin"), ctrl.PublishArticle)
	}
}
