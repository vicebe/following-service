package data_test

import (
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
)

func TestUserExist(ts *testing.T) {

	s, err := data.NewStore("sqlite3", ":memory:")

	if err != nil {
		ts.Fatal(err)
	}

	defer s.Close()

	data.InitializeDB(s)

	ts.Run("tests user found", func(t *testing.T) {
		exists, err := s.UserExists("1")

		if err != nil {
			t.Fatal(err)

		}

		fmt.Printf("exists user 1: %#v", exists)

		if !exists {
			t.Fatalf("User not found")
		}
	})

	ts.Run("tests user not found", func(t *testing.T) {
		exists, err := s.UserExists("4")

		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("exists user 4: %#v", exists)

		if exists {
			t.Fatalf("User found")
		}
	})
}
