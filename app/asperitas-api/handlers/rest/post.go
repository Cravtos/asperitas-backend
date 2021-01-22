package rest

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/posts"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/pkg/errors"
	"net/http"
)

type PostGroup struct {
	Post posts.Setup
}

func (pg PostGroup) Query(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	postsRaw, err := pg.Post.Query(ctx)
	if err != nil {
		return err
	}
	posts := preparePostsToSend(postsRaw)

	return web.Respond(ctx, w, posts, http.StatusOK)
}

func (pg PostGroup) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	pstRaw, err := pg.Post.QueryByID(ctx, params["post_id"])
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case posts.ErrPostNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["post_id"])
		}
	}

	pst := preparePostToSend(pstRaw)

	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg PostGroup) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var np posts.NewPost
	if err := web.Decode(r, &np); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	pstRaw, err := pg.Post.Create(ctx, claims, np, v.Now)
	if err != nil {
		switch err {
		case posts.ErrWrongPostType:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "creating new posts: %+v", np)
		}
	}

	pst := preparePostToSend(pstRaw)
	return web.Respond(ctx, w, pst, http.StatusCreated)
}

func (pg PostGroup) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	if _, err := pg.Post.Delete(ctx, claims, params["post_id"]); err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "ID: %s", params["post_id"])
		}
	}

	success := web.MessageResponse{Msg: "success"}
	return web.Respond(ctx, w, success, http.StatusOK)
}

func (pg PostGroup) QueryByCat(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	pstRaw, err := pg.Post.QueryByCat(ctx, params["category"])
	if err != nil {
		switch err {
		case posts.ErrPostNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["category"])
		}
	}
	pst := preparePostsToSend(pstRaw)
	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg PostGroup) QueryByUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := web.Params(r)
	pstRaw, err := pg.Post.QueryByUser(ctx, params["user"])
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case posts.ErrPostNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["user"])
		}
	}
	pst := preparePostsToSend(pstRaw)
	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg PostGroup) Upvote(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	pstRaw, err := pg.Post.Vote(ctx, claims, params["post_id"], 1)
	if err != nil {
		switch err {
		case posts.ErrPostNotFound:
			return web.NewRequestError(posts.ErrPostNotFound, http.StatusBadRequest)
		case posts.ErrInvalidPostID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "upvoting post with ID: %s", params["post_id"])
		}
	}
	pst := preparePostToSend(pstRaw)
	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg PostGroup) Downvote(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	pstRaw, err := pg.Post.Vote(ctx, claims, params["post_id"], -1)
	if err != nil {
		switch err {
		case posts.ErrPostNotFound:
			return web.NewRequestError(posts.ErrPostNotFound, http.StatusBadRequest)
		case posts.ErrInvalidPostID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "downvoting post with ID: %s", params["post_id"])
		}
	}
	pst := preparePostToSend(pstRaw)
	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg PostGroup) Unvote(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	pstRaw, err := pg.Post.Unvote(ctx, claims, params["post_id"])
	if err != nil {
		switch err {
		case posts.ErrPostNotFound:
			return web.NewRequestError(posts.ErrPostNotFound, http.StatusBadRequest)
		case posts.ErrInvalidPostID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "unvoting posts with ID: %s", params["post_id"])
		}
	}
	pst := preparePostToSend(pstRaw)
	return web.Respond(ctx, w, pst, http.StatusOK)
}

func (pg PostGroup) CreateComment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var nc posts.NewComment
	if err := web.Decode(r, &nc); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	params := web.Params(r)

	pstRaw, err := pg.Post.CreateComment(ctx, claims, nc, params["post_id"], v.Now)
	if err != nil {
		switch err {
		case posts.ErrPostNotFound:
			return web.NewRequestError(posts.ErrPostNotFound, http.StatusBadRequest)
		case posts.ErrInvalidPostID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "creating new comment: %+v", nc)
		}
	}
	pst := preparePostToSend(pstRaw)
	return web.Respond(ctx, w, pst, http.StatusCreated)
}

func (pg PostGroup) DeleteComment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	pstRaw, err := pg.Post.DeleteComment(ctx, claims, params["post_id"], params["comment_id"])
	if err != nil {
		switch err {
		case posts.ErrCommentNotFound:
			return web.NewRequestError(err, http.StatusBadRequest)
		case posts.ErrPostNotFound:
			return web.NewRequestError(err, http.StatusBadRequest)
		case posts.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		case posts.ErrInvalidPostID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case posts.ErrInvalidCommentID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting comment with ID: %s", params["post_id"])
		}
	}
	pst := preparePostToSend(pstRaw)
	return web.Respond(ctx, w, pst, http.StatusOK)
}
