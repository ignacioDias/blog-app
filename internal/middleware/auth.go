package middleware

import (
	"context"
	"net/http"
	"postapi/internal/application"
)

type contextKey string

const UsernameKey contextKey = "username"

type AuthMiddleware struct {
	jwtService application.JWTService
}

func NewAuthMiddleware(jwtService application.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (a *AuthMiddleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			SendResponse(w, r, map[string]string{"error": "Missing authorization header"}, http.StatusUnauthorized)
			return
		}

		const bearerPrefix = "Bearer "
		if len(tokenString) < len(bearerPrefix) || tokenString[:len(bearerPrefix)] != bearerPrefix {
			SendResponse(w, r, map[string]string{"error": "Invalid authorization format"}, http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len(bearerPrefix):]
		username, err := a.jwtService.ValidateToken(tokenString)
		if err != nil {
			SendResponse(w, r, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
