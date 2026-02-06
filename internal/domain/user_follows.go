package domain

type UserFollow struct {
	FollowerUsername string `db:"follower_username"`
	FollowedUsername string `db:"followed_username"`
}

type JsonUserFollow struct {
	FollowerUsername string `json:"follower_username"`
	FollowedUsername string `json:"followed_username"`
}
