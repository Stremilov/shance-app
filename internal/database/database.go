package database

import (
	"fmt"

	"github.com/levstremilov/shance-app/internal/config"
	"github.com/levstremilov/shance-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Tag{},
		&models.Project{},
		&models.ProjectMember{},
		&models.UserTag{},
		&models.ProjectTag{},
		&models.ProjectVacancy{},
		&models.VacancyTechnology{},
		&models.Technology{},
		&models.VacancyResponse{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
