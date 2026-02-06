package domain

type Profile struct {
	Username       string `db:"username"`
	Description    string `db:"description"`
	ProfilePicture string `db:"profile_picture"`
}

type JsonProfile struct {
	Username       string `json:"username"`
	Description    string `json:"description"`
	ProfilePicture string `json:"profile_picture"`
}

type ProfileRequest struct {
	Username       string `json:"username"`
	Description    string `json:"description"`
	ProfilePicture string `json:"profile_picture"`
}
