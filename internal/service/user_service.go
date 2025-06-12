package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/models"
	"github.com/levstremilov/shance-app/internal/repository"
)

type UserServiceInterface interface {
	GetMe(c *gin.Context) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	Delete(id uint) error
	List() ([]models.User, error)
}

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetMe(c *gin.Context) (*models.User, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, fmt.Errorf("user not authenticated")
	}
	return s.userRepo.GetByID(userID.(uint))
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) Update(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) Create(user *models.User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) List() ([]models.User, error) {
	return s.userRepo.List()
}
