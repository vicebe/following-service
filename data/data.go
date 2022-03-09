package data

import "github.com/jmoiron/sqlx"

var Db *sqlx.DB = nil

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
			follow_id
		)
		VALUES (?, ?)`

var GetFollowersSQL = "SELECT follower_id FROM followers WHERE followed_id = ?"

var FindUserSQL = "SELECT id from users WHERE id = ?"
