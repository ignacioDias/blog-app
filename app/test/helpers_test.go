package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"postapi/app/models"
	"testing"
)

func TestMapPostToJson(t *testing.T) {
	// Since mapPostToJson is not exported, we test it indirectly through handler responses
	post := &models.Post{
		ID:      1,
		Title:   "Test Title",
		Content: "Test Content",
		Author:  "testuser",
	}

	mockDB := &MockDB{
		posts: []*models.Post{post},
	}
	a := setupTestApp(mockDB)

	req, err := http.NewRequest("GET", "/api/posts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var jsonPost models.JsonPost
	err = json.NewDecoder(rr.Body).Decode(&jsonPost)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if jsonPost.ID != post.ID {
		t.Errorf("expected ID %d, got %d", post.ID, jsonPost.ID)
	}
	if jsonPost.Title != post.Title {
		t.Errorf("expected Title %s, got %s", post.Title, jsonPost.Title)
	}
	if jsonPost.Content != post.Content {
		t.Errorf("expected Content %s, got %s", post.Content, jsonPost.Content)
	}
	if jsonPost.Author != post.Author {
		t.Errorf("expected Author %s, got %s", post.Author, jsonPost.Author)
	}
}

func TestMapUserToJson(t *testing.T) {
	// Test indirectly through GetProfile handler
	user := &models.User{
		Username: "testuser",
		Password: "hashedpassword", // Should not appear in JSON
		Email:    "test@example.com",
	}

	mockDB := &MockDB{
		getUserReturnUser: user,
	}
	a := setupTestApp(mockDB)

	req, err := http.NewRequest("GET", "/api/users/testuser", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var jsonUser models.JsonUser
	err = json.NewDecoder(rr.Body).Decode(&jsonUser)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if jsonUser.Username != user.Username {
		t.Errorf("expected Username %s, got %s", user.Username, jsonUser.Username)
	}
	if jsonUser.Email != user.Email {
		t.Errorf("expected Email %s, got %s", user.Email, jsonUser.Email)
	}

	// Verify password is not in the response
	var rawMap map[string]interface{}
	json.NewDecoder(bytes.NewReader(rr.Body.Bytes())).Decode(&rawMap)
	if _, exists := rawMap["password"]; exists {
		t.Error("password should not be in JSON response")
	}
}

func TestParseAndSendResponse(t *testing.T) {
	// Test parse and sendResponse indirectly through login handler
	mockDB := &MockDB{
		loginReturnUser: &models.User{
			Username: "testuser",
			Email:    "test@example.com",
		},
	}
	a := setupTestApp(mockDB)

	requestBody := models.UserRequest{
		Username: "testuser",
		Password: "password123",
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if _, exists := response["token"]; !exists {
		t.Error("expected token in response")
	}
	if _, exists := response["user"]; !exists {
		t.Error("expected user in response")
	}
}
