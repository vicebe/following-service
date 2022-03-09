package handlers

type SimpleResponse struct {
	Message string `json:"message"`
}

type FollowersResponse struct {
	Followers []string `json:"followers"`
}
