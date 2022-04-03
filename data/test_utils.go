package data

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func CreateCommunitiesTable(conn *sqlx.DB) {
	const communitySchemaSQL = `CREATE TABLE communities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			external_id TEXT
		)`

	conn.MustExec(communitySchemaSQL)
}

func CreateUsersTable(conn *sqlx.DB) {
	const usersSchemaSQL = `CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			external_id TEXT
		)`

	conn.MustExec(usersSchemaSQL)
}

func CreateUsersFollowersTable(conn *sqlx.DB) {
	const usersFollowersSchemaSQL = `CREATE TABLE users_followers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			followee_id INTEGER,
			follower_id INTEGER,
			FOREIGN KEY (followee_id) REFERENCES users(id),
			FOREIGN KEY (follower_id) REFERENCES users(id)
		)`

	conn.MustExec(usersFollowersSchemaSQL)
}

func CreateCommunitiesFollowersTable(conn *sqlx.DB) {

	const communitiesFollowersSchemaSQL = `CREATE TABLE communities_followers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			community_id INTEGER,
			follower_id INTEGER,
			FOREIGN KEY (community_id) REFERENCES communities(id),
			FOREIGN KEY (follower_id) REFERENCES users(id)
		)`

	conn.MustExec(communitiesFollowersSchemaSQL)
}

func InitializeDB(conn *sqlx.DB) {

	insertUserSQL := "INSERT INTO users (external_id) VALUES (?)"
	addFollowerSQL :=
		"INSERT INTO users_followers (followee_id, follower_id) VALUES (?, ?)"
	insertCommunitySQL := "INSERT INTO communities (external_id) VALUES (?)"
	addFollowerToCommunitySQL :=
		`INSERT
			INTO communities_followers (community_id, follower_id)
			VALUES (?, ?)
		`

	CreateUsersTable(conn)
	CreateUsersFollowersTable(conn)
	CreateCommunitiesTable(conn)
	CreateCommunitiesFollowersTable(conn)

	tx := conn.MustBegin()

	tx.MustExec(insertUserSQL, "1")
	tx.MustExec(insertUserSQL, "2")
	tx.MustExec(insertUserSQL, "3")

	tx.MustExec(addFollowerSQL, "1", "3")
	tx.MustExec(addFollowerSQL, "2", "1")
	tx.MustExec(addFollowerSQL, "3", "1")
	tx.MustExec(addFollowerSQL, "3", "2")

	tx.MustExec(insertCommunitySQL, "1")
	tx.MustExec(addFollowerToCommunitySQL, "1", "1")
	tx.MustExec(addFollowerToCommunitySQL, "1", "2")

	tx.Commit()
}
