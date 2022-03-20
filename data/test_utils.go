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
	followersSchemaSQL :=
		`CREATE TABLE followers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			follower_id TEXT,
			followed_id TEXT
		)`
	insertUserSQL := "INSERT INTO users (user_id) VALUES (?)"
	addFollowerSQL :=
		"INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)"

	s.MustExec(usersSchemaSQL)
	s.MustExec(followersSchemaSQL)

	tx := s.MustBegin()

	tx.MustExec(insertUserSQL, "1")
	tx.MustExec(insertUserSQL, "2")
	tx.MustExec(insertUserSQL, "3")

	tx.MustExec(addFollowerSQL, "1", "3")
	tx.MustExec(addFollowerSQL, "2", "1")
	tx.MustExec(addFollowerSQL, "3", "1")
	tx.MustExec(addFollowerSQL, "3", "2")

	tx.Commit()
}
