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

func IsFollowing(u string, f string) (bool, error) {
	fr := FollerRow{}
	err := Db.Get(
		&fr,
		"SELECT * from followers WHERE follower_id = ? AND followed_id = ?",
		u,
		f,
	)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func HasFollower(u string, f string) (bool, error) {
	return IsFollowing(f, u)
}

// Follow adds user to another user's followers.
func Follow(u string, t string) error {

	found, err := UserExists(u)

	if err != nil {
		return err
	}

	if !found {
		return ErrorUserNotFound
	}

	found, err = UserExists(t)

	if err != nil {
		return err
	}

	if !found {
		return ErrorUserNotFound
	}

	isFollowing, err := IsFollowing(u, t)

	if err != nil {
		return err
	}

	tx, err := Db.Begin()

	if err != nil {
		return err
	}

	if !isFollowing {
		tx.Exec(
			"INSERT INTO followers (follower_id, follow_id) VALUES (?, ?)",
			u,
			t,
		)
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func GetFollowers(u string) ([]string, error) {

	var followers []string

	err := Db.Select(
		&followers,
		"SELECT follower_id FROM followers WHERE followed_id = ?",
		u,
	)

	switch err {
	case nil:
		return followers, nil
	case sql.ErrNoRows:
		return []string{}, nil
	default:
		return nil, err
	}
}

func UserExists(uId string) (bool, error) {
	var u string
	err := Db.Get(
		&u,
		"SELECT id from users WHERE id = ?",
		uId,
	)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}

}
