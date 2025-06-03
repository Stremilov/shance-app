package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/service"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	FirstName string   `json:"first_name" binding:"required" example:"John"`
	LastName  string   `json:"last_name" binding:"required" example:"Doe"`
	Phone     string   `json:"phone" example:"+79001234567"`
	Role      string   `json:"role" example:"user"`
	Tags      []string `json:"tags" example:"['tag1', 'tag2']"`
	Country   string   `json:"country" example:"Russia"`
	City      string   `json:"city" example:"Moscow"`
	Email     string   `json:"email" binding:"required,email" example:"user@example.com"`
	Password  string   `json:"password" binding:"required,min=6" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и возвращает refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.authService.Register(service.RegisterData{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      req.Role,
		Tags:      req.Tags,
		Country:   req.Country,
		City:      req.City,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", tokens.AccessToken, int((time.Hour * 24).Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
}

// Login godoc
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", tokens.AccessToken, int((time.Hour * 24).Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
}

// Refresh godoc
// @Summary Обновление токенов
// @Description Обновляет пару токенов используя refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token" SchemaExample({"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."})
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", tokens.AccessToken, int((time.Hour * 24).Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
}
