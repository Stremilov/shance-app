package repository

import (
	"github.com/levstremilov/shance-app/internal/models"

	"gorm.io/gorm"
)

type TagRepository struct {
	DB *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		DB: db,
	}
}

func (r *TagRepository) Create(tag *models.Tag) error {
	return r.DB.Create(tag).Error
}

func (r *TagRepository) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	if err := r.DB.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) Update(tag *models.Tag) error {
	return r.DB.Save(tag).Error
}

func (r *TagRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Tag{}, id).Error
}

func (r *TagRepository) List() ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) Search(query string) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.DB.Where("name ILIKE ?", "%"+query+"%").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) GetByUserID(userID uint) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.DB.Joins("JOIN user_tags ON user_tags.tag_id = tags.id").
		Where("user_tags.user_id = ?", userID).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) GetByProjectID(projectID uint) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.DB.Joins("JOIN project_tags ON project_tags.tag_id = tags.id").
		Where("project_tags.project_id = ?", projectID).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
