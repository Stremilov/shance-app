package service

import (
	"fmt"

	"github.com/levstremilov/shance-app/internal/models"
	"github.com/levstremilov/shance-app/internal/repository"
	"gorm.io/gorm"
)

type ProjectServiceInterface interface {
	GetAll() ([]models.Project, error)
	GetByID(id uint) (*models.Project, error)
	Create(project *models.Project) error
	Update(project *models.Project) error
	Delete(id uint) error
	Search(query string) ([]models.Project, error)
	IsProjectOwner(projectID, userID uint) (bool, error)
	InviteMember(projectID uint, email, role string) (*models.ProjectMember, error)
	GetProjectMembers(projectID uint) ([]models.ProjectMember, error)
}

type ProjectService struct {
	projectRepo *repository.ProjectRepository
	tagRepo     *repository.TagRepository
}

func NewProjectService(projectRepo *repository.ProjectRepository, tagRepo *repository.TagRepository) ProjectServiceInterface {
	return &ProjectService{
		projectRepo: projectRepo,
		tagRepo:     tagRepo,
	}
}

func (s *ProjectService) GetAll() ([]models.Project, error) {
	return s.projectRepo.List()
}

func (s *ProjectService) GetByID(id uint) (*models.Project, error) {
	return s.projectRepo.GetByID(id)
}

func (s *ProjectService) Create(project *models.Project) error {
	return s.projectRepo.Create(project)
}

func (s *ProjectService) Update(project *models.Project) error {
	return s.projectRepo.Update(project)
}

func (s *ProjectService) Delete(id uint) error {
	return s.projectRepo.Delete(id)
}

func (s *ProjectService) Search(query string) ([]models.Project, error) {
	return s.projectRepo.Search(query)
}

func (s *ProjectService) IsProjectOwner(projectID, userID uint) (bool, error) {
	var member models.ProjectMember
	err := s.projectRepo.DB.Where("project_id = ? AND user_id = ? AND role = ?", projectID, userID, "owner").First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *ProjectService) InviteMember(projectID uint, email, role string) (*models.ProjectMember, error) {
	// Находим пользователя по email
	var user models.User
	if err := s.projectRepo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Проверяем, не является ли пользователь уже участником
	var existingMember models.ProjectMember
	err := s.projectRepo.DB.Where("project_id = ? AND user_id = ?", projectID, user.ID).First(&existingMember).Error
	if err == nil {
		return nil, fmt.Errorf("user is already a member of this project")
	}

	// Создаем запись об участнике
	member := &models.ProjectMember{
		ProjectID: projectID,
		UserID:    user.ID,
		Role:      role,
	}

	if err := s.projectRepo.DB.Create(member).Error; err != nil {
		return nil, err
	}

	// Загружаем данные пользователя
	if err := s.projectRepo.DB.Preload("User").First(member, member.ID).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (s *ProjectService) GetProjectMembers(projectID uint) ([]models.ProjectMember, error) {
	var members []models.ProjectMember
	err := s.projectRepo.DB.Preload("User").Where("project_id = ?", projectID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
