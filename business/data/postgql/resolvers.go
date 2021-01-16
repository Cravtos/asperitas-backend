package postgql

import (
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func Hello(p graphql.ResolveParams) (interface{}, error) {
	return "World", nil
}

func postTitle(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Title, nil
}

func postType(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	if src.Type == "url" {
		return "link", nil
	}
	return src.Type, nil
}

func anyPost(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(Key).(PostGQL)
	if !ok {
		return nil, errors.New("postGQL missing from context")
	}
	posts, err := a.selectAllPosts(p.Context)
	if err != nil {
		return nil, err
	}

	author, err := a.getAuthorByID(p.Context, posts[0].UserID)
	if err != nil {
		return nil, err
	}

	votes, err := a.selectVotesByPostID(p.Context, posts[0].ID)
	if err != nil {
		return nil, err
	}

	comments, err := a.selectCommentsByPostID(p.Context, posts[0].ID)
	if err != nil {
		return nil, err
	}
	posts[0].Author = author
	posts[0].Votes = votes
	posts[0].Comments = comments
	return posts[0], nil
}

func allPosts(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(Key).(PostGQL)
	if !ok {
		return nil, errors.New("postGQL missing from context")
	}
	posts, err := a.selectAllPosts(p.Context)
	if err != nil {
		return nil, err
	}

	for i, post := range posts {
		author, err := a.getAuthorByID(p.Context, post.UserID)
		if err != nil {
			return nil, err
		}

		votes, err := a.selectVotesByPostID(p.Context, post.ID)
		if err != nil {
			return nil, err
		}

		comments, err := a.selectCommentsByPostID(p.Context, post.ID)
		if err != nil {
			return nil, err
		}
		posts[i].Author = author
		posts[i].Votes = votes
		posts[i].Comments = comments
	}
	return posts, nil
}

func postURL(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	if src.Type != "url" {
		return nil, errors.New("provided post is not link post")
	}
	return src.Payload, nil
}

func postText(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	if src.Type != "text" {
		return nil, errors.New("provided post is not text post")
	}
	return src.Payload, nil
}

func postID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.ID, nil
}

func postScore(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	score := 0
	for _, vote := range src.Votes {
		score += vote.Vote
	}
	return score, nil
}

func postViews(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Views, nil
}

func postCategory(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Category, nil
}

func postDateCreated(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.DateCreated, nil
}

func postAuthor(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Author, nil
}

func postUpvotePercentage(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	if len(src.Votes) == 0 {
		return 0, nil
	}

	var positive float32

	for _, vote := range src.Votes {
		if vote.Vote == 1 {
			positive++
		}
	}

	return int(positive / float32(len(src.Votes)) * 100), nil
}

func authorUsername(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Author)
	if !ok {
		return nil, errors.New("author missing from context")
	}
	return src.Username, nil
}

func authorID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Author)
	if !ok {
		return nil, errors.New("author missing from context")
	}
	return src.ID, nil
}

func voteVote(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Vote)
	if !ok {
		return nil, errors.New("vote missing from context")
	}
	return src.Vote, nil
}

func voteUserID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Vote)
	if !ok {
		return nil, errors.New("vote missing from context")
	}
	return src.User, nil
}

func commentID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, errors.New("vote missing from context")
	}
	return src.ID, nil
}

func commentBody(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, errors.New("vote missing from context")
	}
	return src.Body, nil
}

func commentAuthor(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, errors.New("vote missing from context")
	}
	return src.Author, nil
}

func commentDateCreated(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, errors.New("vote missing from context")
	}
	return src.DateCreated, nil
}

func postVotes(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Votes, nil
}

func postComments(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Comments, nil
}
