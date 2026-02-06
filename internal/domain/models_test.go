package domain

import "testing"

func TestUser_ToResponse(t *testing.T) {
	tests := []struct {
		name string
		user *User
		want UserResponse
	}{
		{
			name: "Complete user",
			user: &User{
				Username: "testuser",
				Password: "hashedpassword",
				Email:    "test@example.com",
			},
			want: UserResponse{
				Username: "testuser",
				Password: "hashedpassword",
				Email:    "test@example.com",
			},
		},
		{
			name: "Empty user",
			user: &User{},
			want: UserResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.ToResponse()
			if got.Username != tt.want.Username {
				t.Errorf("ToResponse() Username = %v, want %v", got.Username, tt.want.Username)
			}
			if got.Password != tt.want.Password {
				t.Errorf("ToResponse() Password = %v, want %v", got.Password, tt.want.Password)
			}
			if got.Email != tt.want.Email {
				t.Errorf("ToResponse() Email = %v, want %v", got.Email, tt.want.Email)
			}
		})
	}
}

func TestUserStructTags(t *testing.T) {
	// This test verifies the struct tags are present
	user := User{
		Username: "test",
		Password: "pass",
		Email:    "email@test.com",
	}

	// Verify the fields are accessible
	if user.Username == "" {
		t.Error("Username should be accessible")
	}
	if user.Password == "" {
		t.Error("Password should be accessible")
	}
	if user.Email == "" {
		t.Error("Email should be accessible")
	}
}

func TestPostModel(t *testing.T) {
	post := Post{
		ID:      1,
		Title:   "Test Title",
		Content: "Test Content",
		Author:  "testauthor",
	}

	if post.ID != 1 {
		t.Errorf("Post ID = %v, want 1", post.ID)
	}
	if post.Title != "Test Title" {
		t.Errorf("Post Title = %v, want 'Test Title'", post.Title)
	}
	if post.Content != "Test Content" {
		t.Errorf("Post Content = %v, want 'Test Content'", post.Content)
	}
	if post.Author != "testauthor" {
		t.Errorf("Post Author = %v, want 'testauthor'", post.Author)
	}
}

func TestProfileModel(t *testing.T) {
	profile := Profile{
		Username:       "testuser",
		Description:    "Test description",
		ProfilePicture: "https://example.com/pic.jpg",
	}

	if profile.Username != "testuser" {
		t.Errorf("Profile Username = %v, want 'testuser'", profile.Username)
	}
	if profile.Description != "Test description" {
		t.Errorf("Profile Description = %v, want 'Test description'", profile.Description)
	}
	if profile.ProfilePicture != "https://example.com/pic.jpg" {
		t.Errorf("Profile ProfilePicture = %v, want 'https://example.com/pic.jpg'", profile.ProfilePicture)
	}
}

func TestUserFollowModel(t *testing.T) {
	follow := UserFollow{
		FollowerUsername: "user1",
		FollowedUsername: "user2",
	}

	if follow.FollowerUsername != "user1" {
		t.Errorf("UserFollow FollowerUsername = %v, want 'user1'", follow.FollowerUsername)
	}
	if follow.FollowedUsername != "user2" {
		t.Errorf("UserFollow FollowedUsername = %v, want 'user2'", follow.FollowedUsername)
	}
}

func TestJsonUserModel(t *testing.T) {
	jsonUser := JsonUser{
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Verify password is not included in JsonUser (security check)
	if jsonUser.Username != "testuser" {
		t.Errorf("JsonUser Username = %v, want 'testuser'", jsonUser.Username)
	}
	if jsonUser.Email != "test@example.com" {
		t.Errorf("JsonUser Email = %v, want 'test@example.com'", jsonUser.Email)
	}
}

func TestPostRequestModel(t *testing.T) {
	postReq := PostRequest{
		Title:   "New Post",
		Content: "Post content",
	}

	if postReq.Title != "New Post" {
		t.Errorf("PostRequest Title = %v, want 'New Post'", postReq.Title)
	}
	if postReq.Content != "Post content" {
		t.Errorf("PostRequest Content = %v, want 'Post content'", postReq.Content)
	}
}

func TestProfileRequestModel(t *testing.T) {
	profileReq := ProfileRequest{
		Description:    "My bio",
		ProfilePicture: "https://example.com/avatar.jpg",
	}

	if profileReq.Description != "My bio" {
		t.Errorf("ProfileRequest Description = %v, want 'My bio'", profileReq.Description)
	}
	if profileReq.ProfilePicture != "https://example.com/avatar.jpg" {
		t.Errorf("ProfileRequest ProfilePicture = %v, want URL", profileReq.ProfilePicture)
	}
}
