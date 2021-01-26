package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/cravtos/asperitas-backend/graph/generated"
	"github.com/cravtos/asperitas-backend/graph/model"
)

//todo create resolvers

var inf = model.PostLink{
	PostID:           "id",
	Title:            "title",
	Type:             "type",
	Score:            0,
	Views:            0,
	Category:         model.CategoryFashion,
	DateCreated:      time.Time{},
	UpvotePercentage: 0,
	Author: &model.Author{
		Username: "author",
		AuthorID: "author_id",
	},
	Votes:    make([]*model.Vote, 0),
	Comments: make([]*model.Comment, 0),
	URL:      "URL",
}

var au = &model.AuthData{
	Token: "token",
	User: &model.User{
		Username: "username",
		UserID:   "user_id",
	},
}

func (r *mutationResolver) CreatePost(ctx context.Context, typeArg model.PostType, title string, category model.Category, payload string) (model.Info, error) {
	return inf, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, postID string) (model.Info, error) {
	return inf, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, postID string, text string) (model.Info, error) {
	return inf, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, postID string, commentID string) (model.Info, error) {
	return inf, nil
}

func (r *mutationResolver) Upvote(ctx context.Context, postID string) (model.Info, error) {
	return inf, nil
}

func (r *mutationResolver) Downvote(ctx context.Context, postID string) (model.Info, error) {
	return inf, nil
}

func (r *mutationResolver) Unvote(ctx context.Context, postID string) (model.Info, error) {
	return inf, nil
}

func (r *mutationResolver) Register(ctx context.Context, name string, password string) (*model.AuthData, error) {
	return au, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, name string, password string) (*model.AuthData, error) {
	return au, nil
}

func (r *queryResolver) AnyPost(ctx context.Context) (model.Info, error) {
	return inf, nil
}

func (r *queryResolver) Posts(ctx context.Context, category *model.Category, userID *string) ([]model.Info, error) {
	i := make([]model.Info, 0)
	i = append(i, inf)
	return i, nil
}

func (r *queryResolver) Post(ctx context.Context, postID string) (model.Info, error) {
	return inf, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
