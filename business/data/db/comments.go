package db

import (
	"context"
	"database/sql"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/pkg/errors"
)

// SelectCommentsByPostID return slice of Comment for a single post
func (setup DBset) SelectCommentsByPostID(ctx context.Context, ID string) ([]CommentDB, error) {
	const qComments = `	SELECT * FROM comments WHERE post_id = $1`

	setup.log.Printf("%s: %s", "db.selectCommentsByPostID", database.Log(qComments, ID))

	var comments []CommentDB
	if err := setup.db.SelectContext(ctx, &comments, qComments, ID); err != nil {
		return nil, errors.Wrap(err, "selecting comments")
	}
	return comments, nil
}

// CreateComment creates comment with specified data
func (setup DBset) CreateComment(ctx context.Context, nc CommentDB) error {
	const qComment = `
	INSERT INTO comments
		(comment_id, post_id, user_id, body, date_created)
	VALUES
		($1, $2, $3, $4, $5)`

	setup.log.Printf("%s: %s", "db.createComment",
		database.Log(qComment, nc.ID, nc.PostID, nc.AuthorID, nc.Body, nc.DateCreated))

	if _, err := setup.db.ExecContext(ctx, qComment, nc.ID, nc.PostID, nc.AuthorID, nc.Body, nc.DateCreated); err != nil {
		return errors.Wrap(err, "inserting Comment")
	}
	return nil
}

// GetCommentByID returns comment with ID commentID
func (setup DBset) GetCommentByID(ctx context.Context, commentID string) (CommentDB, error) {
	const qComment = `SELECT * FROM comments WHERE comment_id = $1`

	setup.log.Printf("%s: %s", "db.getCommentByID", database.Log(qComment, commentID))

	var comment CommentDB
	if err := setup.db.GetContext(ctx, &comment, qComment, commentID); err != nil {
		if err == sql.ErrNoRows {
			return CommentDB{}, ErrCommentNotFound
		}
		return CommentDB{}, errors.Wrap(err, "selecting comment by ID")
	}
	return comment, nil
}

// DeleteComment deletes comment.
func (setup DBset) DeleteComment(ctx context.Context, commentID string) error {
	const qDeleteComment = `DELETE FROM comments WHERE comment_id = $1`

	setup.log.Printf("%s: %s", "db.deleteComment", database.Log(qDeleteComment, commentID))

	if _, err := setup.db.ExecContext(ctx, qDeleteComment, commentID); err != nil {
		if err == sql.ErrNoRows {
			return ErrCommentNotFound
		}
		return errors.Wrapf(err, "deleting comment %s", commentID)
	}
	return nil
}
