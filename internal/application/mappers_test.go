package application

import (
	"postapi/internal/domain"
	"testing"
)

func TestMapUserToJson(t *testing.T) {
	tests := []struct {
		name string
		user *domain.User
		want domain.JsonUser
	}{
		{
			name: "Valid user",
			user: &domain.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpassword123",
			},
			want: domain.JsonUser{
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
		{
			name: "User with empty fields",
			user: &domain.User{
				Username: "",
				Email:    "",
				Password: "password",
			},
			want: domain.JsonUser{
				Username: "",
				Email:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapUserToJson(tt.user)
			if got.Username != tt.want.Username {
				t.Errorf("MapUserToJson() Username = %v, want %v", got.Username, tt.want.Username)
			}
			if got.Email != tt.want.Email {
				t.Errorf("MapUserToJson() Email = %v, want %v", got.Email, tt.want.Email)
			}
		})
	}
}

func TestMapPostToJson(t *testing.T) {
	tests := []struct {
		name string
		post *domain.Post
		want domain.JsonPost
	}{
		{
			name: "Valid post",
			post: &domain.Post{
				ID:      1,
				Title:   "Test Post",
				Content: "This is test content",
				Author:  "testuser",
			},
			want: domain.JsonPost{
				ID:      1,
				Title:   "Test Post",
				Content: "This is test content",
				Author:  "testuser",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapPostToJson(tt.post)
			if got.ID != tt.want.ID {
				t.Errorf("MapPostToJson() ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Title != tt.want.Title {
				t.Errorf("MapPostToJson() Title = %v, want %v", got.Title, tt.want.Title)
			}
			if got.Content != tt.want.Content {
				t.Errorf("MapPostToJson() Content = %v, want %v", got.Content, tt.want.Content)
			}
			if got.Author != tt.want.Author {
				t.Errorf("MapPostToJson() Author = %v, want %v", got.Author, tt.want.Author)
			}
		})
	}
}

func TestMapFollowToJson(t *testing.T) {
	follow := &domain.UserFollow{
		FollowerUsername: "user1",
		FollowedUsername: "user2",
	}

	got := MapFollowToJson(follow)

	if got.FollowerUsername != "user1" {
		t.Errorf("MapFollowToJson() FollowerUsername = %v, want user1", got.FollowerUsername)
	}
	if got.FollowedUsername != "user2" {
		t.Errorf("MapFollowToJson() FollowedUsername = %v, want user2", got.FollowedUsername)
	}
}

func TestMapProfileToJson(t *testing.T) {
	profile := &domain.Profile{
		Username:       "testuser",
		Description:    "Test description",
		ProfilePicture: "https://example.com/pic.jpg",
	}

	got := MapProfileToJson(profile)

	if got.Username != "testuser" {
		t.Errorf("MapProfileToJson() Username = %v, want testuser", got.Username)
	}
	if got.Description != "Test description" {
		t.Errorf("MapProfileToJson() Description = %v, want Test description", got.Description)
	}
	if got.ProfilePicture != "https://example.com/pic.jpg" {
		t.Errorf("MapProfileToJson() ProfilePicture = %v, want URL", got.ProfilePicture)
	}
}
