package data

// CommunityRepository is a repository interface to implement common methods for
// accessing community related information.
type CommunityRepository interface {

	// Create a new community.
	Create(community *Community) error

	// Update an existing community.
	Update(community *Community, newCommunity *Community) error

	// Delete a community
	Delete(community *Community) error

	// FindBy finds a community given a key and a value
	FindBy(key string, value interface{}) (*Community, error)

	// IsFollowingCommunity checks if user is following a community
	IsFollowingCommunity(community *Community, user *User) (bool, error)

	// FollowCommunity adds user to the community followers.
	FollowCommunity(community *Community, user *User) error

	// UnfollowCommunity removes user from the community followers.
	UnfollowCommunity(community *Community, user *User) error

	// GetCommunityFollowers retrieves a list of the community followers.
	GetCommunityFollowers(community *Community) ([]User, error)
}
