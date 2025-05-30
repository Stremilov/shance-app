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

func setupRouter(handler *handler.ProjectHandler) *gin.Engine {
	r := gin.Default()

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
		projects := api.Group("/projects")
		{
			projects.POST("", handler.Create)
			projects.GET("", handler.GetAll)
			projects.GET("/:id", handler.GetByID)
			projects.PUT("/:id", handler.Update)
			projects.DELETE("/:id", handler.Delete)
		}
	}
	return r
}

func initDependencies(cfg *config.Config) (*handler.ProjectHandler, error) {
	db, err := config.InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	repo := repository.NewProjectRepository(db)
	svc := service.NewProjectService(repo)
	handler := handler.NewProjectHandler(svc)

	return handler, nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	handler, err := initDependencies(cfg)
	if err != nil {
		log.Fatal("Failed to init dependencies:", err)
	}

	r := setupRouter(handler)
	if err := r.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
