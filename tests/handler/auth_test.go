package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/handler"
	"github.com/levstremilov/shance-app/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(data service.RegisterData) (*domain.TokenPair, error) {
	args := m.Called(data)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenPair), nil
}

func (m *MockAuthService) Login(email, password string) (*domain.TokenPair, error) {
	args := m.Called(email, password)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenPair), nil
}

func (m *MockAuthService) RefreshToken(token string) (*domain.TokenPair, error) {
	args := m.Called(token)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenPair), nil
}

func (m *MockAuthService) ValidateToken(token string) (*service.Claims, error) {
	args := m.Called(token)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Claims), nil
}

func setupAuthTestRouter(handler *handler.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)
	}

	return r
}

func TestRegister(t *testing.T) {
	mockService := new(MockAuthService)
	handler := handler.NewAuthHandler(mockService)
	router := setupAuthTestRouter(handler)

	tests := []struct {
		name    string
		payload struct {
			FirstName string   `json:"first_name"`
			LastName  string   `json:"last_name"`
			Phone     string   `json:"phone"`
			Role      string   `json:"role"`
			Tags      []string `json:"tags"`
			Country   string   `json:"country"`
			City      string   `json:"city"`
			Email     string   `json:"email"`
			Password  string   `json:"password"`
		}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Успешная регистрация",
			payload: struct {
				FirstName string   `json:"first_name"`
				LastName  string   `json:"last_name"`
				Phone     string   `json:"phone"`
				Role      string   `json:"role"`
				Tags      []string `json:"tags"`
				Country   string   `json:"country"`
				City      string   `json:"city"`
				Email     string   `json:"email"`
				Password  string   `json:"password"`
			}{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "password123",
			},
			mockSetup: func() {
				mockService.On("Register", mock.AnythingOfType("service.RegisterData")).Return(
					&domain.TokenPair{
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
					},
					nil,
				)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Ошибка валидации - отсутствует email",
			payload: struct {
				FirstName string   `json:"first_name"`
				LastName  string   `json:"last_name"`
				Phone     string   `json:"phone"`
				Role      string   `json:"role"`
				Tags      []string `json:"tags"`
				Country   string   `json:"country"`
				City      string   `json:"city"`
				Email     string   `json:"email"`
				Password  string   `json:"password"`
			}{
				FirstName: "John",
				LastName:  "Doe",
				Password:  "password123",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				var response struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Error)
			} else {
				var response struct {
					AccessToken  string `json:"access_token"`
					RefreshToken string `json:"refresh_token"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	mockService := new(MockAuthService)
	handler := handler.NewAuthHandler(mockService)
	router := setupAuthTestRouter(handler)

	tests := []struct {
		name    string
		payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Успешный вход",
			payload: struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "john@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockService.On("Login", "john@example.com", "password123").Return(
					&domain.TokenPair{
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
					},
					nil,
				)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Неверные учетные данные",
			payload: struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "john@example.com",
				Password: "wrong-password",
			},
			mockSetup: func() {
				mockService.On("Login", "john@example.com", "wrong-password").Return(
					nil,
					assert.AnError,
				)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				var response struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Error)
			} else {
				var response struct {
					AccessToken  string `json:"access_token"`
					RefreshToken string `json:"refresh_token"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := handler.NewAuthHandler(mockService)
	router := setupAuthTestRouter(handler)

	tests := []struct {
		name    string
		payload struct {
			RefreshToken string `json:"refresh_token"`
		}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Успешное обновление токена",
			payload: struct {
				RefreshToken string `json:"refresh_token"`
			}{
				RefreshToken: "valid-refresh-token",
			},
			mockSetup: func() {
				mockService.On("RefreshToken", "valid-refresh-token").Return(
					&domain.TokenPair{
						AccessToken:  "new-access-token",
						RefreshToken: "new-refresh-token",
					},
					nil,
				)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Неверный refresh token",
			payload: struct {
				RefreshToken string `json:"refresh_token"`
			}{
				RefreshToken: "invalid-refresh-token",
			},
			mockSetup: func() {
				mockService.On("RefreshToken", "invalid-refresh-token").Return(
					nil,
					assert.AnError,
				)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				var response struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Error)
			} else {
				var response struct {
					AccessToken  string `json:"access_token"`
					RefreshToken string `json:"refresh_token"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
			}
		})
	}
}
