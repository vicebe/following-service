package data

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type SqlQuerent struct {
	*sqlx.DB
	l *log.Logger
}

func NewSqlQuerent(db *sqlx.DB, l *log.Logger) *SqlQuerent {
	return &SqlQuerent{DB: db, l: l}
}

func (s *SqlQuerent) DoTransaction(f func(ctx context.Context) error) (retErr error) {
	tx, err := s.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if errC := tx.Commit(); errC != nil {
			s.l.Print(errC)
			retErr = errC
			if errR := tx.Rollback(); errR != nil {
				s.l.Print(errR)
				retErr = errR
			}
		}
	}()

	txCtx := context.WithValue(context.Background(), "tx", tx)

	err = f(txCtx)

	if err != nil {
		s.l.Print(err)
		return err
	}

	return nil
}
