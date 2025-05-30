package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/service"
)

// ProjectHandler представляет обработчик для работы с проектами
type ProjectHandler struct {
	service *service.ProjectService
}

// NewProjectHandler создает новый экземпляр ProjectHandler
func NewProjectHandler(service *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// Create godoc
// @Summary Создание нового проекта
// @Description Создает новый проект в системе
// @Tags projects
// @Accept json
// @Produce json
// @Param project body domain.Project true "Данные проекта"
// @Success 201 {object} domain.Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var project domain.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetAll godoc
// @Summary Получение всех проектов
// @Description Возвращает список всех проектов
// @Tags projects
// @Produce json
// @Success 200 {array} domain.Project
// @Failure 500 {object} map[string]string
// @Router /projects [get]
func (h *ProjectHandler) GetAll(c *gin.Context) {
	projects, err := h.service.GetAllProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetByID godoc
// @Summary Получение проекта по ID
// @Description Возвращает проект по указанному ID
// @Tags projects
// @Produce json
// @Param id path int true "ID проекта"
// @Success 200 {object} domain.Project
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	project, err := h.service.GetProjectByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// Update godoc
// @Summary Обновление проекта
// @Description Обновляет данные существующего проекта
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param project body domain.Project true "Данные проекта"
// @Success 200 {object} domain.Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var project domain.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ID = uint(id)
	if err := h.service.UpdateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// Delete godoc
// @Summary Удаление проекта
// @Description Удаляет проект по указанному ID
// @Tags projects
// @Param id path int true "ID проекта"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteProject(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
