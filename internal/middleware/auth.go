package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/models"
	"github.com/levstremilov/shance-app/internal/service"
)

func AuthMiddleware(authService service.AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			fmt.Printf("Cookie error: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		fmt.Printf("Got access token: %s\n", accessToken)

		// Убираем префикс "Bearer " если он есть
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")

		claims, err := authService.ValidateToken(accessToken)
		if err != nil {
			fmt.Printf("Token validation error: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Устанавливаем данные пользователя в контекст
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user", &models.User{
			ID:    claims.UserID,
			Email: claims.Email,
			Role:  claims.Role,
		})

		c.Next()
	}
}
