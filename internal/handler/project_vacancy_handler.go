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
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description" binding:"required"`
	Technologies []uint   `json:"technologies" binding:"required"` // id технологий
	Questions    []string `json:"questions" binding:"required"`
}

type VacancyResponse struct {
	ID              uint      `json:"id"`
	ProjectID       uint      `json:"project_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Questions       []string  `json:"questions"`
	TechnologyNames []string  `json:"technology_names"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateTechnologyRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateVacancyResponseRequest struct {
	Message string `json:"message"`
}

type ProjectVacancyHandler struct {
	service *service.ProjectVacancyService
}

func NewProjectVacancyHandler(service *service.ProjectVacancyService) *ProjectVacancyHandler {
	return &ProjectVacancyHandler{service: service}
}

// CreateProjectVacancy godoc
// @Summary Создать вакансию для проекта
// @Description Создаёт новую вакансию, привязанную к проекту. Вопросы передаются как массив строк в поле questions.
// @Tags vacancies
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param request body CreateProjectVacancyRequest true "Данные вакансии (title, description, technologies - id технологий, questions - массив строк)"
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

	var questions []models.Question
	for _, qText := range req.Questions {
		var q models.Question
		if err := h.service.DB().Where("description = ?", qText).FirstOrCreate(&q, models.Question{Description: qText}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		questions = append(questions, q)
	}

	vacancy := models.ProjectVacancy{
		ProjectID:    uint(projectID),
		Title:        req.Title,
		Description:  req.Description,
		Technologies: technologies,
		Questions:    questions,
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
		var techQuestions []string
		for _, t := range v.Technologies {
			techNames = append(techNames, t.Name)
		}
		for _, q := range v.Questions {
			techQuestions = append(techQuestions, q.Description)
		}

		resp = append(resp, VacancyResponse{
			ID:              v.ID,
			ProjectID:       v.ProjectID,
			Title:           v.Title,
			Description:     v.Description,
			Questions:       techQuestions,
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

// CreateVacancyResponse godoc
// @Summary Отклик на вакансию
// @Description Создать отклик на вакансию
// @Tags vacancy-responses
// @Accept json
// @Produce json
// @Param id path int true "ID вакансии"
// @Param request body CreateVacancyResponseRequest true "Данные отклика"
// @Success 201 {object} models.SwaggerVacancyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /vacancies/{id} [post]
func (h *ProjectVacancyHandler) CreateVacancyResponse(c *gin.Context) {
	vacancyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vacancy id"})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req CreateVacancyResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := models.VacancyResponse{
		ProjectVacancyID: uint(vacancyID),
		UserID:           userID.(uint),
		Message:          req.Message,
	}

	if err := h.service.DB().Create(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetVacancyResponses godoc
// @Summary Получить отклики на вакансию
// @Description Получить список откликов на вакансию с контактами пользователей
// @Tags vacancy-responses
// @Accept json
// @Produce json
// @Param id path int true "ID вакансии"
// @Success 200 {array} models.SwaggerVacancyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /vacancies/{id}/responses [get]
func (h *ProjectVacancyHandler) GetVacancyResponses(c *gin.Context) {
	vacancyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vacancy id"})
		return
	}

	var responses []models.VacancyResponse
	if err := h.service.DB().Preload("User").Where("project_vacancy_id = ?", vacancyID).Find(&responses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type ResponseInfo struct {
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Phone     string    `json:"phone"`
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"created_at"`
	}

	var result []ResponseInfo
	for _, r := range responses {
		result = append(result, ResponseInfo{
			FirstName: r.User.FirstName,
			LastName:  r.User.LastName,
			Email:     r.User.Email,
			Phone:     r.User.Phone,
			Message:   r.Message,
			CreatedAt: r.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}
