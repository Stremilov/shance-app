package main

// @title Shance API
// @version 1.0
// @description API для управления проектами
// @host localhost:8000
// @BasePath /api/v1
// @schemes http

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/docs"
	"github.com/levstremilov/shance-app/internal/config"
	"github.com/levstremilov/shance-app/internal/handler"
	"github.com/levstremilov/shance-app/internal/middleware"
	"github.com/levstremilov/shance-app/internal/repository"
	"github.com/levstremilov/shance-app/internal/service"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(
	projectHandler *handler.ProjectHandler,
	authHandler *handler.AuthHandler,
	authService service.AuthServiceInterface,
	userHandler *handler.UserHandler,
) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))
	// Swagger документация
	docs.SwaggerInfo.Title = "Shance API"
	docs.SwaggerInfo.Description = "API для управления проектами"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
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

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(authService))
		{
			users.GET("/me", userHandler.GetMe)
			users.GET("/:id", userHandler.GetByID)
			users.PATCH("/update", userHandler.UpdateMe)
		}

		// Защищенные маршруты
		projects := api.Group("/projects")
		projects.Use(middleware.AuthMiddleware(authService))
		{
			projects.POST("", projectHandler.CreateProject)
			projects.GET("", projectHandler.GetProjects)
			projects.GET("/search", projectHandler.SearchProjects)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
		}
	}
	return r
}

func initDependencies(cfg *config.Config) (*handler.ProjectHandler, *handler.AuthHandler, service.AuthServiceInterface, *handler.UserHandler, error) {
	db, err := config.InitDB(cfg)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Очищаем базу данных
	if err := config.CleanDB(db); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to clean db: %w", err)
	}

	// Применяем миграции
	if err := config.RunMigrations(db); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Инициализация репозиториев
	projectRepo := repository.NewProjectRepository(db)
	userRepo := repository.NewUserRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// Инициализация сервисов
	projectService := service.NewProjectService(projectRepo, tagRepo)
	authService := service.NewAuthService(
		userRepo,
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)
	userService := service.NewUserService(userRepo)

	// Инициализация обработчиков
	projectHandler := handler.NewProjectHandler(projectService)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	return projectHandler, authHandler, authService, userHandler, nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	projectHandler, authHandler, authService, userHandler, err := initDependencies(cfg)
	if err != nil {
		log.Fatalf("Failed to init dependencies: %v", err)
	}

	r := setupRouter(projectHandler, authHandler, authService, userHandler)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
