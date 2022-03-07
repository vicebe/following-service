// TODO: test failure
package handlers_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
)

func TestGetFollowers(t *testing.T) {

	from, to := "1", "2"
	rUrl := fmt.Sprintf("/%s/follow/%s", from, to)

	r := chi.NewRouter()
	l := log.New(os.Stdout, "following-service-test", log.LstdFlags)
	us := services.NewUserService(l)
	sh := handlers.NewServiceHandler(l, us)
	r.Post("/{userId}/follow/{toFollowId}", sh.FollowUser)

	req := httptest.NewRequest(http.MethodPost, rUrl, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	res := rr.Result()

	if res.StatusCode != http.StatusNoContent {
		t.Fatal(res.StatusCode)
	}

	user, _ := data.GetUserByID(from)
	newFollowing := user.Following[len(user.Following)-1]
	if newFollowing != to {
		t.Fatalf("user %s is not following %s", from, to)
	}

	user, _ = data.GetUserByID(to)
	newFollower := user.Followers[len(user.Followers)-1]
	if newFollower != from {
		t.Fatalf("user %s is not being followed by %s", to, from)
	}
}
