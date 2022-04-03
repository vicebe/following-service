package data_test

import (
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"sort"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
)

func TestCreateCommunity(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	cr := data.NewCommunityRepositorySQL(
		log.New(os.Stdout, "test", log.LstdFlags),
		db,
	)

	data.InitializeDB(db)

	c := &data.Community{
		ID:         2,
		ExternalID: "NEW-USER",
	}

	err := cr.Create(c)

	if err != nil {
		t.Fatal(err)
	}

	communityCreated, err := cr.FindBy("id", 2)

	if err != nil {
		t.Fatal(err)
	}

	if communityCreated.ExternalID != c.ExternalID {
		t.Fatalf(
			"Community expected:\n%#v\ngot:\n%#v",
			c,
			communityCreated,
		)
	}

}

func TestUpdateCommunity(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	cr := data.NewCommunityRepositorySQL(
		log.New(os.Stdout, "test", log.LstdFlags),
		db,
	)

	data.InitializeDB(db)

	update := &data.Community{
		ExternalID: "3",
	}

	c, err := cr.FindBy("id", 1)

	if err != nil {
		t.Fatal(err)
	}

	if c.ExternalID != "1" {
		t.Fatal("Community 1 does not exist")
	}

	err = cr.Update(c, update)

	if err != nil {
		t.Fatal(err)
	}

	c, err = cr.FindBy("id", 1)

	if c.ExternalID != update.ExternalID {
		t.Fatalf(
			"Update did not work exptected:\n%s\ngot:\n%s",
			update.ExternalID,
			c.ExternalID,
		)
	}
}

func TestDeleteCommunity(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	cr := data.NewCommunityRepositorySQL(
		log.New(os.Stdout, "test", log.LstdFlags),
		db,
	)

	data.InitializeDB(db)

	c, err := cr.FindBy("id", 1)

	if err != nil {
		t.Fatal(err)
	}

	if c.ID != 1 {
		t.Fatal("Community 1 does not exist")
	}

	err = cr.Delete(c)

	if err != nil {
		t.Fatal(err)
	}

	c, err = cr.FindBy("id", 1)

	if err != data.ErrorCommunityNotFound {
		log.Fatal(err)
	}
}

func TestFindByCommunity(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	cr := data.NewCommunityRepositorySQL(
		log.New(os.Stdout, "test", log.LstdFlags),
		db,
	)

	data.InitializeDB(db)

	c, err := cr.FindBy("id", 1)

	if err != nil {
		t.Fatal(err)
	}

	if c.ID != 1 {
		t.Fatal("Community 1 does not exist")
	}

	c, err = cr.FindBy("external_id", 1)

	if c.ID != 1 {
		t.Fatal("Community 1 does not exist")
	}
}

func TestIsFollowingCommunity(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	cr := data.NewCommunityRepositorySQL(
		log.New(os.Stdout, "test", log.LstdFlags),
		db,
	)

	data.InitializeDB(db)

	c := &data.Community{
		ID:         1,
		ExternalID: "1",
	}

	u := &data.User{
		ID:         1,
		ExternalID: "1",
	}

	isFollowing, err := cr.IsFollowingCommunity(c, u)

	if err != nil {
		t.Fatal(err)
	}

	if !isFollowing {
		t.Fatalf("User #%v is not following community #%v", u, c)
	}
}

func TestFollowCommunity(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	cr := data.NewCommunityRepositorySQL(
		log.New(os.Stdout, "test", log.LstdFlags),
		db,
	)

	data.InitializeDB(db)

	c := &data.Community{
		ID:         1,
		ExternalID: "1",
	}

	u := &data.User{
		ID:         3,
		ExternalID: "3",
	}

	isFollowing, err := cr.IsFollowingCommunity(c, u)

	if err != nil {
		t.Fatal(err)
	}

	if isFollowing {
		t.Fatalf("User #%v is following community #%v", u, c)
	}

	err = cr.FollowCommunity(c, u)

	if err != nil {
		t.Fatal(err)
	}

	isFollowing, err = cr.IsFollowingCommunity(c, u)

	if err != nil {
		t.Fatal(err)
	}

	if !isFollowing {
		t.Fatalf("User #%v is not following community #%v", u, c)
	}

}

func TestGetCommunityFollowers(t *testing.T) {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	cr := data.NewCommunityRepositorySQL(
		log.New(os.Stdout, "test", log.LstdFlags),
		db,
	)

	data.InitializeDB(db)

	c := &data.Community{
		ID:         1,
		ExternalID: "1",
	}

	us := []*data.User{
		{
			ID:         1,
			ExternalID: "1",
		},
		{
			ID:         2,
			ExternalID: "2",
		},
	}

	followers, err := cr.GetCommunityFollowers(c)

	if err != nil {
		t.Fatal(err)
	}

	if len(us) != len(followers) {
		t.Fatalf(
			"Followers count not expected. Expected %d got: %d",
			len(us),
			len(followers),
		)
	}

	sort.SliceStable(followers, func(i, j int) bool {
		return followers[i].ID < followers[j].ID
	})

	for i, user := range followers {
		if user.ID != us[i].ID {
			t.Fatalf(
				"Followers not the same. Expected:\n%#v\ngot:\n%#v",
				us[i],
				user,
			)
		}
	}
}
