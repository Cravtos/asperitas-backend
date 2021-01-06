// Package post contains product related CRUD functionality.
package post

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

var (
	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("invalid post id")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")

	//todo ErrNotFound should exist not only for posts (but maybe it should be more serious problem than with Post)

	// ErrNotFound is used when a specific Post is requested but does not exist.
	ErrNotFound = errors.New("post not found")
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
func (p Post) Create(ctx context.Context, claims auth.Claims, np NewPost, now time.Time) (Info, error) {
	post := PostDB{
		ID:          uuid.New().String(),
		Score:       0,
		Views:       0,
		Title:       np.Title,
		Type:        np.Type,
		Category:    np.Category,
		Payload:     np.Text,
		DateCreated: now,
		UserID:      claims.User.ID,
	}
	if post.Type == "url" {
		post.Payload = np.URL
	}

	if err := p.insertPost(ctx, post); err != nil {
		return InfoText{}, err
	}

	if err := p.insertVote(ctx, post.ID, post.UserID, 1); err != nil {
		return InfoText{}, err
	}

	info := infoByPostAndClaims(post, claims)
	return info, nil
}

// Delete removes the product identified by a given ID.
func (p Post) Delete(ctx context.Context, claims auth.Claims, postID string) error {

	if _, err := uuid.Parse(postID); err != nil {
		return ErrInvalidID
	}

	const qSelectAuthor = `SELECT user_id FROM posts WHERE post_id = $1`

	var author Author
	if err := p.db.GetContext(ctx, &author, qSelectAuthor, postID); err != nil {
		return errors.Wrap(err, "selecting post by ID")
	}

	if claims.User.ID == author.ID {
		return ErrForbidden
	}

	const qDeletePost = `DELETE FROM posts WHERE post_id = $1`

	p.log.Printf("%s: %s: %s", "post.Delete", database.Log(qDeletePost, postID))

	if _, err := p.db.ExecContext(ctx, qDeletePost, postID); err != nil {
		return errors.Wrapf(err, "deleting post %s", postID)
	}

	return nil
}

// Query gets all Posts from the database ready to be send to user.
func (p Post) Query(ctx context.Context) ([]Info, error) {

	posts, err := p.selectAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	var info []Info
	for _, post := range posts {
		author, err := p.getAuthorByID(ctx, post.UserID)
		if err != nil {
			return nil, err
		}

		votes, err := p.selectVotesByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}

		comments, err := p.selectCommentsByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}

		info = append(info, infoByDBdata(post, author, votes, comments))
	}

	return info, nil
}

// QueryByID finds the post identified by a given ID ready to be send to user.
func (p Post) QueryByID(ctx context.Context, postID string) (Info, error) {
	if _, err := uuid.Parse(postID); err != nil {
		return InfoText{}, ErrInvalidID
	}

	post, err := p.getPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	author, err := p.getAuthorByID(ctx, post.UserID)
	if err != nil {
		return nil, err
	}

	votes, err := p.selectVotesByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	comments, err := p.selectCommentsByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return infoByDBdata(post, author, votes, comments), nil
}

// QueryByCat finds the post identified by a given Category ready to be send to user.
func (p Post) QueryByCat(ctx context.Context, category string) ([]Info, error) {

	// Get all posts
	const qPost = `SELECT * FROM posts WHERE category = $1`

	p.log.Printf("%s: %s: %s", "post.QueryByCat",
		database.Log(qPost),
	)

	var posts []PostDB
	if err := p.db.SelectContext(ctx, &posts, qPost, category); err != nil {
		return nil, errors.Wrap(err, "selecting posts")
	}

	var info []Info
	for _, post := range posts {
		// todo: divide into functions
		const qAuthor = `SELECT user_id, name FROM users WHERE user_id = $1`

		var author Author
		if err := p.db.GetContext(ctx, &author, qAuthor, post.UserID); err != nil {
			return nil, errors.Wrap(err, "selecting authors")
		}

		// Get posts votes
		const qVotes = `SELECT user_id, vote FROM votes WHERE post_id = $1`

		p.log.Printf("%s: %s: %s", "post.QueryByCat",
			database.Log(qPost),
		)

		var votes []Vote
		if err := p.db.SelectContext(ctx, &votes, qVotes, post.ID); err != nil {
			return nil, errors.Wrap(err, "selecting votes")
		}

		// Get their comments
		const qComments = `
		SELECT 
			name, user_id, cm.date_created, body, comment_id 
		FROM 
			comments cm join users using(user_id) 
		WHERE 
			post_id = $1`

		p.log.Printf("%s: %s: %s", "post.QueryByCat",
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
		p.log.Printf("%s: %s: %s with %v", "post.QueryByCat",
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
				UpvotePercentage: upvotePercentage(votes),
			})
			continue
		}
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
			UpvotePercentage: upvotePercentage(votes),
		})
	}

	return info, nil
}

// QueryByUser finds the post identified by a given user. // Todo: think about how to get post by user from PostDB
func (p Post) QueryByUser(ctx context.Context, user string) ([]Info, error) {
	// Get author
	const qAuthor = `SELECT user_id, name FROM users WHERE name = $1`

	var author Author
	if err := p.db.GetContext(ctx, &author, qAuthor, user); err != nil {
		return nil, errors.Wrap(err, "selecting author")
	}

	// Get all posts
	const qPost = `SELECT * FROM posts WHERE user_id = $1`

	p.log.Printf("%s: %s: %s", "post.QueryByUser",
		database.Log(qPost),
	)

	var posts []PostDB
	if err := p.db.SelectContext(ctx, &posts, qPost, author.ID); err != nil {
		return nil, errors.Wrap(err, "selecting posts")
	}

	var info []Info
	for _, post := range posts {
		// todo: divide into functions

		// Get posts votes
		const qVotes = `SELECT user_id, vote FROM votes WHERE post_id = $1`

		p.log.Printf("%s: %s: %s", "post.QueryByUser",
			database.Log(qPost),
		)

		var votes []Vote
		if err := p.db.SelectContext(ctx, &votes, qVotes, post.ID); err != nil {
			return nil, errors.Wrap(err, "selecting votes")
		}

		// Get their comments
		const qComments = `
		SELECT
			name, user_id, cm.date_created, body, comment_id
		FROM
			comments cm join users using(user_id)
		WHERE
			post_id = $1`

		p.log.Printf("%s: %s: %s", "post.QueryByUser",
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
		p.log.Printf("%s: %s: %s with %v", "post.QueryByUser",
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
				UpvotePercentage: upvotePercentage(votes),
			})
			continue
		}
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
			UpvotePercentage: upvotePercentage(votes),
		})
	}

	return info, nil
}

// Vote adds vote to the post with given postID.
func (p Post) Vote(ctx context.Context, claims auth.Claims, postID string, vote int) (Info, error) {

	// todo: update post score
	// todo: maybe replace qCheckExist by QueryByID with ErrNotFound check
	const qCheckExist = `SELECT COUNT(1) FROM posts WHERE post_id = $1`

	p.log.Printf("%s: %s: %s", "post.Vote",
		database.Log(qCheckExist),
	)

	var exist []int
	if err := p.db.SelectContext(ctx, &exist, qCheckExist, postID); err != nil {
		return nil, errors.Wrap(err, "checking if post exists")
	}

	if exist[0] == 0 {
		return nil, ErrNotFound
	}

	const qDeleteVote = `DELETE FROM votes WHERE post_id = $1 AND user_id = $2`

	p.log.Printf("%s: %s: %s", "post.Vote",
		database.Log(qDeleteVote),
	)

	if _, err := p.db.ExecContext(ctx, qDeleteVote, postID, claims.User.ID); err != nil {
		return nil, errors.Wrap(err, "deleting votes")
	}

	const qPutVote = `
		INSERT INTO votes
			(post_id, user_id, vote)
		VALUES
			($1, $2, $3)`

	p.log.Printf("%s: %s: %s", "post.Vote",
		database.Log(qPutVote),
	)

	if _, err := p.db.ExecContext(ctx, qPutVote, postID, claims.User.ID, vote); err != nil {
		return nil, errors.Wrap(err, "putting vote")
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "getting post after voting")
	}

	return pst, nil
}
