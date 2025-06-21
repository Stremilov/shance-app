package service

import (
	"fmt"
	"strconv"

	"github.com/levstremilov/shance-app/internal/models"
	"github.com/levstremilov/shance-app/internal/repository"
	"gorm.io/gorm"
)

type ProjectServiceInterface interface {
	GetAll() ([]models.Project, error)
	GetByID(id string) (*models.Project, error)
	Create(project *models.Project) error
	Update(project *models.Project) error
	Delete(id string) error
	Search(query string) ([]models.Project, error)
	IsProjectOwner(projectID, userID uint) (bool, error)
	InviteMember(projectID uint, email, role string) (*models.ProjectMember, error)
	GetProjectMembers(projectID uint) ([]models.ProjectMember, error)
}

type ProjectService struct {
	projectRepo *repository.ProjectRepository
	tagRepo     *repository.TagRepository
	userRepo    *repository.UserRepository
}

func NewProjectService(projectRepo *repository.ProjectRepository, tagRepo *repository.TagRepository, userRepo *repository.UserRepository) ProjectServiceInterface {
	return &ProjectService{
		projectRepo: projectRepo,
		tagRepo:     tagRepo,
		userRepo:    userRepo,
	}
}

func (s *ProjectService) GetAll() ([]models.Project, error) {
	return s.projectRepo.List()
}

func (s *ProjectService) GetByID(id string) (*models.Project, error) {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}
	return s.projectRepo.GetByID(uint(idUint))
}

func (s *ProjectService) Create(project *models.Project) error {
	return s.projectRepo.Create(project)
}

func (s *ProjectService) Update(project *models.Project) error {
	return s.projectRepo.Update(project)
}

func (s *ProjectService) Delete(id string) error {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}
	return s.projectRepo.Delete(uint(idUint))
}

func (s *ProjectService) Search(query string) ([]models.Project, error) {
	return s.projectRepo.Search(query)
}

func (s *ProjectService) IsProjectOwner(projectID, userID uint) (bool, error) {
	var member models.ProjectMember
	err := s.projectRepo.GetDB().Where("project_id = ? AND user_id = ? AND role = ?", projectID, userID, "owner").First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *ProjectService) InviteMember(projectID uint, email, role string) (*models.ProjectMember, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	var existingMember models.ProjectMember
	err = s.projectRepo.GetDB().Where("project_id = ? AND user_id = ?", projectID, user.ID).First(&existingMember).Error
	if err == nil {
		return nil, fmt.Errorf("user is already a member of this project")
	}

	err = s.projectRepo.AddMember(projectID, user.ID, role)
	if err != nil {
		return nil, err
	}

	member := &models.ProjectMember{}
	err = s.projectRepo.GetDB().Where("project_id = ? AND user_id = ?", projectID, user.ID).Preload("User").First(member).Error
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (s *ProjectService) GetProjectMembers(projectID uint) ([]models.ProjectMember, error) {
	var members []models.ProjectMember
	err := s.projectRepo.GetDB().Where("project_id = ?", projectID).Preload("User").Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
