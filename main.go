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
	"time"

	"github.com/levstremilov/shance-app/internal/config"
	"github.com/levstremilov/shance-app/internal/database"
	"github.com/levstremilov/shance-app/internal/handler"
	"github.com/levstremilov/shance-app/internal/repository"
	"github.com/levstremilov/shance-app/internal/server"
	"github.com/levstremilov/shance-app/internal/service"
	_ "github.com/lib/pq"

	"gorm.io/gorm"
)

func initDependencies(db *gorm.DB) (*handler.AuthHandler, *handler.UserHandler, *handler.ProjectHandler, *handler.TagHandler, *handler.ProjectVacancyHandler, service.AuthServiceInterface) {
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	tagRepo := repository.NewTagRepository(db)
	vacancyRepo := repository.NewProjectVacancyRepository(db)

	authService := service.NewAuthService(userRepo, "your-secret-key", 24*time.Hour, 168*time.Hour)
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo, tagRepo, userRepo)
	tagService := service.NewTagService(tagRepo)
	vacancyService := service.NewProjectVacancyService(vacancyRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService)
	tagHandler := handler.NewTagHandler(tagService)
	vacancyHandler := handler.NewProjectVacancyHandler(vacancyService)

	return authHandler, userHandler, projectHandler, tagHandler, vacancyHandler, authService
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	authHandler, userHandler, projectHandler, tagHandler, vacancyHandler, authService := initDependencies(db)

	r := server.SetUpRouter(projectHandler, authHandler, userHandler, tagHandler, vacancyHandler, authService, cfg)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
