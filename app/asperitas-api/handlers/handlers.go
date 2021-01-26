// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/cravtos/asperitas-backend/app/asperitas-api/handlers/rest"
	"github.com/cravtos/asperitas-backend/graph"
	"github.com/cravtos/asperitas-backend/graph/generated"
	"github.com/graphql-go/graphql"
	"log"
	"net/http"
	"os"

	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/posts"
	"github.com/cravtos/asperitas-backend/business/data/users"
	"github.com/cravtos/asperitas-backend/business/mid"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(
	build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB, gqlschema graphql.Schema,
) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Panics(log))

	// Register debug check endpoints.
	cg := rest.CheckGroup{
		Build: build,
		Db:    db,
	}
	app.HandleDebug(http.MethodGet, "/readiness", cg.Readiness)
	app.HandleDebug(http.MethodGet, "/liveness", cg.Liveness)

	// Register users endpoints
	ug := rest.UserGroup{
		User: users.New(log, db),
		Auth: a,
	}

	app.Handle(http.MethodPost, "/api/register", ug.Register)
	app.Handle(http.MethodPost, "/api/login", ug.Login)

	// Register posts endpoints
	pg := rest.PostGroup{
		Post: posts.New(log, db),
	}

	app.Handle(http.MethodGet, "/api/posts/", pg.Query)
	app.Handle(http.MethodGet, "/api/posts/:category", pg.QueryByCat)
	app.Handle(http.MethodGet, "/api/post/:post_id", pg.QueryByID)
	app.Handle(http.MethodGet, "/api/user/:user", pg.QueryByUser)
	app.Handle(http.MethodPost, "/api/posts", pg.Create, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/api/post/:post_id", pg.Delete, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/api/post/:post_id", pg.CreateComment, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/api/post/:post_id/:comment_id", pg.DeleteComment, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/api/posts/:post_id/upvote", pg.Upvote, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/api/posts/:post_id/downvote", pg.Downvote, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/api/posts/:post_id/unvote", pg.Unvote, mid.Authenticate(a))

	// Register endpoints for CORS
	cog := rest.CorsGroup{
		Log: log,
	}

	app.Handle(http.MethodOptions, "/api/register", cog.Allow("POST"))
	app.Handle(http.MethodOptions, "/api/login", cog.Allow("POST"))
	app.Handle(http.MethodOptions, "/api/posts", cog.Allow("POST"))
	app.Handle(http.MethodOptions, "/api/post/:post_id", cog.Allow("POST", "DELETE"))
	app.Handle(http.MethodOptions, "/api/post/:post_id/:comment_id", cog.Allow("DELETE"))
	app.Handle(http.MethodOptions, "/api/post/:post_id/upvote", cog.Allow("GET"))
	app.Handle(http.MethodOptions, "/api/post/:post_id/downvote", cog.Allow("GET"))
	app.Handle(http.MethodOptions, "/api/post/:post_id/unvote", cog.Allow("GET"))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Log: log, DB: db, Auth: a}}))

	app.HandleGraphQL(http.MethodPost, "/api/graphql", srv.ServeHTTP, mid.GQLAuth())

	return app
}
