package repository

import (
	"github.com/levstremilov/shance-app/internal/domain"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *domain.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) GetAll() ([]*domain.Project, error) {
	var projects []*domain.Project
	err := r.db.Model(&domain.Project{}).Preload("Tags").Preload("User").Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) GetByID(id uint) (*domain.Project, error) {
	var project domain.Project
	err := r.db.Model(&domain.Project{}).Preload("Tags").First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) Update(project *domain.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Project{}, id).Error
}

func (r *ProjectRepository) Search(query string) ([]*domain.Project, error) {
	var projects []*domain.Project
	err := r.db.Model(&domain.Project{}).Preload("Tags").Preload("User").Where("name ILIKE ? OR title ILIKE ? OR description ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&projects).Error
	return projects, err
}
