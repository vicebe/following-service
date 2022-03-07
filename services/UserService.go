package services

import (
	"log"

	"github.com/vicebe/following-service/data"
)

type UserService struct {
	l *log.Logger
}

func NewUserService(l *log.Logger) *UserService {
	return &UserService{l}
}

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

func (us *UserService) GetFollowers(userId string) ([]string, error) {
	us.l.Printf("[DEBUG] Finding user %s\n", userId)

	user, err := data.GetUserByID(userId)

	if err != nil {
		return nil, err
	}

	followers := user.GetFollowers()

	return followers, nil

}
