package data

import (
	"encoding/json"
	"fmt"
	"io"
)

// ToJson seriealizes the given interface into a string based JSON format
func ToJson(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}

type User struct {
	ID        string   `json:"id"`
	Followers []string `json:"followers"`
	Following []string `json:"following"`
}

var users = []*User{
	{
		ID:        "1",
		Followers: []string{"2", "3"},
		Following: []string{"3"},
	},
	{
		ID:        "2",
		Followers: []string{"3"},
		Following: []string{"1"},
	},
	{
		ID:        "3",
		Followers: []string{"1"},
		Following: []string{"1", "2"},
	},
}

var ErrorUserNotFound = fmt.Errorf("User not found")
var ErrorUserAlreadyFollowed = fmt.Errorf("User being followed")

func elementExists(els []string, key string) bool {
	for _, el := range els {
		if el == key {
			return true
		}
	}

	return false

}

func (u *User) hasFollower(id string) bool {
	return elementExists(u.Followers, id)
}

func (u *User) isFollowing(id string) bool {
	return elementExists(u.Following, id)
}

func (u *User) Follow(toFollowId string) error {

	toFollow, err := GetUserByID(toFollowId)

	if err != nil {
		return ErrorUserNotFound
	}

	if !u.isFollowing(toFollowId) {
		u.Following = append(u.Following, toFollowId)
	}

	if !toFollow.hasFollower(u.ID) {
		toFollow.Followers = append(toFollow.Followers, u.ID)
	}

	return nil
}

func (u *User) GetFollowers() []string {
	return u.Followers
}

func GetUserByID(userId string) (*User, error) {
	// For now we get the users id from a mock db structure.

	for _, p := range users {

		if p.ID == userId {
			return p, nil
		}
	}

	return nil, ErrorUserNotFound
}
