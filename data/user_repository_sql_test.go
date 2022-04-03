package data_test

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
	"log"
	"os"
	"testing"
)

func TestUserRepositorySQL_Create(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	ur := data.NewUserRepositorySQL(
		db,
		log.New(os.Stdout, "test", log.LstdFlags),
	)
	defer db.Close()

	data.CreateUsersTable(db)

	u := &data.User{
		ID:         0,
		ExternalID: "NEW-USER",
	}

	err := ur.Create(u)

	if err != nil {
		t.Fatal(err)
	}

	addedU := &data.User{}
	err = db.Get(addedU, `SELECT * FROM users`)

	if err != nil {
		t.Fatal(err)
	}

	if addedU.ExternalID != u.ExternalID {
		t.Fatalf("Expected %#v got %#v", u, addedU)
	}
}

func TestUserRepositorySQL_Update(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	ur := data.NewUserRepositorySQL(
		db,
		log.New(os.Stdout, "test", log.LstdFlags),
	)
	defer db.Close()

	data.CreateUsersTable(db)

	db.MustExec(`INSERT INTO users (external_id) VALUES ('NEW-USER')`)

	toUpdate := &data.User{
		ID:         1,
		ExternalID: "NEW-USER",
	}

	update := &data.User{
		ID:         1,
		ExternalID: "UPDATED-USER",
	}

	err := ur.Update(toUpdate, update)

	if err != nil {
		t.Fatal(err)
	}

	got := &data.User{}

	err = db.Get(got, `SELECT * FROM users`)

	if err != nil {
		t.Fatal(err)
	}

	if got.ExternalID != update.ExternalID {
		t.Fatalf(
			"expected %#v got %#v",
			update.ExternalID,
			got.ExternalID,
		)
	}
}

func TestUserRepositorySQL_Delete(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	ur := data.NewUserRepositorySQL(
		db,
		log.New(os.Stdout, "test", log.LstdFlags),
	)
	defer db.Close()

	data.CreateUsersTable(db)

	db.MustExec(`INSERT INTO users (external_id) VALUES ('NEW-USER')`)

	u := &data.User{
		ID:         1,
		ExternalID: "TEST-USER",
	}

	err := ur.Delete(u)

	if err != nil {
		t.Fatal(err)
	}

	got := &data.User{}

	err = db.Get(got, `SELECT * FROM users`)

	if err != sql.ErrNoRows {
		t.Fatal(err)
	}
}

func TestUserRepositorySQL_FindBy(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	ur := data.NewUserRepositorySQL(
		db,
		log.New(os.Stdout, "test", log.LstdFlags),
	)
	defer db.Close()

	data.CreateUsersTable(db)

	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-ONE')`)
	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-TWO')`)
	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-THREE')`)

	u, err := ur.FindBy("external_id", "USER-TWO")

	if err != nil {
		t.Fatal(err)
	}

	if u.ExternalID != "USER-TWO" {
		t.Fatalf("expected USER-TWO got: %s", u.ExternalID)
	}
}

func TestUserRepositorySQL_IsFollowingUser(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	ur := data.NewUserRepositorySQL(
		db,
		log.New(os.Stdout, "test", log.LstdFlags),
	)
	defer db.Close()

	data.CreateUsersTable(db)
	data.CreateUsersFollowersTable(db)

	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-ONE')`)
	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-TWO')`)
	db.MustExec(
		`
		INSERT
			INTO users_followers (followee_id, follower_id)
			VALUES ('1', '2')
		`,
	)

	userOne := &data.User{
		ID:         1,
		ExternalID: "USER-ONE",
	}
	userTwo := &data.User{
		ID:         2,
		ExternalID: "USER-TWO",
	}

	isFollowing, err := ur.IsFollowingUser(userTwo, userOne)

	if err != nil {
		t.Fatal(err)
	}

	if !isFollowing {
		t.Fatalf("user %#v is not following user %#v", userTwo, userOne)
	}

	isFollowing, err = ur.IsFollowingUser(userOne, userTwo)

	if err != nil {
		t.Fatal(err)
	}

	if isFollowing {
		t.Fatalf("user %#v is following user %#v", userTwo, userOne)
	}
}

func TestUserRepositorySQL_FollowUser(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	ur := data.NewUserRepositorySQL(
		db,
		log.New(os.Stdout, "test", log.LstdFlags),
	)
	defer db.Close()

	data.CreateUsersTable(db)
	data.CreateUsersFollowersTable(db)

	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-ONE')`)
	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-TWO')`)

	userOne := &data.User{
		ID:         1,
		ExternalID: "USER-ONE",
	}
	userTwo := &data.User{
		ID:         2,
		ExternalID: "USER-TWO",
	}

	err := ur.FollowUser(userTwo, userOne)

	if err != nil {
		t.Fatal(err)
	}

	isFollowing, err := ur.IsFollowingUser(userTwo, userOne)

	if err != nil {
		t.Fatal(err)
	}

	if !isFollowing {
		t.Fatalf("user %#v is not following user %#v", userTwo, userOne)
	}
}

func TestUserRepositorySQL_UnfollowUser(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	ur := data.NewUserRepositorySQL(
		db,
		log.New(os.Stdout, "test", log.LstdFlags),
	)
	defer db.Close()

	data.CreateUsersTable(db)
	data.CreateUsersFollowersTable(db)

	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-ONE')`)
	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-TWO')`)
	db.MustExec(
		`
		INSERT
			INTO users_followers (followee_id, follower_id)
			VALUES ('1', '2')
		`,
	)

	userOne := &data.User{
		ID:         1,
		ExternalID: "USER-ONE",
	}
	userTwo := &data.User{
		ID:         2,
		ExternalID: "USER-TWO",
	}

	err := ur.UnfollowUser(userTwo, userOne)

	if err != nil {
		t.Fatal(err)
	}

	isFollowing, err := ur.IsFollowingUser(userTwo, userOne)

	if err != nil {
		t.Fatal(err)
	}

	if isFollowing {
		t.Fatalf("user %#v is following user %#v", userTwo, userOne)
	}
}

func TestUserRepositorySQL_GetUserFollowers(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	ur := data.NewUserRepositorySQL(
		db,
		log.New(os.Stdout, "test", log.LstdFlags),
	)
	defer db.Close()

	data.CreateUsersTable(db)
	data.CreateUsersFollowersTable(db)

	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-ONE')`)
	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-TWO')`)
	db.MustExec(
		`
		INSERT
			INTO users_followers (followee_id, follower_id)
			VALUES ('1', '2')
		`,
	)

	userOne := &data.User{
		ID:         1,
		ExternalID: "USER-ONE",
	}
	userTwo := &data.User{
		ID:         2,
		ExternalID: "USER-TWO",
	}

	t.Run("test not empty list of followers", func(t *testing.T) {
		followers, err := ur.GetUserFollowers(userOne)

		if err != nil {
			t.Fatal(err)
		}

		if len(followers) != 1 {
			t.Fatalf(
				"amount of followers not expected, expected 1 got %d",
				len(followers),
			)
		}

		if followers[0].ID != 2 {
			t.Fatalf(
				"expected ID: 2 as follower, got %#v instead",
				followers,
			)
		}
	})

	t.Run("test empty list of followers", func(t *testing.T) {
		followers, err := ur.GetUserFollowers(userTwo)

		if err != nil {
			t.Fatal(err)
		}

		if len(followers) != 0 {
			t.Fatalf(
				"expected empty list instead got %d followers",
				len(followers),
			)
		}
	})
}

func TestUserRepositorySQL_GetUserFollowees(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	ur := data.NewUserRepositorySQL(
		db,
		log.New(os.Stdout, "test", log.LstdFlags),
	)
	defer db.Close()

	data.CreateUsersTable(db)
	data.CreateUsersFollowersTable(db)

	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-ONE')`)
	db.MustExec(`INSERT INTO users (external_id) VALUES ('USER-TWO')`)
	db.MustExec(
		`
		INSERT
			INTO users_followers (followee_id, follower_id)
			VALUES ('1', '2')
		`,
	)

	userOne := &data.User{
		ID:         1,
		ExternalID: "USER-ONE",
	}
	userTwo := &data.User{
		ID:         2,
		ExternalID: "USER-TWO",
	}

	t.Run("test not empty list of followees", func(t *testing.T) {
		followees, err := ur.GetUserFollowees(userTwo)

		if err != nil {
			t.Fatal(err)
		}

		if len(followees) != 1 {
			t.Fatalf(
				"amount of followees not expected, expected 1 got %d",
				len(followees),
			)
		}

		if followees[0].ID != 1 {
			t.Fatalf(
				"expected ID: 1 as follower, got %#v instead",
				followees,
			)
		}
	})

	t.Run("test empty list of followees", func(t *testing.T) {
		followees, err := ur.GetUserFollowees(userOne)

		if err != nil {
			t.Fatal(err)
		}

		if len(followees) != 0 {
			t.Fatalf(
				"expected empty list instead got %d followees",
				len(followees),
			)
		}
	})
}
