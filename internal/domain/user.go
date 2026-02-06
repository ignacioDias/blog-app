package domain

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

type JsonUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
	}
}
