package data

import "context"

// Querent is an object that interacts with the database.
type Querent interface {

	// DoTransaction makes applies a transaction to the database.
	DoTransaction(f func(ctx context.Context) error) error
}
