package userproducers

type UserUnfollowedEvent struct {
	// followed user
	FolloweeID string `json:"followee_id"`

	// follower user
	FollowerID string `json:"follower_id"`
}
