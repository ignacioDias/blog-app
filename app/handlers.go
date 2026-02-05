package app

import (
	"fmt"
	"log"
	"net/http"
	"postapi/app/models"

	"github.com/gorilla/mux"
)

func (a *App) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Post API")
	}
}

func (a *App) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(UsernameKey).(string)
		if !ok {
			sendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		req := models.PostRequest{}
		err := parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse body. err = %v \n", err)
			sendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}

		p := &models.Post{
			ID:      0,
			Title:   req.Title,
			Author:  username,
			Content: req.Content,
		}

		err = a.DB.CreatePost(p)
		if err != nil {
			log.Printf("Cannot save post in DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to create post"}, http.StatusInternalServerError)
			return
		}

		resp := mapPostToJson(p)
		sendResponse(w, r, resp, http.StatusOK)
	}
}

func (a *App) GetPostsByUserHandlder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		if username == "" {
			sendResponse(w, r, map[string]string{"error": "Username required"}, http.StatusBadRequest)
			return
		}
		posts, err := a.DB.GetPostsByUser(username)
		if err != nil {
			log.Printf("Cannot get posts, err = %v\n", err)
			sendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		var resp = make([]models.JsonPost, len(posts))
		for idx, post := range posts {
			resp[idx] = mapPostToJson(post)
		}
		sendResponse(w, r, resp, http.StatusOK)
	}
}

func (a *App) RegisterUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := models.UserRequest{}
		err := parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse body. err = %v \n", err)
			sendResponse(w, r, nil, http.StatusBadRequest)
			return
		}
		u := &models.User{
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
		}

		err = a.DB.RegisterUser(u)
		if err != nil {
			log.Printf("Cannot save user in DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		resp := mapUserToJson(u)
		sendResponse(w, r, resp, http.StatusOK)
	}
}

func (a *App) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := models.UserRequest{}
		err := parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse body. err = %v \n", err)
			sendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}

		u := &models.User{
			Username: req.Username,
			Password: req.Password,
		}

		user, err := a.DB.LoginUser(u)
		if err != nil {
			log.Printf("Login failed. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Invalid credentials"}, http.StatusUnauthorized)
			return
		}

		tokenString, err := CreateToken(user.Username)
		if err != nil {
			log.Printf("Cannot create token. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to generate token"}, http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"user":  mapUserToJson(user),
			"token": tokenString,
		}
		sendResponse(w, r, resp, http.StatusOK)
	}
}

func (a *App) GetProfileHandler() http.HandlerFunc {
	return nil
}

func (a *App) GetPostHandlder() http.HandlerFunc {
	return nil
}

func (a *App) UpdatePostHandler() http.HandlerFunc {
	return nil
}
func (a *App) DeletePostHandler() http.HandlerFunc {
	return nil
}
