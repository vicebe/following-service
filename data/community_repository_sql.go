package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Community is a representation of a community in the database.
type Community struct {
	ID int64 `db:"id" json:"id"`

	// ExternalID is the unique identifier for the community, it is not specific to this
	// service but used universally across services to identify the community.
	ExternalID string `db:"external_id" json:"external_id"`
}

type CommunityRepositorySQL struct {
	sq *SqlQuerent
	l  *log.Logger
}

func NewCommunityRepositorySQL(
	l *log.Logger,
	db *sqlx.DB,
) *CommunityRepositorySQL {
	return &CommunityRepositorySQL{
		l:  l,
		sq: NewSqlQuerent(db, l),
	}
}

func (cr *CommunityRepositorySQL) Create(community *Community) error {

	return cr.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)
		const CreateCommunitySQL = `INSERT INTO communities (external_id) VALUES (?)`

		if _, err := tx.Exec(
			CreateCommunitySQL,
			community.ExternalID,
		); err != nil {
			return err
		}

		return nil
	})
}

func (cr *CommunityRepositorySQL) Update(
	community *Community,
	newCommunity *Community,
) error {

	return cr.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)
		const UpdateCommunitySQL = `UPDATE communities SET external_id = ? WHERE external_id = ?`

		if _, err := tx.Exec(
			UpdateCommunitySQL,
			newCommunity.ExternalID,
			community.ExternalID,
		); err != nil {
			return err
		}

		return nil
	})
}

func (cr *CommunityRepositorySQL) Delete(community *Community) error {
	return cr.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)
		DeleteCommunitySQL := `DELETE FROM communities WHERE id = ?`

		if _, err := tx.Exec(
			DeleteCommunitySQL,
			community.ID,
		); err != nil {
			return err
		}

		return nil
	})
}

func (cr *CommunityRepositorySQL) FindBy(
	col string,
	value interface{},
) (*Community, error) {
	FindCommunitySQL := fmt.Sprintf("SELECT * FROM communities WHERE %s = ?", col)
	c := &Community{}
	err := cr.sq.Get(c, FindCommunitySQL, value)

	switch err {
	case nil:
		return c, nil
	case sql.ErrNoRows:
		return nil, ErrorCommunityNotFound
	default:
		return nil, err
	}
}

func (cr *CommunityRepositorySQL) IsFollowingCommunity(
	community *Community,
	user *User,
) (bool, error) {
	const IsFollowingCommunitySQL = `
	SELECT
		1
	FROM communities_followers
	WHERE community_id = ?
		AND follower_id = ?
	`
	var exists int
	err := cr.sq.Get(&exists, IsFollowingCommunitySQL, community.ID, user.ID)

	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (cr *CommunityRepositorySQL) FollowCommunity(
	community *Community,
	user *User,
) error {

	isFollowing, err := cr.IsFollowingCommunity(community, user)

	if err != nil {
		return err
	}

	if isFollowing {
		return nil
	}

	return cr.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)

		const FollowCommunitySQL = `
			INSERT
				INTO communities_followers (community_id, follower_id)
				VALUES (?, ?)
		`

		if _, err := tx.Exec(
			FollowCommunitySQL,
			community.ID,
			user.ID,
		); err != nil {
			return err
		}

		return nil
	})
}

func (cr *CommunityRepositorySQL) UnfollowCommunity(
	community *Community,
	user *User,
) error {

	isFollowing, err := cr.IsFollowingCommunity(community, user)

	if err != nil {
		return err
	}

	if !isFollowing {
		return nil
	}

	return cr.sq.DoTransaction(func(ctx context.Context) error {
		tx := ctx.Value("tx").(*sqlx.Tx)

		const UnfollowCommunitySQL = `
			DELETE
				FROM communities_followers
				WHERE community_id = ?
				AND follower_id = ?
		`

		if _, err := tx.Exec(
			UnfollowCommunitySQL,
			community.ID,
			user.ID,
		); err != nil {
			return err
		}

		return nil
	})
}

func (cr *CommunityRepositorySQL) GetCommunityFollowers(
	community *Community,
) ([]User, error) {
	var followers []User

	const GetCommunityFollowersSQL = `
		SELECT *
		FROM users u
		WHERE
			EXISTS (
				SELECT 1
				FROM communities_followers cf
				WHERE community_id = ?
					AND follower_id = u.id
			)
	`

	err := cr.sq.Select(&followers, GetCommunityFollowersSQL, community.ID)

	if err != nil {
		cr.l.Print(err)
		return nil, err
	}

	if followers == nil {
		followers = []User{}
	}

	return followers, nil
}
