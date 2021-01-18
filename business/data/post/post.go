// Package post contains product related CRUD functionality.
package post

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

var (
	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("invalid post id")

	// ErrInvalidCommentID occurs when an ID is not in a valid form.
	ErrInvalidCommentID = errors.New("invalid comment id")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")

	//ErrCommentNotFound is used when user tries to create post with incorrect type.
	ErrWrongPostType = errors.New("new post should be of type url or text")
)

// PostSet manages the set of API's for product access.
type PostSet struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a PostSet for api access.
func New(log *log.Logger, db *sqlx.DB) PostSet {
	return PostSet{
		log: log,
		db:  db,
	}
}

// Create adds a post to the database. It returns the created post with fields like ID and DateCreated populated.
func (p PostSet) Create(ctx context.Context, claims auth.Claims, np NewPost, now time.Time) (Info, error) {
	dbs := db.NewDBset(p.log, p.db)
	post := db.PostDB{
		ID:          uuid.New().String(),
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

	if err := dbs.InsertPost(ctx, post); err != nil {
		return nil, err
	}

	if err := dbs.InsertVote(ctx, post.ID, post.UserID, 1); err != nil {
		return nil, err
	}

	info := getNewPostInfo(post, claims)
	return info, nil
}

// Delete removes the product identified by a given ID.
func (p PostSet) Delete(ctx context.Context, claims auth.Claims, postID string) error {
	dbs := db.NewDBset(p.log, p.db)
	if _, err := uuid.Parse(postID); err != nil {
		return ErrInvalidID
	}

	post, err := dbs.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}

	if claims.User.ID != post.UserID {
		return ErrForbidden
	}

	return dbs.DeletePost(ctx, postID)
}

// Query gets all Posts from the database ready to be send to user.
func (p PostSet) Query(ctx context.Context) ([]Info, error) {
	dbs := db.NewDBset(p.log, p.db)
	posts, err := dbs.SelectAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	var info []Info
	for _, post := range posts {
		authorDB, err := dbs.GetUserByID(ctx, post.UserID)
		if err != nil {
			return nil, err
		}
		author := convertUser(authorDB)

		votesDB, err := dbs.SelectVotesByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		votes := convertVotes(votesDB)

		commentsWithUserDB, err := dbs.SelectCommentsWithUserByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		comments := convertComments(commentsWithUserDB)

		info = append(info, getInfo(post, author, votes, comments))
	}

	return info, nil
}

// QueryByID finds the post identified by a given ID ready to be send to user.
func (p PostSet) QueryByID(ctx context.Context, postID string) (Info, error) {
	dbs := db.NewDBset(p.log, p.db)
	if _, err := uuid.Parse(postID); err != nil {
		return InfoText{}, ErrInvalidID
	}

	post, err := dbs.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	authorDB, err := dbs.GetUserByID(ctx, post.UserID)
	if err != nil {
		return nil, err
	}
	author := convertUser(authorDB)

	votesDB, err := dbs.SelectVotesByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	votes := convertVotes(votesDB)

	commentsWithUserDB, err := dbs.SelectCommentsWithUserByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	comments := convertComments(commentsWithUserDB)

	return getInfo(post, author, votes, comments), nil
}

// QueryByCat finds the post identified by a given Category ready to be send to user.
func (p PostSet) QueryByCat(ctx context.Context, category string) ([]Info, error) {
	dbs := db.NewDBset(p.log, p.db)
	posts, err := dbs.SelectPostsByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	var info []Info
	for _, post := range posts {
		authorDB, err := dbs.GetUserByID(ctx, post.UserID)
		if err != nil {
			return nil, err
		}
		author := convertUser(authorDB)

		votesDB, err := dbs.SelectVotesByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		votes := convertVotes(votesDB)

		commentsWithUserDB, err := dbs.SelectCommentsWithUserByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		comments := convertComments(commentsWithUserDB)

		info = append(info, getInfo(post, author, votes, comments))
	}

	return info, nil
}

// QueryByUser finds the posts identified by a given user ID.
func (p PostSet) QueryByUser(ctx context.Context, name string) ([]Info, error) {
	dbs := db.NewDBset(p.log, p.db)

	authorDB, err := dbs.GetUserByName(ctx, name)
	if err != nil {
		return nil, err
	}
	author := convertUser(authorDB)

	posts, err := dbs.SelectPostsByUser(ctx, author.ID)
	if err != nil {
		return nil, err
	}

	var info []Info
	for _, post := range posts {
		votesDB, err := dbs.SelectVotesByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		votes := convertVotes(votesDB)

		commentsWithUserDB, err := dbs.SelectCommentsWithUserByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		comments := convertComments(commentsWithUserDB)

		info = append(info, getInfo(post, author, votes, comments))
	}
	return info, nil
}

// Vote adds vote to the post with given postID.
func (p PostSet) Vote(ctx context.Context, claims auth.Claims, postID string, vote int) (Info, error) {
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidID
	}
	dbs := db.NewDBset(p.log, p.db)
	if err := dbs.CheckPost(ctx, postID); err != nil {
		return nil, err
	}

	if err := dbs.CheckVote(ctx, postID, claims.User.ID); err != nil {
		if err != db.ErrVoteNotFound {
			return nil, err
		}
		if err := dbs.InsertVote(ctx, postID, claims.User.ID, vote); err != nil {
			return nil, err
		}
	} else {
		if err := dbs.UpdateVote(ctx, postID, claims.User.ID, vote); err != nil {
			return nil, err
		}
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "getting post after voting")
	}

	return pst, nil
}

// Unvote erases vote to the post from a single user
func (p PostSet) Unvote(ctx context.Context, claims auth.Claims, postID string) (Info, error) {
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidID
	}
	dbs := db.NewDBset(p.log, p.db)
	if err := dbs.CheckPost(ctx, postID); err != nil {
		return nil, err
	}

	if err := dbs.CheckVote(ctx, postID, claims.User.ID); err != nil {
		if err != db.ErrVoteNotFound {
			return nil, err
		} else {
			return nil, nil
		}
	}
	if err := dbs.DeleteVote(ctx, postID, claims.User.ID); err != nil {
		return nil, err
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "getting post after voting")
	}

	return pst, nil
}

// CreateComment creates comment
func (p PostSet) CreateComment(
	ctx context.Context, claims auth.Claims, nc NewComment, postID string, now time.Time) (Info, error) {
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidID
	}
	dbs := db.NewDBset(p.log, p.db)
	if err := dbs.CheckPost(ctx, postID); err != nil {
		return nil, err
	}

	ncDB := db.CommentDB{
		DateCreated: now,
		PostID:      postID,
		AuthorID:    claims.User.ID,
		Body:        nc.Text,
		ID:          uuid.New().String(),
	}

	if err := dbs.CreateComment(ctx, ncDB); err != nil {
		return InfoText{}, err
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "getting post after creating comment")
	}
	return pst, nil

}

// DeleteComment deletes comment
func (p PostSet) DeleteComment(ctx context.Context, claims auth.Claims, postID string, commentID string) (Info, error) {
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidID
	}
	if _, err := uuid.Parse(commentID); err != nil {
		return nil, ErrInvalidCommentID
	}
	dbs := db.NewDBset(p.log, p.db)
	if err := dbs.CheckPost(ctx, postID); err != nil {
		return nil, err
	}
	commentDB, err := dbs.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if claims.User.ID != commentDB.AuthorID {
		return nil, ErrForbidden
	}
	if err := dbs.DeleteComment(ctx, commentID); err != nil {
		return nil, err
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "getting post after voting")
	}

	return pst, nil
}
