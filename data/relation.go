package data

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

// Relation represents a row from the followers table
type Relation struct {

	// relation id
	ID int `db:"id"`

	// User id of the follower in the realtionship
	FollowerId string `db:"follower_id"`

	// User id of the user being followed
	FollewedId string `db:"followed_id"`
}

// RelationStore is a store with a db connection for Relation related operations
type RelationStore struct {
	*sqlx.DB

	l *log.Logger

	// We're adding the user store as dependency because some functions use
	// functions from the UserStore
	us *UserStore
}

func NewRelationStore(db *sqlx.DB, l *log.Logger) *RelationStore {
	return &RelationStore{
		l:  l,
		DB: db,
		us: NewUserStore(db, l),
	}
}

// IsFollowing verifies if user `u` is following user `f` providing their ids
func (rs *RelationStore) IsFollowing(u string, f string) (bool, error) {
	r := &Relation{}
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
