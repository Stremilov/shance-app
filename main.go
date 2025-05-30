package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/config"
	"github.com/levstremilov/shance-app/internal/handler"
	"github.com/levstremilov/shance-app/internal/repository"
	"github.com/levstremilov/shance-app/internal/service"
)

func setupRouter(handler *handler.ProjectHandler) *gin.Engine {
	r := gin.Default()
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
