package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/levstremilov/shance-app/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo        UserRepository
	jwtSecret       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(id uint) (*domain.User, error)
}

func NewAuthService(
	userRepo UserRepository,
	jwtSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		jwtSecret:       []byte(jwtSecret),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"` // "access" или "refresh"
	jwt.RegisteredClaims
}

func (s *AuthService) Register(email, password string) (*domain.TokenPair, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) Login(email, password string) (*domain.TokenPair, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*domain.TokenPair, error) {
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return nil, errors.New("invalid token type")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) generateTokenPair(user *domain.User) (*domain.TokenPair, error) {
	// Генерируем access token
	accessClaims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Генерируем refresh token
	refreshClaims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func (s *AuthService) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
