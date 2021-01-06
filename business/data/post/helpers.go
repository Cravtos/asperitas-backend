package post

import (
	"context"
	"database/sql"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/pkg/errors"
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

//creates new Info using PostDB and auth.Claims given by user
func infoByPostAndClaims(post PostDB, claims auth.Claims) Info {
	var info Info
	if post.Type == "url" {
		info = InfoLink{
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
func infoByDBdata(post PostDB, author Author, votes []Vote, comments []Comment) Info {

	var info Info
	if post.Type == "url" {
		info = InfoLink{
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

	var author Author
	if err := p.db.GetContext(ctx, &author, qAuthor, ID); err != nil {
		return Author{}, errors.Wrap(err, "selecting authors")
	}
	return author, nil
}

//returns slice of Vote for a single post
func (p Post) selectVotesByPostID(ctx context.Context, ID string) ([]Vote, error) {
	const qVotes = `SELECT user_id, vote FROM votes WHERE post_id = $1`

	p.log.Printf("%s: %s: %s", "post.Query",
		database.Log(qVotes),
	)

	var votes []Vote
	if err := p.db.SelectContext(ctx, &votes, qVotes, ID); err != nil {
		return nil, errors.Wrap(err, "selecting votes")
	}
	return votes, nil
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

	p.log.Printf("%s: %s: %s", "post.Query",
		database.Log(qComments),
	)

	var rawComments []CommentWithAuthor
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
	p.log.Printf("%s: %s: %s with %v", "post.Query",
		database.Log(qComments), comments,
	)
	return comments, nil
}

//return all posts stored in DB
func (p Post) selectAllPosts(ctx context.Context) ([]PostDB, error) {
	const qPost = `SELECT * FROM posts`

	p.log.Printf("%s: %s: %s", "post.Query",
		database.Log(qPost),
	)

	var posts []PostDB
	if err := p.db.SelectContext(ctx, &posts, qPost); err != nil {
		return nil, errors.Wrap(err, "selecting posts")
	}
	return posts, nil
}

//getPostByID obtains post from DB using its ID
func (p Post) getPostByID(ctx context.Context, postID string) (PostDB, error) {
	const q = `
	SELECT * FROM
		posts
	WHERE
		post_id = $1`

	p.log.Printf("%s: %s: %s", "product.QueryByID",
		database.Log(q, postID),
	)

	var post PostDB
	if err := p.db.GetContext(ctx, &post, q, postID); err != nil {
		if err == sql.ErrNoRows {
			return PostDB{}, ErrNotFound
		}
		return PostDB{}, errors.Wrap(err, "selecting post by ID")
	}
	return post, nil
}

//insertPost adds one new row to posts DB
func (p Post) insertPost(ctx context.Context, post PostDB) error {
	const qPost = `
	INSERT INTO posts
		(post_id, score, views, type, title, category, payload, date_created, user_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	p.log.Printf("%s: %s: %s", "post.Create",
		database.Log(qPost, post.ID, post.Score, post.Views, post.Type, post.Title, post.Payload, post.Category,
			post.DateCreated, post.UserID),
	)

	if _, err := p.db.ExecContext(ctx, qPost, post.ID, post.Score, post.Views, post.Type, post.Title,
		post.Category, post.Payload, post.DateCreated, post.UserID); err != nil {
		return errors.Wrap(err, "creating post")
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

	p.log.Printf("%s: %s: %s", "post.Create",
		database.Log(qVote, postID, userID, vote),
	)

	if _, err := p.db.ExecContext(ctx, qVote, postID, userID, vote); err != nil {
		return errors.Wrap(err, "upvote created post")
	}
	return nil
}

// todo: function to make InfoLink or InfoText from all fields and payload
