// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/cravtos/asperitas-backend/business/auth" // Import is removed in final PR
	"github.com/cravtos/asperitas-backend/business/data/post"
	"github.com/cravtos/asperitas-backend/business/data/user"
	"github.com/cravtos/asperitas-backend/business/mid"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Panics(log))

	// Register debug check endpoints.
	cg := checkGroup{
		build: build,
		db:    db,
	}
	app.HandleDebug(http.MethodGet, "/readiness", cg.readiness)
	app.HandleDebug(http.MethodGet, "/liveness", cg.liveness)

	ug := userGroup {
		user: user.New(log, db),
		auth: a,
	}

	app.Handle(http.MethodPost, "/api/register", ug.register)
	app.Handle(http.MethodPost, "/api/login", ug.login)

	//Register post and post endpoints
	p := postGroup {
		post: post.New(log, db),
	}

	app.Handle(http.MethodGet, "/api/posts/", p.query)
	app.Handle(http.MethodGet, "/api/posts/:category", p.queryByCat)
	app.Handle(http.MethodGet, "/api/post/:post_id", p.queryByID)
	//todo: app.Handle(http.MethodPost, "/api/post/:post_id", p.createComment, mid.Authenticate(a))
	//todo: app.Handle(http.MethodDelete, "/api/posts/:post_id/:comment_id", p.deleteComment, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/api/posts/", p.create, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/api/posts/:post_id/upvote", p.upvote, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/api/posts/:post_id/downvote", p.downvote, mid.Authenticate(a))
	//todo: app.Handle(http.MethodGet, "/api/posts/:post_id/:comment_id", p.unvote, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/api/post/:post_id", p.delete, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/api/user/:user", p.queryByUser, mid.Authenticate(a))

	return app
}