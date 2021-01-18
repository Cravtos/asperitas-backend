package db

import (
	"context"
	"database/sql"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/pkg/errors"
)

//ObtainPosts returns all posts with a given category from a given user
func (setup DBset) ObtainPosts(ctx context.Context, category, userID string) ([]PostDB, error) {
	if category == "all" && userID == "" {
		return setup.SelectAllPosts(ctx)
	} else if category == "all" {
		return setup.SelectPostsByUser(ctx, userID)
	} else if userID == "" {
		return setup.SelectPostsByCategory(ctx, category)
	} else {
		return setup.SelectPostsByCategoryAndUser(ctx, category, userID)
	}
}

// SelectPostsByCategoryAndUser returns all posts with a given category from user stored in database
func (setup DBset) SelectPostsByCategoryAndUser(ctx context.Context, category string, userID string) ([]PostDB, error) {
	const qPost = `SELECT * FROM posts WHERE category = $1 and user_id = $2`

	setup.log.Printf("%s: %s", "db.selectPostsByCategoryAndUser", database.Log(qPost, category, userID))

	var posts []PostDB
	if err := setup.db.SelectContext(ctx, &posts, qPost, category, userID); err != nil {
		return nil, errors.Wrap(err, "selecting category posts")
	}
	return posts, nil
}

// SelectAllPosts return all posts stored in database
func (setup DBset) SelectAllPosts(ctx context.Context) ([]PostDB, error) {
	const qPost = `SELECT * FROM posts`

	setup.log.Printf("%s: %s", "db.selectAllPosts", database.Log(qPost))

	var posts []PostDB
	if err := setup.db.SelectContext(ctx, &posts, qPost); err != nil {
		return nil, errors.Wrap(err, "selecting all posts")
	}

	return posts, nil
}

// SelectPostsByCategory returns all posts with a given category stored in database
func (setup DBset) SelectPostsByCategory(ctx context.Context, category string) ([]PostDB, error) {
	const qPost = `SELECT * FROM posts WHERE category = $1`

	setup.log.Printf("%s: %s", "db.selectPostsByCategory", database.Log(qPost, category))

	var posts []PostDB
	if err := setup.db.SelectContext(ctx, &posts, qPost, category); err != nil {
		return nil, errors.Wrap(err, "selecting category posts")
	}

	return posts, nil
}

// SelectPostsByUser returns all posts from user stored in database
func (setup DBset) SelectPostsByUser(ctx context.Context, userID string) ([]PostDB, error) {
	const qPost = `SELECT * FROM posts WHERE user_id = $1`

	setup.log.Printf("%s: %s", "db.selectPostsByUser", database.Log(qPost, userID))

	var posts []PostDB
	if err := setup.db.SelectContext(ctx, &posts, qPost, userID); err != nil {
		return nil, errors.Wrap(err, "selecting users posts")
	}

	return posts, nil
}

// GetPostByID obtains post from database using ID
func (setup DBset) GetPostByID(ctx context.Context, postID string) (PostDB, error) {
	const q = `	SELECT * FROM posts WHERE post_id = $1`

	setup.log.Printf("%s: %s", "db.getPostByID", database.Log(q, postID))

	var post PostDB
	if err := setup.db.GetContext(ctx, &post, q, postID); err != nil {
		if err == sql.ErrNoRows {
			return PostDB{}, ErrPostNotFound
		}
		return PostDB{}, errors.Wrap(err, "selecting post by ID")
	}

	return post, nil
}

// CheckPost shows whether post with given ID exist in database or not.
// It returns an error if post doesn't exist.
func (setup DBset) CheckPost(ctx context.Context, postID string) error {
	const qCheckExist = `SELECT COUNT(*) FROM posts WHERE post_id = $1`

	setup.log.Printf("%s: %s", "db.checkPost", database.Log(qCheckExist, postID))

	var exist int
	if err := setup.db.GetContext(ctx, &exist, qCheckExist, postID); err != nil {
		return errors.Wrap(err, "checking if post exists")
	}

	if exist == 0 {
		return ErrPostNotFound
	}
	return nil
}

// InsertPost adds one new row to posts table
func (setup DBset) InsertPost(ctx context.Context, post PostDB) error {
	const qPost = `
	INSERT INTO posts
		(post_id, views, type, title, category, payload, date_created, user_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)`

	setup.log.Printf("%s: %s", "db.insertPost",
		database.Log(qPost, post.ID, post.Views, post.Type, post.Title, post.Category, post.Payload,
			post.DateCreated, post.UserID),
	)

	if _, err := setup.db.ExecContext(ctx, qPost, post.ID, post.Views, post.Type, post.Title,
		post.Category, post.Payload, post.DateCreated, post.UserID); err != nil {
		return errors.Wrap(err, "inserting post")
	}
	return nil
}

// DeletePost deletes post with all its votes and comments
func (setup DBset) DeletePost(ctx context.Context, postID string) error {
	const qDeleteVotes = `DELETE FROM votes WHERE post_id = $1`
	setup.log.Printf("%s: %s", "db.deletePost", database.Log(qDeleteVotes, postID))
	if _, err := setup.db.ExecContext(ctx, qDeleteVotes, postID); err != nil {
		return errors.Wrapf(err, "deleting votes %s", postID)
	}

	const qDeleteComments = `DELETE FROM comments WHERE post_id = $1`
	setup.log.Printf("%s: %s", "db.deletePost", database.Log(qDeleteComments, postID))
	if _, err := setup.db.ExecContext(ctx, qDeleteComments, postID); err != nil {
		return errors.Wrapf(err, "deleting comments %s", postID)
	}

	const qDeletePost = `DELETE FROM posts WHERE post_id = $1`
	setup.log.Printf("%s: %s", "db.deletePost", database.Log(qDeletePost, postID))
	if _, err := setup.db.ExecContext(ctx, qDeletePost, postID); err != nil {
		return errors.Wrapf(err, "deleting post %s", postID)
	}

	return nil
}
