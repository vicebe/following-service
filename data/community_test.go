package data_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
)

func TestCommunityExist(ts *testing.T) {

	s, err := data.NewStore(
		"sqlite3", ":memory:", log.New(os.Stdout, "test", log.LstdFlags),
	)

	if err != nil {
		ts.Fatal(err)
	}

	defer s.Close()

	data.InitializeDB(s)

	ts.Run("tests community found", func(t *testing.T) {
		exists, err := s.CommunityExists("1")

		if err != nil {
			t.Fatal(err)

		}

		fmt.Printf("exists community 1: %#v", exists)

		if !exists {
			t.Fatalf("Community not found")
		}
	})

	ts.Run("tests community not found", func(t *testing.T) {
		exists, err := s.CommunityExists("4")

		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("exists community 4: %#v", exists)

		if exists {
			t.Fatalf("Community found")
		}
	})
}
