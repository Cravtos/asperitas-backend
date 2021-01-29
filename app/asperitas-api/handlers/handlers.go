// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/cravtos/asperitas-backend/graph"
	"github.com/cravtos/asperitas-backend/graph/generated"

	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/mid"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown)

	// Construct GraphQL server
	config := generated.Config{Resolvers: &graph.Resolver{Log: log, DB: db, Auth: a}}
	es := generated.NewExecutableSchema(config)
	srv := handler.NewDefaultServer(es)
	errSetup := graph.NewSetup(log, shutdown)
	srv.SetErrorPresenter(errSetup.ErrorPresenter)

	// Register GraphQL endpoint
	app.Handle(http.MethodPost, "/api/graphql", srv.ServeHTTP, mid.Auth())

	return app
}
