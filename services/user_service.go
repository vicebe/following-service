package services

import (
	"encoding/json"
	"log"

	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/events"
	userproducers "github.com/vicebe/following-service/events/user_producers"
)

type UserServiceI interface {
	GetUser(userID string) (*data.User, error)
	FollowUser(user *data.User, followee *data.User) error
	UnfollowUser(user *data.User, followee *data.User) error
	GetUserFollowers(user *data.User) ([]data.User, error)
	GetUserCommunities(user *data.User) ([]data.Community, error)
	CreateUser(user *data.User) error
}

// UserService is a service that handles common business logic for users.
type UserService struct {
	l                  *log.Logger
	ur                 data.UserRepository
	UserFollowedProd   events.Producer
	UserUnfollowedProd events.Producer
}

func NewUserService(
	l *log.Logger,
	ur data.UserRepository,
	userFollowedProd events.Producer,
	userUnfollowedProd events.Producer,
) *UserService {
	return &UserService{l, ur, userFollowedProd, userUnfollowedProd}
}

func (us *UserService) GetUser(userID string) (*data.User, error) {
	user, err := us.ur.FindBy("external_id", userID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FollowUser adds user to the followers of another user given both ids.
func (us *UserService) FollowUser(user *data.User, followee *data.User) error {

	err := us.ur.FollowUser(user, followee)

	if err != nil {
		return err
	}

	m, err := json.Marshal(&userproducers.UserFollowedEvent{
		FolloweeID: followee.ExternalID,
		FollowerID: user.ExternalID,
	})

	if err != nil {
		us.l.Print("[ERROR]: ", err)
		return err
	}

	if err := us.UserFollowedProd.ProduceEvent(m); err != nil {
		us.l.Print("[ERROR]: ", err)
		return err
	}

	return nil
}

// UnfollowUser removes user from the followers of another user given both ids.
func (us *UserService) UnfollowUser(
	user *data.User, followee *data.User,
) error {

	err := us.ur.UnfollowUser(user, followee)

	if err != nil {
		return err
	}

	m, err := json.Marshal(&userproducers.UserUnfollowedEvent{
		FolloweeID: followee.ExternalID,
		FollowerID: user.ExternalID,
	})

	if err != nil {
		us.l.Print("[ERROR]: ", err)
		return err
	}

	if err := us.UserUnfollowedProd.ProduceEvent(m); err != nil {
		us.l.Print("[ERROR]: ", err)
		return err
	}

	return nil
}

// GetUserFollowers returns a list of a user's followers
func (us *UserService) GetUserFollowers(user *data.User) ([]data.User, error) {

	followers, err := us.ur.GetUserFollowers(user)

	if err != nil {
		return nil, err
	}

	return followers, nil

}

// GetUserCommunities returns a list of communities that the user follows
func (us *UserService) GetUserCommunities(
	user *data.User,
) ([]data.Community, error) {

	followers, err := us.ur.GetUserCommunities(user)

	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (us *UserService) CreateUser(user *data.User) error {

	_, err := us.ur.FindBy("external_id", user.ExternalID)

	if err == data.ErrorUserNotFound {
		if err := us.ur.Create(user); err != nil {
			return err
		}

		return nil
	}

	return err
}
