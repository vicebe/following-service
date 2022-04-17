package communityproducers

type CommunityFollowedEvent struct {
	// followed user
	CommunityID string `json:"community_id"`

	// follower user
	UserID string `json:"user_id"`
}
