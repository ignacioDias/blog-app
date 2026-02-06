package handlers

import (
	"log"
	"net/http"
	"postapi/internal/application"
	models "postapi/internal/domain"
	"postapi/internal/middleware"

	"github.com/gorilla/mux"
)

type FollowHandler struct {
	UserUseCase application.UserUseCase
}

func (fh *FollowHandler) FollowHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(middleware.UsernameKey).(string)
		if !ok {
			middleware.SendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}
		vars := mux.Vars(r)
		followed := vars["username"]
		f := &models.UserFollow{
			FollowerUsername: username,
			FollowedUsername: followed,
		}

		err := fh.UserUseCase.FollowRepo.Create(f)
		if err != nil {
			log.Printf("Cannot save follow in DB. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to create follow"}, http.StatusInternalServerError)
			return
		}

		resp := application.MapFollowToJson(f)
		middleware.SendResponse(w, r, resp, http.StatusOK)
	}
}

func (fh *FollowHandler) UnfollowHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(middleware.UsernameKey).(string)
		if !ok {
			middleware.SendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		unfollowed := vars["username"]

		f := &models.UserFollow{
			FollowerUsername: username,
			FollowedUsername: unfollowed,
		}

		err := fh.UserUseCase.FollowRepo.Delete(f)

		if err != nil {
			log.Printf("Cannot remove follow in DB. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to remove follow"}, http.StatusInternalServerError)
			return
		}

		resp := application.MapFollowToJson(f)
		middleware.SendResponse(w, r, resp, http.StatusOK)
	}
}

func (fh *FollowHandler) GetFollowersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		if username == "" {
			middleware.SendResponse(w, r, map[string]string{"error": "Username required"}, http.StatusBadRequest)
			return
		}
		usernames, err := fh.UserUseCase.FollowRepo.GetFollowers(username)

		if err != nil {
			log.Printf("Cannot get followers from DB. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to get followers"}, http.StatusInternalServerError)
			return
		}
		users := make([]models.JsonUser, len(usernames))
		for i, followerUsername := range usernames {
			user, err := fh.UserUseCase.UserRepo.FindByUsername(followerUsername)

			if err != nil {
				log.Printf("Cannot get follower from DB. err = %v\n", err)
				middleware.SendResponse(w, r, map[string]string{"error": "Failed to get follower details"}, http.StatusInternalServerError)
				return
			}

			users[i] = application.MapUserToJson(user)

		}

		middleware.SendResponse(w, r, users, http.StatusOK)
	}
}

func (fh *FollowHandler) GetFollowingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		if username == "" {
			middleware.SendResponse(w, r, map[string]string{"error": "Username required"}, http.StatusBadRequest)
			return
		}
		usernames, err := fh.UserUseCase.FollowRepo.GetFollowing(username)

		if err != nil {
			log.Printf("Cannot get followings from DB. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to get followings"}, http.StatusInternalServerError)
			return
		}
		users := make([]models.JsonUser, len(usernames))
		for i, followingUsername := range usernames {
			user, err := fh.UserUseCase.UserRepo.FindByUsername(followingUsername)

			if err != nil {
				log.Printf("Cannot get following from DB. err = %v\n", err)
				middleware.SendResponse(w, r, map[string]string{"error": "Failed to get following details"}, http.StatusInternalServerError)
				return
			}

			users[i] = application.MapUserToJson(user)

		}

		middleware.SendResponse(w, r, users, http.StatusOK)
	}
}
