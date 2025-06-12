package handler

import (
	"net/http"
	"strconv"

	_ "gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/models"
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
	Country   string   `json:"tags" example:"РА СИ Я"`
	City      string   `json:"city" example:"Санкт-Петербург"`
}

// GetMe godoc
// @Summary Получение информации о текущем пользователе
// @Description Возвращает информацию о текущем авторизованном пользователе
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.SwaggerUser
// @Failure 401 {object} ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}
	user, err := h.userService.GetMe(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUser godoc
// @Summary Получение информации о пользователе
// @Description Возвращает информацию о пользователе по его ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.SwaggerUser
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}


// UpdateMe godoc
// @Summary Обновление данных текущего пользователя
// @Description Обновляет данные авторизованного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateUserRequest true "Данные пользователя"
// @Success 200 {object} models.SwaggerUser
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	currentUser, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	currentUser.FirstName = req.FirstName
	currentUser.LastName = req.LastName
	currentUser.Phone = req.Phone
	currentUser.Role = req.Role
	currentUser.Country = req.Country
	currentUser.City = req.City

	if len(req.Tags) > 0 {
		tags := make([]models.Tag, len(req.Tags))
		for i, tagName := range req.Tags {
			tags[i] = models.Tag{Name: tagName}
		}
		currentUser.Tags = tags
	} else {
		currentUser.Tags = []models.Tag{}
	}

	if err := h.userService.Update(currentUser); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, currentUser)
}
