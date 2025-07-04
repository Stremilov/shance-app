package config

import (
	"fmt"

	"github.com/levstremilov/shance-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func CleanDB(db *gorm.DB) error {
	// Отключаем проверку внешних ключей
	if err := db.Exec("SET CONSTRAINTS ALL DEFERRED").Error; err != nil {
		return fmt.Errorf("failed to defer constraints: %w", err)
	}

	// Удаляем таблицы в правильном порядке
	tables := []interface{}{
		&models.ProjectTag{},
		&models.Project{},
		&models.User{},
		&models.Tag{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	// Включаем проверку внешних ключей обратно
	if err := db.Exec("SET CONSTRAINTS ALL IMMEDIATE").Error; err != nil {
		return fmt.Errorf("failed to enable constraints: %w", err)
	}

	return nil
}
