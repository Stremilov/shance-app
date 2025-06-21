package repository

import (
	"github.com/levstremilov/shance-app/internal/models"
	"gorm.io/gorm"
)

type ProjectVacancyRepository struct {
	db *gorm.DB
}

func NewProjectVacancyRepository(db *gorm.DB) *ProjectVacancyRepository {
	return &ProjectVacancyRepository{db: db}
}

func (r *ProjectVacancyRepository) Create(vacancy *models.ProjectVacancy) error {
	return r.db.Create(vacancy).Error
}

func (r *ProjectVacancyRepository) FindByProjectID(projectID uint) ([]models.ProjectVacancy, error) {
	var vacancies []models.ProjectVacancy
	err := r.db.Preload("Technologies").Where("project_id = ?", projectID).Find(&vacancies).Error
	if err != nil {
		return nil, err
	}
	return vacancies, nil
}

func (r *ProjectVacancyRepository) DB() *gorm.DB {
	return r.db
}
