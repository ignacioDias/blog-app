package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"postapi/app"
	"postapi/app/models"
	"testing"
)

type MockDB struct {
	posts []*models.Post
	users []*models.User
	// Control error behavior
	shouldFailCreate     bool
	shouldFailGet        bool
	shouldFailGetByUser  bool
	shouldFailRegister   bool
	shouldFailLogin      bool
	loginReturnUser      *models.User
	getByUserReturnPosts []*models.Post
}

func (m *MockDB) Open() error {
	return nil
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) CreatePost(p *models.Post) error {
	if m.shouldFailCreate {
		return errors.New("mock create error")
	}
	p.ID = int64(len(m.posts) + 1)
	m.posts = append(m.posts, p)
	return nil
}

func (m *MockDB) GetPosts() ([]*models.Post, error) {
	if m.shouldFailGet {
		return nil, errors.New("mock get error")
	}
	return m.posts, nil
}

func (m *MockDB) GetPostsByUser(username string) ([]*models.Post, error) {
	if m.shouldFailGetByUser {
		return nil, errors.New("mock get by user error")
	}
	if m.getByUserReturnPosts != nil {
		return m.getByUserReturnPosts, nil
	}
	var userPosts []*models.Post
	for _, post := range m.posts {
		if post.Author == username {
			userPosts = append(userPosts, post)
		}
	}
	return userPosts, nil
}

func (m *MockDB) RegisterUser(u *models.User) error {
	if m.shouldFailRegister {
		return errors.New("mock register error")
	}
	m.users = append(m.users, u)
	return nil
}

func (m *MockDB) LoginUser(u *models.User) (*models.User, error) {
	if m.shouldFailLogin {
		return nil, errors.New("mock login error")
	}
	if m.loginReturnUser != nil {
		return m.loginReturnUser, nil
	}
	return u, nil
}

func setupTestApp(mockDB *MockDB) *app.App {
	a := app.New()
	a.DB = mockDB
	return a
}

func TestIndexHandler(t *testing.T) {
	mockDB := &MockDB{}
	a := setupTestApp(mockDB)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := a.IndexHandler()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Welcome to Post API"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreatePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockShouldFail bool
		expectedStatus int
	}{
		{
			name: "successful post creation",
			requestBody: models.PostRequest{
				Title:   "Test Title",
				Content: "Test Content",
			},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			mockShouldFail: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "database error",
			requestBody: models.PostRequest{
				Title:   "Test Title",
				Content: "Test Content",
			},
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{shouldFailCreate: tt.mockShouldFail}
			a := setupTestApp(mockDB)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req, err := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := a.CreatePostHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.JsonPost
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}
				if response.Title != "Test Title" {
					t.Errorf("expected title 'Test Title', got '%s'", response.Title)
				}
			}
		})
	}
}

func TestGetPostsHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockPosts      []*models.Post
		mockShouldFail bool
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful get posts",
			mockPosts: []*models.Post{
				{ID: 1, Title: "Post 1", Content: "Content 1", Author: "Author 1"},
				{ID: 2, Title: "Post 2", Content: "Content 2", Author: "Author 2"},
			},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "empty posts list",
			mockPosts:      []*models.Post{},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "database error",
			mockPosts:      nil,
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				posts:         tt.mockPosts,
				shouldFailGet: tt.mockShouldFail,
			}
			a := setupTestApp(mockDB)

			req, err := http.NewRequest("GET", "/api/posts", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := a.GetPostsHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response []models.JsonPost
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}
				if len(response) != tt.expectedCount {
					t.Errorf("expected %d posts, got %d", tt.expectedCount, len(response))
				}
			}
		})
	}
}

func TestRegisterUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockShouldFail bool
		expectedStatus int
	}{
		{
			name: "successful user registration",
			requestBody: models.UserRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
			},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			mockShouldFail: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "database error",
			requestBody: models.UserRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
			},
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{shouldFailRegister: tt.mockShouldFail}
			a := setupTestApp(mockDB)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req, err := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := a.RegisterUserHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.JsonUser
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}
				if response.Username != "testuser" {
					t.Errorf("expected username 'testuser', got '%s'", response.Username)
				}
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockShouldFail bool
		mockReturnUser *models.User
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "successful login",
			requestBody: models.UserRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockShouldFail: false,
			mockReturnUser: &models.User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			mockShouldFail: false,
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
		{
			name: "invalid credentials",
			requestBody: models.UserRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockShouldFail: true,
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				shouldFailLogin: tt.mockShouldFail,
				loginReturnUser: tt.mockReturnUser,
			}
			a := setupTestApp(mockDB)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := a.LoginHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectToken {
				var response map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}
				if response["token"] == nil {
					t.Error("expected token in response, got nil")
				}
				if response["user"] == nil {
					t.Error("expected user in response, got nil")
				}
			}
		})
	}
}
