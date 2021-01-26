// Package posts contains product related CRUD functionality.
package posts

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
	// ErrInvalidPostID occurs when an ID is not in a valid form.
	ErrInvalidPostID = errors.New("invalid posts id")

	// ErrInvalidCommentID occurs when an ID is not in a valid form.
	ErrInvalidCommentID = errors.New("invalid comment id")

	// ErrInvalidUserID occurs when an ID is not in a valid form.
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrForbidden occurs when a users tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")

	//ErrWrongPostType is used when users tries to create posts with incorrect type.
	ErrWrongPostType = errors.New("new posts should be of type link or text")

	// ErrPostNotFound is used when a specific Post is requested but does not exist.
	ErrPostNotFound = errors.New("posts not found")

	//ErrCommentNotFound is used when a specific Comment is requested but does not exist
	ErrCommentNotFound = errors.New("comment not found")

	//ErrUserNotFound is used when a specific UserID is requested but does not exist
	ErrUserNotFound = errors.New("users not found")
)

// PostSetup manages the set of API's for product access.
type Setup struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a post Setup for api access.
func New(log *log.Logger, db *sqlx.DB) Setup {
	return Setup{
		log: log,
		db:  db,
	}
}

// Create adds a posts to the database. It returns the created posts with fields like ID and DateCreated populated.
func (p Setup) Create(ctx context.Context, claims auth.Claims, np NewPost, now time.Time) (Info, error) {
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
		return Info{}, ErrWrongPostType
	}
	if post.Type == "link" {
		post.Payload = np.URL
	}

	if err := dbs.InsertPost(ctx, post); err != nil {
		return Info{}, err
	}

	if err := dbs.InsertVote(ctx, post.ID, post.UserID, 1); err != nil {
		return Info{}, err
	}

	info := getNewPostInfo(post, claims)
	return info, nil
}

// Delete removes the product identified by a given ID.
func (p Setup) Delete(ctx context.Context, claims auth.Claims, postID string) (Info, error) {
	dbs := db.NewDBset(p.log, p.db)
	if _, err := uuid.Parse(postID); err != nil {
		return Info{}, ErrInvalidPostID
	}

	post, err := dbs.GetPostByID(ctx, postID)
	if err != nil {
		if err == db.ErrPostNotFound {
			return Info{}, ErrPostNotFound
		}
		return Info{}, err
	}

	if claims.User.ID != post.UserID {
		return Info{}, ErrForbidden
	}
	info := getEmptyInfo(post)
	info, err = p.fillInfo(ctx, info)
	if err != nil {
		return Info{}, err
	}

	return info, dbs.DeletePost(ctx, postID)
}

// Query gets all Posts from the database
func (p Setup) Query(ctx context.Context) ([]Info, error) {
	dbs := db.NewDBset(p.log, p.db)
	posts, err := dbs.SelectAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	var infos []Info
	for _, post := range posts {
		info := getEmptyInfo(post)
		info, err = p.fillInfo(ctx, info)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}

//ObtainPosts returns all posts with a given category from a given users
func (p Setup) ObtainPosts(ctx context.Context, category, userID string) ([]Info, error) {
	dbs := db.NewDBset(p.log, p.db)
	posts, err := dbs.ObtainPosts(ctx, category, userID)
	if err != nil {
		return nil, err
	}
	var infos []Info
	for _, post := range posts {
		info := getEmptyInfo(post)
		info, err = p.fillInfo(ctx, info)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// QueryByID finds the posts identified by a given ID ready to be send to users.
func (p Setup) QueryByID(ctx context.Context, postID string) (Info, error) {
	dbs := db.NewDBset(p.log, p.db)
	if _, err := uuid.Parse(postID); err != nil {
		return Info{}, ErrInvalidPostID
	}

	post, err := dbs.GetPostByID(ctx, postID)
	if err != nil {
		if err == db.ErrPostNotFound {
			return Info{}, ErrPostNotFound
		}
		return Info{}, err
	}

	info := getEmptyInfo(post)
	info, err = p.fillInfo(ctx, info)
	if err != nil {
		return Info{}, err
	}

	return info, nil
}

// QueryByCat finds the posts identified by a given Category ready to be send to users.
func (p Setup) QueryByCat(ctx context.Context, category string) ([]Info, error) {
	dbs := db.NewDBset(p.log, p.db)
	posts, err := dbs.SelectPostsByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	var infos []Info
	for _, post := range posts {
		info := getEmptyInfo(post)
		info, err = p.fillInfo(ctx, info)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}

// QueryByUser finds the posts identified by a given users ID.
func (p Setup) QueryByUser(ctx context.Context, name string) ([]Info, error) {
	dbs := db.NewDBset(p.log, p.db)

	authorDB, err := dbs.GetUserByName(ctx, name)
	if err != nil {
		if err == db.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	author := convertUser(authorDB)

	posts, err := dbs.SelectPostsByUser(ctx, author.ID)
	if err != nil {
		return nil, err
	}

	var infos []Info
	for _, post := range posts {
		info := getEmptyInfo(post)
		info, err = p.fillInfo(ctx, info)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}
	return infos, nil
}

// Vote adds vote to the posts with given postID.
func (p Setup) Vote(ctx context.Context, claims auth.Claims, postID string, vote int) (Info, error) {
	if _, err := uuid.Parse(postID); err != nil {
		return Info{}, ErrInvalidPostID
	}
	dbs := db.NewDBset(p.log, p.db)
	if err := dbs.CheckPost(ctx, postID); err != nil {
		return Info{}, err
	}

	if err := dbs.CheckVote(ctx, postID, claims.User.ID); err != nil {
		if err != db.ErrVoteNotFound {
			return Info{}, err
		}
		if err := dbs.InsertVote(ctx, postID, claims.User.ID, vote); err != nil {
			return Info{}, err
		}
	} else {
		if err := dbs.UpdateVote(ctx, postID, claims.User.ID, vote); err != nil {
			return Info{}, err
		}
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return Info{}, errors.Wrap(err, "getting posts after voting")
	}

	return pst, nil
}

// Unvote erases vote to the posts from a single users
func (p Setup) Unvote(ctx context.Context, claims auth.Claims, postID string) (Info, error) {
	if _, err := uuid.Parse(postID); err != nil {
		return Info{}, ErrInvalidPostID
	}
	dbs := db.NewDBset(p.log, p.db)
	if err := dbs.CheckPost(ctx, postID); err != nil {
		return Info{}, err
	}

	if err := dbs.CheckVote(ctx, postID, claims.User.ID); err != nil {
		if err != db.ErrVoteNotFound {
			return Info{}, err
		} else {
			return Info{}, nil
		}
	}
	if err := dbs.DeleteVote(ctx, postID, claims.User.ID); err != nil {
		return Info{}, err
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return Info{}, errors.Wrap(err, "getting posts after voting")
	}

	return pst, nil
}

// CreateComment creates comment
func (p Setup) CreateComment(
	ctx context.Context, claims auth.Claims, nc NewComment, postID string, now time.Time) (Info, error) {
	if _, err := uuid.Parse(postID); err != nil {
		return Info{}, ErrInvalidPostID
	}
	dbs := db.NewDBset(p.log, p.db)
	if err := dbs.CheckPost(ctx, postID); err != nil {
		if err == db.ErrPostNotFound {
			return Info{}, ErrPostNotFound
		}
		return Info{}, err
	}

	ncDB := db.CommentDB{
		DateCreated: now,
		PostID:      postID,
		AuthorID:    claims.User.ID,
		Body:        nc.Text,
		ID:          uuid.New().String(),
	}

	if err := dbs.CreateComment(ctx, ncDB); err != nil {
		return Info{}, err
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return Info{}, errors.Wrap(err, "getting posts after creating comment")
	}
	return pst, nil

}

// DeleteComment deletes comment
func (p Setup) DeleteComment(ctx context.Context, claims auth.Claims, postID string, commentID string) (Info, error) {
	if _, err := uuid.Parse(postID); err != nil {
		return Info{}, ErrInvalidPostID
	}
	if _, err := uuid.Parse(commentID); err != nil {
		return Info{}, ErrInvalidCommentID
	}
	dbs := db.NewDBset(p.log, p.db)
	if err := dbs.CheckPost(ctx, postID); err != nil {
		if err == db.ErrPostNotFound {
			return Info{}, ErrPostNotFound
		}
		return Info{}, err
	}
	commentDB, err := dbs.GetCommentByID(ctx, commentID)
	if err != nil {
		if err == db.ErrCommentNotFound {
			return Info{}, ErrCommentNotFound
		}
		return Info{}, err
	}

	if claims.User.ID != commentDB.AuthorID {
		return Info{}, ErrForbidden
	}
	if err := dbs.DeleteComment(ctx, commentID); err != nil {
		return Info{}, err
	}

	pst, err := p.QueryByID(ctx, postID)
	if err != nil {
		return Info{}, errors.Wrap(err, "getting posts after voting")
	}

	return pst, nil
}

// AuthorByID finds the user identified by a given ID
func (p Setup) AuthorByID(ctx context.Context, userID string) (Author, error) {
	dbs := db.NewDBset(p.log, p.db)
	if _, err := uuid.Parse(userID); err != nil {
		return Author{}, ErrInvalidUserID
	}

	authorDB, err := dbs.GetUserByID(ctx, userID)
	if err != nil {
		if err == db.ErrUserNotFound {
			return Author{}, ErrUserNotFound
		}
		return Author{}, err
	}

	return convertUser(authorDB), nil
}
