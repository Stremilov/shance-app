package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthServiceInterface interface {
	Register(data RegisterData) (*domain.TokenPair, error)
	Login(email, password string) (*domain.TokenPair, error)
	RefreshToken(token string) (*domain.TokenPair, error)
	ValidateToken(token string) (*Claims, error)
}

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, accessTTL, refreshTTL time.Duration) AuthServiceInterface {
	return &AuthService{
		userRepo:   userRepo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

type RegisterData struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	Role      string
	Tags      []string
	Country   string
	City      string
}

func (s *AuthService) Register(data RegisterData) (*domain.TokenPair, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        data.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		Phone:        data.Phone,
		Role:         data.Role,
		Country:      data.Country,
		City:         data.City,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) Login(email, password string) (*domain.TokenPair, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*domain.TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	user, err := s.userRepo.GetByID(uint(userID))
	if err != nil {
		return nil, err
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) generateTokenPair(user *domain.User) (*domain.TokenPair, error) {
	accessToken, err := s.generateToken(user, s.accessTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user, s.refreshTTL)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateToken(user *domain.User, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
