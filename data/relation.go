package data

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

// Relation represents a row from the followers table
type RelationUsersFollowers struct {

	// relation id
	ID int `db:"id"`

	// User id of the follower in the realtionship
	FollowerId string `db:"follower_id"`

	// User id of the user being followed
	FollewedId string `db:"followed_id"`
}

type RelationCommunityFollowers struct {
	ID int `db:"id"`

	CommunityId string `db:"community_id"`

	FollowerId string `db:"follower_id"`
}

// RelationStore is a store with a db connection for Relation related operations
type RelationStore struct {
	*sqlx.DB

	l *log.Logger

	// We're adding the user store as dependency because some functions use
	// functions from the UserStore
	us *UserStore
	cs *CommunityStore
}

func NewRelationStore(db *sqlx.DB, l *log.Logger) *RelationStore {
	return &RelationStore{
		l:  l,
		DB: db,
		us: NewUserStore(db, l),
		cs: NewCommunityStore(db, l),
	}
}

// IsFollowing verifies if user `u` is following user `f` providing their ids
func (rs *RelationStore) IsFollowing(u string, f string) (bool, error) {
	r := &RelationUsersFollowers{}
	err := rs.Get(r, IsFollowerSQL, u, f)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

// HasFollower verifies if user `u` has user `f` as follower using their ids
func (rs *RelationStore) HasFollower(u string, f string) (bool, error) {
	return rs.IsFollowing(f, u)
}

// Follow adds user to another user's followers.
func (rs *RelationStore) Follow(u string, t string) error {

	found, err := rs.us.UserExists(u)

	if err != nil {
		return err
	}

	if !found {
		return ErrorUserNotFound
	}

	found, err = rs.us.UserExists(t)

	if err != nil {
		return err
	}

	if !found {
		return ErrorUserNotFound
	}

	isFollowing, err := rs.IsFollowing(u, t)

	if err != nil {
		return err
	}

	tx, err := rs.Begin()

	if err != nil {
		return err
	}

	if !isFollowing {
		_, err := tx.Exec(FollowUserSQL, u, t)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

// Unfollow removes user from another user's followers.
func (rs *RelationStore) Unfollow(u string, t string) error {

	found, err := rs.us.UserExists(u)

	if err != nil {
		return err
	}

	if !found {
		return ErrorUserNotFound
	}

	found, err = rs.us.UserExists(t)

	if err != nil {
		return err
	}

	if !found {
		return ErrorUserNotFound
	}

	isFollowing, err := rs.IsFollowing(u, t)

	if err != nil {
		return err
	}

	tx, err := rs.Begin()

	if err != nil {
		return err
	}

	if isFollowing {
		_, err := tx.Exec(UnfollowUserSQL, u, t)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

// GetFollowers returns user's followers ids
func (rs *RelationStore) GetFollowers(u string) ([]string, error) {

	followers := []string{}

	err := rs.Select(&followers, GetFollowersSQL, u)

	if err != nil {
		return nil, err
	}

	return followers, nil
}

// IsFollowingCommunity verifies if user `u` is following community `c`providing
// their ids
func (rs *RelationStore) IsFollowingCommunity(u string, c string) (bool, error) {
	r := &RelationCommunityFollowers{}
	err := rs.Get(r, IsFollowingCommunitySQL, u, c)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

// CommunityHasFollower verifies if community `c` has user `u` as follower using
// their ids.
func (rs *RelationStore) CommunityHasFollower(c string, u string) (bool, error) {
	return rs.IsFollowingCommunity(u, c)
}

// Follow adds user to a community's followers.
func (rs *RelationStore) FollowCommunity(c string, u string) error {
	found, err := rs.cs.CommunityExists(c)

	if err != nil {
		return err
	}

	if !found {
		return ErrorCommunityNotFound
	}

	found, err = rs.us.UserExists(u)

	if err != nil {
		return err
	}

	if !found {
		return ErrorUserNotFound
	}

	isFollowing, err := rs.IsFollowingCommunity(u, c)

	if err != nil {
		return err
	}

	tx, err := rs.Begin()

	if err != nil {
		return err
	}

	if !isFollowing {
		_, err := tx.Exec(FollowCommunitySQL, u, c)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}
