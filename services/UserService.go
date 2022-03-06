package services

import (
	"github.com/vicebe/following-service/data"
)

func FollowUser(userId string, userToFollowId string) error {
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

func GetFollowers(userId string) ([]string, error) {
	user, err := data.GetUserByID(userId)

	if err != nil {
		return nil, err
	}

	followers := user.GetFollowers()

	return followers, nil

}
