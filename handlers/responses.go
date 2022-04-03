package handlers

import "github.com/vicebe/following-service/data"

type SimpleResponse struct {
	Message string `json:"message"`
}

type FollowersResponse struct {
	Followers []data.User `json:"followers"`
}
