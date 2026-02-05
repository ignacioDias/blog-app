package app

import (
	"fmt"
	"log"
	"net/http"
	"postapi/app/models"
	"strconv"

	"github.com/gorilla/mux"
)

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
func (a *App) UpdatePostHandler() http.HandlerFunc {
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
		vars := mux.Vars(r)
		id := vars["post_id"]
		idAsNumber, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			sendResponse(w, r, map[string]string{"error": fmt.Sprintf("Invalid ID %s", id)}, http.StatusBadRequest)
			return
		}
		p := &models.Post{
			ID:      idAsNumber,
			Title:   req.Title,
			Author:  username,
			Content: req.Content,
		}

		err = a.DB.UpdatePost(p)
		if err != nil {
			log.Printf("Cannot save post in DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to update post"}, http.StatusInternalServerError)
			return
		}

		resp := mapPostToJson(p)
		sendResponse(w, r, resp, http.StatusOK)
	}
}
func (a *App) DeletePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(UsernameKey).(string)
		if !ok {
			sendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}
		vars := mux.Vars(r)
		id := vars["post_id"]
		idAsNumber, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			sendResponse(w, r, map[string]string{"error": fmt.Sprintf("Invalid ID %s", id)}, http.StatusBadRequest)
			return
		}

		err = a.DB.DeletePost(idAsNumber, username)
		if err != nil {
			log.Printf("Cannot delete post in DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to delete post"}, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (a *App) GetPostsByUserHandler() http.HandlerFunc { //TODO: PAGINATION
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

func (a *App) GetPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["post_id"]

		idAsNumber, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			sendResponse(w, r, map[string]string{"error": fmt.Sprintf("Invalid ID %s", id)}, http.StatusBadRequest)
			return
		}

		post, err := a.DB.GetPost(idAsNumber)
		if err != nil {
			log.Printf("Cannot get post, err = %v\n", err)
			sendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}

		resp := mapPostToJson(post)
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

		resp := map[string]any{
			"user":  mapUserToJson(user),
			"token": tokenString,
		}
		sendResponse(w, r, resp, http.StatusOK)
	}
}

func (a *App) GetProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		if username == "" {
			sendResponse(w, r, map[string]string{"error": "Username required"}, http.StatusBadRequest)
			return
		}
		user, err := a.DB.GetUserByUsername(username)
		if err != nil {
			log.Printf("Cannot get user, err = %v\n", err)
			sendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		resp := mapUserToJson(user)
		sendResponse(w, r, resp, http.StatusOK)

	}
}

func (a *App) FollowHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(UsernameKey).(string)
		if !ok {
			sendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
		}
		vars := mux.Vars(r)
		followed := vars["username"]
		f := &models.UserFollow{
			FollowerUsername: username,
			FollowedUsername: followed,
		}

		err := a.DB.CreateFollow(f)
		if err != nil {
			log.Printf("Cannot save post in DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to create post"}, http.StatusInternalServerError)
			return
		}

		resp := mapFollowToJson(f)
		sendResponse(w, r, resp, http.StatusOK)
	}
}

func (a *App) UnfollowHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (a *App) GetFollowersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (a *App) GetFollowingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
