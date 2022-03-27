package data

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

// Community represents a row in the users table
type Community struct {

	// ID is the unique identifier for the community, it is not specific to this
	// service but used universally across services to identify the community.
	ID string `db:"id" json:"id"`
}

// CommunityStore is a store with a db connection for Community related
// operations.
type CommunityStore struct {
	*sqlx.DB
	l *log.Logger
}

func NewCommunityStore(db *sqlx.DB, l *log.Logger) *CommunityStore {
	return &CommunityStore{
		l:  l,
		DB: db,
	}
}

// CommunityExists verifies if user exists in the database given its id.
func (cs *CommunityStore) CommunityExists(uId string) (bool, error) {
	var u string
	err := cs.Get(&u, FindCommunitySQL, uId)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}

}
