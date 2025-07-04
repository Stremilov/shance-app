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
	FirstName *string   `json:"name" example:"Новое имя"`
	LastName  *string   `json:"title" example:"Новая фамилия"`
	Phone     *string   `json:"subtitle" example:"Новый номер телефона"`
	Role      *string   `json:"description" example:"Новая роль"`
	Tags      *[]string `json:"photo" example:"new_photo1.jpg, new_photo2.jpg"`
	Country   *string   `json:"tags" example:"РА СИ Я"`
	City      *string   `json:"city" example:"Санкт-Петербург"`
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
// @Router /users/me [patch]
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

	if req.FirstName != nil {
		currentUser.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		currentUser.LastName = *req.LastName
	}
	if req.Phone != nil {
		currentUser.Phone = *req.Phone
	}
	if req.Role != nil {
		currentUser.Role = *req.Role
	}
	if req.Country != nil {
		currentUser.Country = *req.Country
	}
	if req.City != nil {
		currentUser.City = *req.City
	}
	if req.Tags != nil {
		tags := make([]models.Tag, len(*req.Tags))
		for i, tagName := range *req.Tags {
			tags[i] = models.Tag{Name: tagName}
		}
		currentUser.Tags = tags
	}

	if err := h.userService.Update(currentUser); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, currentUser)
}

// GetOwnProjects godoc
// @Summary Получение проектов пользователя
// @Description Возвращает список проектов текущего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} ProjectResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{user_id}/projects [get]
func (h *UserHandler) GetOwnProjects(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid user id type"})
		return
	}

	projects, err := h.userService.GetOwnProjects(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		tags := make([]string, len(p.Tags))
		for j, t := range p.Tags {
			tags[j] = t.Name
		}

		response[i] = ProjectResponse{
			ID:          p.ID,
			Name:        p.Name,
			Title:       p.Title,
			Subtitle:    p.Subtitle,
			Description: p.Description,
			Photo:       []string{p.Photo},
			Tags:        tags,
			UserID:      p.UserID,
		}
	}

	c.JSON(http.StatusOK, response)
}
