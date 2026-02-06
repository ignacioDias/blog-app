package handlers

import (
	"log"
	"net/http"
	"postapi/internal/application"
	models "postapi/internal/domain"
	"postapi/internal/middleware"

	"github.com/gorilla/mux"
)

type ProfileHandler struct {
	ProfileUseCase application.ProfileUseCase
}

func (p *ProfileHandler) GetProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		if username == "" {
			middleware.SendResponse(w, r, map[string]string{"error": "Username required"}, http.StatusBadRequest)
			return
		}

		profile, err := p.ProfileUseCase.ProfileRepository.FindByUsername(username)

		if err != nil {
			log.Printf("Cannot get profile from DB. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to get profile details"}, http.StatusInternalServerError)
			return
		}

		jsonProfile := application.MapProfileToJson(profile)
		middleware.SendResponse(w, r, jsonProfile, http.StatusOK)
	}
}

func (pH *ProfileHandler) CreateProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(middleware.UsernameKey).(string)
		if !ok {
			middleware.SendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		req := models.ProfileRequest{}
		err := middleware.Parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot middleware.Parse body. err = %v \n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}

		p := &models.Profile{
			Username:       username,
			Description:    req.Description,
			ProfilePicture: req.ProfilePicture,
		}

		err = pH.ProfileUseCase.ProfileRepository.Create(p)

		if err != nil {
			log.Printf("Cannot create profile. err = %v \n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to create profile"}, http.StatusInternalServerError)
			return
		}

		jsonProfile := application.MapProfileToJson(p)
		middleware.SendResponse(w, r, jsonProfile, http.StatusOK)
	}
}

func (pH *ProfileHandler) UpdateProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(middleware.UsernameKey).(string)
		if !ok {
			middleware.SendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		req := models.ProfileRequest{}
		err := middleware.Parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot middleware.Parse body. err = %v \n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		currentProfile, err := pH.ProfileUseCase.ProfileRepository.FindByUsername(username)

		if err != nil {
			log.Printf("Cannot get profile. err = %v \n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to get profile"}, http.StatusInternalServerError)
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

		err = pH.ProfileUseCase.ProfileRepository.Update(p)

		if err != nil {
			log.Printf("Cannot update profile. err = %v \n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to update profile"}, http.StatusInternalServerError)
			return
		}

		jsonProfile := application.MapProfileToJson(p)
		middleware.SendResponse(w, r, jsonProfile, http.StatusOK)
	}
}
