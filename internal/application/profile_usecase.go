package application

import (
	models "postapi/internal/domain"
)

type ProfileUseCase struct {
	ProfileRepository models.ProfileRepository
}

func MapProfileToJson(f *models.Profile) models.JsonProfile {
	return models.JsonProfile{
		Username:       f.Username,
		Description:    f.Description,
		ProfilePicture: f.ProfilePicture,
	}
}
