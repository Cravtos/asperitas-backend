package rest

import (
	"context"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"log"
	"net/http"
)

type CorsGroup struct {
	Log *log.Logger
}

// Allow responds with "Access-Control-Allow-Methods" set to methods and
// "Access-Control-Allow-Headers" set to "authorization, content-type"
func (cg CorsGroup) Allow(methods ...string) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return web.AllowMethods(ctx, w, methods)
	}
}
