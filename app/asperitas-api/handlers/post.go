package handlers

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/post"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/pkg/errors"
	"net/http"
)

type postGroup struct {
	post post.Post
}

func (pg postGroup) query(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	posts, err := pg.post.Query(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, posts, http.StatusOK)
}

func (pg postGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	pst, err := pg.post.QueryByID(ctx, params["post_id"])
	if err != nil {
		switch err {
		case post.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case post.ErrPostNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["post_id"])
		}
	}

	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg postGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var np post.NewPost
	if err := web.Decode(r, &np); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	pst, err := pg.post.Create(ctx, claims, np, v.Now)
	if err != nil {
		switch err {
		case post.ErrWrongPostType:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "creating new post: %+v", np)
		}
	}

	return web.Respond(ctx, w, pst, http.StatusCreated)
}

func (pg postGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	if err := pg.post.Delete(ctx, claims, params["post_id"]); err != nil {
		switch err {
		case post.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "ID: %s", params["post_id"])
		}
	}

	success := web.MessageResponse{Msg: "success"}
	return web.Respond(ctx, w, success, http.StatusOK)
}

func (pg postGroup) queryByCat(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	pst, err := pg.post.QueryByCat(ctx, params["category"])
	if err != nil {
		switch err {
		case post.ErrPostNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["category"])
		}
	}

	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg postGroup) queryByUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	pst, err := pg.post.QueryByUser(ctx, params["user"])
	if err != nil {
		switch err {
		case post.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case post.ErrPostNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["user"])
		}
	}

	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg postGroup) upvote(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	pst, err := pg.post.Vote(ctx, claims, params["post_id"], 1)
	if err != nil {
		switch err {
		case post.ErrPostNotFound:
			return web.NewRequestError(post.ErrPostNotFound, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "upvoting post with ID: %s", params["post_id"])
		}
	}

	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg postGroup) downvote(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	pst, err := pg.post.Vote(ctx, claims, params["post_id"], -1)
	if err != nil {
		switch err {
		case post.ErrPostNotFound:
			return web.NewRequestError(post.ErrPostNotFound, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "downvoting post with ID: %s", params["post_id"])
		}
	}

	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg postGroup) unvote(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	pst, err := pg.post.Unvote(ctx, claims, params["post_id"])
	if err != nil {
		switch err {
		case post.ErrPostNotFound:
			return web.NewRequestError(post.ErrPostNotFound, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "unvoting post with ID: %s", params["post_id"])
		}
	}

	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg postGroup) createComment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var nc post.NewComment
	if err := web.Decode(r, &nc); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	params := web.Params(r)

	pst, err := pg.post.CreateComment(ctx, claims, nc, params["post_id"], v.Now)
	if err != nil {
		switch err {
		case post.ErrPostNotFound:
			return web.NewRequestError(post.ErrPostNotFound, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "creating new comment: %+v", nc)
		}
	}

	return web.Respond(ctx, w, pst, http.StatusCreated)
}

func (pg postGroup) deleteComment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	pst, err := pg.post.DeleteComment(ctx, claims, params["post_id"], params["comment_id"])
	if err != nil {
		switch err {
		case post.ErrCommentNotFound:
			return web.NewRequestError(post.ErrCommentNotFound, http.StatusBadRequest)
		case post.ErrPostNotFound:
			return web.NewRequestError(post.ErrPostNotFound, http.StatusBadRequest)
		case post.ErrForbidden:
			return web.NewRequestError(post.ErrForbidden, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "deleting comment with ID: %s", params["post_id"])
		}
	}

	return web.Respond(ctx, w, pst, http.StatusOK)
}
