package test

import (
	"postapi/app"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateToken(t *testing.T) {
	username := "testuser"

	token, err := app.CreateToken(username)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	// Verify the token is valid
	verifiedUsername, err := app.VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}

	if verifiedUsername != username {
		t.Errorf("expected username %s, got %s", username, verifiedUsername)
	}
}

func TestVerifyToken_ValidToken(t *testing.T) {
	username := "testuser"
	token, _ := app.CreateToken(username)

	verifiedUsername, err := app.VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}

	if verifiedUsername != username {
		t.Errorf("expected username %s, got %s", username, verifiedUsername)
	}
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"

	_, err := app.VerifyToken(invalidToken)
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	// Create an expired token
	secretKey := []byte("secret-key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": "testuser",
			"exp":      time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		t.Fatalf("failed to create expired token: %v", err)
	}

	_, err = app.VerifyToken(tokenString)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestVerifyToken_MissingUsername(t *testing.T) {
	// Create a token without username claim
	secretKey := []byte("secret-key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	_, err = app.VerifyToken(tokenString)
	if err == nil {
		t.Error("expected error for token without username")
	}
}
