package services

import (
	"log"

	"github.com/vicebe/following-service/data"
)

// UserService is a service that handles common business logic for users.
type UserService struct {
	l  *log.Logger
	ur data.UserRepository
}

func NewUserService(l *log.Logger, ur data.UserRepository) *UserService {
	return &UserService{l, ur}
}

// FollowUser adds user to the followers of another user given both ids.
func (us *UserService) FollowUser(userId string, userToFollowId string) error {

	follower, err := us.ur.FindBy("external_id", userId)

	if err != nil {
		return err
	}

	followee, err := us.ur.FindBy("external_id", userToFollowId)

	if err != nil {
		return err
	}

	err = us.ur.FollowUser(follower, followee)

	if err != nil {
		return err
	}

	return nil
}

// UnfollowUser removes user from the followers of another user given both ids.
func (us *UserService) UnfollowUser(
	userId string, userToUnfollowId string,
) error {
	follower, err := us.ur.FindBy("external_id", userId)

	if err != nil {
		return err
	}

	followee, err := us.ur.FindBy("external_id", userToUnfollowId)

	if err != nil {
		return err
	}

	err = us.ur.UnfollowUser(follower, followee)

	if err != nil {
		return err
	}

	return nil
}

// GetUserFollowers returns a list of a user's followers
func (us *UserService) GetUserFollowers(userId string) ([]data.User, error) {
	user, err := us.ur.FindBy("external_id", userId)
	followers, err := us.ur.GetUserFollowers(user)

	if err != nil {
		return nil, err
	}

	return followers, nil

}
