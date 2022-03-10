package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
)

// ToJson seriealizes the given interface into a string based JSON format
func ToJson(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}

type FollerRow struct {
	ID         int    `db:"id"`
	FollowerId string `db:"follower_id"`
	FollewedId string `db:"followed_id"`
}

var ErrorUserNotFound = fmt.Errorf("user not found")
var ErrorUserAlreadyFollowed = fmt.Errorf("user being followed")

func (db *DatabaseObject) IsFollowing(u string, f string) (bool, error) {
	fr := FollerRow{}
	err := db.C.Get(&fr, IsFollowerSQL, u, f)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (db *DatabaseObject) HasFollower(u string, f string) (bool, error) {
	return db.IsFollowing(f, u)
}

// Follow adds user to another user's followers.
func (db *DatabaseObject) Follow(u string, t string) error {

	found, err := db.UserExists(u)

	if err != nil {
		return err
	}

	if !found {
		return ErrorUserNotFound
	}

	found, err = db.UserExists(t)

	if err != nil {
		return err
	}

	if !found {
		return ErrorUserNotFound
	}

	isFollowing, err := db.IsFollowing(u, t)

	if err != nil {
		return err
	}

	tx, err := db.C.Begin()

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

func (db *DatabaseObject) GetFollowers(u string) ([]string, error) {

	var followers []string

	err := db.C.Select(&followers, GetFollowersSQL, u)

	switch err {
	case nil:
		return followers, nil
	case sql.ErrNoRows:
		return []string{}, nil
	default:
		return nil, err
	}
}

func (db *DatabaseObject) UserExists(uId string) (bool, error) {
	var u string
	err := db.C.Get(&u, FindUserSQL, uId)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}

}
