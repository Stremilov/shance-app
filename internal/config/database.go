package config

import (
	"database/sql"
	"fmt"

	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/migration"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := migration.RunMigrations(sqlDB); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm instance: %w", err)
	}

	err = db.AutoMigrate(&domain.Project{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
