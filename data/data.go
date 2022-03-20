package data

import (
	"encoding/json"
	"io"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	IsFollowerSQL = `
	SELECT
		*
	FROM Followers
	WHERE follower_id = ?
		AND followed_id = ?
	`

	FollowUserSQL = `
	INSERT
		INTO followers (
			follower_id,
			followed_id
		)
		VALUES (?, ?)`

	UnfollowUserSQL = `
	DELETE
		FROM followers
		WHERE
			follower_id = ?
			AND followed_id = ?
	`

	GetFollowersSQL = `SELECT follower_id FROM followers WHERE followed_id = ?`

	FindUserSQL = `SELECT id from users WHERE id = ?`
)

// ToJson seriealizes the given interface into a string based JSON format
func ToJson(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}

// Store is a register of all implemented Stores.
type Store struct {
	l *log.Logger
	*sqlx.DB
	*UserStore
	*RelationStore
}

func NewStore(driverName string, dataSrcName string, l *log.Logger) (*Store, error) {
	db, err := sqlx.Open(driverName, dataSrcName)

	if err != nil {
		return nil, ErrorOpeningDb
	}

	if err := db.Ping(); err != nil {
		return nil, ErrorCouldntConnectDb
	}

	store := &Store{
		l:             l,
		DB:            db,
		UserStore:     NewUserStore(db, l),
		RelationStore: NewRelationStore(db, l),
	}

	return store, nil
}
