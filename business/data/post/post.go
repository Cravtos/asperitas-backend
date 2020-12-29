// Package post contains product related CRUD functionality.
package post

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

var (
	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")

	// ErrNotFound is used when a specific Post is requested but does not exist.
	ErrNotFound = errors.New("not found")
)

// Post manages the set of API's for product access.
type Post struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a Post for api access.
func New(log *log.Logger, db *sqlx.DB) Post {
	return Post{
		log: log,
		db:  db,
	}
}

// Create adds a Post to the database. It returns the created Post with fields like ID and DateCreated populated.
func (p Post) Create(ctx context.Context, traceID string, claims auth.Claims, np NewPost, now time.Time) (Info, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.create")
	defer span.End()

	// todo: split responses into two (InfoText, InfoLink or smth like that)
	post := Info{
		ID:         uuid.New().String(),
		Author:     Author{Username: claims.User.Username, ID: claims.User.ID},
		Type:		np.Type,
		Title:		np.Title,
		URL:		np.URL,
		Category:	np.Category,
		Text:		np.Text,
		Votes:		[]Vote{{claims.User.Username, 1}},
		Comments:	[]Comment{},
		DateCreated: now.UTC(),
		UpvotePercentage: 100,
	}

	const qPost = `
	INSERT INTO posts
		(post_id, score, views, type, title, url, category, text, date_created, user_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	p.log.Printf("%s: %s: %s", traceID, "post.Create",
		database.Log(qPost, post.ID, post.Score, post.Views, post.Type, post.Title, post.URL, post.Category,
			post.Text, post.DateCreated, post.Author.ID),
	)

	if _, err := p.db.ExecContext(ctx, qPost, post.ID, post.Score, post.Views, post.Type, post.Title,
		post.URL, post.Category, post.Text, post.DateCreated, post.Author.ID); err != nil {
		return Info{}, errors.Wrap(err, "creating post")
	}

	const qVote = `
	INSERT INTO votes
		(post_id, user_id, vote)
	VALUES
		($1, $2, $3)`

	p.log.Printf("%s: %s: %s", traceID, "post.Create",
		database.Log(qVote, post.ID, post.Author.ID, 1),
	)

	if _, err := p.db.ExecContext(ctx, qVote, post.ID, post.Author.ID, 1); err != nil {
		return Info{}, errors.Wrap(err, "upvote created post")
	}

	return post, nil
}

// Delete removes the product identified by a given ID.
func (p Post) Delete(ctx context.Context, traceID string, claims auth.Claims, postID string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.delete")
	defer span.End()

	if _, err := uuid.Parse(postID); err != nil {
		return ErrInvalidID
	}

	post, err := p.QueryByID(ctx, traceID, postID)
	if err != nil {
		errors.Wrap(err, "selecting post by ID")
	}

	if claims.User.ID == post.Author.ID {
		return ErrForbidden
	}

	const q = `
	DELETE FROM
		posts
	WHERE
		post_id = $1`

	p.log.Printf("%s: %s: %s", traceID, "post.Delete",
		database.Log(q, postID),
	)

	if _, err := p.db.ExecContext(ctx, q, postID); err != nil {
		return errors.Wrapf(err, "deleting post %s", postID)
	}

	return nil
}

// Query gets all Posts from the database.
func (p Post) Query(ctx context.Context, traceID string) ([]Info, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.query")
	defer span.End()

	// Get all posts
	const qPost = `SELECT * FROM posts`

	p.log.Printf("%s: %s: %s", traceID, "post.Query",
		database.Log(qPost),
	)

	// todo: uncomment and finish
	//var posts []PostDB
	//if err := p.db.SelectContext(ctx, &posts, qPost); err != nil {
	//	return nil, errors.Wrap(err, "selecting posts")
	//}
	//
	//var info []Info
	//for _, post := range posts {
	//	// todo: divide into functions
	//	// Get posts votes
	//	const qVotes = `SELECT user_id, vote FROM votes WHERE post_id = $1`
	//
	//	p.log.Printf("%s: %s: %s", traceID, "post.Query",
	//		database.Log(qPost),
	//	)
	//
	//	var votes []Vote
	//	if err := p.db.SelectContext(ctx, &votes, qVotes, post.ID); err != nil {
	//		return nil, errors.Wrap(err, "selecting votes")
	//	}
	//
	//	/* todo: think about how to get Authors from the comments
	//	   (1. one more db requests for every comment or keep two values in db about author)
	//	   (2. write your own marshalling function)
	//	*/
	//	 */
	//	// Get their comments
	//	const qComments = `SELECT * FROM comments WHERE post_id = $1`
	//
	//	p.log.Printf("%s: %s: %s", traceID, "post.Query",
	//		database.Log(qPost),
	//	)
	//
	//	var comments []CommentDB
	//	if err := p.db.SelectContext(ctx, &comments, qVotes, post.ID); err != nil {
	//		return nil, errors.Wrap(err, "selecting votes")
	//	}
	//
	//	// Get their authors
	//
	//}



	return posts, nil
}

// QueryByID finds the post identified by a given ID.
func (p Post) QueryByID(ctx context.Context, traceID string, postID string) (Info, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.querybyid")
	defer span.End()

	if _, err := uuid.Parse(postID); err != nil {
		return Info{}, ErrInvalidID
	}

	const q = `
	SELECT * FROM
		posts
	WHERE
		post_id = $1`

	p.log.Printf("%s: %s: %s", traceID, "product.QueryByID",
		database.Log(q, postID),
	)

	var post Info
	if err := p.db.GetContext(ctx, &post, q, postID); err != nil {
		return Info{}, errors.Wrap(err, "selecting post by ID")
	}

	return post, nil
}

// QueryByCat finds the post identified by a given Category.
func (p Post) QueryByCat(ctx context.Context, traceID string, category string) (Info, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.querybycat")
	defer span.End()

	const q = `
	SELECT * FROM
		posts
	WHERE
		category = $1`

	p.log.Printf("%s: %s: %s", traceID, "product.QueryByCat",
		database.Log(q, category),
	)

	var post Info
	if err := p.db.GetContext(ctx, &post, q, category); err != nil {
		if err == sql.ErrNoRows { // Todo: check if error is correct
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrap(err, "selecting posts by category")
	}

	return post, nil
}

// // QueryByUser finds the post identified by a given user. // Todo: think about how to get post by user from PostDB
//func (p Post) QueryByUser(ctx context.Context, traceID string, user string) (Info, error) {
//	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.querybyuser")
//	defer span.End()
//
//	const q = `
//	SELECT * FROM
//		posts
//	WHERE
//		post_id = $1`
//
//	p.log.Printf("%s: %s: %s", traceID, "product.QueryByID",
//		database.Log(q, postID),
//	)
//
//	var prd Info
//	if err := p.db.GetContext(ctx, &prd, q, postID); err != nil {
//		return Info{}, errors.Wrap(err, "selecting single product")
//	}
//
//	return prd, nil
//}
