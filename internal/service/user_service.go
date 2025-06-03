package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/repository"
)

type UserServiceInterface interface {
	GetMe(c *gin.Context) (*domain.User, error)
	GetByID(id uint) (*domain.User, error)
	UpdateByID(id uint, user *domain.User) error
}

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetMe(c *gin.Context) (*domain.User, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, fmt.Errorf("user not authenticated")
	}
	return s.userRepo.GetByID(userID.(uint))
}

func (s *UserService) GetByID(id uint) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) UpdateByID(id uint, user *domain.User) error {
	return s.userRepo.UpdateByID(id, user)
}
