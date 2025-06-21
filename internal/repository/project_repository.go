package repository

import (
	"github.com/levstremilov/shance-app/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := r.db.Preload("Tags").Preload("User").First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}

func (r *ProjectRepository) List() ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.Preload("Tags").Preload("User").Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) Search(query string) ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Preload("Tags").Preload("User").
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) AddMember(projectID, userID uint, role string) error {
	member := models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}
	return r.db.Create(&member).Error
}

func (r *ProjectRepository) RemoveMember(projectID, userID uint) error {
	return r.db.Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&models.ProjectMember{}).Error
}

func (r *ProjectRepository) AddTag(projectID, tagID uint) error {
	projectTag := models.ProjectTag{
		ProjectID: projectID,
		TagID:     tagID,
	}
	return r.db.Create(&projectTag).Error
}

func (r *ProjectRepository) RemoveTag(projectID, tagID uint) error {
	return r.db.Where("project_id = ? AND tag_id = ?", projectID, tagID).
		Delete(&models.ProjectTag{}).Error
}

func (r *ProjectRepository) GetDB() *gorm.DB {
	return r.db
}
