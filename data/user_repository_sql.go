package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// User represents a row in the users table
type User struct {
	ID uint64 `db:"id" json:"id"`

	// ExternalID is the unique identifier for the user, it is not specific to this
	// service but used universally across services to identify the user, in
	// other words the user's UUID.
	ExternalID string `db:"external_id" json:"external_id"`
}

// UserRepositorySQL is a store with a db connection for User related operations
type UserRepositorySQL struct {
	l  *log.Logger
	sq *SqlQuerent
}

func NewUserRepositorySQL(db *sqlx.DB, l *log.Logger) *UserRepositorySQL {
	return &UserRepositorySQL{
		l:  l,
		sq: NewSqlQuerent(db, l),
	}
}

func (u *UserRepositorySQL) Create(user *User) error {
	return u.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)
		const CreateUserSQL = "INSERT INTO users (external_id) VALUES (?)"

		if _, err := tx.Exec(CreateUserSQL, user.ExternalID); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserRepositorySQL) Update(user *User, update *User) error {
	return u.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)
		const UpdateUserSQL = `UPDATE users SET external_id = ? WHERE external_id = ?`

		if _, err := tx.Exec(
			UpdateUserSQL,
			update.ExternalID,
			user.ExternalID,
		); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserRepositorySQL) Delete(user *User) error {
	return u.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)
		DeleteUserSQL := `DELETE FROM users WHERE id = ?`

		if _, err := tx.Exec(DeleteUserSQL, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserRepositorySQL) FindBy(
	key string,
	value interface{},
) (*User, error) {
	FindUserSQL := fmt.Sprintf(
		"SELECT * FROM users WHERE %s = ?",
		key,
	)
	user := &User{}
	err := u.sq.Get(user, FindUserSQL, value)

	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return nil, ErrorCommunityNotFound
	default:
		return nil, err
	}
}

func (u *UserRepositorySQL) IsFollowingUser(
	follower *User,
	followee *User,
) (bool, error) {
	const IsFollowingUserSQL = `
	SELECT
		1
	FROM users_followers
	WHERE follower_id = ?
		AND followee_id = ?
	`
	var exists int
	err := u.sq.Get(&exists, IsFollowingUserSQL, follower.ID, followee.ID)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (u *UserRepositorySQL) FollowUser(user *User, followee *User) error {
	isFollowing, err := u.IsFollowingUser(user, followee)

	if err != nil {
		return err
	}

	if isFollowing {
		return nil
	}

	return u.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)

		const FollowUserSQL = `
			INSERT
				INTO users_followers (follower_id, followee_id)
				VALUES (?, ?)
		`
		_, err := tx.Exec(FollowUserSQL, user.ID, followee.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *UserRepositorySQL) UnfollowUser(user *User, followee *User) error {
	isFollowing, err := u.IsFollowingUser(user, followee)

	if err != nil {
		return err
	}

	if !isFollowing {
		return nil
	}

	return u.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)

		const UnfollowUserSQL = `
			DELETE
				FROM users_followers
				WHERE follower_id = ?
					AND followee_id = ?
		`
		_, err := tx.Exec(UnfollowUserSQL, user.ID, followee.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *UserRepositorySQL) GetUserFollowers(user *User) ([]User, error) {
	var followers []User

	const GetUserFollowersSQL = `
		SELECT *
		FROM users u
		WHERE
			EXISTS (
				SELECT 1
				FROM users_followers cf
				WHERE followee_id = ?
					AND followers_id = u.id
			)
	`

	err := u.sq.Select(&followers, GetUserFollowersSQL, user.ID)

	if err != nil {
		u.l.Print(err)
		return nil, err
	}

	if followers == nil {
		followers = []User{}
	}

	return followers, nil
}

func (u *UserRepositorySQL) GetUserFollowees(user *User) ([]User, error) {
	var followees []User

	const GetUserFolloweesSQL = `
		SELECT *
		FROM users u
		WHERE
			EXISTS (
				SELECT 1
				FROM users_followers cf
				WHERE follower_id = ?
					AND followee_id = u.id
			)
	`

	err := u.sq.Select(&followees, GetUserFolloweesSQL, user.ID)

	if err != nil {
		u.l.Print(err)
		return nil, err
	}

	if followees == nil {
		followees = []User{}
	}

	return followees, nil
}
