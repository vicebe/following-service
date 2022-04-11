package test_utils

import (
	"errors"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
)

type CommunityServiceFollowCommunityMock struct {
	services.CommunityServiceI
}

func (CommunityServiceFollowCommunityMock) FollowCommunity(
	*data.Community,
	*data.User,
) error {
	return nil
}

type CommunityServiceFollowCommunityErrorMock struct {
	services.CommunityServiceI
}

func (CommunityServiceFollowCommunityErrorMock) FollowCommunity(
	*data.Community,
	*data.User,
) error {
	return errors.New("error")
}

type CommunityServiceUnfollowCommunityMock struct {
	services.CommunityServiceI
}

func (CommunityServiceUnfollowCommunityMock) UnfollowCommunity(
	*data.Community,
	*data.User,
) error {
	return nil
}

type CommunityServiceUnfollowCommunityErrorMock struct {
	services.CommunityServiceI
}

func (CommunityServiceUnfollowCommunityErrorMock) UnfollowCommunity(
	*data.Community,
	*data.User,
) error {
	return errors.New("error")
}

type CommunityServiceGetFollowersMock struct {
	services.CommunityServiceI
}

func (CommunityServiceGetFollowersMock) GetCommunityFollowers(
	*data.Community,
) ([]data.User, error) {
	return FollowersList, nil
}

type CommunityServiceGetFollowersErrorMock struct {
	services.CommunityServiceI
}

func (CommunityServiceGetFollowersErrorMock) GetCommunityFollowers(
	*data.Community,
) ([]data.User, error) {
	return nil, errors.New("error")
}
