package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"
)

const (
	communityFollowingsRoutePath    = "/api/communities/{communityID}/followers"
	communityFollowingsRoutePathFmt = "/api/communities/%s/followers"
)

func TestCommunityHandler_FollowCommunity(ts *testing.T) {

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

	cr := data.NewCommunityRepositorySQL(l, db)
	ur := data.NewUserRepositorySQL(l, db)
	cs := services.NewCommunityService(l, cr, ur)
	ch := handlers.NewCommunityHandler(l, cs)

	const URL = communityFollowingsRoutePath + "/{userID}"
	const URLFmt = communityFollowingsRoutePathFmt + "/%s"

	r.Post(URL, ch.FollowCommunity)

	ts.Run("tests ability for user to follow a community", func(t *testing.T) {
		cID, uID := "1", "3"

		rUrl := fmt.Sprintf(URLFmt, cID, uID)

		req := httptest.NewRequest(http.MethodPost, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Fatal(rr.Body.String())
		}

		community, err := cr.FindBy("id", cID)

		if err != nil {
			t.Fatal(err)
		}

		user, err := ur.FindBy("id", uID)

		if err != nil {
			t.Fatal(err)
		}

		isFollowing, err := cr.IsFollowingCommunity(community, user)

		if err != nil {
			t.Fatal(err)
		}

		if !isFollowing {
			t.Fatalf("user %s is not following community %s", uID, cID)
		}
	})

	ts.Run("tests community not found", func(t *testing.T) {

		cID, uID := "4", "3"
		rUrl := fmt.Sprintf(URLFmt, uID, cID)

		req := httptest.NewRequest(http.MethodPost, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Status code returned %d", rr.Code)
		}

		expected := &handlers.SimpleResponse{
			Message: data.ErrorCommunityNotFound.Error(),
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

func TestCommunityHandler_UnfollowCommunity(ts *testing.T) {

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

	cr := data.NewCommunityRepositorySQL(l, db)
	ur := data.NewUserRepositorySQL(l, db)
	cs := services.NewCommunityService(l, cr, ur)
	ch := handlers.NewCommunityHandler(l, cs)

	const URL = communityFollowingsRoutePath + "/{userID}"
	const URLFmt = communityFollowingsRoutePathFmt + "/%s"

	r.Delete(URL, ch.UnfollowCommunity)

	ts.Run("tests ability for user to unfollow a community", func(t *testing.T) {
		cID, uID := "1", "3"

		rUrl := fmt.Sprintf(URLFmt, cID, uID)

		req := httptest.NewRequest(http.MethodDelete, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Fatal(rr.Body.String())
		}

		community, err := cr.FindBy("id", cID)

		if err != nil {
			t.Fatal(err)
		}

		user, err := ur.FindBy("id", uID)

		if err != nil {
			t.Fatal(err)
		}

		isFollowing, err := cr.IsFollowingCommunity(community, user)

		if err != nil {
			t.Fatal(err)
		}

		if isFollowing {
			t.Fatalf("user %s is following community %s", uID, cID)
		}
	})

	ts.Run("tests community not found", func(t *testing.T) {

		cID, uID := "4", "3"
		rUrl := fmt.Sprintf(URLFmt, uID, cID)

		req := httptest.NewRequest(http.MethodDelete, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Status code returned %d", rr.Code)
		}

		expected := &handlers.SimpleResponse{
			Message: data.ErrorCommunityNotFound.Error(),
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

func TestCommunityHandler_GetCommunityFollowers(t *testing.T) {

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

	cr := data.NewCommunityRepositorySQL(l, db)
	ur := data.NewUserRepositorySQL(l, db)
	cs := services.NewCommunityService(l, cr, ur)
	ch := handlers.NewCommunityHandler(l, cs)

	const URL = communityFollowingsRoutePath
	const URLFmt = communityFollowingsRoutePathFmt

	r.Get(URL, ch.GetCommunityFollowers)

	t.Run("test Get community followers", func(t *testing.T) {
		cID := "1"
		rUrl := fmt.Sprintf(URLFmt, cID)

		req := httptest.NewRequest(http.MethodGet, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatalf(
				"http Code status not expected, expected %d got %d",
				res.StatusCode,
				http.StatusOK,
			)
		}

		expected := []data.User{
			{
				ID:         1,
				ExternalID: "1",
			},
			{
				ID:         2,
				ExternalID: "2",
			},
		}

		var got handlers.FollowersResponse

		if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}

		if len(got.Followers) != len(expected) {
			t.Fatalf(
				"followers list not expected, expected %d, got %d",
				len(expected),
				len(got.Followers),
			)
		}

		sort.SliceStable(got.Followers, func(i, j int) bool {
			return got.Followers[i].ID < got.Followers[j].ID
		})

		for i, user := range expected {
			if user.ID != got.Followers[i].ID {
				t.Fatalf(
					"follower not expected.\nexpected:\n%#v\ngot:\n%#v",
					user,
					got.Followers[i],
				)
			}
		}

	})

	t.Run("test Community Not found", func(t *testing.T) {

		cID := "2"
		rUrl := fmt.Sprintf(URLFmt, cID)

		req := httptest.NewRequest(http.MethodGet, rUrl, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf(
				"http Code status not expected, expected %d got %d",
				res.StatusCode,
				http.StatusOK,
			)
		}
	})
}
