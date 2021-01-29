package mid

import (
	"context"
	"net/http"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"authHeader"}

type contextKey struct {
	name string
}

// Auth decodes the share session cookie and packs the session into context
func Auth() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			c := r.Header.Get("authorization")
			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, c)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next(w, r)
		}
	}
}

// GetAuthString finds the user from the context. REQUIRES Middleware to have run.
func GetAuthString(ctx context.Context) string {
	raw, _ := ctx.Value(userCtxKey).(string)
	return raw
}
