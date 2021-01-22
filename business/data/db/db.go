package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
)

var (
	// ErrPostNotFound is used when a specific Post is requested but does not exist.
	ErrPostNotFound = errors.New("posts not found")

	//ErrCommentNotFound is used when a specific Comment is requested but does not exist
	ErrCommentNotFound = errors.New("comment not found")

	//ErrVoteNotFound is used when a specific Vote is requested but does not exist
	ErrVoteNotFound = errors.New("vote not found")

	//ErrUserNotFound is used when a specific UserID is requested but does not exist
	ErrUserNotFound = errors.New("users not found")
)

type DBset struct {
	log *log.Logger
	db  *sqlx.DB
}

func NewDBset(log *log.Logger, db *sqlx.DB) DBset {
	return DBset{log: log, db: db}
}
