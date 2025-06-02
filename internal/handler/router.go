package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/service"
)

func SetupRouter(authService *service.AuthService, projectService *service.ProjectService) *gin.Engine {
	r := gin.Default()

	authHandler := NewAuthHandler(authService)

	// Публичные эндпоинты
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/refresh", authHandler.RefreshToken)

	// Защищенные эндпоинты
	protected := r.Group("/api")
	protected.Use(AuthMiddleware(authService))
	{
		projectHandler := NewProjectHandler(projectService)
		protected.GET("/projects/search", projectHandler.SearchProjects)
	}

	return r
}
