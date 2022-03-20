package data

import (
	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB(s *Store) {

	usersSchemaSQL :=
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT
		)`
	insertUserSQL := "INSERT INTO users (user_id) VALUES (?)"
	addFollowerSQL :=
		"INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)"
	followersSchemaSQL :=
		`CREATE TABLE followers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			follower_id TEXT,
			followed_id TEXT
		)`
	communitySchemaSQL :=
		`CREATE TABLE communities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			community_id TEXT
		)`
	communityFollowersSchemaSQL :=
		`CREATE TABLE community_followers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			community_id TEXT,
			follower_id TEXT
		)`
	insertCommunitySQL := "INSERT INTO communities (community_id) VALUES (?)"
	addFollowerToCommunitySQL :=
		`INSERT
			INTO community_followers (community_id, follower_id)
			VALUES (?, ?)
		`

	s.MustExec(usersSchemaSQL)
	s.MustExec(followersSchemaSQL)
	s.MustExec(communitySchemaSQL)
	s.MustExec(communityFollowersSchemaSQL)

	tx := s.MustBegin()

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
