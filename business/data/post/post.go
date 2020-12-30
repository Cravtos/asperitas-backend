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

// Create adds a post to the database. It returns the created post with fields like ID and DateCreated populated.
func (p Post) Create(ctx context.Context, traceID string, claims auth.Claims, np NewPost, now time.Time) (Info, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.create")
	defer span.End()

	// todo: split responses into two (InfoText, InfoLink or smth like that)
	post := struct {
		ID               string
		Author           Author
		Type             string
		Title            string
		Category         string
		Votes            []Vote
		Score            int
		Views            int
		Comments         []Comment
		DateCreated      time.Time
		UpvotePercentage int
		Payload          string
	}{
		uuid.New().String(),
		Author{Username: claims.User.Username, ID: claims.User.ID},
		np.Type,
		np.Title,
		np.Category,
		[]Vote{{claims.User.Username, 1}},
		1,
		1,
		[]Comment{},
		now.UTC(),
		100,
		np.Text,
	}
	if post.Type == "url" {
		post.Payload = np.URL
	}

	const qPost = `
	INSERT INTO posts
		(post_id, score, views, type, title, category, payload, date_created, user_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	p.log.Printf("%s: %s: %s", traceID, "post.Create",
		database.Log(qPost, post.ID, post.Score, post.Views, post.Type, post.Title, post.Payload, post.Category,
			post.DateCreated, post.Author.ID),
	)

	if _, err := p.db.ExecContext(ctx, qPost, post.ID, post.Score, post.Views, post.Type, post.Title,
		post.Category, post.Payload, post.DateCreated, post.Author.ID); err != nil {
		return InfoText{}, errors.Wrap(err, "creating post")
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
		return InfoText{}, errors.Wrap(err, "upvote created post")
	}

	if post.Type == "url" {
		info := InfoLink{
			ID:               post.ID,
			Score:            post.Score,
			Views:            post.Views,
			Title:            post.Title,
			Payload:          post.Payload,
			Category:         post.Category,
			DateCreated:      post.DateCreated,
			Author:           post.Author,
			Votes:            post.Votes,
			Comments:         post.Comments,
			UpvotePercentage: post.UpvotePercentage,
		}
		return info, nil
	} else {
		info := InfoText{
			ID:               post.ID,
			Score:            post.Score,
			Views:            post.Views,
			Title:            post.Title,
			Payload:          post.Payload,
			Category:         post.Category,
			DateCreated:      post.DateCreated,
			Author:           post.Author,
			Votes:            post.Votes,
			Comments:         post.Comments,
			UpvotePercentage: post.UpvotePercentage,
		}
		return info, nil
	}
}

// Delete removes the product identified by a given ID.
func (p Post) Delete(ctx context.Context, traceID string, claims auth.Claims, postID string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.delete")
	defer span.End()

	if _, err := uuid.Parse(postID); err != nil {
		return ErrInvalidID
	}

	const qPost = `SELECT user_id FROM posts WHERE post_id = $1`

	var author Author
	if err := p.db.GetContext(ctx, &author, qPost, postID); err != nil {
		return errors.Wrap(err, "selecting post by ID")
	}

	if claims.User.ID == author.ID {
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

	var posts []PostDB
	if err := p.db.SelectContext(ctx, &posts, qPost); err != nil {
		return nil, errors.Wrap(err, "selecting posts")
	}

	var info []Info
	for _, post := range posts {
		const qAuthor = `SELECT user_id, name FROM users WHERE user_id = $1`

		var author Author
		if err := p.db.GetContext(ctx, &author, qAuthor, post.UserID); err != nil {
			return nil, errors.Wrap(err, "selecting votes")
		}
		// todo: divide into functions
		// Get posts votes
		const qVotes = `SELECT user_id as User, vote as Vote FROM votes WHERE post_id = $1`

		p.log.Printf("%s: %s: %s", traceID, "post.Query",
			database.Log(qPost),
		)

		var votes []Vote
		if err := p.db.SelectContext(ctx, &votes, qVotes, post.ID); err != nil {
			return nil, errors.Wrap(err, "selecting votes")
		}

		// Get their comments
		const qComments = `SELECT name, user_id, cm.date_created, body, comment_id FROM comments cm join users using(user_id) WHERE post_id = $1`

		p.log.Printf("%s: %s: %s", traceID, "post.Query",
			database.Log(qComments),
		)

		var rawComments []CommentWithAuthor
		if err := p.db.SelectContext(ctx, &rawComments, qComments, post.ID); err != nil {
			return nil, errors.Wrap(err, "selecting comments")
		}

		comments := make([]Comment, 0)
		for _, comment := range rawComments {
			author := Author{
				Username: comment.AuthorName,
				ID:       comment.AuthorID,
			}
			comments = append(comments, Comment{
				DateCreated: comment.DateCreated,
				Author:      author,
				Body:        comment.Body,
				ID:          comment.ID,
			})
		}
		p.log.Printf("%s: %s: %s with %v", traceID, "post.Query",
			database.Log(qComments), comments,
		)
		if post.Type == "url" {
			info = append(info, InfoLink{
				ID:               post.ID,
				Score:            post.Score,
				Views:            post.Views,
				Title:            post.Title,
				Payload:          post.Payload,
				Category:         post.Category,
				DateCreated:      post.DateCreated,
				Author:           author,
				Votes:            votes,
				Comments:         comments,
				UpvotePercentage: 100,
			})
		} else {
			info = append(info, InfoText{
				ID:               post.ID,
				Score:            post.Score,
				Views:            post.Views,
				Title:            post.Title,
				Payload:          post.Payload,
				Category:         post.Category,
				DateCreated:      post.DateCreated,
				Author:           author,
				Votes:            votes,
				Comments:         comments,
				UpvotePercentage: 100,
			})
		}
	}

	return info, nil
}

// todo: divide into functions
// QueryByID finds the post identified by a given ID.
func (p Post) QueryByID(ctx context.Context, traceID string, postID string) (Info, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.post.querybyid")
	defer span.End()

	if _, err := uuid.Parse(postID); err != nil {
		return InfoText{}, ErrInvalidID
	}

	const q = `
	SELECT * FROM
		posts
	WHERE
		post_id = $1`

	p.log.Printf("%s: %s: %s", traceID, "product.QueryByID",
		database.Log(q, postID),
	)

	var post PostDB
	if err := p.db.GetContext(ctx, &post, q, postID); err != nil {
		return InfoText{}, errors.Wrap(err, "selecting post by ID")
	}

	const qAuthor = `SELECT user_id, name FROM users WHERE user_id = $1`

	var author Author
	if err := p.db.GetContext(ctx, &author, qAuthor, post.UserID); err != nil {
		return nil, errors.Wrap(err, "selecting votes")
	}
	const qVotes = `SELECT user_id as User, vote as Vote FROM votes WHERE post_id = $1`

	var votes []Vote
	if err := p.db.SelectContext(ctx, &votes, qVotes, post.ID); err != nil {
		return nil, errors.Wrap(err, "selecting votes")
	}

	// Get their comments
	const qComments = `SELECT name, user_id, cm.date_created, body, comment_id FROM comments cm join users using(user_id) WHERE post_id = $1`

	p.log.Printf("%s: %s: %s", traceID, "post.Query",
		database.Log(qComments),
	)

	var rawComments []CommentWithAuthor
	if err := p.db.SelectContext(ctx, &rawComments, qComments, post.ID); err != nil {
		return nil, errors.Wrap(err, "selecting comments")
	}

	comments := make([]Comment, 0)
	for _, comment := range rawComments {
		author := Author{
			Username: comment.AuthorName,
			ID:       comment.AuthorID,
		}
		comments = append(comments, Comment{
			DateCreated: comment.DateCreated,
			Author:      author,
			Body:        comment.Body,
			ID:          comment.ID,
		})
	}
	if post.Type == "url" {
		return InfoLink{
			ID:               post.ID,
			Score:            post.Score,
			Views:            post.Views,
			Title:            post.Title,
			Payload:          post.Payload,
			Category:         post.Category,
			DateCreated:      post.DateCreated,
			Author:           author,
			Votes:            votes,
			Comments:         comments,
			UpvotePercentage: 100,
		}, nil
	} else {
		return InfoText{
			ID:               post.ID,
			Score:            post.Score,
			Views:            post.Views,
			Title:            post.Title,
			Payload:          post.Payload,
			Category:         post.Category,
			DateCreated:      post.DateCreated,
			Author:           author,
			Votes:            votes,
			Comments:         comments,
			UpvotePercentage: 100,
		}, nil
	}
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
			return InfoText{}, ErrNotFound
		}
		return InfoText{}, errors.Wrap(err, "selecting posts by category")
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
