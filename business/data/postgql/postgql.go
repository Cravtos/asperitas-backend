package postgql

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/data/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
)

var (
	// ErrInvalidPostID occurs when an ID is not in a valid form.
	ErrInvalidPostID = errors.New("invalid post id")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")

	// ErrInvalidCommentID occurs when an ID is not in a valid form.
	ErrInvalidCommentID = errors.New("invalid comment id")
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyPostGQL is used to store/retrieve a Claims value from a context.Context.
const KeyPostGQL ctxKey = 1
const KeyAuthHeader ctxKey = 2
const KeyAuth ctxKey = 3

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

//fillInfos fills slice of Info with votes, comments and author for each Post
func (g PostGQL) fillInfos(ctx context.Context, posts []Info) ([]Info, error) {
	for i, post := range posts {
		filledPost, err := g.fillInfo(ctx, post)
		if err != nil {
			return nil, err
		}
		posts[i] = filledPost
	}
	return posts, nil
}

//fillInfo fills Info with votes, comments and author
func (g PostGQL) fillInfo(ctx context.Context, src Info) (Info, error) {
	dbs := db.NewDBset(g.log, g.db)

	//obtaining votes for post
	votesDB, err := dbs.SelectVotesByPostID(ctx, src.ID)
	if err != nil {
		return Info{}, err
	}
	votes := convertVotes(votesDB)
	src.Votes = votes

	//obtaining author for post
	userDB, err := dbs.GetUserByID(ctx, src.UserID)
	if err != nil {
		return Info{}, err
	}
	author := convertUser(userDB)
	src.Author = &author

	//obtaining comments for post
	commentsWithAuthorDB, err := dbs.SelectCommentsWithUserByPostID(ctx, src.ID)
	if err != nil {
		return Info{}, err
	}
	comments := convertComments(commentsWithAuthorDB)
	src.Comments = comments

	return src, nil
}
