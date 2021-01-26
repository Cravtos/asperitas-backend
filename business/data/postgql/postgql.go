package postgql

import (
	"errors"
	"github.com/cravtos/asperitas-backend/graph"
	"github.com/jmoiron/sqlx"
	"log"
)

var (
	// ErrInvalidPostID occurs when an ID is not in a valid form.
	ErrInvalidPostID = graph.newPublicError(errors.New("invalid postRes id"))

	// ErrForbidden occurs when a users tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = graph.newPublicError(errors.New("attempted action is not allowed"))

	// ErrInvalidCommentID occurs when an ID is not in a valid form.
	ErrInvalidCommentID = graph.newPublicError(errors.New("invalid comment id"))
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyPostGQL is used to store/retrieve a Claims value from a context.Context.
const KeyPostGQL ctxKey = 1
const KeyAuthHeader ctxKey = 2
const KeyAuth ctxKey = 3

type PostGQL struct {
	log *log.Logger
	db  *sqlx.DB
}

func NewPostGQL(log *log.Logger, db *sqlx.DB) PostGQL {
	return PostGQL{
		log: log,
		db:  db,
	}
}
