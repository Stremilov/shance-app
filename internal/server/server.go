package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/docs"
	"github.com/levstremilov/shance-app/internal/config"
	"github.com/levstremilov/shance-app/internal/handler"
	"github.com/levstremilov/shance-app/internal/middleware"
	"github.com/levstremilov/shance-app/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func SetUpRouter(
	projectHandler *handler.ProjectHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	tagHandler *handler.TagHandler,
	vacancyHandler *handler.ProjectVacancyHandler,
	authService service.AuthServiceInterface,
	cfg *config.Config,
) *gin.Engine {
	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.MetricsPath = "/metrics"
	p.Use(r)

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

		protected := api.Group("")

		// Protected routes
		// if cfg.Server.Env == "dev" {
		protected.Use(middleware.AuthMiddleware(authService))
		// }

		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetMe)
				users.PATCH("/me", userHandler.UpdateMe)
				users.GET("/:id", userHandler.GetUser)
				users.GET("/:id/projects", userHandler.GetOwnProjects)
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
				projects.POST("/:id/vacancy", vacancyHandler.CreateProjectVacancy)
				projects.GET("/:id/vacancies", vacancyHandler.GetProjectVacancies)
			}

			vacancies := protected.Group("/vacancies")
			{
				vacancies.POST("/:id", vacancyHandler.CreateVacancyResponse)
				vacancies.GET("/:id/responses", vacancyHandler.GetVacancyResponses)
			}

			// Tag routes
			tags := protected.Group("/tags")
			{
				tags.POST("", tagHandler.CreateTag)
				tags.GET("", tagHandler.ListTags)
				tags.GET("/search", tagHandler.SearchTags)
				tags.PUT("/:id", tagHandler.UpdateTag)
				tags.DELETE("/:id", tagHandler.DeleteTag)
			}

			// Technology route
			protected.POST("/technologies", vacancyHandler.CreateTechnology)
		}
	}

	return r
}
