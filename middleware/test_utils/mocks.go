package test_utils

import (
	"errors"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
)

type FoundUserMock struct {
	services.UserServiceI
}

func (FoundUserMock) GetUser(userID string) (*data.User, error) {
	return &data.User{
		ID:         1,
		ExternalID: "1",
	}, nil
}

type NotFoundUserMock struct {
	services.UserServiceI
}

func (NotFoundUserMock) GetUser(userID string) (*data.User, error) {
	return nil, data.ErrorUserNotFound
}

type UserServiceErrorMock struct {
	services.UserServiceI
}

func (UserServiceErrorMock) GetUser(userID string) (*data.User, error) {
	return nil, errors.New("error")
}
