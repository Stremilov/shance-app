package domain

import (
	"time"
)

type Project struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
	Tags        []Tag     `json:"tags" gorm:"many2many:project_tags;"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProjectTag struct {
	ProjectID uint `json:"project_id"`
	TagID     uint `json:"tag_id"`
}
