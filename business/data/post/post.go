// Package post contains product related CRUD functionality.
package post

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/auth"
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

	// ErrPostNotFound is used when a specific Post is requested but does not exist.
	ErrPostNotFound = errors.New("post not found")

	//ErrCommentNotFound is used when a specific Comment is requested but does not exist
	ErrCommentNotFound = errors.New("comment not found")

	ErrWrongPostType = errors.New("new post should be of type url or text")
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

	// todo: find reason why link posts are not created
	post := postDB{
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
	if post.Type != "link" && post.Type != "text" {
		return nil, ErrWrongPostType
	}
	if post.Type == "link" {
		post.Payload = np.URL
	}

	if err := p.insertPost(ctx, post); err != nil {
		return nil, err
	}

	if err := p.insertVote(ctx, post.ID, post.UserID, 1); err != nil {
		return nil, err
	}

	info := infoByPostAndClaims(post, claims)
	return info, nil
}

// Delete removes the product identified by a given ID.
func (p Post) Delete(ctx context.Context, claims auth.Claims, postID string) error {

	if _, err := uuid.Parse(postID); err != nil {
		return ErrInvalidID
	}

	post, err := p.getPostByID(ctx, postID)
	if err != nil {
		return err
	}

	if claims.User.ID != post.UserID {
		return ErrForbidden
	}

	return p.deletePost(ctx, postID)
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
	posts, err := p.selectPostsByCategory(ctx, category)
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

// QueryByUser finds the posts identified by a given user ID.
func (p Post) QueryByUser(ctx context.Context, name string) ([]Info, error) {
	author, err := p.getAuthorByName(ctx, name)
	if err != nil {
		return nil, err
	}
	posts, err := p.selectPostsByUser(ctx, author.ID)
	if err != nil {
		return nil, err
	}

	var info []Info
	for _, post := range posts {
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

// Vote adds vote to the post with given postID.
func (p Post) Vote(ctx context.Context, claims auth.Claims, postID string, vote int) (Info, error) {
	if err := p.checkPost(ctx, postID); err != nil {
		return nil, err
	}

	if err := p.checkVote(ctx, postID, claims.User.ID); err != nil {
		if err != ErrPostNotFound {
			return nil, err
		}
		if err := p.insertVote(ctx, postID, claims.User.ID, vote); err != nil {
			return nil, err
		}
	} else {
		if err := p.updateVote(ctx, postID, claims.User.ID, vote); err != nil {
			return nil, err
		}
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "getting post after voting")
	}

	return pst, nil
}

//Unvote erases vote to the post from a single user
func (p Post) Unvote(ctx context.Context, claims auth.Claims, postID string) (Info, error) {
	if err := p.checkPost(ctx, postID); err != nil {
		return nil, err
	}

	if err := p.checkVote(ctx, postID, claims.User.ID); err != nil {
		if err != ErrPostNotFound {
			return nil, err
		} else {
			return nil, nil
		}
	}
	if err := p.deleteVote(ctx, postID, claims.User.ID); err != nil {
		return nil, err
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "getting post after voting")
	}

	return pst, nil
}

func (p Post) CreateComment(
	ctx context.Context, claims auth.Claims, nc NewComment, postID string, now time.Time) (Info, error) {
	if err := p.checkPost(ctx, postID); err != nil {
		return nil, err
	}

	if err := p.insertComment(ctx, uuid.New().String(), postID, claims.User.ID, nc.Text, now); err != nil {
		return InfoText{}, err
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "getting post after creating comment")
	}
	return pst, nil

}

func (p Post) DeleteComment(ctx context.Context, claims auth.Claims, postID string, commentID string) (Info, error) {
	if err := p.checkPost(ctx, postID); err != nil {
		return nil, err
	}
	comment, err := p.getCommentByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if claims.User.ID != comment.Author.ID {
		return nil, ErrForbidden
	}
	if err := p.deleteComment(ctx, commentID); err != nil {
		return nil, err
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "getting post after voting")
	}

	return pst, nil
}
