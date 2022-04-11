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

	CommunityOne = data.Community{
		ID:         1,
		ExternalID: "COMMUNITY-ONE",
	}

	CommunitiesList = []data.Community{
		CommunityOne,
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

type UserServiceGetCommunitiesMock struct {
	services.UserServiceI
}

func (UserServiceGetCommunitiesMock) GetUserCommunities(
	*data.User,
) ([]data.Community, error) {
	return CommunitiesList, nil
}

type UserServiceGetCommunitiesErrorMock struct {
	services.UserServiceI
}

func (UserServiceGetCommunitiesErrorMock) GetUserCommunities(
	*data.User,
) ([]data.Community, error) {
	return nil, errors.New("error")
}

type UserServiceUnfollowUserMock struct {
	services.UserServiceI
}

func (UserServiceUnfollowUserMock) GetUser(string) (*data.User, error) {
	return &UserOne, nil
}

func (UserServiceUnfollowUserMock) UnfollowUser(*data.User, *data.User) error {
	return nil
}

type UserServiceUnfollowUserNotFoundMock struct {
	services.UserServiceI
}

func (UserServiceUnfollowUserNotFoundMock) GetUser(string) (*data.User, error) {
	return nil, data.ErrorUserNotFound
}

func (UserServiceUnfollowUserNotFoundMock) UnfollowUser(
	*data.User,
	*data.User,
) error {
	return nil
}

type UserServiceUnfollowUserGetUserErrorMock struct {
	services.UserServiceI
}

func (UserServiceUnfollowUserGetUserErrorMock) GetUser(string) (*data.User, error) {
	return nil, errors.New("error")
}

func (UserServiceUnfollowUserGetUserErrorMock) UnfollowUser(
	*data.User,
	*data.User,
) error {
	return nil
}

type UserServiceUnfollowUserErrorMock struct {
	services.UserServiceI
}

func (UserServiceUnfollowUserErrorMock) GetUser(string) (*data.User, error) {
	return &UserOne, nil
}

func (UserServiceUnfollowUserErrorMock) UnfollowUser(
	*data.User,
	*data.User,
) error {
	return errors.New("error")
}

type UserServiceFollowUserMock struct {
	services.UserServiceI
}

func (UserServiceFollowUserMock) GetUser(string) (*data.User, error) {
	return &UserOne, nil
}

func (UserServiceFollowUserMock) FollowUser(*data.User, *data.User) error {
	return nil
}

type UserServiceFollowUserNotFoundMock struct {
	services.UserServiceI
}

func (UserServiceFollowUserNotFoundMock) GetUser(string) (*data.User, error) {
	return nil, data.ErrorUserNotFound
}

func (UserServiceFollowUserNotFoundMock) FollowUser(
	*data.User,
	*data.User,
) error {
	return nil
}

type UserServiceFollowUserGetUserErrorMock struct {
	services.UserServiceI
}

func (UserServiceFollowUserGetUserErrorMock) GetUser(string) (*data.User, error) {
	return nil, errors.New("error")
}

func (UserServiceFollowUserGetUserErrorMock) FollowUser(
	*data.User,
	*data.User,
) error {
	return nil
}

type UserServiceFollowUserErrorMock struct {
	services.UserServiceI
}

func (UserServiceFollowUserErrorMock) GetUser(string) (*data.User, error) {
	return &UserOne, nil
}

func (UserServiceFollowUserErrorMock) FollowUser(
	*data.User,
	*data.User,
) error {
	return errors.New("error")
}
