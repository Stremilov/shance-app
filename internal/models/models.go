package models

import (
	"time"

	"gorm.io/gorm"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Tag struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Users       []User    `gorm:"many2many:user_tags;"`
	Projects    []Project `gorm:"many2many:project_tags;"`
}

type Project struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Title       string
	Subtitle    string
	Description string
	Photo       string `gorm:"type:jsonb"`
	Status      string `gorm:"default:active"`
	StartDate   time.Time
	EndDate     time.Time
	UserID      uint
	User        User   `gorm:"foreignKey:UserID"`
	Tags        []Tag  `gorm:"many2many:project_tags;"`
	Members     []User `gorm:"many2many:project_members;"`
}

type ProjectMember struct {
	ProjectID uint
	UserID    uint
	Role      string
	JoinedAt  time.Time
	User      User    `gorm:"foreignKey:UserID"`
	Project   Project `gorm:"foreignKey:ProjectID"`
	gorm.Model
}

type UserTag struct {
	UserID uint
	TagID  uint
	gorm.Model
}

type ProjectTag struct {
	ProjectID uint
	TagID     uint
	gorm.Model
}
