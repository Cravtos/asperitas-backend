package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/posts"
	"github.com/cravtos/asperitas-backend/business/data/users"
	"github.com/cravtos/asperitas-backend/business/mid"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/cravtos/asperitas-backend/graph/generated"
	"github.com/cravtos/asperitas-backend/graph/model"
	errs "github.com/pkg/errors"
)

func (r *authorResolver) Posts(ctx context.Context, obj *model.Author, category *model.Category) ([]model.Info, error) {
	ps := posts.New(r.Log, r.DB)
	cat := "all"
	if category != nil {
		cat = category.String()
	}

	infos, err := ps.ObtainPosts(ctx, cat, obj.AuthorID)
	if err != nil {
		return nil, newPrivateError(err)
	}

	return preparePostsToSend(infos), nil
}

func (r *commentResolver) Post(ctx context.Context, obj *model.Comment) (model.Info, error) {
	ps := posts.New(r.Log, r.DB)

	info, err := ps.QueryByID(ctx, obj.PostID)
	if err != nil {
		return nil, newPrivateError(err)
	}

	return preparePostToSend(info), nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, typeArg model.PostType, title string, category model.Category, payload string) (model.Info, error) {
	claims, err := r.Auth.ValidateString(mid.GetAuthString(ctx))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(r.Log, r.DB)
	np := posts.NewPost{
		Title:    title,
		Type:     string(typeArg),
		Category: string(category),
		Text:     payload,
		URL:      payload,
	}
	v, ok := ctx.Value(web.KeyValues).(web.Values)
	if !ok {
		return nil, web.NewShutdownError("web value missing from context")
	}
	info, err := ps.Create(ctx, claims, np, v.Now)
	if err != nil {
		return nil, newPrivateError(err)
	}

	return preparePostToSend(info), nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, postID string) (model.Info, error) {
	ps := posts.New(r.Log, r.DB)

	claims, err := r.Auth.ValidateString(mid.GetAuthString(ctx))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	info, err := ps.Delete(ctx, claims, postID)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, newPublicError(err)
		case posts.ErrPostNotFound:
			return nil, newPublicError(err)
		case posts.ErrForbidden:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(err)
		}
	}
	return preparePostToSend(info), nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, postID string, text string) (model.Info, error) {
	claims, err := r.Auth.ValidateString(mid.GetAuthString(ctx))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	v, ok := ctx.Value(web.KeyValues).(web.Values)
	if !ok {
		return nil, web.NewShutdownError("web value missing from context")
	}

	ps := posts.New(r.Log, r.DB)

	nc := posts.NewComment{Text: text}
	info, err := ps.CreateComment(ctx, claims, nc, postID, v.Now)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, newPublicError(err)
		case posts.ErrPostNotFound:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(err)
		}
	}
	return preparePostToSend(info), nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, postID string, commentID string) (model.Info, error) {
	claims, err := r.Auth.ValidateString(mid.GetAuthString(ctx))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(r.Log, r.DB)

	info, err := ps.DeleteComment(ctx, claims, postID, commentID)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, newPublicError(err)
		case posts.ErrInvalidCommentID:
			return nil, newPublicError(err)
		case posts.ErrPostNotFound:
			return nil, newPublicError(err)
		case posts.ErrCommentNotFound:
			return nil, newPublicError(err)
		case posts.ErrForbidden:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(err)
		}
	}

	return preparePostToSend(info), nil
}

func (r *mutationResolver) Upvote(ctx context.Context, postID string) (model.Info, error) {
	claims, err := r.Auth.ValidateString(mid.GetAuthString(ctx))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(r.Log, r.DB)

	info, err := ps.Vote(ctx, claims, postID, 1)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(err)
		}
	}

	return preparePostToSend(info), nil
}

func (r *mutationResolver) Downvote(ctx context.Context, postID string) (model.Info, error) {
	claims, err := r.Auth.ValidateString(mid.GetAuthString(ctx))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(r.Log, r.DB)

	info, err := ps.Vote(ctx, claims, postID, -1)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(err)
		}
	}
	return preparePostToSend(info), nil
}

func (r *mutationResolver) Unvote(ctx context.Context, postID string) (model.Info, error) {
	claims, err := r.Auth.ValidateString(mid.GetAuthString(ctx))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(r.Log, r.DB)

	info, err := ps.Unvote(ctx, claims, postID)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(err)
		}
	}
	return preparePostToSend(info), nil
}

func (r *mutationResolver) Register(ctx context.Context, name string, password string) (*model.AuthData, error) {
	v, ok := ctx.Value(web.KeyValues).(web.Values)
	if !ok {
		return nil, web.NewShutdownError("web values missing from context")
	}

	us := users.New(r.Log, r.DB)
	nu := users.NewUser{
		Name:     name,
		Password: password,
	}

	_, err := us.Create(ctx, nu, v.Now)
	if err != nil {
		return nil, newPrivateError(errs.Wrapf(err, "unable to create users with name %s", nu.Name))
	}

	claims, err := us.Authenticate(ctx, nu.Name, nu.Password, v.Now)
	if err != nil {
		switch err {
		case users.ErrAuthenticationFailure:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(errs.Wrapf(err, "unable to authenticate users with name %s", nu.Name))
		}
	}

	kid := r.Auth.GetKID()
	Token, err := r.Auth.GenerateToken(kid, claims)
	if err != nil {
		return nil, newPrivateError(errs.Wrapf(err, "generating token"))
	}

	return &model.AuthData{
		Token: Token,
		User:  prepareUser(claims.User),
	}, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, name string, password string) (*model.AuthData, error) {
	v, ok := ctx.Value(web.KeyValues).(web.Values)
	if !ok {
		return nil, web.NewShutdownError("web values missing from context")
	}

	claims, err := users.New(r.Log, r.DB).Authenticate(ctx, name, password, v.Now)
	if err != nil {
		switch err {
		case users.ErrAuthenticationFailure:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(errs.Wrapf(err, "unable to authenticate users with name %s", name))
		}
	}

	kid := r.Auth.GetKID()
	Token, err := r.Auth.GenerateToken(kid, claims)
	if err != nil {
		return nil, newPrivateError(errs.Wrapf(err, "generating token"))
	}

	return &model.AuthData{
		Token: Token,
		User:  prepareUser(claims.User),
	}, nil
}

func (r *queryResolver) AnyPost(ctx context.Context) (model.Info, error) {
	ps := posts.New(r.Log, r.DB)

	infos, err := ps.Query(ctx)
	if err != nil {
		return nil, newPrivateError(err)
	}
	if len(infos) == 0 {
		return nil, newPublicError(errors.New("there is no infos at all"))
	}

	return preparePostToSend(infos[0]), nil
}

func (r *queryResolver) Posts(ctx context.Context, category *model.Category, userID *string) ([]model.Info, error) {
	ps := posts.New(r.Log, r.DB)
	cat := "all"
	us := ""
	if category != nil {
		cat = category.String()
	}
	if userID != nil {
		us = *userID
	}
	infos, err := ps.ObtainPosts(ctx, cat, us)
	if err != nil {
		return nil, newPrivateError(err)
	}
	return preparePostsToSend(infos), nil
}

func (r *queryResolver) Post(ctx context.Context, postID string) (model.Info, error) {
	ps := posts.New(r.Log, r.DB)

	info, err := ps.QueryByID(ctx, postID)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, newPublicError(err)
		case posts.ErrPostNotFound:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(err)
		}
	}
	return preparePostToSend(info), nil
}

func (r *voteResolver) Author(ctx context.Context, obj *model.Vote) (*model.Author, error) {
	ps := posts.New(r.Log, r.DB)

	author, err := ps.AuthorByID(ctx, obj.AuthorID)
	if err != nil {
		return nil, newPrivateError(err)
	}
	return prepareAuthor(&author), nil
}

// Author returns generated.AuthorResolver implementation.
func (r *Resolver) Author() generated.AuthorResolver { return &authorResolver{r} }

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Vote returns generated.VoteResolver implementation.
func (r *Resolver) Vote() generated.VoteResolver { return &voteResolver{r} }

type authorResolver struct{ *Resolver }
type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type voteResolver struct{ *Resolver }
