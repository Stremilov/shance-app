package service

import (
	"github.com/levstremilov/shance-app/internal/models"
	"github.com/levstremilov/shance-app/internal/repository"
	"gorm.io/gorm"
)

type ProjectVacancyService struct {
	repo *repository.ProjectVacancyRepository
}

func NewProjectVacancyService(repo *repository.ProjectVacancyRepository) *ProjectVacancyService {
	return &ProjectVacancyService{repo: repo}
}

func (s *ProjectVacancyService) Create(vacancy *models.ProjectVacancy) error {
	return s.repo.Create(vacancy)
}

func (s *ProjectVacancyService) GetAll(project_id uint) ([]models.ProjectVacancy, error) {
	return s.repo.FindByProjectID(project_id)
}

func (s *ProjectVacancyService) DB() *gorm.DB {
	return s.repo.DB()
}
