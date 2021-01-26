package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cravtos/asperitas-backend/graph/generated"
	"github.com/cravtos/asperitas-backend/graph/model"
)

func (r *mutationResolver) CreatePost(ctx context.Context, typeArg model.PostType, title string, category model.Category, payload string) (model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePost(ctx context.Context, postID string) (model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateComment(ctx context.Context, postID string, text string) (model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteComment(ctx context.Context, postID string, commentID string) (model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Upvote(ctx context.Context, postID string) (model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Downvote(ctx context.Context, postID string) (model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Unvote(ctx context.Context, postID string) (model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Register(ctx context.Context, name string, password string) (*model.AuthData, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SignIn(ctx context.Context, name string, password string) (*model.AuthData, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AnyPost(ctx context.Context) (model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Posts(ctx context.Context, category *model.Category, userID *string) ([]model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Post(ctx context.Context, postID string) (model.Info, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
