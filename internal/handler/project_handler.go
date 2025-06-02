package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/service"
)

// ProjectHandler представляет обработчик для работы с проектами
type ProjectHandler struct {
	projectService *service.ProjectService
}

// CreateProjectRequest представляет запрос на создание проекта
type CreateProjectRequest struct {
	Name        string   `json:"name" binding:"required" example:"Новый проект"`
	Title       string   `json:"title" example:"Заголовок проекта"`
	Subtitle    string   `json:"subtitle" example:"Подзаголовок проекта"`
	Description string   `json:"description" example:"Описание проекта"`
	Photo       []string `json:"photo" example:"['photo1.jpg', 'photo2.jpg']"`
	Tags        []string `json:"tags" example:"['tag1', 'tag2']"`
}

// UpdateProjectRequest представляет запрос на обновление проекта
type UpdateProjectRequest struct {
	Name        string   `json:"name" example:"Обновленный проект"`
	Title       string   `json:"title" example:"Новый заголовок"`
	Subtitle    string   `json:"subtitle" example:"Новый подзаголовок"`
	Description string   `json:"description" example:"Новое описание"`
	Photo       []string `json:"photo" example:"['new_photo1.jpg', 'new_photo2.jpg']"`
	Tags        []string `json:"tags" example:"['new_tag1', 'new_tag2']"`
}

// ProjectResponse представляет ответ API для проекта
type ProjectResponse struct {
	ID          uint         `json:"id" example:"1"`
	Name        string       `json:"name" example:"Новый проект"`
	Title       string       `json:"title" example:"Заголовок проекта"`
	Subtitle    string       `json:"subtitle" example:"Подзаголовок проекта"`
	Description string       `json:"description" example:"Описание проекта"`
	Photo       []string     `json:"photo" example:"['photo1.jpg', 'photo2.jpg']"`
	Tags        []string     `json:"tags" example:"['tag1', 'tag2']"`
	UserID      uint         `json:"user_id" example:"1"`
	User        UserResponse `json:"user"`
	CreatedAt   time.Time    `json:"created_at" example:"2024-03-20T12:00:00Z"`
}

type UserResponse struct {
	ID        uint   `json:"id" example:"1"`
	FirstName string `json:"first_name" example:"Иван"`
	LastName  string `json:"last_name" example:"Иванов"`
	Email     string `json:"email" example:"ivan@example.com"`
}

// NewProjectHandler создает новый экземпляр ProjectHandler
func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// GetProjects godoc
// @Summary Получение всех проектов
// @Description Возвращает список всех проектов
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {array} ProjectResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects [get]
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	projects, err := h.projectService.GetAll()
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

		// Преобразуем строку обратно в массив для ответа
		var photoArray []string
		if err := json.Unmarshal([]byte(p.Photo), &photoArray); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "error processing photo data"})
			return
		}

		response[i] = ProjectResponse{
			ID:          p.ID,
			Name:        p.Name,
			Title:       p.Title,
			Subtitle:    p.Subtitle,
			Description: p.Description,
			Photo:       photoArray,
			Tags:        tags,
			UserID:      p.UserID,
			User: UserResponse{
				ID:        p.User.ID,
				FirstName: p.User.FirstName,
				LastName:  p.User.LastName,
				Email:     p.User.Email,
			},
			CreatedAt: p.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProject godoc
// @Summary Получение проекта по ID
// @Description Возвращает проект по указанному ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Success 200 {object} ProjectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid project ID"})
		return
	}

	project, err := h.projectService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Project not found"})
		return
	}

	tags := make([]string, len(project.Tags))
	for i, t := range project.Tags {
		tags[i] = t.Name
	}

	// Преобразуем строку обратно в массив для ответа
	var photoArray []string
	if err := json.Unmarshal([]byte(project.Photo), &photoArray); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "error processing photo data"})
		return
	}

	response := ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Title:       project.Title,
		Subtitle:    project.Subtitle,
		Description: project.Description,
		Photo:       photoArray,
		Tags:        tags,
		UserID:      project.UserID,
		User: UserResponse{
			ID:        project.User.ID,
			FirstName: project.User.FirstName,
			LastName:  project.User.LastName,
			Email:     project.User.Email,
		},
		CreatedAt: project.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// CreateProject godoc
// @Summary Создание нового проекта
// @Description Создает новый проект в системе
// @Tags projects
// @Accept json
// @Produce json
// @Param request body CreateProjectRequest true "Данные проекта"
// @Success 201 {object} ProjectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Преобразуем массив фотографий в JSON строку
	photoJSON, err := json.Marshal(req.Photo)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid photo format"})
		return
	}

	project := &domain.Project{
		Name:        req.Name,
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		Photo:       string(photoJSON),
		UserID:      userID.(uint),
	}

	if err := h.projectService.Create(project); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	tags := make([]string, len(project.Tags))
	for i, t := range project.Tags {
		tags[i] = t.Name
	}

	// Преобразуем строку обратно в массив для ответа
	var photoArray []string
	if err := json.Unmarshal([]byte(project.Photo), &photoArray); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "error processing photo data"})
		return
	}

	response := ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Title:       project.Title,
		Subtitle:    project.Subtitle,
		Description: project.Description,
		Photo:       photoArray,
		Tags:        tags,
		UserID:      project.UserID,
		User: UserResponse{
			ID:        project.User.ID,
			FirstName: project.User.FirstName,
			LastName:  project.User.LastName,
			Email:     project.User.Email,
		},
		CreatedAt: project.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateProject godoc
// @Summary Обновление проекта
// @Description Обновляет данные существующего проекта
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param request body UpdateProjectRequest true "Данные проекта"
// @Success 200 {object} ProjectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid project ID"})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Преобразуем массив фотографий в JSON строку
	photoJSON, err := json.Marshal(req.Photo)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid photo format"})
		return
	}

	project := &domain.Project{
		ID:          uint(id),
		Name:        req.Name,
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		Photo:       string(photoJSON),
	}

	if err := h.projectService.Update(project); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	tags := make([]string, len(project.Tags))
	for i, t := range project.Tags {
		tags[i] = t.Name
	}

	// Преобразуем строку обратно в массив для ответа
	var photoArray []string
	if err := json.Unmarshal([]byte(project.Photo), &photoArray); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "error processing photo data"})
		return
	}

	response := ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Title:       project.Title,
		Subtitle:    project.Subtitle,
		Description: project.Description,
		Photo:       photoArray,
		Tags:        tags,
		UserID:      project.UserID,
		User: UserResponse{
			ID:        project.User.ID,
			FirstName: project.User.FirstName,
			LastName:  project.User.LastName,
			Email:     project.User.Email,
		},
		CreatedAt: project.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteProject godoc
// @Summary Удаление проекта
// @Description Удаляет проект по указанному ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid project ID"})
		return
	}

	if err := h.projectService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchProjects godoc
// @Summary Поиск проектов по названию
// @Description Возвращает список проектов, названия которых содержат поисковый запрос
// @Tags projects
// @Accept json
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Success 200 {array} ProjectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/search [get]
func (h *ProjectHandler) SearchProjects(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Search query is required"})
		return
	}

	projects, err := h.projectService.Search(query)
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

		// Преобразуем строку обратно в массив для ответа
		var photoArray []string
		if err := json.Unmarshal([]byte(p.Photo), &photoArray); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "error processing photo data"})
			return
		}

		response[i] = ProjectResponse{
			ID:          p.ID,
			Name:        p.Name,
			Title:       p.Title,
			Subtitle:    p.Subtitle,
			Description: p.Description,
			Photo:       photoArray,
			Tags:        tags,
			UserID:      p.UserID,
			User: UserResponse{
				ID:        p.User.ID,
				FirstName: p.User.FirstName,
				LastName:  p.User.LastName,
				Email:     p.User.Email,
			},
			CreatedAt: p.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}
