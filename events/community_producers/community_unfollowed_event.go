package communityproducers

type CommunityUnfollowedEvent struct {
	// unfollowed community
	CommunityID string `json:"community_id"`

	// follower user
	UserID string `json:"user_id"`
}
