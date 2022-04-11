package test_utils

import (
	"errors"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/services"
)

type FoundCommunityMock struct {
	services.CommunityServiceI
}

func (FoundCommunityMock) GetCommunity(_ string) (*data.Community, error) {
	return &data.Community{
		ID:         1,
		ExternalID: "1",
	}, nil
}

type NotFoundCommunityMock struct {
	services.CommunityServiceI
}

func (NotFoundCommunityMock) GetCommunity(_ string) (*data.Community, error) {
	return nil, data.ErrorCommunityNotFound
}

type CommunityServiceErrorMock struct {
	services.CommunityServiceI
}

func (CommunityServiceErrorMock) GetCommunity(_ string) (*data.Community, error) {
	return nil, errors.New("error")
}
