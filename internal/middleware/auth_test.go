package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock JWT Service for testing
type mockJWTService struct {
	validateFunc func(token string) (string, error)
}

func (m *mockJWTService) GenerateToken(username string) (string, error) {
	return "mock-token", nil
}

func (m *mockJWTService) ValidateToken(token string) (string, error) {
	if m.validateFunc != nil {
		return m.validateFunc(token)
	}
	return "", errors.New("not implemented")
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	mockService := &mockJWTService{}
	authMiddleware := NewAuthMiddleware(mockService)

	handler := authMiddleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when auth fails")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	mockService := &mockJWTService{}
	authMiddleware := NewAuthMiddleware(mockService)

	tests := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "token-without-bearer"},
		{"Wrong prefix", "Basic token123"},
		{"Only Bearer", "Bearer"},
		{"Bearer with space only", "Bearer "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := authMiddleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				t.Error("Handler should not be called")
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
			}
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockService := &mockJWTService{
		validateFunc: func(token string) (string, error) {
			return "", errors.New("invalid token")
		},
	}
	authMiddleware := NewAuthMiddleware(mockService)

	handler := authMiddleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	expectedUsername := "testuser"
	mockService := &mockJWTService{
		validateFunc: func(token string) (string, error) {
			if token == "valid-token" {
				return expectedUsername, nil
			}
			return "", errors.New("invalid token")
		},
	}
	authMiddleware := NewAuthMiddleware(mockService)

	handlerCalled := false
	handler := authMiddleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		// Check if username is in context
		username, ok := r.Context().Value(UsernameKey).(string)
		if !ok {
			t.Error("Username not found in context")
			return
		}
		if username != expectedUsername {
			t.Errorf("Expected username %s, got %s", expectedUsername, username)
		}

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	handler(w, req)

	if !handlerCalled {
		t.Error("Handler should have been called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_ContextKey(t *testing.T) {
	// Test that the context key is of the correct type
	if UsernameKey != contextKey("username") {
		t.Error("UsernameKey should be 'username'")
	}

	// Test setting and getting from context
	ctx := context.WithValue(context.Background(), UsernameKey, "testuser")
	username, ok := ctx.Value(UsernameKey).(string)
	if !ok {
		t.Error("Failed to retrieve username from context")
	}
	if username != "testuser" {
		t.Errorf("Expected 'testuser', got %s", username)
	}
}

func TestAuthMiddleware_DifferentTokens(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		expectedUser   string
		expectedStatus int
	}{
		{"User1", "token1", "user1", http.StatusOK},
		{"User2", "token2", "user2", http.StatusOK},
		{"Admin", "admin-token", "admin", http.StatusOK},
		{"Invalid", "invalid", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockJWTService{
				validateFunc: func(token string) (string, error) {
					switch token {
					case "token1":
						return "user1", nil
					case "token2":
						return "user2", nil
					case "admin-token":
						return "admin", nil
					default:
						return "", errors.New("invalid token")
					}
				},
			}

			authMiddleware := NewAuthMiddleware(mockService)

			handler := authMiddleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				username, _ := r.Context().Value(UsernameKey).(string)
				if username != tt.expectedUser {
					t.Errorf("Expected username %s, got %s", tt.expectedUser, username)
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestNewAuthMiddleware(t *testing.T) {
	mockService := &mockJWTService{}
	authMiddleware := NewAuthMiddleware(mockService)

	if authMiddleware == nil {
		t.Error("NewAuthMiddleware should not return nil")
	}

	if authMiddleware.jwtService != mockService {
		t.Error("JWT service not properly assigned")
	}
}
