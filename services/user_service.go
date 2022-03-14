package services

import (
	"log"

	"github.com/vicebe/following-service/data"
)

// UserService is a service that handles common business logic for users.
type UserService struct {
	l  *log.Logger
	ds *data.Store
}

func NewUserService(l *log.Logger, ds *data.Store) *UserService {
	return &UserService{l, ds}
}

// FollowUser adds user to the followers of another user given both ids.
func (us *UserService) FollowUser(userId string, userToFollowId string) error {
	us.l.Printf(
		"[DEBUG] starting follow process for user (%s -> %s)\n",
		userId,
		userToFollowId,
	)

	err := us.ds.Follow(userId, userToFollowId)

	if err != nil {
		return err
	}

	return nil
}

// GetFollowers returns a list of a user's followers
func (us *UserService) GetFollowers(userId string) ([]string, error) {
	us.l.Printf("[DEBUG] Finding user %s\n", userId)

	followers, err := us.ds.GetFollowers(userId)

	if err != nil {
		return nil, err
	}

	return followers, nil

}
