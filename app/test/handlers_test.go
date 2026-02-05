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
	posts    []*models.Post
	users    []*models.User
	follows  []*models.UserFollow
	profiles map[string]*models.Profile
	// Control error behavior
	shouldFailCreate        bool
	shouldFailUpdate        bool
	shouldFailGet           bool
	shouldFailGetByUser     bool
	shouldFailDelete        bool
	shouldFailRegister      bool
	shouldFailLogin         bool
	shouldFailGetUser       bool
	shouldFailCreateFollow  bool
	shouldFailRemoveFollow  bool
	shouldFailCreateProfile bool
	shouldFailUpdateProfile bool
	shouldFailGetProfile    bool
	loginReturnUser         *models.User
	getByUserReturnPosts    []*models.Post
	getUserReturnUser       *models.User
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

func (m *MockDB) UpdatePost(p *models.Post) error {
	if m.shouldFailUpdate {
		return errors.New("mock update error")
	}
	for i, post := range m.posts {
		if post.ID == p.ID && post.Author == p.Author {
			m.posts[i] = p
			return nil
		}
	}
	return errors.New("post not found or unauthorized")
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

func (m *MockDB) CreateFollow(f *models.UserFollow) error {
	if m.shouldFailCreateFollow {
		return errors.New("mock create follow error")
	}
	m.follows = append(m.follows, f)
	return nil
}

func (m *MockDB) Removefollow(f *models.UserFollow) error {
	if m.shouldFailRemoveFollow {
		return errors.New("mock remove follow error")
	}
	for i, follow := range m.follows {
		if follow.FollowerUsername == f.FollowerUsername && follow.FollowedUsername == f.FollowedUsername {
			m.follows = append(m.follows[:i], m.follows[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockDB) GetFollowers(username string) ([]string, error) {
	if m.shouldFailGetUser {
		return nil, errors.New("mock get followers error")
	}
	var followers []string
	for _, follow := range m.follows {
		if follow.FollowedUsername == username {
			followers = append(followers, follow.FollowerUsername)
		}
	}
	return followers, nil
}

func (m *MockDB) GetFollowings(username string) ([]string, error) {
	if m.shouldFailGetUser {
		return nil, errors.New("mock get following error")
	}
	var following []string
	for _, follow := range m.follows {
		if follow.FollowerUsername == username {
			following = append(following, follow.FollowedUsername)
		}
	}
	return following, nil
}

func (m *MockDB) GetProfile(username string) (*models.Profile, error) {
	if m.shouldFailGetProfile {
		return nil, errors.New("mock get profile error")
	}
	if m.profiles != nil {
		if profile, exists := m.profiles[username]; exists {
			return profile, nil
		}
	}
	return &models.Profile{
		Username:       username,
		Description:    "Test description",
		ProfilePicture: "test.jpg",
	}, nil
}

func (m *MockDB) CreateProfile(p *models.Profile) error {
	if m.shouldFailCreateProfile {
		return errors.New("mock create profile error")
	}
	if m.profiles == nil {
		m.profiles = make(map[string]*models.Profile)
	}
	m.profiles[p.Username] = p
	return nil
}

func (m *MockDB) UpdateProfile(p *models.Profile) error {
	if m.shouldFailUpdateProfile {
		return errors.New("mock update profile error")
	}
	if m.profiles == nil {
		m.profiles = make(map[string]*models.Profile)
	}
	if _, exists := m.profiles[p.Username]; !exists {
		return errors.New("profile not found")
	}
	m.profiles[p.Username] = p
	return nil
}

func setupTestApp(mockDB *MockDB) *app.App {
	a := app.New()
	a.DB = mockDB
	return a
}

func TestCreatePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
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

func TestUpdatePostHandler(t *testing.T) {
	tests := []struct {
		name                 string
		postID               string
		requestBody          interface{}
		withAuth             bool
		username             string
		mockPosts            []*models.Post
		mockShouldFailGet    bool
		mockShouldFailUpdate bool
		expectedStatus       int
	}{
		{
			name:   "successful update - all fields",
			postID: "1",
			requestBody: models.PostRequest{
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			withAuth: true,
			username: "testuser",
			mockPosts: []*models.Post{
				{ID: 1, Title: "Original Title", Content: "Original Content", Author: "testuser"},
			},
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusOK,
		},
		{
			name:   "successful update - partial fields (only title)",
			postID: "1",
			requestBody: models.PostRequest{
				Title:   "Updated Title",
				Content: "",
			},
			withAuth: true,
			username: "testuser",
			mockPosts: []*models.Post{
				{ID: 1, Title: "Original Title", Content: "Original Content", Author: "testuser"},
			},
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusOK,
		},
		{
			name:   "successful update - partial fields (only content)",
			postID: "1",
			requestBody: models.PostRequest{
				Title:   "",
				Content: "Updated Content",
			},
			withAuth: true,
			username: "testuser",
			mockPosts: []*models.Post{
				{ID: 1, Title: "Original Title", Content: "Original Content", Author: "testuser"},
			},
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusOK,
		},
		{
			name:   "missing authorization",
			postID: "1",
			requestBody: models.PostRequest{
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			withAuth:             false,
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusUnauthorized,
		},
		{
			name:                 "invalid request body",
			postID:               "1",
			requestBody:          "invalid json",
			withAuth:             true,
			username:             "testuser",
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusBadRequest,
		},
		{
			name:   "invalid post id",
			postID: "invalid",
			requestBody: models.PostRequest{
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			withAuth:             true,
			username:             "testuser",
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusBadRequest,
		},
		{
			name:   "database error on get post",
			postID: "1",
			requestBody: models.PostRequest{
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			withAuth:             true,
			username:             "testuser",
			mockShouldFailGet:    true,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusInternalServerError,
		},
		{
			name:   "database error on update post",
			postID: "1",
			requestBody: models.PostRequest{
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			withAuth: true,
			username: "testuser",
			mockPosts: []*models.Post{
				{ID: 1, Title: "Original Title", Content: "Original Content", Author: "testuser"},
			},
			mockShouldFailGet:    false,
			mockShouldFailUpdate: true,
			expectedStatus:       http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				posts:            tt.mockPosts,
				shouldFailGet:    tt.mockShouldFailGet,
				shouldFailUpdate: tt.mockShouldFailUpdate,
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

			req, err := http.NewRequest("PATCH", "/api/posts/"+tt.postID, bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			if tt.withAuth {
				ctx := context.WithValue(req.Context(), app.UsernameKey, tt.username)
				req = req.WithContext(ctx)
			}

			req = mux.SetURLVars(req, map[string]string{"post_id": tt.postID})

			rr := httptest.NewRecorder()
			handler := a.UpdatePostHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.JsonPost
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}

				// Check if partial update worked correctly
				reqParsed := tt.requestBody.(models.PostRequest)
				if reqParsed.Title != "" && response.Title != reqParsed.Title {
					t.Errorf("expected title %s, got %s", reqParsed.Title, response.Title)
				}
				if reqParsed.Content != "" && response.Content != reqParsed.Content {
					t.Errorf("expected content %s, got %s", reqParsed.Content, response.Content)
				}
				// Check that empty fields preserved original values
				if reqParsed.Title == "" && len(tt.mockPosts) > 0 && response.Title != tt.mockPosts[0].Title {
					t.Errorf("expected title to be preserved as %s, got %s", tt.mockPosts[0].Title, response.Title)
				}
				if reqParsed.Content == "" && len(tt.mockPosts) > 0 && response.Content != tt.mockPosts[0].Content {
					t.Errorf("expected content to be preserved as %s, got %s", tt.mockPosts[0].Content, response.Content)
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
		mockShouldFail bool
		expectedStatus int
	}{
		{
			name:           "successful get profile",
			username:       "testuser",
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "database error",
			username:       "testuser",
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				shouldFailGetProfile: tt.mockShouldFail,
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
func TestFollowHandler(t *testing.T) {
	tests := []struct {
		name                string
		withAuth            bool
		username            string
		followedUser        string
		mockShouldFail      bool
		expectedStatus      int
		expectedFollowCount int
	}{
		{
			name:                "successful follow",
			withAuth:            true,
			username:            "user1",
			followedUser:        "user2",
			mockShouldFail:      false,
			expectedStatus:      http.StatusOK,
			expectedFollowCount: 1,
		},
		{
			name:           "missing authorization",
			withAuth:       false,
			username:       "",
			followedUser:   "user2",
			mockShouldFail: false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "database error",
			withAuth:       true,
			username:       "user1",
			followedUser:   "user2",
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{shouldFailCreateFollow: tt.mockShouldFail}
			a := setupTestApp(mockDB)

			req, err := http.NewRequest("POST", "/api/follow/"+tt.followedUser, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.withAuth {
				ctx := context.WithValue(req.Context(), app.UsernameKey, tt.username)
				req = req.WithContext(ctx)
			}

			req = mux.SetURLVars(req, map[string]string{"username": tt.followedUser})

			rr := httptest.NewRecorder()
			handler := a.FollowHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.JsonUserFollow
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}

				if response.FollowerUsername != tt.username {
					t.Errorf("expected follower username %s, got %s", tt.username, response.FollowerUsername)
				}

				if response.FollowedUsername != tt.followedUser {
					t.Errorf("expected followed username %s, got %s", tt.followedUser, response.FollowedUsername)
				}

				if len(mockDB.follows) != tt.expectedFollowCount {
					t.Errorf("expected %d follows in database, got %d", tt.expectedFollowCount, len(mockDB.follows))
				}
			}
		})
	}
}

func TestUnfollowHandler(t *testing.T) {
	tests := []struct {
		name                string
		withAuth            bool
		username            string
		unfollowedUser      string
		mockShouldFail      bool
		expectedStatus      int
		initialFollows      []*models.UserFollow
		expectedFollowCount int
	}{
		{
			name:           "successful unfollow",
			withAuth:       true,
			username:       "user1",
			unfollowedUser: "user2",
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			initialFollows: []*models.UserFollow{
				{FollowerUsername: "user1", FollowedUsername: "user2"},
				{FollowerUsername: "user1", FollowedUsername: "user3"},
			},
			expectedFollowCount: 1,
		},
		{
			name:           "missing authorization",
			withAuth:       false,
			username:       "",
			unfollowedUser: "user2",
			mockShouldFail: false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "database error",
			withAuth:       true,
			username:       "user1",
			unfollowedUser: "user2",
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unfollow non-existent follow",
			withAuth:       true,
			username:       "user1",
			unfollowedUser: "user5",
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			initialFollows: []*models.UserFollow{
				{FollowerUsername: "user1", FollowedUsername: "user2"},
			},
			expectedFollowCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				shouldFailRemoveFollow: tt.mockShouldFail,
				follows:                tt.initialFollows,
			}
			a := setupTestApp(mockDB)

			req, err := http.NewRequest("DELETE", "/api/follow/"+tt.unfollowedUser, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.withAuth {
				ctx := context.WithValue(req.Context(), app.UsernameKey, tt.username)
				req = req.WithContext(ctx)
			}

			req = mux.SetURLVars(req, map[string]string{"username": tt.unfollowedUser})

			rr := httptest.NewRecorder()
			handler := a.UnfollowHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.JsonUserFollow
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}

				if response.FollowerUsername != tt.username {
					t.Errorf("expected follower username %s, got %s", tt.username, response.FollowerUsername)
				}

				if response.FollowedUsername != tt.unfollowedUser {
					t.Errorf("expected followed username %s, got %s", tt.unfollowedUser, response.FollowedUsername)
				}

				if tt.initialFollows != nil && len(mockDB.follows) != tt.expectedFollowCount {
					t.Errorf("expected %d follows remaining in database, got %d", tt.expectedFollowCount, len(mockDB.follows))
				}
			}
		})
	}
}

func TestGetFollowersHandler(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		initialFollows []*models.UserFollow
		initialUsers   []*models.User
		mockShouldFail bool
		expectedStatus int
		expectedCount  int
	}{
		{
			name:     "successful get followers",
			username: "user1",
			initialFollows: []*models.UserFollow{
				{FollowerUsername: "user2", FollowedUsername: "user1"},
				{FollowerUsername: "user3", FollowedUsername: "user1"},
			},
			initialUsers: []*models.User{
				{Username: "user2", Email: "user2@example.com"},
				{Username: "user3", Email: "user3@example.com"},
			},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "no followers",
			username:       "user1",
			initialFollows: []*models.UserFollow{},
			initialUsers:   []*models.User{},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "database error on GetFollowers",
			username:       "user1",
			initialFollows: []*models.UserFollow{},
			initialUsers:   []*models.User{},
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
		{
			name:     "database error on GetUserByUsername",
			username: "user1",
			initialFollows: []*models.UserFollow{
				{FollowerUsername: "user2", FollowedUsername: "user1"},
			},
			initialUsers:   []*models.User{},
			mockShouldFail: false,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				follows:           tt.initialFollows,
				users:             tt.initialUsers,
				shouldFailGetUser: tt.mockShouldFail,
			}
			a := setupTestApp(mockDB)

			req, err := http.NewRequest("GET", "/api/"+tt.username+"/followers", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"username": tt.username})

			rr := httptest.NewRecorder()
			handler := a.GetFollowersHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response []models.JsonUser
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}

				if len(response) != tt.expectedCount {
					t.Errorf("expected %d followers, got %d", tt.expectedCount, len(response))
				}
			}
		})
	}
}

func TestGetFollowingHandler(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		initialFollows []*models.UserFollow
		initialUsers   []*models.User
		mockShouldFail bool
		expectedStatus int
		expectedCount  int
	}{
		{
			name:     "successful get following",
			username: "user1",
			initialFollows: []*models.UserFollow{
				{FollowerUsername: "user1", FollowedUsername: "user2"},
				{FollowerUsername: "user1", FollowedUsername: "user3"},
			},
			initialUsers: []*models.User{
				{Username: "user2", Email: "user2@example.com"},
				{Username: "user3", Email: "user3@example.com"},
			},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "no following",
			username:       "user1",
			initialFollows: []*models.UserFollow{},
			initialUsers:   []*models.User{},
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "database error on GetFollowing",
			username:       "user1",
			initialFollows: []*models.UserFollow{},
			initialUsers:   []*models.User{},
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
		{
			name:     "database error on GetUserByUsername",
			username: "user1",
			initialFollows: []*models.UserFollow{
				{FollowerUsername: "user1", FollowedUsername: "user2"},
			},
			initialUsers:   []*models.User{},
			mockShouldFail: false,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				follows:           tt.initialFollows,
				users:             tt.initialUsers,
				shouldFailGetUser: tt.mockShouldFail,
			}
			a := setupTestApp(mockDB)

			req, err := http.NewRequest("GET", "/api/"+tt.username+"/following", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"username": tt.username})

			rr := httptest.NewRecorder()
			handler := a.GetFollowingHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response []models.JsonUser
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}

				if len(response) != tt.expectedCount {
					t.Errorf("expected %d following, got %d", tt.expectedCount, len(response))
				}
			}
		})
	}
}

func TestCreateProfileHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		withAuth       bool
		username       string
		mockShouldFail bool
		expectedStatus int
	}{
		{
			name: "successful profile creation",
			requestBody: models.ProfileRequest{
				Description:    "Test description",
				ProfilePicture: "test.jpg",
			},
			withAuth:       true,
			username:       "testuser",
			mockShouldFail: false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing authorization",
			requestBody: models.ProfileRequest{
				Description:    "Test description",
				ProfilePicture: "test.jpg",
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
			requestBody: models.ProfileRequest{
				Description:    "Test description",
				ProfilePicture: "test.jpg",
			},
			withAuth:       true,
			username:       "testuser",
			mockShouldFail: true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{shouldFailCreateProfile: tt.mockShouldFail}
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

			req, err := http.NewRequest("POST", "/api/profile", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			if tt.withAuth {
				ctx := context.WithValue(req.Context(), app.UsernameKey, tt.username)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()
			handler := a.CreateProfileHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.JsonProfile
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}

				if response.Username != tt.username {
					t.Errorf("expected username %s, got %s", tt.username, response.Username)
				}
			}
		})
	}
}

func TestUpdateProfileHandler(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          interface{}
		withAuth             bool
		username             string
		existingProfile      *models.Profile
		mockShouldFailGet    bool
		mockShouldFailUpdate bool
		expectedStatus       int
	}{
		{
			name: "successful profile update - all fields",
			requestBody: models.ProfileRequest{
				Description:    "Updated description",
				ProfilePicture: "updated.jpg",
			},
			withAuth: true,
			username: "testuser",
			existingProfile: &models.Profile{
				Username:       "testuser",
				Description:    "Old description",
				ProfilePicture: "old.jpg",
			},
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusOK,
		},
		{
			name: "successful profile update - partial fields",
			requestBody: models.ProfileRequest{
				Description:    "Updated description",
				ProfilePicture: "",
			},
			withAuth: true,
			username: "testuser",
			existingProfile: &models.Profile{
				Username:       "testuser",
				Description:    "Old description",
				ProfilePicture: "old.jpg",
			},
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusOK,
		},
		{
			name: "missing authorization",
			requestBody: models.ProfileRequest{
				Description:    "Updated description",
				ProfilePicture: "updated.jpg",
			},
			withAuth:             false,
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusUnauthorized,
		},
		{
			name:                 "invalid request body",
			requestBody:          "invalid json",
			withAuth:             true,
			username:             "testuser",
			mockShouldFailGet:    false,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusBadRequest,
		},
		{
			name: "database error on get profile",
			requestBody: models.ProfileRequest{
				Description:    "Updated description",
				ProfilePicture: "updated.jpg",
			},
			withAuth:             true,
			username:             "testuser",
			mockShouldFailGet:    true,
			mockShouldFailUpdate: false,
			expectedStatus:       http.StatusInternalServerError,
		},
		{
			name: "database error on update profile",
			requestBody: models.ProfileRequest{
				Description:    "Updated description",
				ProfilePicture: "updated.jpg",
			},
			withAuth: true,
			username: "testuser",
			existingProfile: &models.Profile{
				Username:       "testuser",
				Description:    "Old description",
				ProfilePicture: "old.jpg",
			},
			mockShouldFailGet:    false,
			mockShouldFailUpdate: true,
			expectedStatus:       http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				shouldFailGetProfile:    tt.mockShouldFailGet,
				shouldFailUpdateProfile: tt.mockShouldFailUpdate,
				profiles:                make(map[string]*models.Profile),
			}

			if tt.existingProfile != nil {
				mockDB.profiles[tt.username] = tt.existingProfile
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

			req, err := http.NewRequest("PATCH", "/api/profile", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			if tt.withAuth {
				ctx := context.WithValue(req.Context(), app.UsernameKey, tt.username)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()
			handler := a.UpdateProfileHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.JsonProfile
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("could not decode response: %v", err)
				}

				if response.Username != tt.username {
					t.Errorf("expected username %s, got %s", tt.username, response.Username)
				}

				// Check if partial update worked correctly
				reqParsed := tt.requestBody.(models.ProfileRequest)
				if reqParsed.Description != "" && response.Description != reqParsed.Description {
					t.Errorf("expected description %s, got %s", reqParsed.Description, response.Description)
				}
				if reqParsed.Description == "" && tt.existingProfile != nil && response.Description != tt.existingProfile.Description {
					t.Errorf("expected description to be preserved as %s, got %s", tt.existingProfile.Description, response.Description)
				}
			}
		})
	}
}
