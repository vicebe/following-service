package data

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"testing"
)

func TestSqlQuerent_DoTransaction(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	l := log.New(os.Stdout, "test", log.LstdFlags)
	type fields struct {
		DB *sqlx.DB
		l  *log.Logger
	}
	type args struct {
		f func(ctx context.Context) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test error not thrown",
			fields: fields{
				DB: db,
				l:  l,
			},
			args: args{f: func(ctx context.Context) error {
				tx, ok := ctx.Value("tx").(*sqlx.Tx)

				if !ok {
					return fmt.Errorf("couldn't get tx object")
				}

				if _, err := tx.Exec(
					`CREATE TABLE test (id INTEGER)`,
				); err != nil {
					return err
				}

				return nil
			}},
			wantErr: false,
		},
		{
			name: "test function error returns",
			fields: fields{
				DB: db,
				l:  l,
			},
			args: args{
				f: func(ctx context.Context) error {
					return fmt.Errorf("Error returned!")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SqlQuerent{
				DB: tt.fields.DB,
				l:  tt.fields.l,
			}
			if err := s.DoTransaction(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("DoTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
