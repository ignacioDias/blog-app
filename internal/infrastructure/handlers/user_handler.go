package handlers

import (
	"log"
	"net/http"
	"postapi/internal/application"
	models "postapi/internal/domain"
	"postapi/internal/middleware"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserUseCase application.UserUseCase
	JWTService  application.JWTService
}

func (uh *UserHandler) RegisterUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := models.UserResponse{}
		err := middleware.Parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot middleware.Parse body. err = %v \n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		u := &models.User{
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
		}

		err = uh.UserUseCase.UserRepo.Create(u)
		if err != nil {
			log.Printf("Cannot save user in DB. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		resp := application.MapUserToJson(u)
		middleware.SendResponse(w, r, resp, http.StatusOK)
	}
}

func (uh *UserHandler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := models.UserResponse{}
		err := middleware.Parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot middleware.Parse body. err = %v \n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}

		u := &models.User{
			Username: req.Username,
			Password: req.Password,
		}

		user, err := uh.UserUseCase.UserRepo.LoginUser(u)
		if err != nil {
			log.Printf("Login failed. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Invalid credentials"}, http.StatusUnauthorized)
			return
		}
		tokenString, err := uh.JWTService.GenerateToken(user.Username)
		if err != nil {
			log.Printf("Cannot create token. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to generate token"}, http.StatusInternalServerError)
			return
		}

		resp := map[string]any{
			"user":  application.MapUserToJson(user),
			"token": tokenString,
		}
		middleware.SendResponse(w, r, resp, http.StatusOK)
	}
}

func (uh *UserHandler) GetUserByUsernameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		if username == "" {
			middleware.SendResponse(w, r, map[string]string{"error": "Username required"}, http.StatusBadRequest)
			return
		}
		user, err := uh.UserUseCase.UserRepo.FindByUsername(username)
		if err != nil {
			log.Printf("Cannot get user, err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to get user"}, http.StatusInternalServerError)
			return
		}
		resp := application.MapUserToJson(user)
		middleware.SendResponse(w, r, resp, http.StatusOK)

	}
}
