package config

import (
	"os"
	"time"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	Server struct {
		Host string
		Port string
		Env  string
	}
	JWT struct {
		Secret          string
		AccessTokenTTL  time.Duration
		RefreshTokenTTL time.Duration
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			DBName   string
		}{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "levstremilov"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "shance_db"),
		},
		Server: struct {
			Host string
			Port string
			Env  string
		}{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8000"),
			Env:  getEnv("ENV", "prod"),
		},
		JWT: struct {
			Secret          string
			AccessTokenTTL  time.Duration
			RefreshTokenTTL time.Duration
		}{
			Secret:          getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 720 * time.Hour,
		},
	}

	return config, nil
}
