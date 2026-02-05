package test

import (
	"encoding/json"
	"postapi/app/models"
	"testing"
)

func TestPostModel(t *testing.T) {
	post := models.Post{
		ID:      1,
		Title:   "Test Title",
		Content: "Test Content",
		Author:  "Test Author",
	}

	if post.ID != 1 {
		t.Errorf("expected ID 1, got %d", post.ID)
	}
	if post.Title != "Test Title" {
		t.Errorf("expected Title 'Test Title', got '%s'", post.Title)
	}
}

func TestPostRequestJSON(t *testing.T) {
	jsonData := `{
        "title": "Test Title",
        "content": "Test Content",
        "author": "Test Author"
    }`

	var postReq models.PostRequest
	err := json.Unmarshal([]byte(jsonData), &postReq)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if postReq.Title != "Test Title" {
		t.Errorf("expected Title 'Test Title', got '%s'", postReq.Title)
	}
	if postReq.Content != "Test Content" {
		t.Errorf("expected Content 'Test Content', got '%s'", postReq.Content)
	}
	if postReq.Author != "Test Author" {
		t.Errorf("expected Author 'Test Author', got '%s'", postReq.Author)
	}
}

func TestJsonPostSerialization(t *testing.T) {
	jsonPost := models.JsonPost{
		ID:      1,
		Title:   "Test Title",
		Content: "Test Content",
		Author:  "Test Author",
	}

	data, err := json.Marshal(jsonPost)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	var result models.JsonPost
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if result.ID != jsonPost.ID {
		t.Errorf("expected ID %d, got %d", jsonPost.ID, result.ID)
	}
	if result.Title != jsonPost.Title {
		t.Errorf("expected Title '%s', got '%s'", jsonPost.Title, result.Title)
	}
}

func TestUserModel(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	if user.Username != "testuser" {
		t.Errorf("expected Username 'testuser', got '%s'", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected Email 'test@example.com', got '%s'", user.Email)
	}
}

func TestUserRequestJSON(t *testing.T) {
	jsonData := `{
        "username": "testuser",
        "password": "password123",
        "email": "test@example.com"
    }`

	var userReq models.UserRequest
	err := json.Unmarshal([]byte(jsonData), &userReq)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if userReq.Username != "testuser" {
		t.Errorf("expected Username 'testuser', got '%s'", userReq.Username)
	}
	if userReq.Password != "password123" {
		t.Errorf("expected Password 'password123', got '%s'", userReq.Password)
	}
	if userReq.Email != "test@example.com" {
		t.Errorf("expected Email 'test@example.com', got '%s'", userReq.Email)
	}
}

func TestJsonUserOmitsPassword(t *testing.T) {
	jsonUser := models.JsonUser{
		Username: "testuser",
		Email:    "test@example.com",
	}

	data, err := json.Marshal(jsonUser)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	// Check that password field doesn't exist in JSON
	jsonString := string(data)
	if contains(jsonString, "password") {
		t.Error("JsonUser should not contain password field")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
