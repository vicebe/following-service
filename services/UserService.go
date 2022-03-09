package services

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vicebe/following-service/data"
)

// UserService is a service that handles common business logic for users.
type UserService struct {
	l  *log.Logger
	db *sqlx.DB
}

func NewUserService(l *log.Logger, db *sqlx.DB) *UserService {
	return &UserService{l, db}
}

// FollowUser adds user to the followers of another user given both ids.
func (us *UserService) FollowUser(userId string, userToFollowId string) error {
	us.l.Printf(
		"[DEBUG] starting follow process for user (%s -> %s)\n",
		userId,
		userToFollowId,
	)
	user, err := data.GetUserByID(userId)

	if err != nil {
		return err
	}

	err = user.Follow(userToFollowId)

	if err != nil {
		return err
	}

	return nil
}

// GetFollowers returns a list of a user's followers
func (us *UserService) GetFollowers(userId string) ([]string, error) {
	us.l.Printf("[DEBUG] Finding user %s\n", userId)

	user, err := data.GetUserByID(userId)

	if err != nil {
		return nil, err
	}

	followers := user.GetFollowers()

	return followers, nil

}
