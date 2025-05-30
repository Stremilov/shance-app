package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	// Устанавливаем префиксы для переменных окружения
	viper.SetEnvPrefix("SHANCE")
	viper.BindEnv("server.port", "SHANCE_SERVER_PORT")
	viper.BindEnv("database.host", "SHANCE_DB_HOST")
	viper.BindEnv("database.port", "SHANCE_DB_PORT")
	viper.BindEnv("database.user", "SHANCE_DB_USER")
	viper.BindEnv("database.password", "SHANCE_DB_PASSWORD")
	viper.BindEnv("database.dbname", "SHANCE_DB_NAME")
	viper.BindEnv("jwt.secret", "SHANCE_JWT_SECRET")
	viper.BindEnv("jwt.access_token_ttl", "SHANCE_JWT_ACCESS_TTL")
	viper.BindEnv("jwt.refresh_token_ttl", "SHANCE_JWT_REFRESH_TTL")

	// Устанавливаем значения по умолчанию
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "shance_db")
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.access_token_ttl", "15m")
	viper.SetDefault("jwt.refresh_token_ttl", "720h")

	// Пытаемся прочитать конфигурационный файл
	if err := viper.ReadInConfig(); err != nil {
		// Если файл не найден, это не ошибка, так как у нас есть значения по умолчанию
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("jwt.secret"),
			AccessTokenTTL:  viper.GetDuration("jwt.access_token_ttl"),
			RefreshTokenTTL: viper.GetDuration("jwt.refresh_token_ttl"),
		},
	}

	return config, nil
}
