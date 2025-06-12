package models

// SwaggerUser представляет пользователя для Swagger документации
type SwaggerUser struct {
	ID        uint   `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Role      string `json:"role" example:"user"`
	Phone     string `json:"phone" example:"+1234567890"`
	Country   string `json:"country" example:"USA"`
	City      string `json:"city" example:"New York"`
	CreatedAt string `json:"created_at" example:"2024-03-12T15:04:05Z"`
	UpdatedAt string `json:"updated_at" example:"2024-03-12T15:04:05Z"`
}

// SwaggerProject представляет проект для Swagger документации
type SwaggerProject struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"my-project"`
	Title       string `json:"title" example:"My Project"`
	Subtitle    string `json:"subtitle" example:"A great project"`
	Description string `json:"description" example:"This is a project description"`
	Status      string `json:"status" example:"active"`
	UserID      uint   `json:"user_id" example:"1"`
	UserEmail   string `json:"user_email" example:"user@example.com"`
	CreatedAt   string `json:"created_at" example:"2024-03-12T15:04:05Z"`
	UpdatedAt   string `json:"updated_at" example:"2024-03-12T15:04:05Z"`
}

// SwaggerTag представляет тег для Swagger документации
type SwaggerTag struct {
	ID        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"important"`
	CreatedAt string `json:"created_at" example:"2024-03-12T15:04:05Z"`
	UpdatedAt string `json:"updated_at" example:"2024-03-12T15:04:05Z"`
}

// SwaggerListResponse представляет ответ со списком для Swagger документации
type SwaggerListResponse struct {
	Count    int64       `json:"count" example:"100"`
	Next     string      `json:"next" example:"/api/v1/users?page=2&page_size=10"`
	Previous string      `json:"previous" example:"/api/v1/users?page=1&page_size=10"`
	Results  interface{} `json:"results"`
}
