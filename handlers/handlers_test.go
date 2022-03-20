// TODO: improve error responses
package handlers_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
)

func TestFollowUser(ts *testing.T) {

	r := chi.NewRouter()
	l := log.New(os.Stdout, "following-service-test", log.LstdFlags)

	ds, err := data.NewStore("sqlite3", ":memory:")

	if err != nil {
		ts.Fatal(err)
	}

	defer ds.Close()

	data.InitializeDB(ds)

	as := services.NewAppService(l, ds)
	sh := handlers.NewHandler(l, as)

	r.Post("/{userId}/follow/{toFollowId}", sh.FollowUser)

	ts.Run("tests ability for user to follow", func(t *testing.T) {
		from, to := "1", "2"
		rUrl := fmt.Sprintf("/%s/follow/%s", from, to)

		req := httptest.NewRequest(http.MethodPost, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Fatal(rr.Body.String())
		}

		isFollowing, err := ds.IsFollowing(from, to)

		if err != nil {
			t.Fatal(err)
		}

		if !isFollowing {
			t.Fatalf("user %s is not following %s", from, to)
		}

		hasFollower, err := ds.HasFollower(to, from)

		if err != nil {
			t.Fatal(err)
		}

		if !hasFollower {
			t.Fatalf("user %s has no follower %s", to, from)
		}
	})

	ts.Run("tests user not found", func(t *testing.T) {

		from, to := "4", "2"
		rUrl := fmt.Sprintf("/%s/follow/%s", from, to)

		req := httptest.NewRequest(http.MethodPost, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Status code returned %d", rr.Code)
		}

		expected := &handlers.SimpleResponse{
			Message: data.ErrorUserNotFound.Error(),
		}

		var expectedRes bytes.Buffer

		data.ToJson(expected, &expectedRes)

		jsonRes := rr.Body.String()
		expectedResStr := expectedRes.String()

		if jsonRes != expectedResStr {
			t.Fatalf(
				"responses are not equal.\nexpected: %s\ngiven %s",
				expectedResStr,
				jsonRes,
			)
		}
	})

}

func TestUnfollowUser(ts *testing.T) {

	r := chi.NewRouter()
	l := log.New(os.Stdout, "following-service-test", log.LstdFlags)

	ds, err := data.NewStore("sqlite3", ":memory:")

	if err != nil {
		ts.Fatal(err)
	}

	defer ds.Close()

	data.InitializeDB(ds)

	as := services.NewAppService(l, ds)
	sh := handlers.NewHandler(l, as)

	r.Delete("/{userId}/unfollow/{toUnfollowId}", sh.UnfollowUser)

	ts.Run("tests ability for user to unfollow", func(t *testing.T) {
		from, to := "1", "3"
		rUrl := fmt.Sprintf("/%s/unfollow/%s", from, to)

		req := httptest.NewRequest(http.MethodDelete, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Fatal(rr.Body.String())
		}

		isFollowing, err := ds.IsFollowing(from, to)

		if err != nil {
			t.Fatal(err)
		}

		if isFollowing {
			t.Fatalf("user %s is following %s", from, to)
		}

		hasFollower, err := ds.HasFollower(to, from)

		if err != nil {
			t.Fatal(err)
		}

		if hasFollower {
			t.Fatalf("user %s has follower %s", to, from)
		}
	})

	ts.Run("tests user not found", func(t *testing.T) {

		from, to := "4", "2"
		rUrl := fmt.Sprintf("/%s/unfollow/%s", from, to)

		req := httptest.NewRequest(http.MethodDelete, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Status code returned %d", rr.Code)
		}

		expected := &handlers.SimpleResponse{
			Message: data.ErrorUserNotFound.Error(),
		}

		var expectedRes bytes.Buffer

		data.ToJson(expected, &expectedRes)

		jsonRes := rr.Body.String()
		expectedResStr := expectedRes.String()

		if jsonRes != expectedResStr {
			t.Fatalf(
				"responses are not equal.\nexpected: %s\ngiven %s",
				expectedResStr,
				jsonRes,
			)
		}
	})

}
