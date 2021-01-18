package postgql

import (
	"github.com/cravtos/asperitas-backend/business/data/db"
)

//convertPosts turns slice of db.PostDB to slice of Info
func convertPosts(postsDB []db.PostDB) []Info {
	posts := make([]Info, 0)
	for _, post := range postsDB {
		posts = append(posts, convertPost(post))
	}
	return posts
}

//convertPost turns db.PostDB to Info
func convertPost(postDB db.PostDB) Info {
	return Info{
		ID:          postDB.ID,
		Views:       postDB.Views,
		Type:        postDB.Type,
		Title:       postDB.Title,
		Category:    postDB.Category,
		Payload:     postDB.Payload,
		DateCreated: postDB.DateCreated,
		UserID:      postDB.UserID,
		Author:      nil,
		Votes:       nil,
		Comments:    nil,
	}
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
	return Vote{UserID: voteDB.User, Vote: voteDB.Vote}
}

//convertUser turns db.UserDB to Author
func convertUser(userDB db.UserDB) Author {
	return Author{Username: userDB.Username, ID: userDB.ID}
}

func parseCatAndUser(category interface{}, userID interface{}) (string, string) {
	if category == nil && userID == nil {
		return "all", ""
	} else if category == nil {
		if user, ok := userID.(string); ok {
			return "all", user
		}
		return "all", ""
	} else if userID == nil {
		if cat, ok := category.(string); ok {
			return cat, ""
		}
		return "all", ""
	} else {
		if cat, ok := category.(string); ok {
			if user, ok := userID.(string); ok {
				return cat, user
			}
			return cat, ""
		}
		if user, ok := userID.(string); ok {
			return "all", user
		}
		return "all", ""
	}
}
