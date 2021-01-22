package rest

import (
	"github.com/cravtos/asperitas-backend/business/data/posts"
)

func preparePostsToSend(postsRaw []posts.Info) []Info {
	var posts []Info
	for _, postRaw := range postsRaw {
		posts = append(posts, preparePostToSend(postRaw))
	}
	return posts
}

func preparePostToSend(postRaw posts.Info) Info {
	var info Info
	if postRaw.Type == "link" || postRaw.Type == "url" {
		info = InfoLink{
			Type:             "link",
			ID:               postRaw.ID,
			Score:            postRaw.Score,
			Views:            postRaw.Views,
			Title:            postRaw.Title,
			Payload:          postRaw.Payload,
			Category:         postRaw.Category,
			DateCreated:      postRaw.DateCreated,
			Author:           prepareAuthor(postRaw.Author),
			Votes:            prepareVotes(postRaw.Votes),
			Comments:         prepareComments(postRaw.Comments),
			UpvotePercentage: postRaw.UpvotePercentage,
		}
	} else {
		info = InfoLink{
			Type:             "link",
			ID:               postRaw.ID,
			Score:            postRaw.Score,
			Views:            postRaw.Views,
			Title:            postRaw.Title,
			Payload:          postRaw.Payload,
			Category:         postRaw.Category,
			DateCreated:      postRaw.DateCreated,
			Author:           prepareAuthor(postRaw.Author),
			Votes:            prepareVotes(postRaw.Votes),
			Comments:         prepareComments(postRaw.Comments),
			UpvotePercentage: postRaw.UpvotePercentage,
		}
	}
	return info
}

func prepareComments(commentsRaw []posts.Comment) []Comment {
	comments := make([]Comment, 0)
	for _, raw := range commentsRaw {
		comments = append(comments, prepareComment(raw))
	}
	return comments
}

func prepareComment(raw posts.Comment) Comment {
	return Comment{
		DateCreated: raw.DateCreated,
		Author:      prepareAuthor(&raw.Author),
		Body:        raw.Body,
		ID:          raw.ID,
	}
}

func prepareVotes(votesRaw []posts.Vote) []Vote {
	votes := make([]Vote, 0)
	for _, voteRaw := range votesRaw {
		votes = append(votes, prepareVote(voteRaw))
	}
	return votes
}

func prepareVote(raw posts.Vote) Vote {
	return Vote{
		User: raw.UserID,
		Vote: raw.Vote,
	}
}

func prepareAuthor(author *posts.Author) Author {
	return Author{
		Username: author.Username,
		ID:       author.ID,
	}
}
