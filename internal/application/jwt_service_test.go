package application

import (
	"strings"
	"testing"
	"time"
)

func TestJWTService_GenerateToken(t *testing.T) {
	secretKey := "test-secret-key"
	jwtService := NewJWTService(secretKey)

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "Valid username",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "Empty username",
			username: "",
			wantErr:  false,
		},
		{
			name:     "Username with special characters",
			username: "test@user.com",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwtService.GenerateToken(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateToken() returned empty token")
			}
			if !tt.wantErr && len(strings.Split(token, ".")) != 3 {
				t.Error("GenerateToken() returned invalid JWT format")
			}
		})
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	secretKey := "test-secret-key"
	jwtService := NewJWTService(secretKey)
	username := "testuser"

	validToken, err := jwtService.GenerateToken(username)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		wantUser  string
		wantErr   bool
		setupFunc func() string
	}{
		{
			name:     "Valid token",
			token:    validToken,
			wantUser: username,
			wantErr:  false,
		},
		{
			name:     "Empty token",
			token:    "",
			wantUser: "",
			wantErr:  true,
		},
		{
			name:     "Invalid token format",
			token:    "invalid.token.format",
			wantUser: "",
			wantErr:  true,
		},
		{
			name:     "Token with wrong secret",
			wantUser: "",
			wantErr:  true,
			setupFunc: func() string {
				wrongService := NewJWTService("wrong-secret")
				token, _ := wrongService.GenerateToken(username)
				return token
			},
		},
		{
			name:     "Malformed token",
			token:    "not.a.token",
			wantUser: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.token
			if tt.setupFunc != nil {
				token = tt.setupFunc()
			}

			gotUser, err := jwtService.ValidateToken(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUser != tt.wantUser {
				t.Errorf("ValidateToken() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestJWTService_RoundTrip(t *testing.T) {
	secretKey := "test-secret-key-roundtrip"
	jwtService := NewJWTService(secretKey)

	usernames := []string{
		"user1",
		"user2",
		"admin",
		"test@example.com",
	}

	for _, username := range usernames {
		t.Run(username, func(t *testing.T) {
			token, err := jwtService.GenerateToken(username)
			if err != nil {
				t.Fatalf("GenerateToken() failed: %v", err)
			}

			time.Sleep(10 * time.Millisecond)

			gotUsername, err := jwtService.ValidateToken(token)
			if err != nil {
				t.Fatalf("ValidateToken() failed: %v", err)
			}

			if gotUsername != username {
				t.Errorf("Round trip failed: got %v, want %v", gotUsername, username)
			}
		})
	}
}
