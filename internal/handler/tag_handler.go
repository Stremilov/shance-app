package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/models"
	"github.com/levstremilov/shance-app/internal/service"
	"gorm.io/gorm"
)

type TagHandler struct {
	tagService service.TagServiceInterface
}

func NewTagHandler(tagService service.TagServiceInterface) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

type CreateTagRequest struct {
	Name        string `json:"name" binding:"required" example:"Тег"`
}

type UpdateTagRequest struct {
	Name        string `json:"name" example:"Обновленный тег"`
}

type TagResponse struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Тег"`
}

// CreateTag godoc
// @Summary Создание нового тега
// @Description Создает новый тег в системе
// @Tags tags
// @Accept json
// @Produce json
// @Param request body CreateTagRequest true "Данные тега"
// @Success 201 {object} TagResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tag := &models.Tag{
		Name: req.Name,
	}

	if err := h.tagService.Create(tag); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}

	c.JSON(http.StatusCreated, response)
}


// UpdateTag godoc
// @Summary Обновление тега
// @Description Обновляет информацию о теге
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "ID тега"
// @Param request body UpdateTagRequest true "Данные тега"
// @Success 200 {object} TagResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tag ID"})
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tag := &models.Tag{
		Model: gorm.Model{ID: uint(id)},
		Name:  req.Name,
	}

	if err := h.tagService.Update(tag); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTag godoc
// @Summary Удаление тега
// @Description Удаляет тег по его ID
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "ID тега"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid tag ID"})
		return
	}

	if err := h.tagService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTags godoc
// @Summary Получение списка тегов
// @Description Возвращает список всех тегов
// @Tags tags
// @Accept json
// @Produce json
// @Success 200 {array} TagResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	tags, err := h.tagService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := make([]TagResponse, len(tags))
	for i, tag := range tags {
		response[i] = TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	c.JSON(http.StatusOK, response)
}

// SearchTags godoc
// @Summary Поиск тегов
// @Description Возвращает список тегов, соответствующих поисковому запросу
// @Tags tags
// @Accept json
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Success 200 {array} TagResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags/search [get]
func (h *TagHandler) SearchTags(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Search query is required"})
		return
	}

	tags, err := h.tagService.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := make([]TagResponse, len(tags))
	for i, tag := range tags {
		response[i] = TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	c.JSON(http.StatusOK, response)
}
