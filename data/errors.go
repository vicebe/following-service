package data

import "errors"

var (
	ErrorUserNotFound             = errors.New("user not found")
	ErrorUserAlreadyFollowed      = errors.New("user being followed")
	ErrorCommunityNotFound        = errors.New("community not found")
	ErrorCommunityAlreadyFollowed = errors.New("community being followed")
	ErrorCouldntConnectDb         = errors.New("could not connect to database")
	ErrorOpeningDb                = errors.New("could not open database")
)
