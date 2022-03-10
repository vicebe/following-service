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
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
)

func initializeDB() *data.DatabaseObject {
	c := sqlx.MustConnect("sqlite3", ":memory:")
	usersSchemaSQL :=
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT
		)`
	followersSchemaSQL :=
		`CREATE TABLE followers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			follower_id TEXT,
			followed_id TEXT
		)`
	insertUserSQL := "INSERT INTO users (user_id) VALUES (?)"
	addFollowerSQL :=
		"INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)"

	c.MustExec(usersSchemaSQL)
	c.MustExec(followersSchemaSQL)

	tx := c.MustBegin()

	tx.MustExec(insertUserSQL, "1")
	tx.MustExec(insertUserSQL, "2")
	tx.MustExec(insertUserSQL, "3")

	tx.MustExec(addFollowerSQL, "1", "3")
	tx.MustExec(addFollowerSQL, "2", "1")
	tx.MustExec(addFollowerSQL, "3", "1")
	tx.MustExec(addFollowerSQL, "3", "2")

	tx.Commit()

	db := data.NewDatabaseObject(c)

	return db
}
func TestFollowUser(ts *testing.T) {

	r := chi.NewRouter()
	l := log.New(os.Stdout, "following-service-test", log.LstdFlags)
	db := initializeDB()
	us := services.NewUserService(l, db)
	sh := handlers.NewServiceHandler(l, us)
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

		isFollowing, err := db.IsFollowing(from, to)

		if err != nil {
			t.Fatal(err)
		}

		if !isFollowing {
			t.Fatalf("user %s is not following %s", from, to)
		}

		hasFollower, err := db.HasFollower(to, from)

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
