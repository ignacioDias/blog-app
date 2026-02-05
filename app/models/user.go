package models

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

type JsonUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
