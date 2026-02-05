package app

import (
	"encoding/json"
	"log"
	"net/http"
	"postapi/app/models"
)

func parse(_ http.ResponseWriter, r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func sendResponse(w http.ResponseWriter, _ *http.Request, data any, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Cannot format json. err = %v\n", err)
	}
}

func mapPostToJson(p *models.Post) models.JsonPost {
	return models.JsonPost{
		ID:      p.ID,
		Author:  p.Author,
		Content: p.Content,
		Title:   p.Title,
	}
}

func mapUserToJson(u *models.User) models.JsonUser {
	return models.JsonUser{
		Username: u.Username,
		Email:    u.Email,
	}
}

func mapFollowToJson(f *models.UserFollow) models.JsonUserFollow {
	return models.JsonUserFollow{
		FollowerUsername: f.FollowerUsername,
		FollowedUsername: f.FollowedUsername,
	}
}

func mapProfileToJson(f *models.Profile) models.JsonProfile {
	return models.JsonProfile{
		Username:       f.Username,
		Description:    f.Description,
		ProfilePicture: f.ProfilePicture,
	}
}
