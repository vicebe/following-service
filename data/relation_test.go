package data_test

import (
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
)

// TODO: test apparently doesn't work as expected, verify
func TestFollow(t *testing.T) {

	s, err := data.NewStore(
		"sqlite3", ":memory:", log.New(os.Stdout, "test", log.LstdFlags),
	)

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
func TestUnfollow(t *testing.T) {

	s, err := data.NewStore(
		"sqlite3", ":memory:", log.New(os.Stdout, "test", log.LstdFlags),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer s.Close()

	data.InitializeDB(s)

	u, v := "1", "3"

	if err = s.Unfollow(u, v); err != nil {
		t.Fatal(err)
	}

	isFollowing, err := s.IsFollowing(u, v)

	if err != nil {
		t.Fatal(err)
	}

	if isFollowing {
		t.Fatalf("User %s is following %s", u, v)
	}

	hasFollower, err := s.HasFollower(v, u)

	if err != nil {
		t.Fatal(err)
	}

	if hasFollower {
		t.Fatalf("User %s sill has as follower %s", v, u)
	}

}

func TestFollowCommunity(t *testing.T) {

	s, err := data.NewStore(
		"sqlite3", ":memory:", log.New(os.Stdout, "test", log.LstdFlags),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer s.Close()

	data.InitializeDB(s)

	u, c := "3", "1"

	if err = s.FollowCommunity(c, u); err != nil {
		t.Fatal(err)
	}

	isFollowing, err := s.IsFollowingCommunity(u, c)

	if err != nil {
		t.Fatal(err)
	}

	if !isFollowing {
		t.Fatalf("User %s is not following community %s", u, c)
	}

	hasFollower, err := s.CommunityHasFollower(c, u)

	if err != nil {
		t.Fatal(err)
	}

	if !hasFollower {
		t.Fatalf("Community %s does not have as follower user %s", u, c)
	}

}
