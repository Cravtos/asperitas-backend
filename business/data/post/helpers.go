package post

import (
	"context"
	"database/sql"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/pkg/errors"
	"time"
)

func upvotePercentage(votes []Vote) int {
	var positive float32

	for _, vote := range votes {
		if vote.Vote == 1 {
			positive++
		}
	}

	return int(positive / float32(len(votes)) * 100)
}

//creates new Info using postDB and auth.Claims given by user
func infoByPostAndClaims(post postDB, claims auth.Claims) Info {
	var info Info
	if post.Type == "url" {
		info = InfoLink{
			Type:        "link",
			ID:          post.ID,
			Score:       post.Score,
			Views:       post.Views,
			Title:       post.Title,
			Payload:     post.Payload,
			Category:    post.Category,
			DateCreated: post.DateCreated,
			Author: Author{
				Username: claims.User.Username,
				ID:       claims.User.ID,
			},
			Votes: []Vote{
				{User: claims.User.ID, Vote: 1},
			},
			Comments:         []Comment{},
			UpvotePercentage: 100,
		}
	} else {
		info = InfoText{
			Type:        "text",
			ID:          post.ID,
			Score:       post.Score,
			Views:       post.Views,
			Title:       post.Title,
			Payload:     post.Payload,
			Category:    post.Category,
			DateCreated: post.DateCreated,
			Author: Author{
				Username: claims.User.Username,
				ID:       claims.User.ID,
			},
			Votes: []Vote{
				{User: claims.User.ID, Vote: 1},
			},
			Comments:         []Comment{},
			UpvotePercentage: 100,
		}
	}
	return info
}

//creates new Info using data got from DB
func infoByDBdata(post postDB, author Author, votes []Vote, comments []Comment) Info {
	var info Info
	if post.Type == "url" {
		info = InfoLink{
			Type:             "link",
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
		}
	} else {
		info = InfoText{
			Type:             "text",
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
		}
	}
	return info
}

//obtains Author using its ID in DB
func (p Post) getAuthorByID(ctx context.Context, ID string) (Author, error) {
	const qAuthor = `SELECT user_id, name FROM users WHERE user_id = $1`

	p.log.Printf("%s: %s", "post.helpers.getAuthorByID", database.Log(qAuthor))

	var author Author
	if err := p.db.GetContext(ctx, &author, qAuthor, ID); err != nil {
		return Author{}, errors.Wrap(err, "selecting authors")
	}
	return author, nil
}

//returns slice of Vote for a single post
func (p Post) selectVotesByPostID(ctx context.Context, ID string) ([]Vote, error) {
	const qVotes = `SELECT user_id, vote FROM votes WHERE post_id = $1`

	p.log.Printf("%s: %s", "post.helpers.selectVotesByPostID", database.Log(qVotes))

	var votes []Vote
	if err := p.db.SelectContext(ctx, &votes, qVotes, ID); err != nil {
		return nil, errors.Wrap(err, "selecting votes")
	}
	return votes, nil
}

//returns score for a single post
func (p Post) getPostScore(ctx context.Context, ID string) (int, error) {
	const qScore = `SELECT SUM(vote) as score FROM votes WHERE post_id = $1 HAVING SUM(vote) is not null`

	p.log.Printf("%s: %s", "post.helpers.getPostScore", database.Log(qScore))

	var score int
	if err := p.db.GetContext(ctx, &score, qScore, ID); err != nil {
		return 0, nil
	}
	return score, nil
}

//return slice of Comment for a single post
func (p Post) selectCommentsByPostID(ctx context.Context, ID string) ([]Comment, error) {
	const qComments = `
		SELECT 
			name, user_id, cm.date_created, body, comment_id 
		FROM 
			comments cm join users using(user_id) 
		WHERE 
			post_id = $1`

	p.log.Printf("%s: %s", "post.helpers.selectCommentsByPostID", database.Log(qComments))

	var rawComments []struct {
		DateCreated time.Time `db:"date_created"`
		AuthorName  string    `db:"name"`
		AuthorID    string    `db:"user_id"`
		Body        string    `db:"body"`
		ID          string    `db:"comment_id"`
	}
	if err := p.db.SelectContext(ctx, &rawComments, qComments, ID); err != nil {
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
	return comments, nil
}

//return all posts stored in DB
func (p Post) selectAllPosts(ctx context.Context) ([]postDB, error) {
	const qPost = `SELECT * FROM posts`

	p.log.Printf("%s: %s", "post.helpers.selectAllPosts", database.Log(qPost))

	var posts []postDB
	if err := p.db.SelectContext(ctx, &posts, qPost); err != nil {
		return nil, errors.Wrap(err, "selecting all posts")
	}

	for i := range posts {
		score, err := p.getPostScore(ctx, posts[i].ID)
		if err != nil {
			return nil, err
		}
		posts[i].Score = score
	}
	return posts, nil
}

//selectPostsByCategory returns all posts with a given category stored in DB
func (p Post) selectPostsByCategory(ctx context.Context, category string) ([]postDB, error) {
	const qPost = `SELECT * FROM posts WHERE category = $1`

	p.log.Printf("%s: %s", "post.helpers.selectPostsByCategory", database.Log(qPost))

	var posts []postDB
	if err := p.db.SelectContext(ctx, &posts, qPost, category); err != nil {
		return nil, errors.Wrap(err, "selecting category posts")
	}
	for i := range posts {
		score, err := p.getPostScore(ctx, posts[i].ID)
		if err != nil {
			return nil, err
		}
		posts[i].Score = score
	}
	return posts, nil
}

//selectPostsByUser returns all posts from user stored in DB
func (p Post) selectPostsByUser(ctx context.Context, userID string) ([]postDB, error) {
	const qPost = `SELECT * FROM posts WHERE user_id = $1`

	p.log.Printf("%s: %s", "post.helpers.selectPostsByUser", database.Log(qPost))

	var posts []postDB
	if err := p.db.SelectContext(ctx, &posts, qPost, userID); err != nil {
		return nil, errors.Wrap(err, "selecting users posts")
	}
	for i := range posts {
		score, err := p.getPostScore(ctx, posts[i].ID)
		if err != nil {
			return nil, err
		}
		posts[i].Score = score
	}
	return posts, nil
}

//getPostByID obtains post from DB using its ID
func (p Post) getPostByID(ctx context.Context, postID string) (postDB, error) {
	const q = `	SELECT * FROM posts WHERE post_id = $1`

	p.log.Printf("%s: %s", "post.helpers.getPostByID", database.Log(q, postID))

	var post postDB
	if err := p.db.GetContext(ctx, &post, q, postID); err != nil {
		if err == sql.ErrNoRows {
			return postDB{}, ErrPostNotFound
		}
		return postDB{}, errors.Wrap(err, "selecting post by ID")
	}
	score, err := p.getPostScore(ctx, post.ID)
	if err != nil {
		return postDB{}, err
	}
	post.Score = score
	return post, nil
}

//checkPost shows whether post with given ID exist in DB or not
//it returns an error if post does not exist or nil if does
func (p Post) checkPost(ctx context.Context, postID string) error {
	const qCheckExist = `SELECT COUNT(*) FROM posts WHERE post_id = $1`

	p.log.Printf("%s: %s", "post.helpers.checkPost", database.Log(qCheckExist))

	var exist int
	if err := p.db.GetContext(ctx, &exist, qCheckExist, postID); err != nil {
		return errors.Wrap(err, "checking if post exists")
	}

	if exist == 0 {
		return ErrPostNotFound
	}
	return nil
}

//insertPost adds one new row to posts DB
func (p Post) insertPost(ctx context.Context, post postDB) error {
	const qPost = `
	INSERT INTO posts
		(post_id, views, type, title, category, payload, date_created, user_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)`

	p.log.Printf("%s: %s", "post.helpers.insertPost",
		database.Log(qPost, post.ID, post.Views, post.Type, post.Title, post.Payload, post.Category,
			post.DateCreated, post.UserID),
	)

	if _, err := p.db.ExecContext(ctx, qPost, post.ID, post.Views, post.Type, post.Title,
		post.Category, post.Payload, post.DateCreated, post.UserID); err != nil {
		return errors.Wrap(err, "inserting post")
	}
	return nil
}

//insertVote adds one row to votes DB
func (p Post) insertVote(ctx context.Context, postID string, userID string, vote int) error {
	const qVote = `
	INSERT INTO votes
		(post_id, user_id, vote)
	VALUES
		($1, $2, $3)`

	p.log.Printf("%s: %s", "post.helpers.insertVote", database.Log(qVote, postID, userID, vote))

	if _, err := p.db.ExecContext(ctx, qVote, postID, userID, vote); err != nil {
		return errors.Wrap(err, "inserting Vote")
	}
	return nil
}

func (p Post) deletePost(ctx context.Context, postID string) error {

	const qDeletePost = `DELETE FROM posts WHERE post_id = $1`

	p.log.Printf("%s: %s", "post.helpers.deletePost", database.Log(qDeletePost, postID))

	if _, err := p.db.ExecContext(ctx, qDeletePost, postID); err != nil {
		return errors.Wrapf(err, "deleting post %s", postID)
	}
	return nil
}

//checkVote shows whether vote to given post by user exist in DB or not
//it returns an error if vote does not exist or nil if does
func (p Post) checkVote(ctx context.Context, postID string, userID string) error {
	const qCheckExist = `SELECT COUNT(*) FROM votes WHERE post_id = $1 AND user_id = $2`

	p.log.Printf("%s: %s", "post.helpers.checkVote", database.Log(qCheckExist))

	var exist int
	if err := p.db.GetContext(ctx, &exist, qCheckExist, postID, userID); err != nil {
		return errors.Wrap(err, "checking if vote exists")
	}

	if exist == 0 {
		return ErrPostNotFound
	}
	return nil
}

func (p Post) updateVote(ctx context.Context, postID string, userID string, vote int) error {
	const qUpdateVote = `UPDATE votes SET vote = $3 WHERE post_id = $1 AND user_id = $2`

	p.log.Printf("%s: %s", "post.helpers.updateVote", database.Log(qUpdateVote))

	if _, err := p.db.ExecContext(ctx, qUpdateVote, postID, userID, vote); err != nil {
		return errors.Wrap(err, "updating vote")
	}
	return nil
}

func (p Post) deleteVote(ctx context.Context, postID string, userID string) error {

	const qDeleteVote = `DELETE FROM votes WHERE post_id = $1 AND user_id = $2`

	p.log.Printf("%s: %s", "post.helpers.deleteVote", database.Log(qDeleteVote, postID, userID))

	if _, err := p.db.ExecContext(ctx, qDeleteVote, postID, userID); err != nil {
		return errors.Wrapf(err, "deleting vote on %s from %s", postID, userID)
	}
	return nil
}

func (p Post) insertComment(
	ctx context.Context, commentID string, postID string, userID string, text string, now time.Time) error {
	const qComment = `
	INSERT INTO comments
		(comment_id, post_id, user_id, body, date_created)
	VALUES
		($1, $2, $3, $4, $5)`

	p.log.Printf("%s: %s", "post.helpers.insertComment",
		database.Log(qComment, commentID, postID, userID, text, now))

	if _, err := p.db.ExecContext(ctx, qComment, commentID, postID, userID, text, now); err != nil {
		return errors.Wrap(err, "inserting Comment")
	}
	return nil
}

func (p Post) getCommentByID(ctx context.Context, commentID string) (Comment, error) {
	const qComment = `
		SELECT name, user_id, cm.date_created, body, comment_id 
		FROM comments cm join users using(user_id) 
		WHERE comment_id = $1`

	p.log.Printf("%s: %s", "post.helpers.getCommentByID", database.Log(qComment, commentID))

	var rawComment struct {
		DateCreated time.Time `db:"date_created"`
		AuthorName  string    `db:"name"`
		AuthorID    string    `db:"user_id"`
		Body        string    `db:"body"`
		ID          string    `db:"comment_id"`
	}
	if err := p.db.GetContext(ctx, &rawComment, qComment, commentID); err != nil {
		if err == sql.ErrNoRows {
			return Comment{}, ErrCommentNotFound
		}
		return Comment{}, errors.Wrap(err, "selecting comment by ID")
	}

	author := Author{
		Username: rawComment.AuthorName,
		ID:       rawComment.AuthorID,
	}
	comment := Comment{
		DateCreated: rawComment.DateCreated,
		Author:      author,
		Body:        rawComment.Body,
		ID:          rawComment.ID,
	}
	return comment, nil
}

func (p Post) deleteComment(ctx context.Context, commentID string) error {
	const qDeleteComment = `DELETE FROM comments WHERE comment_id = $1`

	p.log.Printf("%s: %s", "post.helpers.deleteComment", database.Log(qDeleteComment, commentID))

	if _, err := p.db.ExecContext(ctx, qDeleteComment, commentID); err != nil {
		if err == sql.ErrNoRows {
			return ErrCommentNotFound
		}
		return errors.Wrapf(err, "deleting comment %s", commentID)
	}
	return nil
}
