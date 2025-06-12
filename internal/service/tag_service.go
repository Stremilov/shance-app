package service

import (
	"github.com/levstremilov/shance-app/internal/models"
	"github.com/levstremilov/shance-app/internal/repository"
)

type TagServiceInterface interface {
	Create(tag *models.Tag) error
	GetByID(id uint) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(id uint) error
	List() ([]models.Tag, error)
	Search(query string) ([]models.Tag, error)
}

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

func (s *TagService) Create(tag *models.Tag) error {
	return s.tagRepo.Create(tag)
}

func (s *TagService) GetByID(id uint) (*models.Tag, error) {
	return s.tagRepo.GetByID(id)
}

func (s *TagService) Update(tag *models.Tag) error {
	return s.tagRepo.Update(tag)
}

func (s *TagService) Delete(id uint) error {
	return s.tagRepo.Delete(id)
}

func (s *TagService) List() ([]models.Tag, error) {
	return s.tagRepo.List()
}

func (s *TagService) Search(query string) ([]models.Tag, error) {
	return s.tagRepo.Search(query)
}
