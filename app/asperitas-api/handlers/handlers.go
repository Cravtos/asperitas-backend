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

	// Register user endpoints
	ug := userGroup{
		user: user.New(log, db),
		auth: a,
	}

	app.Handle(http.MethodPost, "/api/register", ug.register)
	app.Handle(http.MethodPost, "/api/login", ug.login)

	// Register post endpoints
	pg := postGroup{
		post: post.New(log, db),
	}

	app.Handle(http.MethodGet, "/api/posts/", pg.query)
	app.Handle(http.MethodGet, "/api/posts/:category", pg.queryByCat)
	app.Handle(http.MethodGet, "/api/post/:post_id", pg.queryByID)
	app.Handle(http.MethodGet, "/api/user/:user", pg.queryByUser)
	app.Handle(http.MethodPost, "/api/posts", pg.create, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/api/post/:post_id", pg.delete, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/api/post/:post_id", pg.createComment, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/api/post/:post_id/:comment_id", pg.deleteComment, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/api/post/:post_id/upvote", pg.upvote, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/api/post/:post_id/downvote", pg.downvote, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/api/post/:post_id/unvote", pg.unvote, mid.Authenticate(a))

	// Register endpoints for CORS
	cog := corsGroup{
		log: log,
	}

	app.Handle(http.MethodOptions, "/api/register", cog.allow("POST"))
	app.Handle(http.MethodOptions, "/api/login", cog.allow("POST"))
	app.Handle(http.MethodOptions, "/api/posts", cog.allow("POST")) // todo: smhw it redirects to /api/posts/. need fix
	app.Handle(http.MethodOptions, "/api/post/:post_id", cog.allow("POST", "DELETE"))
	app.Handle(http.MethodOptions, "/api/post/:post_id/:comment_id", cog.allow("DELETE"))
	app.Handle(http.MethodOptions, "/api/post/:post_id/upvote", cog.allow("GET"))
	app.Handle(http.MethodOptions, "/api/post/:post_id/downvote", cog.allow("GET"))
	app.Handle(http.MethodOptions, "/api/post/:post_id/unvote", cog.allow("GET"))

	return app
}