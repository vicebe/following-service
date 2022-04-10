package data

import "errors"

var (
	ErrorUserNotFound      = errors.New("user not found")
	ErrorCommunityNotFound = errors.New("community not found")
)
