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

// Create adds a Post to the database. It returns the created Post with
// fields like ID and DateCreated populated.
func (p Post) Create(ctx context.Context, traceID string, claims auth.Claims, np NewPost, now time.Time) (Info, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.create")
	defer span.End()

	post := Info{
		ID:         uuid.New().String(),
		//Author:     Author{Username: claims.User.Username, ID: claims.User.ID},
		Type:		np.Type,
		Title:		np.Title,
		URL:		np.URL,
		Category:	np.Category,
		Text:		np.Text,
		// Votes:		[]Votes{},
		// Comments:	[]Comment{},
		DateCreated: now.UTC(),
	}

	//const q = `
	//INSERT INTO posts
	//	(score, views, type, title, url, author, category, votes, comments, date_created, text, post_id)
	//VALUES
	//	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	//
	//p.log.Printf("%s: %s: %s", traceID, "post.Create",
	//	database.Log(q, post.Score, post.Views, post.Type, post.Title, post.URL, post.Author, post.Category,
	//		post.Votes, post.Comments, post.DateCreated, post.Text, post.ID),
	//)
	//
	//if _, err := p.db.ExecContext(ctx, q, post.Score, post.Views, post.Type, post.Title, post.URL, post.Author,
	//	post.Category, post.Votes, post.Comments, post.DateCreated, post.Text, post.ID); err != nil {
	//	return Info{}, errors.Wrap(err, "inserting product")
	//}

	const q = `
	INSERT INTO posts
		(post_id, score, views, type, title, url, category, text, date_created)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	p.log.Printf("%s: %s: %s", traceID, "post.Create",
		database.Log(q, post.ID, post.Score, post.Views, post.Type, post.Title, post.URL, post.Category,
			post.Text, post.DateCreated),
	)

	if _, err := p.db.ExecContext(ctx, q, post.ID, post.Score, post.Views, post.Type, post.Title, post.URL,
		post.Category, post.Text, post.DateCreated); err != nil {
		return Info{}, errors.Wrap(err, "inserting product")
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


	//// If you are not an admin. // Todo: check if auth.Claims.Subjects == post.Author.ID
	//if !claims.Authorized(auth.RoleAdmin) {
	//	return ErrForbidden
	//}

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

	const q = `SELECT * FROM posts`

	p.log.Printf("%s: %s: %s", traceID, "post.Query",
		database.Log(q),
	)

	var posts []Info
	if err := p.db.SelectContext(ctx, &posts, q); err != nil {
		return nil, errors.Wrap(err, "selecting posts")
	}

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

// // QueryByUser finds the post identified by a given user. // Todo: think about how to get post by user from DB
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
