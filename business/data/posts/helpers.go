package posts

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/db"
)

// getNewPostInfo creates Info using postDB and auth.Claims given by users
func getNewPostInfo(post db.PostDB, claims auth.Claims) Info {
	var info Info
	info = Info{
		Type:        "link",
		ID:          post.ID,
		Score:       1,
		Views:       post.Views,
		Title:       post.Title,
		Payload:     post.Payload,
		Category:    post.Category,
		DateCreated: post.DateCreated,
		Author: &Author{
			Username: claims.User.Username,
			ID:       claims.User.ID,
		},
		Votes: []Vote{
			{UserID: claims.User.ID, Vote: 1},
		},
		Comments:         []Comment{},
		UpvotePercentage: 100,
	}
	if post.Type == "text" {
		info.Type = post.Type
	}
	return info
}

// getInfo creates new Info using data from DB
func getInfo(post db.PostDB, author *Author, votes []Vote, comments []Comment) Info {
	var info Info
	info = Info{
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

	return info
}

// getEmptyInfo creates new Info filled only with db.posts data
func getEmptyInfo(post db.PostDB) Info {
	var info Info
	info = Info{
		Type:        "link",
		ID:          post.ID,
		Views:       post.Views,
		Title:       post.Title,
		Payload:     post.Payload,
		Category:    post.Category,
		DateCreated: post.DateCreated,
		UserID:      post.UserID,
	}
	if post.Type == "text" {
		info.Type = post.Type
	}
	return info
}

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
	return Vote{UserID: voteDB.UserID, Vote: voteDB.Vote}
}

//convertUser turns db.UserDB to Author
func convertUser(userDB db.UserDB) Author {
	return Author{Username: userDB.Username, ID: userDB.ID}
}

//fillInfos fills slice of Info with votes, comments and author for each Post
func (p Setup) fillInfos(ctx context.Context, posts []Info) ([]Info, error) {
	for i, post := range posts {
		filledPost, err := p.fillInfo(ctx, post)
		if err != nil {
			return nil, err
		}
		posts[i] = filledPost
	}
	return posts, nil
}

//fillInfo fills Info with votes, comments and author
func (p Setup) fillInfo(ctx context.Context, src Info) (Info, error) {
	dbs := db.NewDBset(p.log, p.db)

	//obtaining votes for postRes
	votesDB, err := dbs.SelectVotesByPostID(ctx, src.ID)
	if err != nil {
		return Info{}, err
	}
	votes := convertVotes(votesDB)
	src.Votes = votes
	//obtaining author for postRes
	userDB, err := dbs.GetUserByID(ctx, src.UserID)
	if err != nil {
		return Info{}, err
	}
	author := convertUser(userDB)
	src.Author = &author
	//obtaining comments for postRes
	commentsWithAuthorDB, err := dbs.SelectCommentsWithUserByPostID(ctx, src.ID)
	if err != nil {
		return Info{}, err
	}
	comments := convertComments(commentsWithAuthorDB)
	src.Comments = comments
	src.Score = score(votes)
	src.UpvotePercentage = upvotePercentage(votes)
	return src, nil
}
