package db

import (
	"time"
)

// postDB represents an individual post in database.
type PostDB struct {
	ID          string    `db:"post_id"`
	Views       int       `db:"views"`
	Type        string    `db:"type"`
	Title       string    `db:"title"`
	Category    string    `db:"category"`
	Payload     string    `db:"payload"`
	DateCreated time.Time `db:"date_created"`
	UserID      string    `db:"user_id"`
}

// UserDB represents author in database
type UserDB struct {
	Username string `db:"name"`
	ID       string `db:"user_id"`
}

// VoteDB represents info about user vote.
type VoteDB struct {
	User string `db:"user_id"`
	Vote int    `db:"vote"`
}

// CommentDB represents an individual comment in database
type CommentDB struct {
	DateCreated time.Time `db:"date_created"`
	PostID      string    `db:"post_id"`
	AuthorID    string    `db:"user_id"`
	Body        string    `db:"body"`
	ID          string    `db:"comment_id"`
}

// CommentDB represents an individual comment in database
type CommentWithUserDB struct {
	AuthorName  string    `db:"name"`
	DateCreated time.Time `db:"date_created"`
	PostID      string    `db:"post_id"`
	AuthorID    string    `db:"user_id"`
	Body        string    `db:"body"`
	ID          string    `db:"comment_id"`
}