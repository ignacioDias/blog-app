package app

import (
	"context"
	"net/http"
)

type contextKey string

const UsernameKey contextKey = "username"

func (a *App) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			sendResponse(w, r, map[string]string{"error": "Missing authorization header"}, http.StatusUnauthorized)
			return
		}

		const bearerPrefix = "Bearer "
		if len(tokenString) < len(bearerPrefix) || tokenString[:len(bearerPrefix)] != bearerPrefix {
			sendResponse(w, r, map[string]string{"error": "Invalid authorization format"}, http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len(bearerPrefix):]
		username, err := VerifyToken(tokenString)
		if err != nil {
			sendResponse(w, r, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
			return
		}

		// Agregar username al contexto
		ctx := context.WithValue(r.Context(), UsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
