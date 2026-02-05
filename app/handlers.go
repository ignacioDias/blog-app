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

		content := req.Content
		title := req.Title

		oldPost, err := a.DB.GetPost(idAsNumber)
		if err != nil {
			log.Printf("Cannot get post, err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to get post"}, http.StatusInternalServerError)
			return
		}
		if req.Content == "" {
			content = oldPost.Content
		}
		if req.Title == "" {
			title = oldPost.Title
		}
		p := &models.Post{
			ID:      idAsNumber,
			Title:   title,
			Author:  username,
			Content: content,
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
			sendResponse(w, r, map[string]string{"error": "Failed to get posts"}, http.StatusInternalServerError)
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
			sendResponse(w, r, map[string]string{"error": "Failed to get post"}, http.StatusInternalServerError)
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
			sendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
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

func (a *App) GetUserByUsernameHandler() http.HandlerFunc {
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
			sendResponse(w, r, map[string]string{"error": "Failed to get user"}, http.StatusInternalServerError)
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
			return
		}
		vars := mux.Vars(r)
		followed := vars["username"]
		f := &models.UserFollow{
			FollowerUsername: username,
			FollowedUsername: followed,
		}

		err := a.DB.CreateFollow(f)
		if err != nil {
			log.Printf("Cannot save follow in DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to create follow"}, http.StatusInternalServerError)
			return
		}

		resp := mapFollowToJson(f)
		sendResponse(w, r, resp, http.StatusOK)
	}
}

func (a *App) UnfollowHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(UsernameKey).(string)
		if !ok {
			sendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		unfollowed := vars["username"]

		f := &models.UserFollow{
			FollowerUsername: username,
			FollowedUsername: unfollowed,
		}

		err := a.DB.Removefollow(f)

		if err != nil {
			log.Printf("Cannot remove follow in DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to remove follow"}, http.StatusInternalServerError)
			return
		}

		resp := mapFollowToJson(f)
		sendResponse(w, r, resp, http.StatusOK)
	}
}

func (a *App) GetFollowersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		if username == "" {
			sendResponse(w, r, map[string]string{"error": "Username required"}, http.StatusBadRequest)
			return
		}
		usernames, err := a.DB.GetFollowers(username)

		if err != nil {
			log.Printf("Cannot get followers from DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to get followers"}, http.StatusInternalServerError)
			return
		}
		users := make([]models.JsonUser, len(usernames))
		for i, followerUsername := range usernames {
			user, err := a.DB.GetUserByUsername(followerUsername)

			if err != nil {
				log.Printf("Cannot get follower from DB. err = %v\n", err)
				sendResponse(w, r, map[string]string{"error": "Failed to get follower details"}, http.StatusInternalServerError)
				return
			}

			users[i] = mapUserToJson(user)

		}

		sendResponse(w, r, users, http.StatusOK)
	}
}

func (a *App) GetFollowingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		if username == "" {
			sendResponse(w, r, map[string]string{"error": "Username required"}, http.StatusBadRequest)
			return
		}
		usernames, err := a.DB.GetFollowings(username)

		if err != nil {
			log.Printf("Cannot get followings from DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to get followings"}, http.StatusInternalServerError)
			return
		}
		users := make([]models.JsonUser, len(usernames))
		for i, followingUsername := range usernames {
			user, err := a.DB.GetUserByUsername(followingUsername)

			if err != nil {
				log.Printf("Cannot get following from DB. err = %v\n", err)
				sendResponse(w, r, map[string]string{"error": "Failed to get following details"}, http.StatusInternalServerError)
				return
			}

			users[i] = mapUserToJson(user)

		}

		sendResponse(w, r, users, http.StatusOK)
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

		profile, err := a.DB.GetProfile(username)

		if err != nil {
			log.Printf("Cannot get profile from DB. err = %v\n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to get profile details"}, http.StatusInternalServerError)
			return
		}

		sendResponse(w, r, mapProfileToJson(profile), http.StatusOK)
	}
}

func (a *App) CreateProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(UsernameKey).(string)
		if !ok {
			sendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		req := models.ProfileRequest{}
		err := parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse body. err = %v \n", err)
			sendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}

		p := &models.Profile{
			Username:       username,
			Description:    req.Description,
			ProfilePicture: req.ProfilePicture,
		}

		err = a.DB.CreateProfile(p)

		if err != nil {
			log.Printf("Cannot create profile. err = %v \n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to create profile"}, http.StatusInternalServerError)
			return
		}

		jsonProfile := mapProfileToJson(p)
		sendResponse(w, r, jsonProfile, http.StatusOK)
	}
}

func (a *App) UpdateProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(UsernameKey).(string)
		if !ok {
			sendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		req := models.ProfileRequest{}
		err := parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse body. err = %v \n", err)
			sendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		currentProfile, err := a.DB.GetProfile(username)

		if err != nil {
			log.Printf("Cannot get profile. err = %v \n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to get profile"}, http.StatusInternalServerError)
			return
		}
		description := req.Description
		profilePicture := req.ProfilePicture
		if req.Description == "" {
			description = currentProfile.Description
		}
		if req.ProfilePicture == "" {
			profilePicture = currentProfile.ProfilePicture
		}
		p := &models.Profile{
			Username:       username,
			Description:    description,
			ProfilePicture: profilePicture,
		}

		err = a.DB.UpdateProfile(p)

		if err != nil {
			log.Printf("Cannot update profile. err = %v \n", err)
			sendResponse(w, r, map[string]string{"error": "Failed to update profile"}, http.StatusInternalServerError)
			return
		}

		jsonProfile := mapProfileToJson(p)
		sendResponse(w, r, jsonProfile, http.StatusOK)
	}
}
