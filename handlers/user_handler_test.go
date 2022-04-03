// TODO: improve error responses
package handlers_test

import (
	"bytes"
	"fmt"
	"github.com/jmoiron/sqlx"
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

const (
	usersFollowersRoutePath    = "/api/users/{userID}/followers"
	usersFollowersRoutePathFmt = "/api/users/%s/followers"
)

func TestUserHandler_FollowUser(ts *testing.T) {

	r := chi.NewRouter()
	l := log.New(os.Stdout, "following-service-test", log.LstdFlags)
	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	data.InitializeDB(db)

	ur := data.NewUserRepositorySQL(l, db)
	us := services.NewUserService(l, ur)
	uh := handlers.NewUserHandler(l, us)

	const URL = usersFollowersRoutePath + "/{followerID}"
	const URLFmt = usersFollowersRoutePathFmt + "/%s"

	r.Post(URL, uh.FollowUser)

	ts.Run("tests ability for user to follow", func(t *testing.T) {
		from, to := "1", "2"
		rUrl := fmt.Sprintf(URLFmt, to, from)

		req := httptest.NewRequest(http.MethodPost, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Fatal(rr.Body.String())
		}

		follower, err := ur.FindBy("id", from)

		if err != nil {
			t.Fatal(err)
		}

		followee, err := ur.FindBy("id", to)

		if err != nil {
			t.Fatal(err)
		}

		isFollowing, err := ur.IsFollowingUser(follower, followee)

		if err != nil {
			t.Fatal(err)
		}

		if !isFollowing {
			t.Fatalf("user %s is not following %s", from, to)
		}

	})

	ts.Run("tests user not found", func(t *testing.T) {

		from, to := "4", "2"
		rUrl := fmt.Sprintf(URLFmt, to, from)

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

func TestUserHandler_UnfollowUser(ts *testing.T) {

	r := chi.NewRouter()
	l := log.New(os.Stdout, "following-service-test", log.LstdFlags)
	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	data.InitializeDB(db)

	ur := data.NewUserRepositorySQL(l, db)
	us := services.NewUserService(l, ur)
	uh := handlers.NewUserHandler(l, us)

	const URL = usersFollowersRoutePath + "/{followerID}"
	const URLFmt = usersFollowersRoutePathFmt + "/%s"
	r.Delete(URL, uh.UnfollowUser)

	ts.Run("tests ability for user to unfollow", func(t *testing.T) {
		from, to := "1", "3"
		rUrl := fmt.Sprintf(URLFmt, to, from)

		req := httptest.NewRequest(http.MethodDelete, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Fatal(rr.Body.String())
		}

		follower, err := ur.FindBy("id", from)

		if err != nil {
			t.Fatal(err)
		}

		followee, err := ur.FindBy("id", to)

		if err != nil {
			t.Fatal(err)
		}

		isFollowing, err := ur.IsFollowingUser(follower, followee)

		if err != nil {
			t.Fatal(err)
		}

		if isFollowing {
			t.Fatalf("user %s is following %s", from, to)
		}

	})

	ts.Run("tests user not found", func(t *testing.T) {

		from, to := "4", "2"
		rUrl := fmt.Sprintf(URLFmt, to, from)

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
