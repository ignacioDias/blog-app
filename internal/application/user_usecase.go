package application

import (
	models "postapi/internal/domain"
)

type UserUseCase struct {
	UserRepo   models.UserRepository
	FollowRepo models.UserFollowRepository
}

func MapUserToJson(u *models.User) models.JsonUser {
	return models.JsonUser{
		Username: u.Username,
		Email:    u.Email,
	}
}

func MapFollowToJson(f *models.UserFollow) models.JsonUserFollow {
	return models.JsonUserFollow{
		FollowerUsername: f.FollowerUsername,
		FollowedUsername: f.FollowedUsername,
	}
}
