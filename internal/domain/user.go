package domain

import (
	"time"
)

type Tag struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
}

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FirstName    string    `json:"first_name" gorm:"not null"`
	LastName     string    `json:"last_name" gorm:"not null"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"`
	Tags         []Tag     `json:"tags" gorm:"many2many:user_tags;"`
	Country      string    `json:"country"`
	City         string    `json:"city"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
}

type TokenPair struct {
	AccessToken  string `json:"-"`
	RefreshToken string `json:"refresh_token"`
}
