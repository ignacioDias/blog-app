package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"postapi/app"
	"postapi/app/models"
	"testing"

	"github.com/gorilla/mux"
)

// MockDB implements database.PostDB interface for testing
type MockDB struct {
	posts []*models.Post
	users []*models.User
	// Control error behavior
	shouldFailCreate     bool
	shouldFailGet        bool
	shouldFailGetByUser  bool
	shouldFailDelete     bool
	shouldFailRegister   bool
	shouldFailLogin      bool
	shouldFailGetUser    bool
	loginReturnUser      *models.User
	getByUserReturnPosts []*models.Post
	getUserReturnUser    *models.User
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

func (m *MockDB) GetPost(id int64) (*models.Post, error) {
	if m.shouldFailGet {
		return nil, errors.New("mock get error")
	}
	for _, post := range m.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}

func (m *MockDB) DeletePost(id int64, username string) error {
	if m.shouldFailDelete {
		return errors.New("mock delete error")
	}
	for i, post := range m.posts {
		if post.ID == id {
			m.posts = append(m.posts[:i], m.posts[i+1:]...)
			return nil
		}
	}
	return errors.New("post not found")
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

func (m *MockDB) GetUserByUsername(username string) (*models.User, error) {
	if m.shouldFailGetUser {
		return nil, errors.New("mock get user error")
	}
	if m.getUserReturnUser != nil {
		return m.getUserReturnUser, nil
	}
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
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
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Welcome to Post API"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreatePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		withAuth       bool
		username       string
		mockShouldFail bool
		expectedStatus int
	}{
		{
			name: "successful post creation",
			requestBody: models.PostRequest{
				Title:   "Test Post",
				Content: "Test Content",
			},
			withAuth:       true,
			username:       "testuser",
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing authorization",
			requestBody: models.PostRequest{
				Title:   "Test Post",
				Content: "Test Content",
			},
			withAuth:       false,
			mockShouldFail: false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			withAuth:       true,
			username:       "testuser",
			mockShouldFail: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "database error",
			requestBody: models.PostRequest{
				Title:   "Test Post",
				Content: "Test Content",
			},
			withAuth:       true,
			username:       "testuser",
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{shouldFailCreate: tt.mockShouldFail}
			a := setupTestApp(mockDB)

			var reqBody []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			req, err := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			if tt.withAuth {
				ctx := context.WithValue(req.Context(), app.UsernameKey, tt.username)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()
			handler := a.CreatePostHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestGetPostsByUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		mockPosts      []*models.Post
		mockShouldFail bool
		expectedStatus int
		expectedCount  int
	}{
		{
			name:     "successful get posts",
			username: "testuser",
			mockPosts: []*models.Post{
				{ID: 1, Title: "Post 1", Content: "Content 1", Author: "testuser"},
				{ID: 2, Title: "Post 2", Content: "Content 2", Author: "testuser"},
			},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "empty posts list",
			username:       "testuser",
			mockPosts:      []*models.Post{},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "database error",
			username:       "testuser",
			mockPosts:      nil,
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				posts:               tt.mockPosts,
				shouldFailGetByUser: tt.mockShouldFail,
			}
			a := setupTestApp(mockDB)

			req, err := http.NewRequest("GET", "/api/"+tt.username+"/posts", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"username": tt.username})

			rr := httptest.NewRecorder()
			handler := a.GetPostsByUserHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
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

func TestGetPostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockPosts      []*models.Post
		mockShouldFail bool
		expectedStatus int
	}{
		{
			name:   "successful get post",
			postID: "1",
			mockPosts: []*models.Post{
				{ID: 1, Title: "Test Post", Content: "Test Content", Author: "testuser"},
			},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid post id",
			postID:         "invalid",
			mockPosts:      []*models.Post{},
			mockShouldFail: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "database error",
			postID:         "1",
			mockPosts:      nil,
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				posts:         tt.mockPosts,
				shouldFailGet: tt.mockShouldFail,
			}
			a := setupTestApp(mockDB)

			req, err := http.NewRequest("GET", "/api/posts/"+tt.postID, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"post_id": tt.postID})

			rr := httptest.NewRecorder()
			handler := a.GetPostHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestDeletePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		withAuth       bool
		username       string
		mockPosts      []*models.Post
		mockShouldFail bool
		expectedStatus int
	}{
		{
			name:     "successful delete",
			postID:   "1",
			withAuth: true,
			username: "testuser",
			mockPosts: []*models.Post{
				{ID: 1, Title: "Test Post", Content: "Test Content", Author: "testuser"},
			},
			mockShouldFail: false,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing authorization",
			postID:         "1",
			withAuth:       false,
			mockShouldFail: false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid post id",
			postID:         "invalid",
			withAuth:       true,
			username:       "testuser",
			mockShouldFail: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "database error",
			postID:         "1",
			withAuth:       true,
			username:       "testuser",
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				posts:            tt.mockPosts,
				shouldFailDelete: tt.mockShouldFail,
			}
			a := setupTestApp(mockDB)

			req, err := http.NewRequest("DELETE", "/api/posts/"+tt.postID, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.withAuth {
				ctx := context.WithValue(req.Context(), app.UsernameKey, tt.username)
				req = req.WithContext(ctx)
			}

			req = mux.SetURLVars(req, map[string]string{"post_id": tt.postID})

			rr := httptest.NewRecorder()
			handler := a.DeletePostHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
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
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{shouldFailRegister: tt.mockShouldFail}
			a := setupTestApp(mockDB)

			var reqBody []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := a.RegisterUserHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
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
			name: "login failed - invalid credentials",
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

			var reqBody []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := a.LoginHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectToken && rr.Code == http.StatusOK {
				var response map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}
				if _, exists := response["token"]; !exists {
					t.Error("expected token in response")
				}
			}
		})
	}
}

func TestGetProfileHandler(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		mockUser       *models.User
		mockShouldFail bool
		expectedStatus int
	}{
		{
			name:     "successful get profile",
			username: "testuser",
			mockUser: &models.User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "database error",
			username:       "testuser",
			mockUser:       nil,
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				getUserReturnUser: tt.mockUser,
				shouldFailGetUser: tt.mockShouldFail,
			}
			a := setupTestApp(mockDB)

			req, err := http.NewRequest("GET", "/api/"+tt.username, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"username": tt.username})

			rr := httptest.NewRecorder()
			handler := a.GetProfileHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
