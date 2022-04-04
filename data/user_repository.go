package data

type UserRepository interface {

	// Create a new user.
	Create(user *User) error

	// Update an existing user.
	Update(user *User, newUser *User) error

	// Delete a user
	Delete(user *User) error

	// FindBy finds a user given a key and a value
	FindBy(key string, value interface{}) (*User, error)

	// IsFollowingUser checks if user is following a user
	IsFollowingUser(follower *User, followee *User) (bool, error)

	// FollowUser adds user to the user followers.
	FollowUser(user *User, followee *User) error

	// UnfollowUser removes follower from the followee's followers.
	UnfollowUser(follower *User, followee *User) error

	// GetUserFollowers retrieves a list of the user followers.
	GetUserFollowers(user *User) ([]User, error)

	// GetUserFollowees retrieves a list of the users that this user is
	// following.
	GetUserFollowees(user *User) ([]User, error)

	// GetUserCommunities retrieves a list of the communities that the user
	// follows
	GetUserCommunities(user *User) ([]Community, error)
}
