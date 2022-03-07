package data_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/vicebe/following-service/data"
)

func assertEqualStringSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func assertEqualUsers(a *data.User, b *data.User) bool {
	return a.ID == b.ID &&
		assertEqualStringSlices(a.Followers, b.Followers) &&
		assertEqualStringSlices(a.Following, b.Following)
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

func TestUserByID(t *testing.T) {
	id := "1"
	user, err := data.GetUserByID(id)

	if err != nil {
		t.Fatalf("Could not get user by id %v", id)
	}

	wanted := &data.User{
		ID:        "1",
		Followers: []string{"2", "3"},
		Following: []string{"3"},
	}

	if !assertEqualUsers(user, wanted) {
		t.Fatalf("user found not expected: %#v != %#v", user, wanted)
	}
}

func TestFollow(t *testing.T) {
	u := &data.User{
		ID:        "1",
		Followers: []string{},
		Following: []string{"2"},
	}

	id := "3"

	if err := u.Follow(id); err != nil {
		t.Fatalf("Error following user %s", id)
	}

	if u.Following[len(u.Following)-1] != id {
		t.Fatalf("User %s is not following %s", u.ID, id)
	}

	followed, _ := data.GetUserByID(id)

	if followed.Followers[len(followed.Followers)-1] != u.ID {
		t.Fatalf("User %s is not following %s", id, u.ID)
	}

	fmt.Printf("user %#v, following %#v", u, followed)

}
