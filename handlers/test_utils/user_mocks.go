package test_utils

import (
	"errors"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
)

var (
	UserOne = data.User{
		ID:         1,
		ExternalID: "USER-ONE",
	}

	UserTwo = data.User{
		ID:         2,
		ExternalID: "USER-TWO",
	}

	FollowersList = []data.User{
		UserOne,
		UserTwo,
	}
)

type UserServiceGetFollowersMock struct {
	services.UserServiceI
}

func (UserServiceGetFollowersMock) GetUserFollowers(
	*data.User,
) ([]data.User, error) {
	return FollowersList, nil
}

type UserServiceGetFollowersErrorMock struct {
	services.UserServiceI
}

func (UserServiceGetFollowersErrorMock) GetUserFollowers(
	*data.User,
) ([]data.User, error) {
	return nil, errors.New("error")
}
