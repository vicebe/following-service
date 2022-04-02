package data

import "context"

type Querent interface {
	DoTransaction(f func(ctx context.Context) error) error
}
