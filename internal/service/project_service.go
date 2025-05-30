package service

import (
	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(project *domain.Project) error {
	return s.repo.Create(project)
}

func (s *ProjectService) GetAllProjects() ([]domain.Project, error) {
	return s.repo.GetAll()
}

func (s *ProjectService) GetProjectByID(id uint) (*domain.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectService) UpdateProject(project *domain.Project) error {
	return s.repo.Update(project)
}

func (s *ProjectService) DeleteProject(id uint) error {
	return s.repo.Delete(id)
} 