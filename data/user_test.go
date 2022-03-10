package data_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
)

func initializeDB() *data.DatabaseObject {
	c := sqlx.MustConnect("sqlite3", ":memory:")
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

	c.MustExec(usersSchemaSQL)
	c.MustExec(followersSchemaSQL)

	tx := c.MustBegin()

	tx.MustExec(insertUserSQL, "1")
	tx.MustExec(insertUserSQL, "2")
	tx.MustExec(insertUserSQL, "3")

	tx.MustExec(addFollowerSQL, "1", "3")
	tx.MustExec(addFollowerSQL, "2", "1")
	tx.MustExec(addFollowerSQL, "3", "1")
	tx.MustExec(addFollowerSQL, "3", "2")

	tx.Commit()

	db := data.NewDatabaseObject(c)

	return db
}

func TestToJson(t *testing.T) {

	type simpleResponse struct {
		Message string `json:"message"`
	}

	sr := &simpleResponse{Message: "test"}

	var b bytes.Buffer
	if err := data.ToJson(sr, &b); err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(b.String())
	wanted := fmt.Sprintf("{\"message\":\"%s\"}", sr.Message)

	if got != wanted {
		t.Fatalf("wanted \"%v\" got \"%v\"", wanted, got)
	}
}

func TestUserExist(ts *testing.T) {

	dbo := initializeDB()
	defer dbo.C.Close()

	ts.Run("tests user found", func(t *testing.T) {
		exists, err := dbo.UserExists("1")

		if err != nil {
			t.Fatal(err)

		}

		fmt.Printf("exists user 1: %#v", exists)

		if !exists {
			t.Fatalf("User not found")
		}
	})

	ts.Run("tests user not found", func(t *testing.T) {
		exists, err := dbo.UserExists("4")

		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("exists user 4: %#v", exists)

		if exists {
			t.Fatalf("User found")
		}
	})
}

func TestFollow(t *testing.T) {
	dbo := initializeDB()
	defer dbo.C.Close()

	u, v := "1", "2"

	if err := dbo.Follow(u, v); err != nil {
		t.Fatal(err)
	}

	isFollowing, err := dbo.IsFollowing(u, v)

	if err != nil {
		t.Fatal(err)
	}

	if !isFollowing {
		t.Fatalf("User %s is not following %s", u, v)
	}

	hasFollower, err := dbo.HasFollower(v, u)

	if err != nil {
		t.Fatal(err)
	}

	if !hasFollower {
		t.Fatalf("User %s does not have as follower %s", v, u)
	}

}
