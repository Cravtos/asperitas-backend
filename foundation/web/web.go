// Package web contains a small web framework extension.
package web

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dimfeld/httptreemux/v5"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values are stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	Now time.Time
}

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	mux      *httptreemux.ContextMux
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	mux := httptreemux.NewContextMux()
	mux.RedirectTrailingSlash = false

	return &App{
		mux:      mux,
		shutdown: shutdown,
	}
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// HandleGraphQL sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) HandleGraphQL(method string, path string, handler http.HandlerFunc, mw ...GQLMiddleware) {
	a.handleGraphQL(method, path, handler, mw...)
}

func (a *App) handleGraphQL(method string, path string, h http.HandlerFunc, mw ...GQLMiddleware) {
	// First wrap handler specific middleware around this handler.
	h = wrapGQLMiddleware(mw, h)
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Set the context with the required values to
		// process the request.
		v := Values{
			Now: time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, v)

		// and call the next with our new context
		r = r.WithContext(ctx)
		// Call the wrapped handler functions.
		h(w, r)
	}
	a.mux.Handle(method, path, handler)
}
