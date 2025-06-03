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

type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) Create(project *domain.Project) error {
	args := m.Called(project)
	if args.Error(0) == nil {
		project.ID = 1
		project.CreatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockProjectService) GetAll() ([]*domain.Project, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), nil
}

func (m *MockProjectService) GetByID(id uint) (*domain.Project, error) {
	args := m.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), nil
}

func (m *MockProjectService) Update(project *domain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProjectService) Search(query string) ([]*domain.Project, error) {
	args := m.Called(query)
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func setupTestRouter(handler *handler.ProjectHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	projects := r.Group("/api/v1/projects")
	{
		projects.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Next()
		})
		projects.POST("", handler.CreateProject)
		projects.GET("", handler.GetProjects)
		projects.GET("/search", handler.SearchProjects)
		projects.GET("/:id", handler.GetProject)
		projects.PUT("/:id", handler.UpdateProject)
		projects.DELETE("/:id", handler.DeleteProject)
	}

	return r
}

func TestCreateProject(t *testing.T) {
	mockService := new(MockProjectService)
	handler := handler.NewProjectHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name    string
		payload struct {
			Name        string   `json:"name"`
			Title       string   `json:"title"`
			Subtitle    string   `json:"subtitle"`
			Description string   `json:"description"`
			Photo       []string `json:"photo"`
			Tags        []string `json:"tags"`
		}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Успешное создание проекта",
			payload: struct {
				Name        string   `json:"name"`
				Title       string   `json:"title"`
				Subtitle    string   `json:"subtitle"`
				Description string   `json:"description"`
				Photo       []string `json:"photo"`
				Tags        []string `json:"tags"`
			}{
				Name:        "Test Project",
				Title:       "Test Title",
				Subtitle:    "Test Subtitle",
				Description: "Test Description",
				Photo:       []string{"photo1.jpg"},
				Tags:        []string{"tag1"},
			},
			mockSetup: func() {
				mockService.On("Create", mock.AnythingOfType("*domain.Project")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "Ошибка валидации - отсутствует имя",
			payload: struct {
				Name        string   `json:"name"`
				Title       string   `json:"title"`
				Subtitle    string   `json:"subtitle"`
				Description string   `json:"description"`
				Photo       []string `json:"photo"`
				Tags        []string `json:"tags"`
			}{
				Title:       "Test Title",
				Description: "Test Description",
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{
				Name:  "access_token",
				Value: "test-token",
			})

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
			}
		})
	}
}

func TestGetProjects(t *testing.T) {
	mockService := new(MockProjectService)
	handler := handler.NewProjectHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Успешное получение списка проектов",
			mockSetup: func() {
				mockService.On("GetAll").Return([]*domain.Project{
					{
						ID:          1,
						Name:        "Test Project",
						Title:       "Test Title",
						Description: "Test Description",
						Photo:       "[]",
						Tags:        []domain.Tag{},
						UserID:      1,
						CreatedAt:   time.Now(),
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
			req.AddCookie(&http.Cookie{
				Name:  "access_token",
				Value: "test-token",
			})
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
			}
		})
	}
}

func TestGetProject(t *testing.T) {
	mockService := new(MockProjectService)
	handler := handler.NewProjectHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name           string
		projectID      string
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name:      "Успешное получение проекта",
			projectID: "1",
			mockSetup: func() {
				mockService.On("GetByID", uint(1)).Return(&domain.Project{
					ID:          1,
					Name:        "Test Project",
					Title:       "Test Title",
					Description: "Test Description",
					Photo:       "[]",
					Tags:        []domain.Tag{},
					UserID:      1,
					CreatedAt:   time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Неверный ID проекта",
			projectID:      "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:      "Проект не найден",
			projectID: "999",
			mockSetup: func() {
				mockService.On("GetByID", uint(999)).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/"+tt.projectID, nil)
			req.AddCookie(&http.Cookie{
				Name:  "access_token",
				Value: "test-token",
			})
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
			}
		})
	}
}
