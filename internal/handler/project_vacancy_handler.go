package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/models"
	"github.com/levstremilov/shance-app/internal/service"
)

// CreateProjectVacancyRequest описывает тело запроса на создание вакансии
type CreateProjectVacancyRequest struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Technologies []uint `json:"technologies" binding:"required"` // id технологий
}

type VacancyResponse struct {
	ID              uint      `json:"id"`
	ProjectID       uint      `json:"project_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	TechnologyNames []string  `json:"technology_names"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateTechnologyRequest struct {
	Name string `json:"name" binding:"required"`
}

type ProjectVacancyHandler struct {
	service *service.ProjectVacancyService
}

func NewProjectVacancyHandler(service *service.ProjectVacancyService) *ProjectVacancyHandler {
	return &ProjectVacancyHandler{service: service}
}

// CreateProjectVacancy godoc
// @Summary Создать вакансию для проекта
// @Description Создаёт новую вакансию, привязанную к проекту
// @Tags vacancies
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param request body CreateProjectVacancyRequest true "Данные вакансии"
// @Success 201 {object} models.SwaggerProjectVacancy
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{id}/vacancy [post]
func (h *ProjectVacancyHandler) CreateProjectVacancy(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req CreateProjectVacancyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var technologies []models.Technology
	if len(req.Technologies) > 0 {
		if err := h.service.DB().Find(&technologies, req.Technologies).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	vacancy := models.ProjectVacancy{
		ProjectID:    uint(projectID),
		Title:        req.Title,
		Description:  req.Description,
		Technologies: technologies,
	}

	if err := h.service.Create(&vacancy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vacancy)
}

// GetProjectVacancies godoc
// @Summary Получить вакансии проекта
// @Description Получает список вакансий, привязанных к проекту
// @Tags vacancies
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Success 200 {array} models.SwaggerProjectVacancy
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{id}/vacancies [get]
func (h *ProjectVacancyHandler) GetProjectVacancies(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	vacancies, err := h.service.GetAll(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []VacancyResponse
	for _, v := range vacancies {
		var techNames []string
		for _, t := range v.Technologies {
			techNames = append(techNames, t.Name)
		}
		resp = append(resp, VacancyResponse{
			ID:              v.ID,
			ProjectID:       v.ProjectID,
			Title:           v.Title,
			Description:     v.Description,
			TechnologyNames: techNames,
			CreatedAt:       v.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// CreateTechnology godoc
// @Summary Создать технологию
// @Description Создаёт новую технологию
// @Tags technologies
// @Accept json
// @Produce json
// @Param request body CreateTechnologyRequest true "Данные технологии"
// @Success 201 {object} models.Technology
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /technologies [post]
func (h *ProjectVacancyHandler) CreateTechnology(c *gin.Context) {
	var req CreateTechnologyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tech := models.Technology{Name: req.Name}
	if err := h.service.DB().Create(&tech).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tech)
}
