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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/docs"
	"github.com/levstremilov/shance-app/internal/config"
	"github.com/levstremilov/shance-app/internal/database"
	"github.com/levstremilov/shance-app/internal/handler"
	"github.com/levstremilov/shance-app/internal/middleware"
	"github.com/levstremilov/shance-app/internal/repository"
	"github.com/levstremilov/shance-app/internal/service"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func setupRouter(
	projectHandler *handler.ProjectHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	tagHandler *handler.TagHandler,
	authService service.AuthServiceInterface,
	cfg *config.Config,
) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger
	docs.SwaggerInfo.Title = "Shance API"
	docs.SwaggerInfo.Description = "API для управления проектами"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetMe)
				users.PUT("/me", userHandler.UpdateMe)
				users.GET("/:id", userHandler.GetUser)
			}

			// Project routes
			projects := protected.Group("/projects")
			{
				projects.POST("", projectHandler.CreateProject)
				projects.GET("", projectHandler.GetProjects)
				projects.GET("/:id", projectHandler.GetProject)
				projects.PUT("/:id", projectHandler.UpdateProject)
				projects.DELETE("/:id", projectHandler.DeleteProject)
				projects.GET("/search", projectHandler.SearchProjects)
				projects.POST("/:id/invite", projectHandler.InviteMember)
				projects.GET("/:id/members", projectHandler.GetProjectMembers)
			}

			// Tag routes
			tags := protected.Group("/tags")
			{
				tags.POST("", tagHandler.CreateTag)
				tags.GET("", tagHandler.ListTags)
				tags.GET("/search", tagHandler.SearchTags)
				tags.GET("/:id", tagHandler.GetTag)
				tags.PUT("/:id", tagHandler.UpdateTag)
				tags.DELETE("/:id", tagHandler.DeleteTag)
			}
		}
	}

	return r
}

func initDependencies(db *gorm.DB) (*handler.AuthHandler, *handler.UserHandler, *handler.ProjectHandler, *handler.TagHandler, service.AuthServiceInterface) {
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	tagRepo := repository.NewTagRepository(db)

	authService := service.NewAuthService(userRepo, "your-secret-key", 24*time.Hour, 168*time.Hour)
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo, tagRepo)
	tagService := service.NewTagService(tagRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService)
	tagHandler := handler.NewTagHandler(tagService)

	return authHandler, userHandler, projectHandler, tagHandler, authService
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

	authHandler, userHandler, projectHandler, tagHandler, authService := initDependencies(db)

	r := setupRouter(projectHandler, authHandler, userHandler, tagHandler, authService, cfg)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
