package repository

import (
	"github.com/levstremilov/shance-app/internal/domain"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) GetOrCreate(name string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.Model(&domain.Tag{}).Where("name = ?", name).First(&tag).Error
	if err == nil {
		return &tag, nil
	}

	tag = domain.Tag{Name: name}
	err = r.db.Create(&tag).Error
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *TagRepository) GetByNames(names []string) ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.db.Model(&domain.Tag{}).Where("name IN ?", names).Find(&tags).Error
	return tags, err
}
