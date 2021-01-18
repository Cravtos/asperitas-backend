package db

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

// SelectCommentsWithUserByPostID return slice of Comment for a single post
func (setup DBset) SelectCommentsWithUserByPostID(ctx context.Context, ID string) ([]CommentWithUserDB, error) {
	const qComments = `
		SELECT name, user_id, cm.date_created, body, comment_id 
		FROM comments cm join users using(user_id) 
		WHERE post_id = $1`

	//setup.log.Printf("%s: %s", "db.selectCommentsByPostID", database.Log(qComments, ID))

	var comments []CommentWithUserDB
	if err := setup.db.SelectContext(ctx, &comments, qComments, ID); err != nil {
		return nil, errors.Wrap(err, "selecting comments")
	}
	return comments, nil
}

// GetCommentWithUserByID returns comment with ID commentID
func (setup DBset) GetCommentWithUserByID(ctx context.Context, commentID string) (CommentWithUserDB, error) {
	const qComment = `
		SELECT name, user_id, cm.date_created, body, comment_id 
		FROM comments cm join users using(user_id) 
		WHERE comment_id = $1`

	//setup.log.Printf("%s: %s", "db.getCommentByID", database.Log(qComment, commentID))

	var comment CommentWithUserDB
	if err := setup.db.GetContext(ctx, &comment, qComment, commentID); err != nil {
		if err == sql.ErrNoRows {
			return CommentWithUserDB{}, ErrCommentNotFound
		}
		return CommentWithUserDB{}, errors.Wrap(err, "selecting comment by ID")
	}
	return comment, nil
}
