package main

// @title Shance API
// @version 1.0
// @description API для управления проектами
// @host localhost:8080
// @BasePath /api/v1
// @schemes http

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/docs"
	"github.com/levstremilov/shance-app/internal/config"
	"github.com/levstremilov/shance-app/internal/handler"
	"github.com/levstremilov/shance-app/internal/repository"
	"github.com/levstremilov/shance-app/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(
	projectHandler *handler.ProjectHandler,
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
) *gin.Engine {
	r := gin.Default()

	// Swagger документация
	docs.SwaggerInfo.Title = "Shance API"
	docs.SwaggerInfo.Description = "API для управления проектами"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		// Публичные маршруты
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Защищенные маршруты
		projects := api.Group("/projects")
		projects.Use(handler.AuthMiddleware(authService))
		{
			projects.POST("", projectHandler.Create)
			projects.GET("", projectHandler.GetAll)
			projects.GET("/:id", projectHandler.GetByID)
			projects.PUT("/:id", projectHandler.Update)
			projects.DELETE("/:id", projectHandler.Delete)
		}
	}
	return r
}

func initDependencies(cfg *config.Config) (*handler.ProjectHandler, *handler.AuthHandler, *service.AuthService, error) {
	db, err := config.InitDB(cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Инициализация репозиториев
	projectRepo := repository.NewProjectRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервисов
	projectService := service.NewProjectService(projectRepo)
	authService := service.NewAuthService(
		userRepo,
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// Инициализация обработчиков
	projectHandler := handler.NewProjectHandler(projectService)
	authHandler := handler.NewAuthHandler(authService)

	return projectHandler, authHandler, authService, nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	projectHandler, authHandler, authService, err := initDependencies(cfg)
	if err != nil {
		log.Fatalf("Failed to init dependencies: %v", err)
	}

	r := setupRouter(projectHandler, authHandler, authService)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
