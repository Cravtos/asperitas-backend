package postgql

import (
	"context"
	"database/sql"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

//todo take common helpers out of postgql and post

var (
	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("invalid post id")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")

	// ErrPostNotFound is used when a specific Post is requested but does not exist.
	ErrPostNotFound = errors.New("post not found")

	//ErrCommentNotFound is used when a specific Comment is requested but does not exist
	ErrCommentNotFound = errors.New("comment not found")

	//ErrCommentNotFound is used when user tries to create post with incorrect type.
	ErrWrongPostType = errors.New("new post should be of type url or text")
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// Key is used to store/retrieve a Claims value from a context.Context.
const Key ctxKey = 1

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

// postDB represents an individual post in database. (with additional field "score" counted using votes table)
type postDB struct {
	ID          string    `db:"post_id"`
	Views       int       `db:"views"`
	Type        string    `db:"type"`
	Title       string    `db:"title"`
	Category    string    `db:"category"`
	Payload     string    `db:"payload"`
	DateCreated time.Time `db:"date_created"`
	UserID      string    `db:"user_id"`
	Author      Author
	Votes       []Vote
	Comments    []Comment
}

// Author represents info about authorType
type Author struct {
	Username string `db:"name"`
	ID       string `db:"user_id"`
}

// Vote represents info about user vote.
type Vote struct {
	UserID string `db:"user_id" `
	Vote   int    `db:"vote"`
}

// Comment represents info about comments for the post prepared to be sent to user.
type Comment struct {
	PostID      string
	DateCreated time.Time
	Author      Author
	Body        string
	ID          string
}

// getPostScore returns score of a single post
func (p PostGQL) getPostScore(ctx context.Context, ID string) (int, error) {
	const qScore = `SELECT SUM(vote) as score FROM votes WHERE post_id = $1 HAVING SUM(vote) is not null`

	p.log.Printf("%s: %s", "post.helpers.getPostScore", database.Log(qScore))

	var score int
	if err := p.db.GetContext(ctx, &score, qScore, ID); err != nil {
		return 0, nil
	}
	return score, nil
}

func (p PostGQL) obtainPosts(ctx context.Context, category interface{}) ([]postDB, error) {
	if category == nil || category == 1 {
		return p.selectAllPosts(ctx)
	}
	_, ok := category.(string)
	if !ok {
		return nil, errors.Errorf("got wrong category %v\n", category)
	}
	return p.selectPostsByCategory(ctx, category.(string))
}

// selectPostsByCategory returns all posts with a given category stored in database
func (p PostGQL) selectPostsByCategory(ctx context.Context, category string) ([]postDB, error) {
	const qPost = `SELECT * FROM posts WHERE category = $1`

	p.log.Printf("%s: %s", "post.helpers.selectPostsByCategory", database.Log(qPost))

	var posts []postDB
	if err := p.db.SelectContext(ctx, &posts, qPost, category); err != nil {
		return nil, errors.Wrap(err, "selecting category posts")
	}
	return posts, nil
}

// selectAllPosts return all posts stored in database
func (p PostGQL) selectAllPosts(ctx context.Context) ([]postDB, error) {
	const qPost = `SELECT * FROM posts`

	p.log.Printf("%s: %s", "post.helpers.selectAllPosts", database.Log(qPost))

	var posts []postDB
	if err := p.db.SelectContext(ctx, &posts, qPost); err != nil {
		return nil, errors.Wrap(err, "selecting all posts")
	}

	return posts, nil
}

// getAuthorByID obtains Author using ID from database
func (p PostGQL) getAuthorByID(ctx context.Context, ID string) (Author, error) {
	const qAuthor = `SELECT user_id, name FROM users WHERE user_id = $1`

	p.log.Printf("%s: %s", "post.helpers.getAuthorByID", database.Log(qAuthor))

	var author Author
	if err := p.db.GetContext(ctx, &author, qAuthor, ID); err != nil {
		return Author{}, errors.Wrap(err, "selecting authors")
	}
	return author, nil
}

// selectVotesByPostID returns slice of Vote for a single post
func (p PostGQL) selectVotesByPostID(ctx context.Context, ID string) ([]Vote, error) {
	const qVotes = `SELECT user_id, vote FROM votes WHERE post_id = $1`

	p.log.Printf("%s: %s", "post.helpers.selectVotesByPostID", database.Log(qVotes))

	var votes []Vote
	if err := p.db.SelectContext(ctx, &votes, qVotes, ID); err != nil {
		return nil, errors.Wrap(err, "selecting votes")
	}
	return votes, nil
}

// selectCommentsByPostID return slice of Comment for a single post
func (p PostGQL) selectCommentsByPostID(ctx context.Context, ID string) ([]Comment, error) {
	const qComments = `
		SELECT 
			name, user_id, cm.date_created, body, comment_id 
		FROM 
			comments cm join users using(user_id) 
		WHERE 
			post_id = $1`

	p.log.Printf("%s: %s", "post.helpers.selectCommentsByPostID", database.Log(qComments))

	var rawComments []struct {
		DateCreated time.Time `db:"date_created"`
		AuthorName  string    `db:"name"`
		AuthorID    string    `db:"user_id"`
		Body        string    `db:"body"`
		ID          string    `db:"comment_id"`
	}
	if err := p.db.SelectContext(ctx, &rawComments, qComments, ID); err != nil {
		return nil, errors.Wrap(err, "selecting comments")
	}

	comments := make([]Comment, 0)
	for _, comment := range rawComments {
		author := Author{
			Username: comment.AuthorName,
			ID:       comment.AuthorID,
		}
		comments = append(comments, Comment{
			PostID:      ID,
			DateCreated: comment.DateCreated,
			Author:      author,
			Body:        comment.Body,
			ID:          comment.ID,
		})
	}
	return comments, nil
}

func (p PostGQL) selectPostsByUser(ctx context.Context, userID string) ([]postDB, error) {
	const qPost = `SELECT * FROM posts WHERE user_id = $1`

	p.log.Printf("%s: %s", "post.helpers.selectPostsByUser", database.Log(qPost))

	var posts []postDB
	if err := p.db.SelectContext(ctx, &posts, qPost, userID); err != nil {
		return nil, errors.Wrap(err, "selecting users posts")
	}

	return posts, nil
}

// getPostByID obtains post from database using ID
func (p PostGQL) getPostByID(ctx context.Context, postID string) (postDB, error) {
	const q = `	SELECT * FROM posts WHERE post_id = $1`

	p.log.Printf("%s: %s", "post.helpers.getPostByID", database.Log(q, postID))

	var post postDB
	if err := p.db.GetContext(ctx, &post, q, postID); err != nil {
		if err == sql.ErrNoRows {
			return postDB{}, ErrPostNotFound
		}
		return postDB{}, errors.Wrap(err, "selecting post by ID")
	}
	return post, nil
}
