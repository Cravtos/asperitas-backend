package post

import (
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/db"
)

func upvotePercentage(votes []Vote) int {
	if len(votes) == 0 {
		return 0
	}

	var positive float32

	for _, vote := range votes {
		if vote.Vote == 1 {
			positive++
		}
	}

	return int(positive / float32(len(votes)) * 100)
}

func score(votes []Vote) int {
	score := 0
	for _, vote := range votes {
		score += vote.Vote
	}
	return score
}

// getNewPostInfo creates Info using postDB and auth.Claims given by user
func getNewPostInfo(post db.PostDB, claims auth.Claims) Info {
	var info Info
	if post.Type == "link" {
		info = InfoLink{
			Type:        "link",
			ID:          post.ID,
			Score:       1,
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
			Score:       1,
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

// getInfo creates new Info using data from DB
func getInfo(post db.PostDB, author Author, votes []Vote, comments []Comment) Info {
	var info Info
	if post.Type == "url" {
		info = InfoLink{
			Type:             "link",
			ID:               post.ID,
			Score:            score(votes),
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
			Score:            score(votes),
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

//convertComments turns slice of db.CommentWithUserDB to slice of Comments
func convertComments(commentsDB []db.CommentWithUserDB) []Comment {
	comments := make([]Comment, 0)
	for _, comment := range commentsDB {
		comments = append(comments, convertComment(comment))
	}
	return comments
}

//convertComment turns db.CommentWithUserDB to Comment
func convertComment(commentDB db.CommentWithUserDB) Comment {
	author := Author{
		Username: commentDB.AuthorName,
		ID:       commentDB.AuthorID,
	}
	comment := Comment{
		DateCreated: commentDB.DateCreated,
		Author:      author,
		Body:        commentDB.Body,
		ID:          commentDB.ID,
	}
	return comment
}

//convertVotes turns slice of db.VoteDB to slice of Vote
func convertVotes(votesDB []db.VoteDB) []Vote {
	votes := make([]Vote, 0)
	for _, vote := range votesDB {
		votes = append(votes, convertVote(vote))
	}
	return votes
}

//convertVote turns db.CommentWithUserDB to Comment
func convertVote(voteDB db.VoteDB) Vote {
	return Vote{User: voteDB.User, Vote: voteDB.Vote}
}

//convertUser turns db.UserDB to Author
func convertUser(userDB db.UserDB) Author {
	return Author{Username: userDB.Username, ID: userDB.ID}
}
