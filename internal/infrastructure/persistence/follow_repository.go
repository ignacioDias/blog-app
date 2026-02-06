package persistence

import (
	models "postapi/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserFollowRepositoryImpl struct {
	db *sqlx.DB
}

func (u *UserFollowRepositoryImpl) Create(follow *models.UserFollow) error {
	_, err := u.db.Exec(insertFollowSchema, follow.FollowerUsername, follow.FollowedUsername)
	return err
}

func (u *UserFollowRepositoryImpl) Delete(follow *models.UserFollow) error {
	_, err := u.db.Exec(removeFollowSchema, follow.FollowerUsername, follow.FollowedUsername)
	return err
}

func (u *UserFollowRepositoryImpl) GetFollowers(username string) ([]string, error) {
	var followers []string
	err := u.db.Select(&followers, getFollowersSchema, username)

	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (u *UserFollowRepositoryImpl) GetFollowing(username string) ([]string, error) {
	var following []string
	err := u.db.Select(&following, getFollowingSchema, username)

	if err != nil {
		return nil, err
	}

	return following, nil
}
