package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "shance_db")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("SERVER_PORT", "8000")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: %v", err)
	}

	config := &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBPort:     viper.GetString("DB_PORT"),
		ServerPort: viper.GetString("SERVER_PORT"),
	}

	return config, nil
}
