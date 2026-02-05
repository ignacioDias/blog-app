package test

import (
	"net/http"
	"net/http/httptest"
	"postapi/app"
	"testing"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	mockDB := &MockDB{}
	a := setupTestApp(mockDB)

	// Create a valid token
	username := "testuser"
	token, err := app.CreateToken(username)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	// Create a simple handler that checks if username is in context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUsername, ok := r.Context().Value(app.UsernameKey).(string)
		if !ok {
			t.Error("username not found in context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if ctxUsername != username {
			t.Errorf("expected username %s, got %s", username, ctxUsername)
		}
		w.WriteHeader(http.StatusOK)
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	middleware := a.AuthMiddleware(nextHandler)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	mockDB := &MockDB{}
	a := setupTestApp(mockDB)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
		w.WriteHeader(http.StatusOK)
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middleware := a.AuthMiddleware(nextHandler)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	mockDB := &MockDB{}
	a := setupTestApp(mockDB)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name  string
		token string
	}{
		{"no bearer prefix", "invalidtoken"},
		{"wrong prefix", "Basic token123"},
		{"empty after bearer", "Bearer "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", tt.token)

			rr := httptest.NewRecorder()
			middleware := a.AuthMiddleware(nextHandler)
			middleware.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
			}
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockDB := &MockDB{}
	a := setupTestApp(mockDB)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
		w.WriteHeader(http.StatusOK)
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	rr := httptest.NewRecorder()
	middleware := a.AuthMiddleware(nextHandler)
	middleware.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}
