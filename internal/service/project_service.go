package service

import (
	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/repository"
)

type ProjectServiceInterface interface {
	GetAll() ([]*domain.Project, error)
	GetByID(id uint) (*domain.Project, error)
	Create(project *domain.Project) error
	Update(project *domain.Project) error
	Delete(id uint) error
	Search(query string) ([]*domain.Project, error)
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

func (s *ProjectService) GetAll() ([]*domain.Project, error) {
	return s.projectRepo.GetAll()
}

func (s *ProjectService) GetByID(id uint) (*domain.Project, error) {
	return s.projectRepo.GetByID(id)
}

func (s *ProjectService) Create(project *domain.Project) error {
	if len(project.Tags) > 0 {
		tagNames := make([]string, len(project.Tags))
		for i, tag := range project.Tags {
			tagNames[i] = tag.Name
		}
		tags, err := s.tagRepo.GetByNames(tagNames)
		if err != nil {
			return err
		}
		project.Tags = tags
	}
	return s.projectRepo.Create(project)
}

func (s *ProjectService) Update(project *domain.Project) error {
	if len(project.Tags) > 0 {
		tagNames := make([]string, len(project.Tags))
		for i, tag := range project.Tags {
			tagNames[i] = tag.Name
		}
		tags, err := s.tagRepo.GetByNames(tagNames)
		if err != nil {
			return err
		}
		project.Tags = tags
	}
	return s.projectRepo.Update(project)
}

func (s *ProjectService) Delete(id uint) error {
	return s.projectRepo.Delete(id)
}

func (s *ProjectService) Search(query string) ([]*domain.Project, error) {
	return s.projectRepo.Search(query)
}
