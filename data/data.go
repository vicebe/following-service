package data

import "github.com/jmoiron/sqlx"

type DatabaseObject struct {
	C *sqlx.DB
}

func NewDatabaseObject(c *sqlx.DB) *DatabaseObject {
	return &DatabaseObject{C: c}
}

var IsFollowerSQL = `
	SELECT
		*
	FROM Followers
	WHERE follower_id = ?
		AND followed_id = ?
	`

var FollowUserSQL = `
	INSERT
		INTO followers (
			follower_id,
			followed_id
		)
		VALUES (?, ?)`

var GetFollowersSQL = "SELECT follower_id FROM followers WHERE followed_id = ?"

var FindUserSQL = "SELECT id from users WHERE id = ?"
