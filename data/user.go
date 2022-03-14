package data

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// User represents a row in the users table
type User struct {

	// ID is the unique identifier for the user, it is not specific to this
	// service but used universally across services to identified the user, in
	// other words the user's UUID.
	ID string `db:"id" json:"id"`
}

// UserStore is a store with a db connection for User related operations
type UserStore struct {
	*sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

// UserExists verifies if user exists in the database given its id
func (us *UserStore) UserExists(uId string) (bool, error) {
	var u string
	err := us.Get(&u, FindUserSQL, uId)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}

}
