package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Email        string         `json:"email" gorm:"unique;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Role         string         `json:"role" gorm:"default:user"`
	LastLogin    time.Time      `json:"last_login"`
	Phone        string         `json:"phone"`
	Country      string         `json:"country"`
	City         string         `json:"city"`
	Tags         []Tag          `json:"tags" gorm:"many2many:user_tags;"`
	Projects     []Project      `json:"projects" gorm:"many2many:project_members;"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
