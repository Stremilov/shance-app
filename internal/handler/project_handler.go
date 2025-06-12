package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/models"
	"github.com/levstremilov/shance-app/internal/service"
	"gorm.io/gorm"
)

// ProjectHandler представляет обработчик для работы с проектами
type ProjectHandler struct {
	projectService service.ProjectServiceInterface
}

// CreateProjectRequest представляет запрос на создание проекта
type CreateProjectRequest struct {
	Name        string                  `form:"name" binding:"required" example:"Новый проект"`
	Title       string                  `form:"title" example:"Заголовок проекта"`
	Subtitle    string                  `form:"subtitle" example:"Подзаголовок проекта"`
	Description string                  `form:"description" example:"Описание проекта"`
	Tags        []string                `form:"tags" example:"tag1,tag2"`
	Photos      []*multipart.FileHeader `form:"photos" binding:"required"`
}

// UpdateProjectRequest представляет запрос на обновление проекта
type UpdateProjectRequest struct {
	Name        string   `json:"name" example:"Обновленный проект"`
	Title       string   `json:"title" example:"Новый заголовок"`
	Subtitle    string   `json:"subtitle" example:"Новый подзаголовок"`
	Description string   `json:"description" example:"Новое описание"`
	Photo       []string `json:"photo" example:"['new_photo1.jpg', 'new_photo2.jpg']"`
	Tags        []string `json:"tags" example:"new_tag1,new_tag2"`
}

// ProjectResponse представляет ответ API для проекта
type ProjectResponse struct {
	ID          uint         `json:"id" example:"1"`
	Name        string       `json:"name" example:"Новый проект"`
	Title       string       `json:"title" example:"Заголовок проекта"`
	Subtitle    string       `json:"subtitle" example:"Подзаголовок проекта"`
	Description string       `json:"description" example:"Описание проекта"`
	Photo       []string     `json:"photo" example:"['photo1.jpg', 'photo2.jpg']"`
	Tags        []string     `json:"tags" example:"tag1,tag2"`
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

// InviteMemberRequest представляет запрос на приглашение участника
type InviteMemberRequest struct {
	Email string `json:"email" binding:"required" example:"user@example.com"`
	Role  string `json:"role" binding:"required" example:"member"`
}

// ProjectMemberResponse представляет ответ с информацией об участнике
type ProjectMemberResponse struct {
	ID        uint   `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	FirstName string `json:"first_name" example:"Иван"`
	LastName  string `json:"last_name" example:"Иванов"`
	Role      string `json:"role" example:"member"`
}

// NewProjectHandler создает новый экземпляр ProjectHandler
func NewProjectHandler(projectService service.ProjectServiceInterface) *ProjectHandler {
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
// @Summary Получение информации о проекте
// @Description Возвращает информацию о проекте по его ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Success 200 {object} models.SwaggerProject
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
// @Description Создает новый проект в системе с фотографиями
// @Tags projects
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Название проекта"
// @Param title formData string false "Заголовок проекта"
// @Param subtitle formData string false "Подзаголовок проекта"
// @Param description formData string false "Описание проекта"
// @Param tags formData []string false "Теги проекта"
// @Param photos formData file true "Фотографии проекта"
// @Success 201 {object} ProjectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Обработка загруженных фотографий
	photoPaths := make([]string, 0, len(req.Photos))
	for _, photo := range req.Photos {
		// Создаем уникальное имя файла
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), photo.Filename)
		filepath := fmt.Sprintf("uploads/projects/%s", filename)

		// Сохраняем файл
		if err := c.SaveUploadedFile(photo, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to save photo"})
			return
		}

		photoPaths = append(photoPaths, filepath)
	}

	photoJSON, err := json.Marshal(photoPaths)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid photo format"})
		return
	}

	project := &models.Project{
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

	response := ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Title:       project.Title,
		Subtitle:    project.Subtitle,
		Description: project.Description,
		Photo:       photoPaths,
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

	project := &models.Project{
		Model:       gorm.Model{ID: uint(id)},
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
		CreatedAt:   project.CreatedAt,
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

// ListProjects godoc
// @Summary Получение списка проектов
// @Description Возвращает список всех проектов
// @Tags projects
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы"
// @Param page_size query int false "Размер страницы"
// @Success 200 {object} models.SwaggerListResponse{results=[]models.SwaggerProject}
// @Failure 400 {object} ErrorResponse
// @Router /projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	// ... existing code ...
}

// InviteMember godoc
// @Summary Приглашение участника в проект
// @Description Приглашает пользователя в проект по email
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param request body InviteMemberRequest true "Данные приглашения"
// @Success 200 {object} ProjectMemberResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /projects/{id}/invite [post]
func (h *ProjectHandler) InviteMember(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid project ID"})
		return
	}

	var req InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Проверяем права доступа
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Проверяем, является ли пользователь владельцем проекта
	isOwner, err := h.projectService.IsProjectOwner(uint(projectID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "only project owner can invite members"})
		return
	}

	// Приглашаем пользователя
	member, err := h.projectService.InviteMember(uint(projectID), req.Email, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := ProjectMemberResponse{
		ID:        member.User.ID,
		Email:     member.User.Email,
		FirstName: member.User.FirstName,
		LastName:  member.User.LastName,
		Role:      member.Role,
	}

	c.JSON(http.StatusOK, response)
}

// GetProjectMembers godoc
// @Summary Получение списка участников проекта
// @Description Возвращает список всех участников проекта
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Success 200 {array} ProjectMemberResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /projects/{id}/members [get]
func (h *ProjectHandler) GetProjectMembers(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid project ID"})
		return
	}

	// Проверяем аутентификацию
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Получаем список участников
	members, err := h.projectService.GetProjectMembers(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := make([]ProjectMemberResponse, len(members))
	for i, member := range members {
		response[i] = ProjectMemberResponse{
			ID:        member.User.ID,
			Email:     member.User.Email,
			FirstName: member.User.FirstName,
			LastName:  member.User.LastName,
			Role:      member.Role,
		}
	}

	c.JSON(http.StatusOK, response)
}
