package repository

import (
	"github.com/levstremilov/shance-app/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Project{},
	)
}
