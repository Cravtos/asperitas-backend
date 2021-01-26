package graph

import (
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/posts"
	"github.com/cravtos/asperitas-backend/graph/model"
)

func preparePostsToSend(postsRaw []posts.Info) []model.Info {
	var posts []model.Info
	for _, postRaw := range postsRaw {
		posts = append(posts, preparePostToSend(postRaw))
	}
	return posts
}

func preparePostToSend(postRaw posts.Info) model.Info {
	var info model.Info
	if postRaw.Type == "link" || postRaw.Type == "url" {
		info = model.PostLink{
			Type:             "link",
			PostID:           postRaw.ID,
			Score:            postRaw.Score,
			Views:            postRaw.Views,
			Title:            postRaw.Title,
			URL:              postRaw.Payload,
			Category:         model.Category(postRaw.Category),
			DateCreated:      postRaw.DateCreated,
			Author:           prepareAuthor(postRaw.Author),
			Votes:            prepareVotes(postRaw.Votes),
			Comments:         prepareComments(postRaw.Comments),
			UpvotePercentage: postRaw.UpvotePercentage,
		}
	} else {
		info = model.PostText{
			Type:             "text",
			PostID:           postRaw.ID,
			Score:            postRaw.Score,
			Views:            postRaw.Views,
			Title:            postRaw.Title,
			Text:             postRaw.Payload,
			Category:         model.Category(postRaw.Category),
			DateCreated:      postRaw.DateCreated,
			Author:           prepareAuthor(postRaw.Author),
			Votes:            prepareVotes(postRaw.Votes),
			Comments:         prepareComments(postRaw.Comments),
			UpvotePercentage: postRaw.UpvotePercentage,
		}
	}
	return info
}

func prepareComments(commentsRaw []posts.Comment) []*model.Comment {
	comments := make([]*model.Comment, 0)
	for _, raw := range commentsRaw {
		comments = append(comments, prepareComment(raw))
	}
	return comments
}

func prepareComment(raw posts.Comment) *model.Comment {
	return &model.Comment{
		DateCreated: raw.DateCreated,
		Author:      prepareAuthor(&raw.Author),
		Body:        raw.Body,
		CommentID:   raw.ID,
	}
}

func prepareVotes(votesRaw []posts.Vote) []*model.Vote {
	votes := make([]*model.Vote, 0)
	for _, voteRaw := range votesRaw {
		votes = append(votes, prepareVote(voteRaw))
	}
	return votes
}

func prepareVote(raw posts.Vote) *model.Vote {
	return &model.Vote{
		AuthorID: raw.UserID,
		Vote:     raw.Vote,
	}
}

func prepareAuthor(author *posts.Author) *model.Author {
	return &model.Author{
		Username: author.Username,
		AuthorID: author.ID,
	}
}

func prepareAuthData(data auth.Data) *model.AuthData {
	return &model.AuthData{
		Token: data.Token,
		User:  prepareUser(data.User),
	}
}

func prepareUser(data auth.User) *model.User {
	return &model.User{
		Username: data.Username,
		UserID:   data.ID,
	}
}
