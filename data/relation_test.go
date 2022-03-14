package data_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
)

// TODO: test apparently doesn't work as expected, verify
func TestFollow(t *testing.T) {

	s, err := data.NewStore("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	defer s.Close()

	data.InitializeDB(s)

	u, v := "1", "2"

	if err = s.Follow(u, v); err != nil {
		t.Fatal(err)
	}

	isFollowing, err := s.IsFollowing(u, v)

	if err != nil {
		t.Fatal(err)
	}

	if !isFollowing {
		t.Fatalf("User %s is not following %s", u, v)
	}

	hasFollower, err := s.HasFollower(v, u)

	if err != nil {
		t.Fatal(err)
	}

	if !hasFollower {
		t.Fatalf("User %s does not have as follower %s", v, u)
	}

}
