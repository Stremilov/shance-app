package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/service"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type UpdateUserRequest struct {
	FirstName string   `json:"name" example:"Новое имя"`
	LastName  string   `json:"title" example:"Новая фамилия"`
	Phone     string   `json:"subtitle" example:"Новый номер телефона"`
	Role      string   `json:"description" example:"Новая роль"`
	Tags      []string `json:"photo" example:"new_photo1.jpg, new_photo2.jpg"`
	Country   string   `json:"tags" example:"new_tag1, new_tag2"`
	City      string   `json:"city" example:"Санкт-Петербург"`
}

// GetMe godoc
// @Summary Получение данных текущего пользователя
// @Description Возвращает данные авторизованного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	user, err := h.userService.GetMe(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetByID godoc
// @Summary Получение пользователя по ID
// @Description Возвращает данные пользователя по указанному ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, user)
}

// UpdateMe godoc
// @Summary Обновление данных текущего пользователя
// @Description Обновляет данные авторизованного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param request body UpdateUserRequest true "Данные пользователя"
// @Success 200 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tags := make([]domain.Tag, len(req.Tags))
	for i, tagName := range req.Tags {
		tags[i] = domain.Tag{Name: tagName}
	}

	user := &domain.User{
		ID:        userID.(uint),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      req.Role,
		Tags:      tags,
		Country:   req.Country,
		City:      req.City,
	}

	if err := h.userService.UpdateByID(userID.(uint), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
