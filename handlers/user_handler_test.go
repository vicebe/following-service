// TODO: improve error responses
package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/handlers/test_utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
)

const (
	usersFollowersRoutePath    = "/api/users/{userID}/followers"
	usersFollowersRoutePathFmt = "/api/users/%s/followers"

	usersCommunitiesRoutePath    = "/api/users/{userID}/communities"
	usersCommunitiesRoutePathFmt = "/api/users/%s/communities"
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

func TestUserHandler_GetCommunities(t *testing.T) {
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

	const URL = usersCommunitiesRoutePath
	const URLFmt = usersCommunitiesRoutePathFmt
	r.Get(URL, uh.GetCommunities)

	uID := "1"
	rUrl := fmt.Sprintf(URLFmt, uID)

	req := httptest.NewRequest(http.MethodGet, rUrl, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	res := rr.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	type CommunityResponse struct {
		Communities []data.Community `json:"communities"`
	}

	var communityResponse CommunityResponse
	err := json.Unmarshal(rr.Body.Bytes(), &communityResponse)
	if err != nil {
		t.Fatal(err)
	}

	communities := communityResponse.Communities

	if len(communities) != 1 {
		t.Fatalf("expected one community got %d", len(communities))
	}

	if communities[0].ID != 1 {
		t.Fatalf("expected community 1 got %d", communities[0].ID)
	}

}

func TestUserHandler_GetFollowers(t *testing.T) {
	// responses
	followersResponse, _ := json.Marshal(
		&handlers.FollowersResponse{Followers: test_utils.FollowersList},
	)

	internalError := handlers.MakeInternalErrorResponse()
	internalErrorResponse, _ := json.Marshal(&internalError)

	// tests cases
	tests := []struct {
		name         string
		l            *log.Logger
		userService  services.UserServiceI
		method       string
		statusCode   int
		responseBody string
		middleware   func(http.Handler) http.Handler
	}{
		{
			name: "test get followers success",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceGetFollowersMock{},
			method:       http.MethodGet,
			statusCode:   http.StatusOK,
			responseBody: string(followersResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
		{
			name: "test user not passed in context",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceGetFollowersMock{},
			method:       http.MethodGet,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.IdentityMiddleware,
		},
		{
			name: "test get user errors",
			l: log.New(
				os.Stdout,
				"following-service-test",
				log.LstdFlags,
			),
			userService:  test_utils.UserServiceGetFollowersErrorMock{},
			method:       http.MethodGet,
			statusCode:   http.StatusInternalServerError,
			responseBody: string(internalErrorResponse),
			middleware:   test_utils.AddUserToRequestContext(&test_utils.UserOne),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				tt.method,
				"/get-followers",
				nil,
			)

			rr := httptest.NewRecorder()

			uh := handlers.NewUserHandler(tt.l, tt.userService)

			r := chi.NewRouter()

			r.Use(tt.middleware)

			r.Get("/get-followers", uh.GetFollowers)

			r.ServeHTTP(rr, request)

			if rr.Code != tt.statusCode {
				t.Errorf(
					"expected http status code %d got %d",
					rr.Code,
					tt.statusCode,
				)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.responseBody {
				t.Errorf(
					"expected response '%s' got '%s",
					rr.Body.String(),
					tt.responseBody,
				)
			}
		})
	}
}
