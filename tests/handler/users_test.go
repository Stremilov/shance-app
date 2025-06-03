package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/levstremilov/shance-app/internal/domain"
	"github.com/levstremilov/shance-app/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetMe(c *gin.Context) (*domain.User, error) {
	args := m.Called(c)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), nil
}

func (m *MockUserService) GetByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), nil
}

func (m *MockUserService) UpdateByID(id uint, user *domain.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func setupUserTestRouter(handler *handler.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	users := r.Group("/api/v1/users")
	{
		users.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Next()
		})
		users.GET("/me", handler.GetMe)
		users.GET("/:id", handler.GetByID)
		users.PATCH("/update", handler.UpdateMe)
	}

	return r
}

func TestGetMe(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)
	router := setupUserTestRouter(handler)

	tests := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Успешное получение данных пользователя",
			mockSetup: func() {
				mockService.On("GetMe", mock.AnythingOfType("*gin.Context")).Return(&domain.User{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
					CreatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
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
				var response domain.User
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "John", response.FirstName)
				assert.Equal(t, "Doe", response.LastName)
				assert.Equal(t, "john@example.com", response.Email)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)
	router := setupUserTestRouter(handler)

	tests := []struct {
		name           string
		userID         string
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "Успешное получение пользователя",
			userID: "1",
			mockSetup: func() {
				mockService.On("GetByID", uint(1)).Return(&domain.User{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
					CreatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Неверный ID пользователя",
			userID:         "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:   "Пользователь не найден",
			userID: "999",
			mockSetup: func() {
				mockService.On("GetByID", uint(999)).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.userID, nil)
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
				var response domain.User
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "John", response.FirstName)
				assert.Equal(t, "Doe", response.LastName)
				assert.Equal(t, "john@example.com", response.Email)
			}
		})
	}
}

func TestUpdateMe(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)
	router := setupUserTestRouter(handler)

	tests := []struct {
		name    string
		payload struct {
			FirstName string   `json:"name"`
			LastName  string   `json:"title"`
			Phone     string   `json:"subtitle"`
			Role      string   `json:"description"`
			Tags      []string `json:"photo"`
			Country   string   `json:"tags"`
			City      string   `json:"city"`
		}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Успешное обновление всех данных",
			payload: struct {
				FirstName string   `json:"name"`
				LastName  string   `json:"title"`
				Phone     string   `json:"subtitle"`
				Role      string   `json:"description"`
				Tags      []string `json:"photo"`
				Country   string   `json:"tags"`
				City      string   `json:"city"`
			}{
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "+1234567890",
				Role:      "user",
				Tags:      []string{"tag1", "tag2"},
				Country:   "USA",
				City:      "New York",
			},
			mockSetup: func() {
				mockService.On("GetByID", uint(1)).Return(&domain.User{
					ID:        1,
					FirstName: "Old",
					LastName:  "Name",
					Phone:     "old_phone",
					Role:      "old_role",
					Tags:      []domain.Tag{{Name: "old_tag"}},
					Country:   "Old Country",
					City:      "Old City",
				}, nil)
				mockService.On("UpdateByID", uint(1), mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		// TODO добавить тесты на неполные данные
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/update", bytes.NewBuffer(payload))
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
				var response domain.User
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Проверяем все поля
				assert.Equal(t, tt.payload.FirstName, response.FirstName)
				assert.Equal(t, tt.payload.LastName, response.LastName)
				assert.Equal(t, tt.payload.Phone, response.Phone)
				assert.Equal(t, tt.payload.Role, response.Role)
				assert.Equal(t, tt.payload.Country, response.Country)
				assert.Equal(t, tt.payload.City, response.City)

				if len(tt.payload.Tags) > 0 {
					assert.Equal(t, len(tt.payload.Tags), len(response.Tags))
					for i, tag := range response.Tags {
						assert.Equal(t, tt.payload.Tags[i], tag.Name)
					}
				} else {
					assert.Empty(t, response.Tags)
				}
			}
		})
	}
}
