package gql

import (
	"context"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

//todo think about names for everything
//todo take common helpers out of gql and post

// ctxKey represents the type of value for the context key.
type ctxKey int

// Key is used to store/retrieve a Claims value from a context.Context.
const Key ctxKey = 1

type Access struct {
	log *log.Logger
	db  *sqlx.DB
}

func NewAccess(log *log.Logger, db *sqlx.DB) Access {
	return Access{
		log: log,
		db:  db,
	}
}

// postDB represents an individual post in database. (with additional field "score" counted using votes table)
type postDB struct {
	ID          string    `db:"post_id"`
	Score       int       `db:"score"`
	Views       int       `db:"views"`
	Type        string    `db:"type"`
	Title       string    `db:"title"`
	Category    string    `db:"category"`
	Payload     string    `db:"payload"`
	DateCreated time.Time `db:"date_created"`
	UserID      string    `db:"user_id"`
}

// Author represents info about author
type Author struct {
	Username string `db:"name" json:"username"`
	ID       string `db:"user_id" json:"id"`
}

// Vote represents info about user vote.
type Vote struct {
	User string `db:"user_id" json:"user"`
	Vote int    `db:"vote" json:"vote"`
}

// Comment represents info about comments for the post prepared to be sent to user.
type Comment struct {
	DateCreated time.Time `json:"created"`
	Author      Author    `json:"author"`
	Body        string    `json:"body"`
	ID          string    `json:"id"`
}

// getPostScore returns score of a single post
func (p Access) getPostScore(ctx context.Context, ID string) (int, error) {
	const qScore = `SELECT SUM(vote) as score FROM votes WHERE post_id = $1 HAVING SUM(vote) is not null`

	p.log.Printf("%s: %s", "post.helpers.getPostScore", database.Log(qScore))

	var score int
	if err := p.db.GetContext(ctx, &score, qScore, ID); err != nil {
		return 0, nil
	}
	return score, nil
}

// selectAllPosts return all posts stored in database
func (p Access) selectAllPosts(ctx context.Context) ([]postDB, error) {
	const qPost = `SELECT * FROM posts`

	p.log.Printf("%s: %s", "post.helpers.selectAllPosts", database.Log(qPost))

	var posts []postDB
	if err := p.db.SelectContext(ctx, &posts, qPost); err != nil {
		return nil, errors.Wrap(err, "selecting all posts")
	}

	for i := range posts {
		score, err := p.getPostScore(ctx, posts[i].ID)
		if err != nil {
			return nil, err
		}
		posts[i].Score = score
	}
	return posts, nil
}
