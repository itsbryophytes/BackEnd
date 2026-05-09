package providers

import (
	"github.com/Caknoooo/go-gin-clean-starter/config"
	articleController "github.com/Caknoooo/go-gin-clean-starter/modules/article/controller"
	authController "github.com/Caknoooo/go-gin-clean-starter/modules/auth/controller"
	authRepo "github.com/Caknoooo/go-gin-clean-starter/modules/auth/repository"
	authService "github.com/Caknoooo/go-gin-clean-starter/modules/auth/service"
	chatController "github.com/Caknoooo/go-gin-clean-starter/modules/chat/controller"
	chatRepo "github.com/Caknoooo/go-gin-clean-starter/modules/chat/repository"
	chatService "github.com/Caknoooo/go-gin-clean-starter/modules/chat/service"
	documentController "github.com/Caknoooo/go-gin-clean-starter/modules/document/controller"
	fileController "github.com/Caknoooo/go-gin-clean-starter/modules/file/controller"
	fileRepo "github.com/Caknoooo/go-gin-clean-starter/modules/file/repository"
	fileService "github.com/Caknoooo/go-gin-clean-starter/modules/file/service"
	metricController "github.com/Caknoooo/go-gin-clean-starter/modules/metric/controller"
	processingController "github.com/Caknoooo/go-gin-clean-starter/modules/processing/controller"
	profileController "github.com/Caknoooo/go-gin-clean-starter/modules/profile/controller"
	profileRepo "github.com/Caknoooo/go-gin-clean-starter/modules/profile/repository"
	profileService "github.com/Caknoooo/go-gin-clean-starter/modules/profile/service"
	userController "github.com/Caknoooo/go-gin-clean-starter/modules/user/controller"
	"github.com/Caknoooo/go-gin-clean-starter/modules/user/repository"
	userService "github.com/Caknoooo/go-gin-clean-starter/modules/user/service"
	healthLogController "github.com/Caknoooo/go-gin-clean-starter/modules/health_log/controller"
	healthLogRepo "github.com/Caknoooo/go-gin-clean-starter/modules/health_log/repository"
	healthLogService "github.com/Caknoooo/go-gin-clean-starter/modules/health_log/service"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/constants"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/rag"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func InitDatabase(injector *do.Injector) {
	do.ProvideNamed(injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
		return config.SetUpDatabaseConnection(), nil
	})
}

func RegisterDependencies(injector *do.Injector) {
	InitDatabase(injector)

	do.ProvideNamed(injector, constants.JWTService, func(i *do.Injector) (authService.JWTService, error) {
		return authService.NewJWTService(), nil
	})

	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)
	ragClient := rag.NewClient()

	// Repositories
	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := authRepo.NewRefreshTokenRepository(db)
	fileRepository := fileRepo.NewFileRepository(db)
	chatRepository := chatRepo.NewChatRepository()
	profileRepository := profileRepo.NewProfileRepository()
	healthLogRepository := healthLogRepo.NewHealthLogRepository(db)

	// Services
	userSvc := userService.NewUserService(userRepository, db)
	authSvc := authService.NewAuthService(userRepository, refreshTokenRepository, jwtService, db)
	fileSvc := fileService.NewFileService(fileRepository)
	chatSvc := chatService.NewChatService(ragClient, chatRepository, db)
	profileSvc := profileService.NewProfileService(profileRepository, db)
	healthLogSvc := healthLogService.NewHealthLogService(healthLogRepository)

	// Controllers
	do.Provide(injector, func(i *do.Injector) (userController.UserController, error) {
		return userController.NewUserController(i, userSvc), nil
	})
	do.Provide(injector, func(i *do.Injector) (authController.AuthController, error) {
		return authController.NewAuthController(i, authSvc), nil
	})
	do.Provide(injector, func(i *do.Injector) (fileController.FileController, error) {
		return fileController.NewFileController(fileSvc), nil
	})
	do.Provide(injector, func(i *do.Injector) (documentController.DocumentController, error) {
		return documentController.NewDocumentController(db, ragClient), nil
	})
	do.Provide(injector, func(i *do.Injector) (processingController.ProcessingController, error) {
		return processingController.NewProcessingController(db, ragClient), nil
	})
	do.Provide(injector, func(i *do.Injector) (metricController.MetricController, error) {
		return metricController.NewMetricController(db, ragClient), nil
	})
	do.Provide(injector, func(i *do.Injector) (chatController.ChatController, error) {
		return chatController.NewChatController(chatSvc), nil
	})
	do.Provide(injector, func(i *do.Injector) (profileController.ProfileController, error) {
		return profileController.NewProfileController(profileSvc), nil
	})
	do.Provide(injector, func(i *do.Injector) (articleController.ArticleController, error) {
		return articleController.NewArticleController(db), nil
	})
	do.Provide(injector, func(i *do.Injector) (healthLogController.HealthLogController, error) {
		return healthLogController.NewHealthLogController(healthLogSvc), nil
	})

	_ = jwtService
}
