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
		Author:  "testuser",
	}

	if post.ID != 1 {
		t.Errorf("expected ID 1, got %d", post.ID)
	}
	if post.Title != "Test Title" {
		t.Errorf("expected Title 'Test Title', got '%s'", post.Title)
	}
	if post.Content != "Test Content" {
		t.Errorf("expected Content 'Test Content', got '%s'", post.Content)
	}
	if post.Author != "testuser" {
		t.Errorf("expected Author 'testuser', got '%s'", post.Author)
	}
}

func TestPostRequestJSON(t *testing.T) {
	jsonData := `{
        "title": "Test Title",
        "content": "Test Content"
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
}

func TestJsonPostSerialization(t *testing.T) {
	jsonPost := models.JsonPost{
		ID:      1,
		Title:   "Test Title",
		Content: "Test Content",
		Author:  "testuser",
	}

	data, err := json.Marshal(jsonPost)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	expected := `{"id":1,"title":"Test Title","content":"Test Content","author":"testuser"}`
	var expectedMap, actualMap map[string]interface{}

	json.Unmarshal([]byte(expected), &expectedMap)
	json.Unmarshal(data, &actualMap)

	if len(expectedMap) != len(actualMap) {
		t.Errorf("JSON structure mismatch")
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

func TestJsonUserSerialization(t *testing.T) {
	jsonUser := models.JsonUser{
		Username: "testuser",
		Email:    "test@example.com",
	}

	data, err := json.Marshal(jsonUser)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	var result models.JsonUser
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if result.Username != jsonUser.Username {
		t.Errorf("expected Username '%s', got '%s'", jsonUser.Username, result.Username)
	}
	if result.Email != jsonUser.Email {
		t.Errorf("expected Email '%s', got '%s'", jsonUser.Email, result.Email)
	}
}

func TestJsonUserDoesNotContainPassword(t *testing.T) {
	jsonUser := models.JsonUser{
		Username: "testuser",
		Email:    "test@example.com",
	}

	data, err := json.Marshal(jsonUser)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	jsonString := string(data)
	if len(jsonString) > 0 && (jsonString[0] != '{' || jsonString[len(jsonString)-1] != '}') {
		t.Error("expected valid JSON object")
	}

	// Ensure password field is not present
	var rawMap map[string]interface{}
	json.Unmarshal(data, &rawMap)
	if _, exists := rawMap["password"]; exists {
		t.Error("JsonUser should not contain password field")
	}
}
